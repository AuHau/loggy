package state

import (
	"github.com/antonmedv/expr/vm"
	"github.com/auhau/allot"
)

type State struct {
	// Logs view
	IsFollowing bool
	IsFilterOn  bool
	InputName   string

	IsLogsFirstLine bool
	Logs            string
	TotalLines      int
	MatchingLines   int
	NonPatternLines int

	// ErrorModal
	DisplayError bool
	ErrorMessage string

	// HelpModal
	DisplayHelp bool

	// InputFilter
	FilterString       string
	FilterExpression   *vm.Program
	DisplayFilterInput bool

	// PatternFilter
	DisplayPatternInput  bool
	ParsingPatternString string
	ParsingPattern       allot.Command
}
