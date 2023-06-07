package main

import (
	"context"
	"log"

	"github.com/ronaudinho/drt/internal/client/opendota"
	"github.com/ronaudinho/drt/internal/client/valve"
	"github.com/ronaudinho/drt/internal/ui"
)

func main() {
	ctx := context.Background()
	openDotaAPI := opendota.NewDefaultAPI()
	replayAPI := valve.NewDefaultReplay()

	app := ui.NewApp(openDotaAPI, replayAPI)
	if err := app.MainPage(ctx); err != nil {
		log.Fatal(err)
	}
}
