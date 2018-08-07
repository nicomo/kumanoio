package mailers

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
)

// SendInvitation sends an invitation to the invited person
// when sponsor registers her email in /users/new
// called from actions/users.go Create
func SendInvitation(data map[string]string) error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "Invitation to Kumano"
	m.From = "nicolas.kumanoio@gmail.com"
	m.To = []string{data["emailTo"]}
	err := m.AddBody(r.HTML("invitation_send.html"), render.Data{
		"sponsorID":       data["sponsorID"],
		"invitationURL":   data["invitationURL"],
		"sponsorName":     data["sponsorName"],
		"sponsorNickname": data["sponsor.Nickname"],
	})
	if err != nil {
		return errors.WithStack(err)
	}

	err = smtp.Send(m)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
