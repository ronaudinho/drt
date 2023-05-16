package main

import (
	"log"

	"github.com/rivo/tview"
	"github.com/ronaudinho/drt/internal/ui"
)

func main() {
	application := tview.NewApplication()

	app := ui.NewApp(application)

	if err := app.MainPage(); err != nil {
		log.Fatal(err)
	}
}
