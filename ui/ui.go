package ui

import (
	"errors"
	"github.com/auhau/loggy/store"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
)

// Key bindings
const (
	FILTER_KEY        = '/'
	TOGGLE_FILTER_KEY = 'f'
	PATTERN_KEY       = 'p'
	FOLLOW_KEY        = 'b'
	HELP_KEY          = 'h'
)

// Pages names
const (
	MAIN_PAGE_NAME  = "main"
	HELP_PAGE_NAME  = "helpModal"
	ERROR_PAGE_NAME = "errorModal"
)

// UI components
var (
	app          *tview.Application
	layout       *tview.Flex
	logsView     *tview.TextView
	errorModal   *tview.Modal
	filterInput  *tview.InputField
	patternInput *tview.InputField
	helpModal    *tview.Modal
	pages        *tview.Pages
)

// State
var (
	pattern    string
	filter     string
	isFilterOn bool
)

// TODO: Status bar that
//  - displays when new data are present but you are not following
//  - displays when you are detached
//  - displays number of (filtered) lines
//  - displays file name / stdin
//  - displays if filter is applied / set or not
// 	- display number of lines that did and did not match the parsing pattern

// TODO: Coloring based on type of log (error | warning | debug | info)
// TODO: Follow should not be default (maybe only for STDIN?)

// handleLogsViewInput is main handler for key inputs as it serves the LogsView
// That is the reason why it is placed here and not in `view_logs.go` file.
func handleLogsViewInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case TOGGLE_FILTER_KEY:
		logsView.Clear()

		var err error
		if isFilterOn {
			err = store.Filter("")
			isFilterOn = false
		} else {
			err = store.Filter(filter)
			isFilterOn = true
		}

		if err != nil {
			ShowError(err)
		}
	case FILTER_KEY:
		layout.AddItem(filterInput, 1, 1, true)
		app.SetFocus(filterInput)
	case PATTERN_KEY:
		layout.AddItem(patternInput, 1, 1, true)
		app.SetFocus(patternInput)
	case FOLLOW_KEY:
		logsView.ScrollToEnd()
	case HELP_KEY:
		pages.ShowPage(HELP_PAGE_NAME)
		app.SetFocus(helpModal)
	}

	return event
}

// Bootstrap setup the tview App and bootstraps all its components
// It returns also io.Writer that is used to pass logs into the LogsView
func Bootstrap(bufferSize int, pattern string) (*tview.Application, io.Writer, error) {
	if app != nil {
		return nil, nil, errors.New("application initialized")
	}

	app = tview.NewApplication()

	logsView = makeLogsView(bufferSize)
	errorModal = makeErrorModal()
	helpModal = makeHelpModal()
	filterInput = makeFilterInput()
	patternInput = makePatternInput(pattern)

	layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(logsView, 0, 1, true)

	pages = tview.NewPages().
		AddPage(MAIN_PAGE_NAME, layout, true, true).
		AddPage(HELP_PAGE_NAME, helpModal, true, false).
		AddPage(ERROR_PAGE_NAME, errorModal, true, false)

	app.SetRoot(pages, true)
	app.SetFocus(logsView)

	return app, tview.ANSIWriter(logsView), nil
}
