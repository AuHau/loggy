package ui

import (
	"github.com/auhau/loggy/store"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: Autocomplete the pattern's variables
// TODO: Possibility to save patterns (maybe also per directory?)

func handleFilterInput(key tcell.Key) {
	if key == tcell.KeyEsc {
		// User canceled entering the filter string
		// lets revert the filter to original string.
		filterInput.SetText(filter)
	} else {
		isFilterOn = true
		logsView.Clear()
		err := store.Filter(filterInput.GetText())
		if err != nil {
			ShowError(err)

			// There was an error with the new filter string
			// lets try to use the old one to populate the screen.
			// We gonna ignore any more errors though.
			logsView.Clear()
			store.Filter(filter)
			filterInput.SetText(filter)
			layout.RemoveItem(filterInput)
			return
		}

		filter = filterInput.GetText()
	}

	layout.RemoveItem(filterInput)
	app.SetFocus(logsView)
}

func makeFilterInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Filter: ").
		SetDoneFunc(handleFilterInput)
}
