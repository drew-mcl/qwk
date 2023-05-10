package substitution

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ErrArgOutOfRange is returned when a placeholder refers to an argument index that is out of range
var ErrArgOutOfRange = errors.New("argument index out of range")

// ErrInvalidArgIndex is returned when a placeholder contains an invalid argument index
var ErrInvalidArgIndex = errors.New("invalid argument index")

// ExpandVariables replaces variable placeholders with their values in a command string
// If a placeholder is not matched to any argument or variable, it is left as is
// If a placeholder contains an invalid argument index, an error is returned
func ExpandVariables(command string, variables map[string]string, args []string) (string, error) {
	// Replace global variables
	for key, value := range variables {
		placeholder := fmt.Sprintf("$(%s)", key)
		command = strings.ReplaceAll(command, placeholder, value)
	}

	// Replace inline variables and arguments
	argPattern := regexp.MustCompile(`\$\d+|\$\*`)
	command = argPattern.ReplaceAllStringFunc(command, func(s string) string {
		switch {
		case s == "$*":
			// Replace with all arguments as a single string
			return strings.Join(args, " ")
		default:
			// Replace with the corresponding argument
			argIndex, err := strconv.Atoi(s[1:])
			if err != nil {
				// The argument index is not a valid integer
				return s
			}
			if argIndex > 0 && argIndex <= len(args) {
				return args[argIndex-1]
			}
			// The argument index is out of range
			return s
		}
	})

	return command, nil
}
