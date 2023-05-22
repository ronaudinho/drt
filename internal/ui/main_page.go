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
	"github.com/ronaudinho/drt/internal/util"
)

type App struct {
	*tview.Application
	openDotaAPI *opendota.OpenDotaAPI
	replayAPI   *valve.Replay
}

func NewApp(
	app *tview.Application,
	openDotaAPI *opendota.OpenDotaAPI,
	replayAPI *valve.Replay,
) *App {
	return &App{
		Application: app,
		openDotaAPI: openDotaAPI,
		replayAPI:   replayAPI,
	}
}

func (a *App) MainPage(ctx context.Context) error {
	var err error

	layoutListDownload := tview.NewFlex().
		SetDirection(tview.FlexRow)

	txtDownloads := tview.NewTextView().
		SetText("Downloads").
		SetTextAlign(tview.AlignLeft)

	listDownload := tview.NewList()

	listReplay := a.getListDownload()

	for _, replay := range listReplay {
		// TODO:
		// add binding to `selected func` redirect to new page
		// which will display a replay
		listDownload.AddItem(replay, "", '-', nil)
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("DotA2 Replay")
	textView.SetBorder(true)

	var matchID string

	form := tview.NewForm().
		AddInputField("Match ID", "", 20, nil, func(text string) {
			matchID = text
		}).
		AddButton("Find", func() {
			defer func() {
				if err != nil {
					log.Println("Failed to find match")
					os.Exit(1)
				} else {
					// TODO:
					// redirect to new page
				}
			}()

			var matchIDInt int64

			matchIDInt, err = strconv.ParseInt(matchID, 10, 64)
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

			// TODO
			// If match id is already exists in local folder then
			// 	- binding to new page which will display a replay
			// otherwise
			// 	- download first then redirect to new page which will display a replay
		}).
		SetHorizontal(true)

	form.SetBorder(true).
		SetTitle("Find DotA2 Match ID").
		SetTitleAlign(tview.AlignCenter)

	layoutListDownload.
		AddItem(txtDownloads, 0, 1, false).
		AddItem(listDownload, 0, 2, false)
	layoutListDownload.SetBorder(true)

	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(form, 0, 2, false).
		AddItem(layoutListDownload, 0, 2, false)

	err = a.Application.SetRoot(rootLayout, true).SetFocus(form).Run()

	if err != nil {
		return err
	}

	return nil
}

func (App) getListDownload() []string {
	listReplay := make([]string, 0)

	files, err := os.ReadDir(valve.DefaultDestination)
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), util.DefaultReplayExtensionSuffix) {
			// NOTE
			// using `strings.Cut` so we not iterate to all string value
			replayName, _, _ := strings.Cut(file.Name(), util.DefaultReplayExtensionSuffix)
			listReplay = append(listReplay, replayName)
		}
	}

	return listReplay
}
