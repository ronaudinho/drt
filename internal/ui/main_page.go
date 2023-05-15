package ui

import (
	"github.com/rivo/tview"
)

func MainPage(app *tview.Application) error {
	layoutListDownload := tview.NewFlex().
		SetDirection(tview.FlexRow)

	txtDownloads := tview.NewTextView().
		SetText("Downloads").
		SetTextAlign(tview.AlignLeft)

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
			app.SetFocus(listDownload)
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

	err := app.SetRoot(rootLayout, true).SetFocus(form).Run()

	if err != nil {
		return err
	}

	return nil
}
