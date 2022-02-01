package ui

import (
	"github.com/auhau/loggy/store"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func handlePatternInput(key tcell.Key) {
	if key == tcell.KeyEsc {
		// User canceled entering the pattern string
		// lets revert the filter to original string.
		patternInput.SetText(pattern)
	} else {
		err := store.SetParsePattern(patternInput.GetText())

		if err != nil {
			ShowError(err)

			// There was an error with the new parsing pattern string
			// lets try to use the old one to populate the screen.
			// We gonna ignore any more errors though.
			store.SetParsePattern(pattern)
			patternInput.SetText(pattern)
			layout.RemoveItem(patternInput)
			return
		}

		pattern = patternInput.GetText()

		// Lets apply the filter to the new parsing pattern
		handleFilterInput(tcell.KeyEnter)
	}

	layout.RemoveItem(patternInput)
	app.SetFocus(logsView)
}

func makePatternInput(bootstrappingPattern string) *tview.InputField {
	pattern = bootstrappingPattern
	err := store.SetParsePattern(bootstrappingPattern)

	if err != nil {
		ShowError(err)
	}

	return tview.NewInputField().
		SetText(bootstrappingPattern).
		SetLabel("Parsing pattern: ").
		SetDoneFunc(handlePatternInput)
}
