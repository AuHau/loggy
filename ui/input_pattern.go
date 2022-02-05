package ui

import (
	"errors"
	"fmt"
	"github.com/auhau/allot"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/auhau/loggy/store"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// makeParsingPattern is a main entry point for UI to set a new parsing pattern
func makeParsingPattern(pattern string) (parsingPatternInstance allot.Command, err error) {
	// There might be "so invalid syntax" that regex starts panicking
	// and not returning error, so we catch everything here just to be sure.
	defer func() {
		receivedErr := recover()
		if receivedErr != nil {
			err = errors.New("invalid syntax of parsing pattern")
		}
	}()

	parsingPatternInstance, err = allot.NewWithEscaping(pattern, store.Types)

	if err != nil {
		return parsingPatternInstance, fmt.Errorf("invalid syntax of parsing pattern: %s", err)
	}

	return parsingPatternInstance, nil
}

func patternInputReducer(s gredux.State, action gredux.Action) gredux.State {
	st := s.(state.State)

	switch action.ID {
	case actions.ActionNameDisplayPatternInput:
		st.DisplayPatternInput = true
		return st
	case actions.ActionNameSetPattern:
		pattern, err := makeParsingPattern(action.Data.(string))
		if err != nil {
			st.DisplayError = true
			st.ErrorMessage = fmt.Sprint(err)
			return st
		}

		st.ParsingPattern = pattern
		st.ParsingPatternString = action.Data.(string)
		st.DisplayPatternInput = false
		return st
	case actions.ActionNameToggleNonPatternLines:
		st.DisplayNonPatternLines = !st.DisplayNonPatternLines
		return st
	}

	return st
}

func makePatternInput(stateStore *gredux.Store) *tview.InputField {
	initPattern := stateStore.State().(state.State).ParsingPatternString
	input := tview.NewInputField().
		SetText(initPattern).
		SetLabel("Parsing pattern: ")

	var handlePatternInput = func(key tcell.Key) {
		if key == tcell.KeyEsc {
			// User canceled entering the filter string
			// lets revert the filter to original string.
			currentPattern := stateStore.State().(state.State).ParsingPatternString
			input.SetText(currentPattern)
		} else {
			stateStore.Dispatch(actions.SetPattern(input.GetText()))
		}
	}

	input.SetDoneFunc(handlePatternInput)
	stateStore.AddReducer(patternInputReducer)

	return input
}
