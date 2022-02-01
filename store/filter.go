package store

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/auhau/allot"
	"reflect"
)

// parsePattern is current parsing pattern used on every log line
var parsePattern allot.Command

// filterExpression is compiled expression that is used to filter out the logs
var filterExpression *govaluate.EvaluableExpression

// types are pattern's types usable in declaring parameter's type and its regex shape
// TODO: Allow custom types in config
var types = map[string]string{
	"string":  "[^\\s]+",
	"integer": "[0-9]+",
	"rest":    ".*",
}

// TODO: Option to define if non-matching lines should be printed or not.

// SetParsePattern is a main entry point for UI to set a new parsing pattern
func SetParsePattern(pattern string) (err error) {
	// There might be "so invalid syntax" that regex starts panicking
	// and not returning error, so we catch everything here just to be sure.
	defer func() {
		receivedErr := recover()
		if receivedErr != nil {
			err = errors.New("invalid syntax of parsing pattern")
		}
	}()

	parsePattern, err = allot.NewWithEscaping(pattern, types)

	if err != nil {
		return fmt.Errorf("invalid syntax of parsing pattern: %s", err)
	}

	return nil
}

// Filter goes through buffer and writes out logs that match
// It assumes that display was cleared out
func Filter(filter string) error {
	mu.Lock()
	defer mu.Unlock()
	var err error
	firstLine := true

	// Empty string means no filtering
	if filter == "" {
		filterExpression = nil
	} else {
		filterExpression, err = govaluate.NewEvaluableExpression(filter)
		if err != nil {
			return err
		}
	}

	for element := buffer.Front(); element != nil; element = element.Next() {
		line := fmt.Sprintf("%s", element.Value)

		result, err := isLineMatching(line)
		if err != nil {
			return err
		}

		if result {
			// We don't want empty line in beginning nor end
			if firstLine {
				_, err = fmt.Fprint(writer, line)
				firstLine = false
			} else {
				_, err = fmt.Fprint(writer, "\n"+line)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// buildParameters parses the logs line using parsing pattern,
// extracts specified parameters and parse them into type if needed.
func buildParameters(match allot.MatchInterface) (map[string]interface{}, error) {
	parameters := make(map[string]interface{}, 10)
	patternParams := parsePattern.Parameters()

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

func isLineMatching(line string) (bool, error) {
	// if no filter is configured that for sure it is matching
	if filterExpression == nil {
		return true, nil
	}

	match, err := parsePattern.Match(line)

	// error is returned if the line does not match so we return only false
	if err != nil {
		return false, nil
	}

	parameters, err := buildParameters(match)
	if err != nil {
		return false, err
	}

	result, err := filterExpression.Evaluate(parameters)
	if err != nil {
		return false, err
	}

	if reflect.TypeOf(result).String() != "bool" {
		return false, errors.New("filter expression did not return boolean")
	}

	return result.(bool), nil
}
