package store

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/auhau/allot"
	"strings"
)

const PATTERN_MATCHING_PARAMETER_NAME = "patternMatches"

// Filter goes through buffer and writes out logs that match
// It assumes that display was cleared out
func Filter(filter *vm.Program, pattern allot.Command, nonPatternMatchingDecorator func(string) string) (totalLines, matchingLines, nonPatternLines int, logs string, err error) {
	mu.Lock()
	defer mu.Unlock()
	builder := strings.Builder{}
	firstLine := true

	// We walk the buffer from back to front
	for element := buffer.Back(); element != nil; element = element.Prev() {
		line := fmt.Sprintf("%s", element.Value)
		totalLines += 1

		filterMatched, patternMatch, err := IsLineMatching(line, filter, pattern)
		if err != nil {
			return 0, 0, 0, "", err
		}

		if filterMatched {
			if !firstLine {
				line = "\n" + line
			} else {
				firstLine = false
			}

			if !patternMatch {
				line = nonPatternMatchingDecorator(line)
			}

			builder.WriteString(line)
			matchingLines += 1
		}

		if !patternMatch {
			nonPatternLines += 1
		}
	}

	return totalLines, matchingLines, nonPatternLines, builder.String(), nil
}

// buildParameters parses the logs line using parsing pattern,
// extracts specified parameters and parse them into type if needed.
func buildParameters(match allot.MatchInterface, pattern allot.Command) (map[string]interface{}, error) {
	parameters := make(map[string]interface{}, 10)
	patternParams := pattern.Parameters()

	for _, parameter := range patternParams {
		switch parameter.Data() {
		case "string", "rest":
			if match != nil {
				value, err := match.String(parameter.Name())
				if err != nil {
					return nil, err
				}

				parameters[parameter.Name()] = value
			} else {
				parameters[parameter.Name()] = ""
			}

		case "integer":
			if match != nil {
				value, err := match.Integer(parameter.Name())
				if err != nil {
					return nil, err
				}

				parameters[parameter.Name()] = value
			} else {
				parameters[parameter.Name()] = 0
			}
		default:
			return nil, fmt.Errorf("unknown data type %s", parameter.Data())
		}
	}

	return parameters, nil
}

// IsLineMatching check if for given filter and pattern the line matches.
// It returns enum values FILTER_MATCH, PARSE_PATTERN_NO_MATCH or FILTER_NO_MATCH according it matching result
func IsLineMatching(line string, filter *vm.Program, pattern allot.Command) (filterMatched bool, patternMatched bool, err error) {
	match, matchError := pattern.Match(line)

	// if no filter is configured than we don't have to do filterExpression evaluation
	if filter == nil {
		// there was a parse pattern error though, which signals that the line is not valid
		if matchError != nil {
			return true, false, nil
		}

		return true, true, nil
	}

	var (
		parameters map[string]interface{}
	)

	// We did not match the line against the parsing pattern so not expected parameters are available
	// but we will expose this information as parameter itself to the user.
	parameters, err = buildParameters(match, pattern)
	parameters[PATTERN_MATCHING_PARAMETER_NAME] = matchError == nil
	if err != nil {
		return false, false, err
	}

	result, err := expr.Run(filter, parameters)
	if err != nil {
		return false, false, err
	}

	if result.(bool) {
		return true, matchError == nil, nil
	}

	return false, matchError == nil, nil
}
