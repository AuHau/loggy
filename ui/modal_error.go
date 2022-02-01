package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowError is public function that triggers the ErrorModal with given error
func ShowError(err error) {
	errorModal.SetText(fmt.Sprintf(`There was an error!

%s`, err))
	pages.ShowPage(ERROR_PAGE_NAME)
	app.SetFocus(errorModal)
}

func makeErrorModal() *tview.Modal {
	modal := tview.NewModal().
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage(ERROR_PAGE_NAME)
			app.SetFocus(logsView)
		})

	modal.SetBackgroundColor(tcell.ColorDarkRed)
	return modal
}
