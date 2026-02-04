package main

import (
	"budget_tracker/actions"
	"log"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
