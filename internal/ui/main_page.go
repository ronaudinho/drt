package ui

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rivo/tview"

	"github.com/ronaudinho/drt/internal/client/opendota"
	"github.com/ronaudinho/drt/internal/client/valve"
)

type App struct {
	*tview.Application
	openDotaAPI *opendota.API
	replayAPI   *valve.Replay
}

func NewApp(openDotaAPI *opendota.API, replayAPI *valve.Replay) *App {
	return &App{
		Application: tview.NewApplication(),
		openDotaAPI: openDotaAPI,
		replayAPI:   replayAPI,
	}
}

func (a *App) MainPage(ctx context.Context) error {
	var err error

	title := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("dota 2 replay").
		SetBorder(true)

	var matchID string
	form := tview.NewForm().
		AddInputField("match id", "", 20, nil, func(text string) {
			matchID = text
		}).
		AddButton("find", func() {
			defer func() {
				if err != nil {
					log.Println("failed to find match")
					os.Exit(1)
				} else {
					// TODO: redirect to new page
				}
			}()

			matchIDInt, err := strconv.ParseInt(matchID, 10, 64)
			if err != nil {
				return
			}

			matchDetail, err := a.openDotaAPI.FetchMatchDetail(ctx, matchIDInt)
			if err != nil {
				return
			}

			err = a.replayAPI.Download(ctx, matchDetail.ReplayURL, valve.DefaultDestination)
			if err != nil {
				return
			}

			// TODO: If match id is already exists in local folder then
			// 	- binding to new page which will display a replay
			// otherwise
			// 	- download first then redirect to new page which will display a replay
		}).
		SetHorizontal(true)

	form.SetBorder(true).
		SetTitle("find dota 2 match id").
		SetTitleAlign(tview.AlignCenter)

	downloads := tview.NewFlex().
		SetDirection(tview.FlexRow)
	downloadsTitle := tview.NewTextView().
		SetText("downloads").
		SetTextAlign(tview.AlignLeft)
	downloadsList := tview.NewList()

	replays := a.getDownloadedReplays()
	for _, replay := range replays {
		// TODO: add binding to `selected func` redirect to new page
		// which will display a replay
		downloadsList.AddItem(replay, "", '-', nil)
	}

	downloads.
		AddItem(downloadsTitle, 0, 1, false).
		AddItem(downloadsList, 0, 2, false).
		SetBorder(true)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 0, 1, false).
		AddItem(form, 0, 2, false).
		AddItem(downloads, 0, 2, false)

	return a.Application.SetRoot(root, true).SetFocus(form).Run()
}

func (App) getDownloadedReplays() []string {
	var replays []string

	files, err := os.ReadDir(valve.DefaultDestination)
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), valve.DefaultReplayExtensionSuffix) {
			// NOTE
			// using `strings.Cut` so we not iterate to all string value
			replay, _, _ := strings.Cut(file.Name(), valve.DefaultReplayExtensionSuffix)
			replays = append(replays, replay)
		}
	}

	return replays
}
