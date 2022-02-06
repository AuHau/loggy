package ui

import (
	"fmt"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/auhau/loggy/store"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func logsViewReducer(s gredux.State, action gredux.Action) gredux.State {
	st := s.(state.State)

	switch action.ID {
	case actions.ActionNameTurnOnFollowing:
		st.IsFollowing = true
		return st
	case actions.ActionNameTurnOffFollowing:
		st.IsFollowing = false
		return st
	case actions.ActionNameAddLogLine:
		line := action.Data.(string)
		var lineWithNL string

		if st.IsLogsFirstLine {
			lineWithNL = line
			st.IsLogsFirstLine = false
		} else {
			lineWithNL += "\n" + line
		}

		result, err := store.IsLineMatching(line, st.FilterExpression, st.ParsingPattern)
		if err != nil {
			st.DisplayError = true
			st.ErrorMessage = fmt.Sprint(err)
			return st
		}

		if st.FilterString != "" {
			if result == store.MATCH {
				st.TotalLines += 1
				st.MatchingLines += 1
				st.Logs += lineWithNL
			} else if result == store.FILTER_NO_MATCH {
				st.TotalLines += 1
			} else if result == store.PARSE_PATTERN_NO_MATCH {
				if st.DisplayNonPatternLines {
					st.Logs += lineWithNL
				}
				st.TotalLines += 1
				st.NonPatternLines += 1
			}
		} else {
			st.TotalLines += 1
			st.MatchingLines += 1

			if result == store.PARSE_PATTERN_NO_MATCH {
				st.NonPatternLines += 1
			}

			st.Logs += lineWithNL
		}
		return st
	case actions.ActionNameDropLogLine:
		line := action.Data.(string)
		// We gonna check if the line matches filter only there is some filter set
		if st.FilterString != "" {
			result, err := store.IsLineMatching(line, st.FilterExpression, st.ParsingPattern)
			if err != nil {
				st.DisplayError = true
				st.ErrorMessage = fmt.Sprint(err)
				return st
			}

			if result == store.MATCH {
				st.TotalLines -= 1
				st.MatchingLines -= 1
			} else if result == store.FILTER_NO_MATCH {
				st.TotalLines -= 1
			} else if result == store.PARSE_PATTERN_NO_MATCH {
				st.TotalLines -= 1
				st.NonPatternLines -= 1
			}
		} else {
			st.TotalLines -= 1
			st.MatchingLines -= 1
		}
		return st
	}

	return st
}

// TODO: Add indicator that there is a new log line that was not yet seen by the user
func makeLogsView(bufferSize int, stateStore *gredux.Store) *tview.TextView {
	view := tview.NewTextView().
		SetMaxLines(bufferSize).
		SetWrap(true).
		SetWordWrap(true).
		SetDynamicColors(true)

	isFollowing := stateStore.State().(state.State).IsFollowing
	if isFollowing {
		view.ScrollToEnd() // This makes sure that we follow the end if user requested
	}

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()

		switch key {
		case tcell.KeyRune:
			switch event.Rune() {
			case TOGGLE_FILTER_KEY:
				stateStore.Dispatch(actions.ToggleFilter())
			case SET_FILTER_KEY:
				stateStore.Dispatch(actions.DisplayFilterInput())
			case TOGGLE_NON_PATTERN_LINES_KEY:
				stateStore.Dispatch(actions.ToggleNonPatternLines())
			case SET_PATTERN_KEY:
				stateStore.Dispatch(actions.DisplayPatternInput())
			case HELP_KEY:
				stateStore.Dispatch(actions.DisplayHelp())
			case 'g', 'j', 'k': // Movement is happening => breaking of following
				stateStore.Dispatch(actions.TurnOffFollowing())
			case 'G':
				stateStore.Dispatch(actions.TurnOnFollowing())
			}
		case tcell.KeyEnd:
			stateStore.Dispatch(actions.TurnOnFollowing())
		case tcell.KeyHome, tcell.KeyUp, tcell.KeyDown, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyCtrlF, tcell.KeyCtrlB:
			stateStore.Dispatch(actions.TurnOffFollowing())
		}

		return event
	})

	stateStore.AddReducer(logsViewReducer)
	stateStore.AddHook(func(s gredux.State) {
		st := s.(state.State)

		view.SetText(st.Logs)
	}, []string{
		actions.ActionNameFilter,
		actions.ActionNameToggleFilter,
		actions.ActionNameToggleNonPatternLines,
		actions.ActionNameSetPattern,
		actions.ActionNameAddLogLine,
	})
	return view
}
