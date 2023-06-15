package pkg

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/drew-mcl/logz"
	"github.com/mattn/go-shellwords"
)

type Executor struct {
	root      *Node
	variables map[string]string
	args      []string
}

// NewExecutor initializes a new Executor instance
func NewExecutor(root *Node, variables map[string]string, args []string) *Executor {
	return &Executor{
		root:      root,
		variables: variables,
		args:      args,
	}
}

// Run starts executing the parsed makefile
func (e *Executor) Run() error {
	for _, child := range e.root.Children {
		if child.Type == TargetLine {
			err := e.RunTarget(child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// runTarget runs a specific target and all of its dependencies
func (e *Executor) RunTarget(target *Node) error {
	logz.Trace("Running target:", target.Key)

	// Resolve and run dependencies first
	dependencies := strings.Split(target.Value, " ")
	for _, dependency := range dependencies {
		if dependency != "" {
			dependencyNode := e.FindTarget(dependency)
			if dependencyNode == nil {
				return errors.New("Unknown dependency: " + dependency)
			}

			err := e.RunTarget(dependencyNode)
			if err != nil {
				return err
			}
		}
	}

	// Run the commands associated with this target
	for _, child := range target.Children {
		if child.Type == CommandLine {
			err := e.runCommand(child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// runCommand executes a command line
func (e *Executor) runCommand(command *Node) error {
	// Substitute variables in the command line
	cmdLine := e.substituteVariables(command.Value)

	// Parse the command line into arguments and inline environment vars
	args, envVars, err := parseCmdLine(cmdLine)
	if err != nil {
		logz.Error("Failed to parse command:", err)
		return err
	}

	// Run the command directly, without /bin/sh
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = append(os.Environ(), envVars...)

	logz.Info("Command to execute is:", cmd)
	logz.Debug("Env vars are:", cmd.Env)

	var stderr bytes.Buffer

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmtedErr := strings.TrimSpace(stderr.String())
		errLines := strings.Split(fmtedErr, "\n")
		for _, line := range errLines {
			logz.Error(line)
		}
		return err
	}
	return nil
}

// findTarget finds a target node by its key
func (e *Executor) FindTarget(key string) *Node {
	for _, child := range e.root.Children {
		if child.Type == TargetLine && child.Key == key {
			return child
		}
	}
	return nil
}

// substituteVariables substitutes variable references in a string
func (e *Executor) substituteVariables(input string) string {
	output := input
	for key, value := range e.variables {
		// Assumes the variable reference is in the form of $(VAR)
		output = strings.ReplaceAll(output, "$("+key+")", value)
	}

	//substitutes positional arguments
	for i, arg := range e.args {
		placeHolder := "$" + strconv.Itoa(i+1)
		output = strings.ReplaceAll(output, placeHolder, arg)
	}

	// Substitute $@ with all arguments
	output = strings.ReplaceAll(output, "$@", strings.Join(e.args, " "))
	return output
}

func parseCmdLine(cmdLine string) (args []string, envVars []string, err error) {
	// Split the command line into words
	words, err := shellwords.Parse(cmdLine)
	if err != nil {
		return nil, nil, err
	}

	// Separate words into args and environment variables
	for _, word := range words {
		if strings.Contains(word, "=") {
			if isValidEnvVarAssignment(word) {
				logz.Debug("Found inline variable:", word, "in command:", cmdLine)
				envVars = append(envVars, word)
			} else {
				logz.Error("Invalid inline variable assignment:", word)
			}
		} else {
			args = append(args, word)
		}
	}

	return args, envVars, nil
}

func isValidEnvVarAssignment(word string) bool {
	// Validating that the word has the form "key=value"
	parts := strings.SplitN(word, "=", 2)
	if len(parts) != 2 {
		return false
	}

	// Validating that the key doesn't contain any slashes or invalid characters
	key := strings.TrimSpace(parts[0])
	if strings.ContainsAny(key, "/") {
		return false
	}

	// Add any other checks you consider necessary for a valid key

	return true
}
