package ui

import (
	"fmt"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func errorModalReducer(s gredux.State, action gredux.Action) gredux.State {
	st := s.(state.State)

	switch action.ID {
	case actions.ActionNameHideError:
		st.ErrorMessage = ""
		st.DisplayError = false
		return st
	}

	return st
}

func makeErrorModal(stateStore *gredux.Store) *tview.Modal {
	modal := tview.NewModal().
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			stateStore.Dispatch(actions.HideError())
		})
	modal.SetBackgroundColor(tcell.ColorDarkRed)

	stateStore.AddReducer(errorModalReducer)
	stateStore.AddHook(func(s gredux.State) {
		st := s.(state.State)

		modal.SetText(fmt.Sprintf(`There was an error!

%s`, st.ErrorMessage))
	}, []string{
		actions.ActionNameDropLogLine,
		actions.ActionNameFilter,
		actions.ActionNameToggleFilter,
		actions.ActionNameSetPattern,
		actions.ActionNameAddLogLine,
	})

	return modal
}
