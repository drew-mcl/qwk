package pkg

import (
	"strings"

	"github.com/drew-mcl/logz"
)

func Lex(input string) []*Node {
	lines := strings.Split(input, "\n")

	nodes := []*Node{}

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		if line == "" {
			continue
		}

		// Handle variable assignment
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "\t") {
			logz.Trace("Found variable line:", line)

			parts := strings.SplitN(line, "=", 2)
			varName := parts[0]
			varValue := parts[1]

			logz.Debug("Variable deceleration:", varName)
			logz.Debug("Variable assignment:", varValue)

			variable := &Node{
				Type:  VariableLine,
				Key:   strings.TrimSpace(varName),
				Value: strings.TrimSpace(varValue),
			}
			nodes = append(nodes, variable)
			continue
		}

		// Handle targets
		if strings.Contains(line, ":") {
			logz.Trace("Found target line:", line)

			parts := strings.SplitN(line, ":", 2)

			logz.Trace("Target deceleration:", parts[0])
			logz.Trace("Target optional recursions:", "[", parts[1], "]")

			target := &Node{
				Type:  TargetLine,
				Key:   strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			}
			nodes = append(nodes, target)
			continue
		}

		// Handle commands
		if strings.HasPrefix(line, "\t") {
			logz.Trace("Found command line:", line)

			command := &Node{
				Type:  CommandLine,
				Value: strings.TrimPrefix(line, "\t"),
			}
			nodes = append(nodes, command)
			continue
		}

		// Anything else is an error
		logz.Error("Unexpected line:", line)
	}

	return nodes
}
