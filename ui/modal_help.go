package ui

import (
	"fmt"
	"github.com/auhau/gredux"
	"github.com/auhau/loggy/state"
	"github.com/auhau/loggy/state/actions"
	"github.com/auhau/loggy/store"
	"github.com/rivo/tview"
)
import "github.com/chonla/format"

// Version string set in release time
var Version string = ""

var KEY_MAP = map[string]interface{}{
	"filter":       string(SET_FILTER_KEY),
	"pattern":      string(SET_PATTERN_KEY),
	"toggleFilter": string(TOGGLE_FILTER_KEY),
	"help":         string(HELP_KEY),
	"version":      Version,
}

var HelpHomeText = format.Sprintf(`
Status bar on top displays several helpful information. Describing from left to right:
 - Input name
 - Optional "F" indicator that shows if loggy is following the end of the logs
 - Filter status that displays "<number of filter matching lines>/<number of total lines>". If it has green background than filter is applied otherwise is turned off or not set.
 - Optional number of lines that were not possible to match against the parsing pattern.s

Main key shortcuts:
 - "%<filter>s" for setting filter
 - "%<toggleFilter>s" for toggling filter
 - "%<pattern>s" for setting parsing pattern input
 - "%<help>s" for displaying help`, KEY_MAP)

var HelpNavigationText = `Logs navigation:
 - "j", "k" or arrow keys for scrolling by one line 
 - "g" to move to top
 - "G" to move to bottom and follow bottom
 - "Ctrl-F", "page down" to move down by one page
 - "Ctrl-B", "page up" to move up by one page
 - Mouse scrolling also scrolls the view accordingly`

var HelpInputsText = `Input fields:
 - Left arrow: Move left by one character.
 - Right arrow: Move right by one character.
 - Home, Ctrl-A, Alt-a: Move to the beginning of the line.
 - End, Ctrl-E, Alt-e: Move to the end of the line.
 - Alt-left, Alt-b: Move left by one word.
 - Alt-right, Alt-f: Move right by one word.
 - Backspace: Delete the character before the cursor.
 - Delete: Delete the character after the cursor.
 - Ctrl-K: Delete from the cursor to the end of the line.
 - Ctrl-W: Delete the last word before the cursor.
 - Ctrl-U: Delete the entire line.`

var HelpParsingPatternText = fmt.Sprintf(`The logs are parsed using parsing pattern that you have to configure in order to use filters. The lines are tokenized using space character. Internally regex is used for parsing, but the input pattern is escaped by default for special characters so you don't have to worry about special characters. You define parameters using syntax "<name:type>", where name is the name of parameter that you can refer to in filters and type is predefined type used to correctly find and parse the parameter.

Lines that were not possible to parsed are colored with red color. Moreover counter of how many lines were not possible to parse is displayed in the status bar on the right end of it. It is only present if there are some lines that were not possible to parse.  

There is built-in bool parameter called "%s" reserved for marking log lines that were or were not possible to match against the parsing pattern. So you can use that to debug your parsing pattern with expressions like "!%[1]s"

Supported types:
 - "string" defines string containing non-whitespace characters: [^\s]+
 - "integer" defines a integer: [0-9]+
 - "rest" collects the rest of the line: .*

Example log and bellow its parsing pattern:
[2022-09-11T15:04:22](authorization) DEBUG 200 We have received login information
[<timestamp:string>](<component:string>) <level:string> <code:integer> <message:rest>`, store.PATTERN_MATCHING_PARAMETER_NAME)

var HelpFilterText = `In order to use filter for the logs you have to define parsing pattern in which you define parameters that are extracted from the log lines. Then you can write filter expressions that will be applied on the logs. Filter has to return bool otherwise error will be shown.

loggy uses internally "expr" which has very rich set of arithmetic/string operators that you can use for your filters. Brief overview:
 - modifiers: + - / * % **
 - comparators: > >= < <= == !=
 - logical ops: not ! or || and &&
 - numeric constants, as 64-bit floating point (12345.678)
 - numeric range: '..' (18..45) 
 - string constants (single or double quotes)
 - string operators: + matches contains startsWith endsWith
 - boolean constants: true false
 - parenthesis to control order of evaluation ( )
 - arrays e.g. [1, 2, 3]
 - maps - e.g. {foo: "bar"}
 - ternary conditional: ? :
 - built in functions: len() all() none() any() one() filter() map() count()

For more details see: https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md

Example of filter for the parsing pattern log above:
level == "DEBUG" - display only debug messages
code in 400..500 - display logs with code in range from 400 to 500 (inclusive)`

func helpModalReducer(s gredux.State, action gredux.Action) gredux.State {
	st := s.(state.State)

	switch action.ID {
	case actions.ActionNameDisplayHelp:
		st.DisplayHelp = true
		return st
	case actions.ActionNameHideHelp:
		st.DisplayHelp = false
		return st
	}

	return st
}

func makeHelpModal(stateStore *gredux.Store) *tview.Modal {
	helpPages := map[string]string{
		"General":         HelpHomeText,
		"Navigation":      HelpNavigationText,
		"Inputs":          HelpInputsText,
		"Parsing pattern": HelpParsingPatternText,
		"Filter inputs":   HelpFilterText,
	}
	pagesButton := []string{"General", "Navigation", "Inputs", "Parsing pattern", "Filter inputs"}

	modal := tview.NewModal()
	modal.
		SetText(fmt.Sprintf(`loggy %s

%s

<< press ESC to close >>`, Version, HelpHomeText)).
		AddButtons(pagesButton).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// Empty buttonLabel means ESC key pressed
			if buttonLabel == "Close" || buttonLabel == "" {
				stateStore.Dispatch(actions.HideHelp())
			} else {
				modal.SetText(fmt.Sprintf(`loggy %s

%s

<<press ESC to close>>`, Version, helpPages[buttonLabel]))
			}
		})

	stateStore.AddReducer(helpModalReducer)

	return modal
}
