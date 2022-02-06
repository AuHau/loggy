package ui

import (
	"fmt"
)

func makeNonPatternMatchedDecorator() func(string) string {
	return func(s string) string {
		return fmt.Sprintf("%s%s[-:-:-]", NON_PATTERN_MATCHING_LINES_FORMAT, s)
	}
}
