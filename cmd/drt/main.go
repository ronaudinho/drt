package main

import (
	"context"
	"log"
	"net/http"

	"github.com/rivo/tview"

	"github.com/ronaudinho/drt/internal/client/opendota"
	"github.com/ronaudinho/drt/internal/client/valve"
	"github.com/ronaudinho/drt/internal/ui"
)

func main() {
	ctx := context.Background()

	application := tview.NewApplication()

	httpDefault := http.DefaultClient

	openDotaAPI := opendota.NewMatchAPI(httpDefault, "https://api.opendota.com/api")
	replayAPI := valve.NewReplay(httpDefault)

	app := ui.NewApp(application, openDotaAPI, replayAPI)

	if err := app.MainPage(ctx); err != nil {
		log.Fatal(err)
	}
}
