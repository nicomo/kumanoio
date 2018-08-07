package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/nicomo/kumano/models"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// uncomment for non-Bootstrap form helpers:
			// "form":     plush.FormHelper,
			// "form_for": plush.FormForHelper,
			"is_admin":   isAdmin,
			"is_self":    isSelf,
			"can_invite": canInvite,
		},
	})
}

func isAdmin(help plush.HelperContext) bool {
	if help.Value("current_user") != nil {
		return help.Value("current_user").(*models.User).IsAdmin
	}
	return false
}

// the user is looking at her own content (profile, text, etc)
func isSelf(help plush.HelperContext) bool {
	return help.Value("self").(bool)
}

func canInvite(help plush.HelperContext) bool {
	if help.Value("current_user") != nil {
		return help.Value("current_user").(*models.User).IsAdmin || help.Value("current_user").(*models.User).SponsorshipsCount > 0
	}
	return false
}
