package mailers

import (
	"log"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
)

var smtp mail.Sender
var r *render.Engine

func init() {

	// Pulling config from the env.
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	var err error
	sender, err := mail.NewSMTPSender(host, port, user, password)

	// FIXME: switch to TLS/SSL 
	// see https://support.google.com/accounts/answer/6010255
	// https://gobuffalo.io/en/docs/mail/
	// port 587 with TLS
	// sender.Dialer.TLSConfig = &tls.Config{...}
	sender.Dialer.SSL = true
	if err != nil {
		log.Fatal(err)
	}

	smtp = sender

	r = render.New(render.Options{
		HTMLLayout:   "layout.html",
		TemplatesBox: packr.NewBox("../templates/mail"),
		Helpers:      render.Helpers{},
	})
}
