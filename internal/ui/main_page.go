package ui

import (
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
}

func NewApp(app *tview.Application) *App {
	return &App{
		Application: app,
	}
}

func (a *App) MainPage() error {
	layoutListDownload := tview.NewFlex().
		SetDirection(tview.FlexRow)

	txtDownloads := tview.NewTextView().
		SetText("Downloads").
		SetTextAlign(tview.AlignLeft)

	// TODO
	// Call function to get list of replay DotA2 from
	// local folder in `~/.config/drt/`
	//
	// And add binding to `selected func` redirect to new page
	// which will display a replay
	listDownload := tview.NewList().
		AddItem("123456781", "Replay 1", '1', nil).
		AddItem("123456781", "Replay 2", '2', nil)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("DotA2 Replay")
	textView.SetBorder(true)

	form := tview.NewForm().
		AddInputField("Match ID", "", 20, nil, nil).
		AddButton("Find", func() {
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

	err := a.Application.SetRoot(rootLayout, true).SetFocus(form).Run()

	if err != nil {
		return err
	}

	return nil
}
