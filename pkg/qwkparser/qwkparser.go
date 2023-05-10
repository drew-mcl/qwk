package qwkparser

import (
	"qwk/pkg/logger"
	"strings"
)

type Target struct {
	Name       string
	Commands   []string
	DependsOn  []string
	BuildAfter []*Target
}

type Qwkfile struct {
	Variables map[string]string
	Targets   map[string]*Target
}

func isCommentOrEmpty(line string) bool {
	return strings.HasPrefix(line, "#") || line == ""
}

func parseVariable(line string, mf *Qwkfile) {
	logger.Debug("Parsing variable:", line)
	parts := strings.Split(line, "=")
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	mf.Variables[key] = value
	logger.Info("Added variable:", key, "=", value)
}

func parseTarget(line string, mf *Qwkfile) *string {
	logger.Debug("Parsing target:", line)
	parts := strings.Split(line, ":")
	targetName := strings.TrimSpace(parts[0])

	var dependencies []string
	if len(parts) > 1 {
		dependencies = strings.Split(parts[1], " ")
		for i, dep := range dependencies {
			dependencies[i] = strings.TrimSpace(dep)
		}
	}

	mf.Targets[targetName] = &Target{
		Name:     targetName,
		Commands: []string{},
	}
	ParseDependencies(mf.Targets[targetName], dependencies)

	logger.Info("Added target:", targetName, "with dependencies:", dependencies)
	return &targetName
}

func ParseDependencies(target *Target, dependencies []string) {
	logger.Debug("Adding dependencies to target:", target.Name)
	for _, dep := range dependencies {
		if dep != "" {
			target.DependsOn = append(target.DependsOn, dep)
		}
	}
}

func parseCommand(line string, mf *Qwkfile, currentTarget *string) {
	logger.Debug("Parsing command:", line)
	command := line
	mf.Targets[*currentTarget].Commands = append(mf.Targets[*currentTarget].Commands, command)
	logger.Info("Added command to target:", *currentTarget, ":", command)
}

func processLine(line string, mf *Qwkfile, currentTarget *string) {
	switch {
	case strings.Contains(line, "="):
		parseVariable(line, mf)
	case strings.Contains(line, ":"):
		*currentTarget = *parseTarget(line, mf)
	default:
		if *currentTarget != "" {
			parseCommand(line, mf, currentTarget)
		}
	}
}

func ParseQwkfile(content string) (*Qwkfile, error) {
	lines := strings.Split(content, "\n")
	mf := &Qwkfile{
		Variables: make(map[string]string),
		Targets:   make(map[string]*Target),
	}

	currentTarget := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if isCommentOrEmpty(line) {
			continue
		}

		processLine(line, mf, &currentTarget)
	}
	return mf, nil
}
