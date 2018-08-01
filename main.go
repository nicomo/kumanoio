package main

import (
	"log"

	"github.com/nicomo/kumano/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
