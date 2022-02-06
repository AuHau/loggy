package ui

import (
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/gdamore/tcell/v2"
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

// Color theme
const (
	STATUS_BAR_BACKGROUND_COLOR = tcell.ColorGray

	NON_PATTERN_LINES_STATUS_BACKGROUND_COLOR = tcell.ColorDarkRed
	NON_PATTERN_LINES_STATUS_TEXT_COLOR       = tcell.ColorLightGray

	FILTER_STATUS_NONACTIVE_BACKGROUND_COLOR = tcell.ColorGray
	FILTER_STATUS_ACTIVE_BACKGROUND_COLOR    = tcell.ColorLimeGreen
	FILTER_STATUS_TEXT_COLOR                 = tcell.ColorLightGray

	INPUT_NAME_BACKGROUND_COLOR = tcell.ColorDarkSlateGray
	INPUT_NAME_TEXT_COLOR       = tcell.ColorLightGray

	FOLLOWING_STATUS_TEXT_COLOR                 = tcell.ColorLightGray
	FOLLOWING_STATUS_ACTIVE_BACKGROUND_COLOR    = tcell.ColorNavy
	FOLLOWING_STATUS_NONACTIVE_BACKGROUND_COLOR = tcell.ColorGray
)

// TODO: Coloring based on type of log (error | warning | debug | info)
// TODO: Follow should not be default (maybe only for STDIN?)
// TODO: List of pre-configured patterns available in the UI

// Bootstrap setup the tview App and bootstraps all its components
// It returns also io.Writer that is used to pass logs into the LogsView
func Bootstrap(stateStore *gredux.Store, bufferSize int) (*tview.Application, error) {
	app := tview.NewApplication()

	logsView := makeLogsView(bufferSize, stateStore)
	statusBar := makeStatusBar(stateStore)
	helpModal := makeHelpModal(stateStore)
	errorModal := makeErrorModal(stateStore)
	filterInput := makeFilterInput(stateStore)
	patternInput := makePatternInput(stateStore)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(statusBar, 1, 0, false).
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
			app.SetFocus(errorModal)
			return
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
			if layout.GetItemCount() == 2 {
				layout.AddItem(filterInput, 1, 0, true)
			}
			focusPrimitive = filterInput
		} else {
			layout.RemoveItem(filterInput)
		}

		if st.DisplayPatternInput {
			if layout.GetItemCount() == 2 {
				layout.AddItem(patternInput, 1, 0, true)
			}
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
		actions.ActionNameSetPattern,            // Can display errors
		actions.ActionNameToggleFilter,          // Can display errors
		actions.ActionNameToggleNonPatternLines, // Can display errors
	})

	app.SetRoot(pages, true)
	app.SetFocus(logsView)

	return app, nil
}
