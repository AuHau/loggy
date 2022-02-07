package ui

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/auhau/loggy/store"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: Autocomplete the pattern's variables
// TODO: Possibility to save patterns (maybe also per directory?)

func filterInputReducer(s gredux.State, action gredux.Action) gredux.State {
	st := s.(state.State)

	switch action.ID {
	case actions.ActionNameDisplayFilterInput:
		st.DisplayFilterInput = true
		return st

	case actions.ActionNameHideFilterInput:
		st.DisplayFilterInput = false
		return st

	case actions.ActionNameFilter:
		var (
			filterExpression *vm.Program
			err              error
			filterString     = action.Data.(string)
		)

		if filterString == "" {
			filterExpression = nil
		} else {
			filterExpression, err = expr.Compile(filterString, expr.AsBool())
			if err != nil {
				st.DisplayError = true
				st.ErrorMessage = fmt.Sprint(err)
				return st
			}
		}

		totalLines, matchingLines, nonPatternLines, logs, err := store.Filter(filterExpression, st.ParsingPattern, makeNonPatternMatchedDecorator())
		if err != nil {
			st.DisplayError = true
			st.ErrorMessage = fmt.Sprint(err)
			return st
		}

		st.Logs = logs
		st.FilterExpression = filterExpression
		st.TotalLines = totalLines
		st.MatchingLines = matchingLines
		st.NonPatternLines = nonPatternLines
		st.IsFilterOn = filterString != ""
		st.DisplayFilterInput = false
		st.FilterString = action.Data.(string)

		return st

	case actions.ActionNameToggleFilter:
		// No filter is set, so we don't have to bother with changing anything
		if st.FilterString == "" {
			return st
		}

		var (
			logs string
			err  error
		)

		if st.IsFilterOn { // Filter is On, so we are turning it off so displaying all logs
			_, _, _, logs, err = store.Filter(nil, st.ParsingPattern, makeNonPatternMatchedDecorator())
		} else { // Filter is Off, so we are turning it on so displaying filtered logs, the number of lines shouldn't have changed.
			_, _, _, logs, err = store.Filter(st.FilterExpression, st.ParsingPattern, makeNonPatternMatchedDecorator())
		}

		if err != nil {
			st.DisplayError = true
			st.ErrorMessage = fmt.Sprint(err)
			return st
		}

		st.Logs = logs
		st.IsFilterOn = !st.IsFilterOn
		return st
	}

	return st
}

func makeFilterInput(stateStore *gredux.Store) *tview.InputField {
	input := tview.NewInputField().
		SetLabel("Filter: ")

	var handleFilterInput = func(key tcell.Key) {
		if key == tcell.KeyEsc {
			// User canceled entering the filter string
			// lets revert the filter to original string.
			currentFilter := stateStore.State().(state.State).FilterString
			input.SetText(currentFilter)

			stateStore.Dispatch(actions.HideFilterInput())
		} else {
			stateStore.Dispatch(actions.Filter(input.GetText()))
		}
	}

	input.SetDoneFunc(handleFilterInput)
	stateStore.AddReducer(filterInputReducer)

	return input
}
