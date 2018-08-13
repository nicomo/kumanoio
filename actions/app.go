package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/markbates/goth/gothic"
	"github.com/nicomo/kumano/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// T is used to generate translation string throughout the UI
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_kumano_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// setting the user in the session
		app.Use(SetCurrentUser)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		//ROUTING
		app.GET("/", HomeHandler)

		// authentication of users
		auth := app.Group("/auth")
		auth.GET("/invitation/{invitation_token}", InvitationRedeem)
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)
		auth.DELETE("", AuthDestroy)

		//
		// texts routes
		//
		// single pages, not linked to text model directly
		app.POST("/texts/{text_id}/star", StarHandler)

		// texts group routes
		tr := &TextsResource{}
		textsGroup := app.Group("/texts")
		textsGroup.Use(LoginRequired)
		textsGroup.Middleware.Skip(LoginRequired, tr.Show, tr.List)
		textsGroup.GET("/", tr.List)
		textsGroup.POST("/", tr.Create)
		textsGroup.GET("/new", tr.New)
		textsGroup.GET("/drafts", tr.ListDrafts)
		textsGroup.GET("/user/{user_id}", tr.ListUserTexts)
		textsGroup.GET("/{text_id}", tr.Show)
		textsGroup.GET("/{text_id}/edit", tr.Edit)
		textsGroup.PUT("/{text_id}", tr.Update)
		textsGroup.DELETE("/{text_id}", tr.Destroy)

		// users routes
		ur := &UsersResource{}
		usersGroup := app.Group("/users")
		usersGroup.Use(LoginRequired)
		usersGroup.Middleware.Skip(LoginRequired, ur.Show)
		usersGroup.GET("/", ur.List)                // GET /users => ur.List
		usersGroup.GET("/new", ur.New)              // GET /users/new => ur.New
		usersGroup.GET("/{user_id}", ur.Show)       // GET /users/{user_id} => ur.Show
		usersGroup.GET("/{user_id}/edit", ur.Edit)  // GET /users/{user_id}/edit => ur.Edit
		usersGroup.POST("/", ur.Create)             // POST /users => ur.Create
		usersGroup.PUT("/{user_id}", ur.Update)     // PUT /users/{user_id} => ur.Update
		usersGroup.DELETE("/{user_id}", ur.Destroy) //  DELETE /users/{user_id} => ur.Destroy

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
