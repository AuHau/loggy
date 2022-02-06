package ui

import (
	"fmt"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/rivo/tview"
)

var (
	statusInputName       *tview.TextView
	statusFollowingStatus *tview.TextView
	statusFilterStatus    *tview.TextView
	statusNonPatternLines *tview.TextView
)

func makeStatusBar(stateStore *gredux.Store) *tview.Flex {
	statusFilterStatus = tview.NewTextView().
		SetText("0 / 0").
		SetTextColor(FILTER_STATUS_TEXT_COLOR).
		SetTextAlign(tview.AlignCenter)
	statusFilterStatus.SetBackgroundColor(FILTER_STATUS_NONACTIVE_BACKGROUND_COLOR)

	statusNonPatternLines = tview.NewTextView().
		SetText("0").
		SetTextColor(NON_PATTERN_LINES_STATUS_TEXT_COLOR).
		SetTextAlign(tview.AlignCenter)
	statusNonPatternLines.SetBackgroundColor(NON_PATTERN_LINES_STATUS_BACKGROUND_COLOR)

	statusInputName = tview.NewTextView().
		SetText(" " + stateStore.State().(state.State).InputName).
		SetTextColor(INPUT_NAME_TEXT_COLOR).
		SetDynamicColors(true)
	statusInputName.SetBackgroundColor(INPUT_NAME_BACKGROUND_COLOR)

	statusFollowingStatus = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(FOLLOWING_STATUS_TEXT_COLOR)

	if stateStore.State().(state.State).IsFollowing {
		statusFollowingStatus.SetText("F")
		statusFollowingStatus.SetBackgroundColor(FOLLOWING_STATUS_ACTIVE_BACKGROUND_COLOR)
	} else {
		statusFollowingStatus.SetBackgroundColor(FOLLOWING_STATUS_NONACTIVE_BACKGROUND_COLOR)
	}

	grid := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusInputName, 0, 2, false).
		AddItem(statusFollowingStatus, 5, 0, false).
		AddItem(nil, 0, 8, false).
		AddItem(statusFilterStatus, 17, 0, false)

	grid.Box = tview.NewBox().SetBackgroundColor(STATUS_BAR_BACKGROUND_COLOR)

	stateStore.AddHook(func(s gredux.State) {
		st := s.(state.State)

		if st.IsFilterOn {
			statusFilterStatus.SetBackgroundColor(FILTER_STATUS_ACTIVE_BACKGROUND_COLOR)
		} else {
			statusFilterStatus.SetBackgroundColor(FILTER_STATUS_NONACTIVE_BACKGROUND_COLOR)
		}
		statusFilterStatus.SetText(fmt.Sprintf("%d / %d", st.MatchingLines, st.TotalLines))

		if st.IsFollowing {
			statusInputName.SetText(" " + st.InputName)
			statusFollowingStatus.SetText("F")
			statusFollowingStatus.SetBackgroundColor(FOLLOWING_STATUS_ACTIVE_BACKGROUND_COLOR)
		} else {
			statusFollowingStatus.SetText("")
			statusFollowingStatus.SetBackgroundColor(FOLLOWING_STATUS_NONACTIVE_BACKGROUND_COLOR)
		}

		if st.NonPatternLines > 0 {
			if grid.GetItemCount() == 4 {
				grid.AddItem(statusNonPatternLines, 7, 0, false)
			}

			statusNonPatternLines.SetText(fmt.Sprint(st.NonPatternLines))
		} else {
			if grid.GetItemCount() == 5 {
				grid.RemoveItem(statusNonPatternLines)
			}
		}
	}, []string{
		actions.ActionNameAddLogLine,
		actions.ActionNameDropLogLine,
		actions.ActionNameFilter,
		actions.ActionNameToggleFilter,
		actions.ActionNameSetPattern,
		actions.ActionNameTurnOnFollowing,
		actions.ActionNameTurnOffFollowing,
	})
	return grid
}
