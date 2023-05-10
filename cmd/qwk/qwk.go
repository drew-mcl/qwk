package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"qwk/pkg/logger"
	"qwk/pkg/qwkparser"
	"qwk/pkg/substitution"
)

// ResolveDependencies resolves the dependencies between targets in a Qwkfile
func ResolveDependencies(qwkfile *qwkparser.Qwkfile) error {
	visited := make(map[string]bool)
	for _, target := range qwkfile.Targets {
		err := visit(target, qwkfile, visited)
		if err != nil {
			return err
		}
	}
	return nil
}

// visit is a recursive helper function used by ResolveDependencies to traverse the dependency graph
func visit(target *qwkparser.Target, qwkfile *qwkparser.Qwkfile, visited map[string]bool) error {
	if visited[target.Name] {
		return nil
	}
	visited[target.Name] = true

	for _, depName := range target.DependsOn {
		dep, ok := qwkfile.Targets[depName]
		if !ok {
			return fmt.Errorf("dependency target not found: %s", depName)
		}
		err := visit(dep, qwkfile, visited)
		if err != nil {
			return err
		}
		target.BuildAfter = append(target.BuildAfter, dep)
	}

	return nil
}

// DetermineTargetsToBuild returns a list of targets that need to be built in order to build the given target
func DetermineTargetsToBuild(target *qwkparser.Target) ([]*qwkparser.Target, error) {
	targetsToBuild := []*qwkparser.Target{}

	// Create a map to keep track of visited targets
	visited := make(map[string]bool)

	// Define a recursive function to walk the target dependency graph and add the targets that need to be built to targetsToBuild
	var walk func(*qwkparser.Target) error
	walk = func(t *qwkparser.Target) error {
		if visited[t.Name] {
			return nil
		}
		visited[t.Name] = true

		// Recursively walk the dependencies of the target
		for _, dep := range t.BuildAfter {
			err := walk(dep)
			if err != nil {
				return err
			}
		}

		shouldBuild, err := shouldBuildTarget(t)
		if err != nil {
			return err
		}
		if shouldBuild {
			targetsToBuild = append(targetsToBuild, t)
		}

		return nil
	}

	// Walk the dependency graph starting from the given target
	err := walk(target)
	if err != nil {
		return nil, err
	}

	return targetsToBuild, nil
}

// shouldBuildTarget determines if the given target needs to be built by checking the modification time of its dependencies.
func shouldBuildTarget(target *qwkparser.Target) (bool, error) {
	targetInfo, err := os.Stat(target.Name)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat target: %v", err)
	}

	for _, dep := range target.BuildAfter {
		depInfo, err := os.Stat(dep.Name)
		if err != nil {
			return false, fmt.Errorf("failed to stat dependency: %v", err)
		}
		if depInfo.ModTime().After(targetInfo.ModTime()) {
			return true, nil
		}
	}

	return false, nil
}

// ExecuteTargets executes the commands associated with each of the given targets.
func ExecuteTargets(targets []*qwkparser.Target, qwkfile *qwkparser.Qwkfile, args []string) error {
	for _, target := range targets {
		for _, command := range target.Commands {
			expandedCommand, err := substitution.ExpandVariables(command, qwkfile.Variables, args)
			if err != nil {
				return fmt.Errorf("variable expansion failed: %w", err)
			}
			cmd := exec.Command("sh", "-c", expandedCommand)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("command execution failed: %v", err)
			}
		}
	}
	return nil
}

func main() {
	var versionFlag = flag.Bool("v", false, "Display the version of this program")
	var helpFlag = flag.Bool("h", false, "Display usage information")
	var qwkfileFlag = flag.String("f", "qwkfile", "Specify the qwkfile to use")

	logger.SetLevel(logger.WARN)
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version 1.0.0")
		os.Exit(0)
	}

	if *helpFlag {
		fmt.Println("Usage: qwk [OPTIONS] TARGET [ARGS...]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Println("No target specified. Available targets:")
		content, err := ioutil.ReadFile(*qwkfileFlag)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		mf, err := qwkparser.ParseQwkfile(string(content))
		if err != nil {
			log.Fatalf("Error: %v", err)

		}
		for key := range mf.Targets {
			log.Println(key)
		}
		os.Exit(1)
	}

	targetName := args[0]

	content, err := ioutil.ReadFile(*qwkfileFlag)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	mf, err := qwkparser.ParseQwkfile(string(content))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = ResolveDependencies(mf)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	target, ok := mf.Targets[targetName]
	if !ok {
		log.Fatalf("Error: Target not found: %s", targetName)
	}

	targetsToBuild, err := DetermineTargetsToBuild(target)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = ExecuteTargets(targetsToBuild, mf, args[1:])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
