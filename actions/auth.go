package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/pop/nulls"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitter"
	"github.com/nicomo/kumano/models"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/twitter/callback")),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
	)
}

// AuthCallback manages callbacks from Authentication providers
func AuthCallback(c buffalo.Context) error {
	gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	// check user already exists or not
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("provider = ? and provider_id = ?", gothUser.Provider, gothUser.UserID)
	exists, err := q.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}

	// provision empty user
	u := &models.User{}

	// just login in: populate user from DB
	if exists {
		if err = q.First(u); err != nil {
			return errors.WithStack(err)
		}

		// user logged in
		// minus 1 point for days not logged in + 1 for logging in today
		diff := int(time.Since(u.LastLoggedAt) / time.Minute)
		u.Score += -(diff / 1440) + models.PointsLogsIn
		u.LastLoggedAt = time.Now()

		verrs, err := tx.ValidateAndUpdate(u)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template
			c.Set("errors", verrs)
			fmt.Printf("\nVerrs: %v\n", verrs)
			return c.Redirect(302, "/")
		}

		// set session user to logged in user and redirect to home

		// FIXME: either user current_user_id or current_user
		// currently doing both in different places

		c.Session().Set("current_user_id", u.ID)
		if err = c.Session().Save(); err != nil {
			return errors.WithStack(err)
		}

		return c.Redirect(302, "/")

	}

	// Signing up
	// validate user account from invitation
	if cUserID := c.Session().Get("current_user_id"); cUserID != nil {

		err = tx.Find(u, cUserID)
		if err != nil {
			// FIXME: couldn't find user, manage error
			fmt.Println(err)
		}

		// populate user from oauth info
		u.Name = nulls.NewString(gothUser.Name)
		u.Provider = nulls.NewString(gothUser.Provider)
		u.ProviderID = nulls.NewString(gothUser.UserID)

		//		u.Provider = ToNullString(gothUser.Provider)
		//		u.ProviderID = ToNullString(gothUser.UserID)
		u.AvatarURL = nulls.NewString(gothUser.AvatarURL)

		// retrieve nickname and check if it's unique
		// generate a random one if it's not
		nick := models.NickValidate(gothUser.NickName, tx)
		u.Nickname = nulls.NewString(nick)
		if nick != gothUser.NickName {
			mssg := fmt.Sprintf("@%s was already taken, we used @%s. Hope you like it. You can change it in your profile.", gothUser.NickName, nick)
			c.Flash().Add("success", mssg)
		}

		// new user gets points on creation
		u.Score += models.PointsCreatesAccount
		u.SignedUpAt = time.Now()
		u.LastLoggedAt = time.Now()

		// clean up invitation token from account
		u.InvitationToken = ""

		//FIXME: duplicate key value violates unique constraint "users_email_idx"
		// u.Email = ""

		err = tx.Save(u)
		if err != nil {
			return errors.WithStack(err)
		}

		return c.Redirect(302, "/")

	}

	c.Flash().Add("danger", T.Translate(c, "auth.callback.failure"))
	return c.Redirect(302, "/")

}

// AuthDestroy logs the user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", T.Translate(c, "auth.destroy.success"))
	return c.Redirect(302, "/")
}

// InvitationRedeem gives access to the
func InvitationRedeem(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	user := models.User{}
	q := tx.Where("invitation_token = ?", c.Param("invitation_token"))
	if err := q.First(&user); err != nil {
		// Either user already has working account (invitation redeemed)
		// or no invitation at all
		// either way, redirect to home with message
		c.Flash().Add("danger", T.Translate(c, "auth.invitation.failure"))
		return c.Render(403, r.HTML("/"))
	}

	// set current user in session
	c.Session().Set("current_user_id", user.ID)
	if err := c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("users/signup"))
}
