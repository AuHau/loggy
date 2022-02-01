package ui

import "github.com/rivo/tview"

func makeLogsView(bufferSize int) *tview.TextView {
	view := tview.NewTextView().
		SetMaxLines(bufferSize).
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		}).
		ScrollToEnd() // This makes sure that we follow the end

	view.SetInputCapture(handleLogsViewInput)
	return view
}
