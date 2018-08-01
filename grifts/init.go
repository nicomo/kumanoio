package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/nicomo/kumano/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
