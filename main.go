package main

import (
	"log"

	"github.com/rivo/tview"
	"github.com/ronaudinho/drt/internal/ui"
)

func main() {
	app := tview.NewApplication()

	err := ui.MainPage(app)

	if err != nil {
		log.Fatal(err)
	}
}
