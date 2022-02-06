package store

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/auhau/allot"
	"reflect"
	"strings"
)

// Filter goes through buffer and writes out logs that match
// It assumes that display was cleared out
func Filter(filter *govaluate.EvaluableExpression, pattern allot.Command, displayNonPatternLines bool) (totalLines, matchingLines, nonPatternLines int, logs string, err error) {
	mu.Lock()
	defer mu.Unlock()
	builder := strings.Builder{}
	firstLine := true

	// We walk the buffer from back to front
	for element := buffer.Back(); element != nil; element = element.Prev() {
		line := fmt.Sprintf("%s", element.Value)
		totalLines += 1

		result, err := IsLineMatching(line, filter, pattern)
		if err != nil {
			return 0, 0, 0, "", err
		}

		switch result {
		case MATCH:
			if !firstLine {
				line = "\n" + line
			} else {
				firstLine = false
			}

			builder.WriteString(line)
			matchingLines += 1
		case PARSE_PATTERN_NO_MATCH:
			nonPatternLines += 1

			if displayNonPatternLines {
				if !firstLine {
					line = "\n" + line
				} else {
					firstLine = false
				}

				builder.WriteString(line)
			}
		case FILTER_NO_MATCH:
			// NO OP
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
			value, err := match.String(parameter.Name())
			if err != nil {
				return nil, err
			}

			parameters[parameter.Name()] = value
		case "integer":
			value, err := match.Integer(parameter.Name())
			if err != nil {
				return nil, err
			}

			parameters[parameter.Name()] = value
		default:
			return nil, fmt.Errorf("unknown data type %s", parameter.Data())
		}
	}

	return parameters, nil
}

const (
	MATCH = iota
	PARSE_PATTERN_NO_MATCH
	FILTER_NO_MATCH
)

// IsLineMatching check if for given filter and pattern the line matches.
// It returns enum values MATCH, PARSE_PATTERN_NO_MATCH or FILTER_NO_MATCH according it matching result
func IsLineMatching(line string, filter *govaluate.EvaluableExpression, pattern allot.Command) (int, error) {
	match, err := pattern.Match(line)

	// error is returned if the line does not match but that is not an error but valid state
	if err != nil {
		return PARSE_PATTERN_NO_MATCH, nil
	}

	// if no filter is configured that for sure it is matching
	if filter == nil {
		return MATCH, nil
	}

	parameters, err := buildParameters(match, pattern)
	if err != nil {
		return -1, err
	}

	result, err := filter.Evaluate(parameters)
	if err != nil {
		return -1, err
	}

	if reflect.TypeOf(result).String() != "bool" {
		return -1, errors.New("filter expression did not return boolean")
	}

	if result.(bool) {
		return MATCH, nil
	}

	return FILTER_NO_MATCH, nil
}
