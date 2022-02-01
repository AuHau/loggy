package ui

import (
	"github.com/rivo/tview"
)
import "github.com/chonla/format"

// Version string set in release time
var Version string = ""

var KEY_MAP = map[string]interface{}{
	"filter":       string(FILTER_KEY),
	"pattern":      string(PATTERN_KEY),
	"follow":       string(FOLLOW_KEY),
	"toggleFilter": string(TOGGLE_FILTER_KEY),
	"help":         string(HELP_KEY),
	"version":      Version,
}

var HELP_TEXT = format.Sprintf(`loggy %<version>s

loggy is a swiss-knife for analyzing and reading logs. It parses the logs using a pattern that you define and allows you to filter the logs based on the parsed values.

Following keys are supported:
 - "%<filter>s" for setting filter
 - "%<toggleFilter>s" for toggling filter
 - "%<pattern>s" for setting parsing pattern input
 - "%<follow>s" for scroll to bottom and follow new data
 - "%<help>s" for displaying help

Navigation:
 - "j", "k" or arrow keys for scrolling by one line 
 - "g", "G" to move to top / bottom
 - "Ctrl-F", "page down" to move down by one page
 - "Ctrl-B", "page up" to move up by one page
`, KEY_MAP)

func makeHelpModal() *tview.Modal {
	return tview.NewModal().
		SetText(HELP_TEXT).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage(HELP_PAGE_NAME)
			app.SetFocus(logsView)
		})
}
