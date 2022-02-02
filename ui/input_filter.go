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

			// There was an error in filter string
			// But redraw the screen only if filter is actually turned on
			if isFilterOn {
				logsView.Clear()
				store.Filter(filter)
			}

			// Lets revert the filter to previous one
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
