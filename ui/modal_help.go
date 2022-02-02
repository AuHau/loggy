package ui

import (
	"fmt"
	"github.com/rivo/tview"
)
import "github.com/chonla/format"

// Version string set in release time
var Version string = ""

var KEY_MAP = map[string]interface{}{
	"filter":       string(FILTER_KEY),
	"pattern":      string(PATTERN_KEY),
	"follow":       string(FOLLOW_KEY),
	"toggleFilter": string(TOGGLE_FILTER_KEY),
	"help":         string(HELP_KEY),
	"version":      Version,
}

var HelpHomeText = format.Sprintf(`General keys:
 - "%<filter>s" for setting filter
 - "%<toggleFilter>s" for toggling filter
 - "%<pattern>s" for setting parsing pattern input
 - "%<follow>s" for scroll to bottom and follow new data
 - "%<help>s" for displaying help`, KEY_MAP)

var HelpNavigationText = `Logs navigation:
 - "j", "k" or arrow keys for scrolling by one line 
 - "g", "G" to move to top / bottom
 - "Ctrl-F", "page down" to move down by one page
 - "Ctrl-B", "page up" to move up by one page`

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

var HelpParsingPatternText = `The logs are parsed using parsing pattern that you have to configure in order to use filters. The lines are tokenized using space character so you have to define section of the line. Internally regex is used for parsing, but the input pattern is escaped by default for special characters so you don't have to worry about that. You define parameters using syntax "<name:type>", where name is the name of parameter that you can refer to in filters and type is predefined type used to correctly find and parse the parameter.

Supported types:
 - "string" defines string containing non-whitespace characters: [^\s]+
 - "integer" defines a integer: [0-9]+
 - "rest" collects the rest of the line: .*

Example log and bellow its parsing pattern:
[2022-09-11T15:04:22](authorization) DEBUG 200 We have received login information
[<timestamp:string>](<component:string>) <level:string> <code:integer> <message:rest>`

var HelpFilterText = `In order to use filter for the logs you have to define parsing pattern in which you define parameters that are extracted from the log lines. Then you can write filter expressions that will be applied on the logs. Filter has to return bool otherwise error will be shown.

loggy uses internally "govaluate" which has very rich set of C-like arithmetic/string expressions that you can use for your filters. Brief overview:
 - modifiers: + - / * & | ^ ** %% >> <<
 - comparators: > >= < <= == != =~ !~
 - logical ops: || &&
 - numeric constants, as 64-bit floating point (12345.678)
 - string constants (single quotes: 'foobar')
 - date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
 - boolean constants: true false
 - parenthesis to control order of evaluation ( )
 - arrays (anything separated by , within parenthesis: (1, 2, 'foo'))
 - prefixes: ! - ~
 - ternary conditional: ? :
 - null coalescence: ??

For more details see: https://github.com/Knetic/govaluate/blob/master/MANUAL.md

Example of filter for the parsing pattern log above:
level == "DEBUG" - display only debug messages
code > 400 - display logs with code higher then 400`

func makeHelpModal() *tview.Modal {
	helpPages := map[string]string{
		"General":         HelpHomeText,
		"Navigation":      HelpNavigationText,
		"Inputs":          HelpInputsText,
		"Parsing pattern": HelpParsingPatternText,
		"Filter inputs":   HelpFilterText,
	}
	pagesButton := []string{"Close", "General", "Navigation", "Inputs", "Parsing pattern", "Filter inputs"}

	modal := tview.NewModal()
	modal.
		SetText(fmt.Sprintf(`loggy %s

%s`, Version, HelpHomeText)).
		AddButtons(pagesButton).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// Empty buttonLabel means ESC key pressed
			if buttonLabel == "Close" || buttonLabel == "" {
				pages.HidePage(HELP_PAGE_NAME)
				app.SetFocus(logsView)
			} else {
				modal.SetText(fmt.Sprintf(`loggy %s

%s`, Version, helpPages[buttonLabel]))
			}
		})

	return modal
}
