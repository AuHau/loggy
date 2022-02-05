package actions

import "github.com/auhau/gredux"

const (
	ActionNameHideError             = "hideError"
	ActionNameDisplayHelp           = "displayHelp"
	ActionNameHideHelp              = "hideHelp"
	ActionNameDisplayFilterInput    = "displayFilterInput"
	ActionNameHideFilterInput       = "hideFilterInput"
	ActionNameFilter                = "filter"
	ActionNameToggleFilter          = "toggleFilter"
	ActionNameDisplayPatternInput   = "displayPatternInput"
	ActionNameToggleNonPatternLines = "toggleNonPatternLines"
	ActionNameSetPattern            = "setPattern"
	ActionNameTurnOnFollowing       = "turnOnFollowing"
	ActionNameTurnOffFollowing      = "turnOffFollowing"
	ActionNameDropLogLine           = "dropLogLine"
	ActionNameAddLogLine            = "addLogLine"
)

func HideError() gredux.Action {
	return gredux.Action{ID: ActionNameHideError, Data: nil}
}

func DisplayHelp() gredux.Action {
	return gredux.Action{ID: ActionNameDisplayHelp, Data: nil}
}
func HideHelp() gredux.Action {
	return gredux.Action{ID: ActionNameHideHelp, Data: nil}
}

func DisplayFilterInput() gredux.Action {
	return gredux.Action{ID: ActionNameDisplayFilterInput, Data: nil}
}
func HideFilterInput() gredux.Action {
	return gredux.Action{ID: ActionNameHideFilterInput, Data: nil}
}
func Filter(filterString string) gredux.Action {
	return gredux.Action{ID: ActionNameFilter, Data: filterString}
}
func ToggleFilter() gredux.Action {
	return gredux.Action{ID: ActionNameToggleFilter, Data: nil}
}

func DisplayPatternInput() gredux.Action {
	return gredux.Action{ID: ActionNameDisplayPatternInput, Data: nil}
}
func ToggleNonPatternLines() gredux.Action {
	return gredux.Action{ID: ActionNameToggleNonPatternLines, Data: nil}
}
func SetPattern(pattern string) gredux.Action {
	return gredux.Action{ID: ActionNameSetPattern, Data: pattern}
}

func TurnOnFollowing() gredux.Action {
	return gredux.Action{ID: ActionNameTurnOnFollowing, Data: nil}
}
func TurnOffFollowing() gredux.Action {
	return gredux.Action{ID: ActionNameTurnOffFollowing, Data: nil}
}

func AddLogLine(line string) gredux.Action {
	return gredux.Action{ID: ActionNameAddLogLine, Data: line}
}
func DropLogLine(line string) gredux.Action {
	return gredux.Action{ID: ActionNameDropLogLine, Data: line}
}
