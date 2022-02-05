package ui

import (
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/rivo/tview"
)

// Key bindings
const (
	SET_FILTER_KEY               = '/'
	TOGGLE_FILTER_KEY            = 'f'
	TOGGLE_NON_PATTERN_LINES_KEY = 'a'
	SET_PATTERN_KEY              = 'p'
	HELP_KEY                     = 'h'
)

// Pages names
const (
	MAIN_PAGE_NAME  = "main"
	HELP_PAGE_NAME  = "helpModal"
	ERROR_PAGE_NAME = "errorModal"
)

// TODO: Status bar that
//  - displays number of (filtered) lines
//  - displays if filter is applied / set or not
// 	- display number of lines that did not match the parsing pattern

// TODO: Coloring based on type of log (error | warning | debug | info)
// TODO: Follow should not be default (maybe only for STDIN?)
// TODO: List of pre-configured patterns available

// Bootstrap setup the tview App and bootstraps all its components
// It returns also io.Writer that is used to pass logs into the LogsView
func Bootstrap(stateStore *gredux.Store, bufferSize int) (*tview.Application, error) {
	app := tview.NewApplication()

	logsView := makeLogsView(bufferSize, stateStore)
	//statusBar := makeStatusBar(stateStore)
	helpModal := makeHelpModal(stateStore)
	errorModal := makeErrorModal(stateStore)
	filterInput := makeFilterInput(stateStore)
	patternInput := makePatternInput(stateStore)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		//AddItem(statusBar, 1, 0, false).
		AddItem(logsView, 0, 10, true)

	pages := tview.NewPages().
		AddPage(MAIN_PAGE_NAME, layout, true, true).
		AddPage(HELP_PAGE_NAME, helpModal, true, false).
		AddPage(ERROR_PAGE_NAME, errorModal, true, false)

	stateStore.AddHook(func(s gredux.State) {
		st := s.(state.State)
		var focusPrimitive tview.Primitive = logsView

		if st.DisplayError {
			pages.ShowPage(ERROR_PAGE_NAME)
			focusPrimitive = errorModal
		} else {
			pages.HidePage(ERROR_PAGE_NAME)
		}

		if st.DisplayHelp {
			pages.ShowPage(HELP_PAGE_NAME)
			focusPrimitive = helpModal
		} else {
			pages.HidePage(HELP_PAGE_NAME)
		}

		if st.DisplayFilterInput {
			layout.AddItem(filterInput, 1, 0, true)
			focusPrimitive = filterInput
		} else {
			layout.RemoveItem(filterInput)
		}

		if st.DisplayPatternInput {
			layout.AddItem(patternInput, 1, 0, true)
			focusPrimitive = patternInput
		} else {
			layout.RemoveItem(patternInput)
		}

		app.SetFocus(focusPrimitive)
	}, []string{
		actions.ActionNameHideError,
		actions.ActionNameAddLogLine,  // Can display errors
		actions.ActionNameDropLogLine, // Can display errors
		actions.ActionNameDisplayHelp,
		actions.ActionNameHideHelp,
		actions.ActionNameDisplayFilterInput,
		actions.ActionNameHideFilterInput,
		actions.ActionNameFilter, // Can display errors
		actions.ActionNameDisplayPatternInput,
		actions.ActionNameSetPattern,   // Can display errors
		actions.ActionNameToggleFilter, // Can display errors
	})

	stateStore.AfterUpdate(func(s gredux.State) {
		// The components should have set their state correctly using Hooks
		// now lets render the app
		//app.QueueUpdateDraw(func() {})
	})

	app.SetRoot(pages, true)
	app.SetFocus(logsView)

	return app, nil
}
