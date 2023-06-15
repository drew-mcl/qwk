package main

import (
	"bufio"
	"compile/pkg"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/drew-mcl/logz"
)

var (
	logLevel = flag.String("l", "warn", "set log level")
)

func main() {

	flag.Parse()

	logz.SetLevel(*logLevel)

	if len(flag.Args()) == 0 {
		logz.Error("Incorrect usage")
		fmt.Fprintf(os.Stderr, "Usage: %s [-l loglevel] <command>\n", os.Args[0])
		os.Exit(1)
	}

	qwkfilePath := "./qwkfile"

	file, err := os.Open(qwkfilePath)
	if err != nil {
		logz.Error("Could not open qwkfile:", err)
		panic(err)
	}
	defer file.Close()

	var data strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		logz.Trace("scanner:", line)

		data.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		logz.Error("Could not scan file:", err)
		panic(err)
	}

	nodes := pkg.Lex(data.String())

	mf, err := pkg.Parse(nodes)
	if err != nil {
		logz.Error("Could not parse nodes:", err)
		panic(err)
	}

	//pkg.PrintTree(mf)
	logz.InfoWithSuccess("Parsed qwkfile")

	// Gather all variables in a map
	variables := make(map[string]string)
	for _, node := range nodes {
		if node.Type == pkg.VariableLine {
			variables[node.Key] = node.Value
		}
	}

	args := flag.Args()[1:]

	// Create a new Executor and run it
	executor := pkg.NewExecutor(mf, variables, args)

	target := flag.Arg(0)

	node := executor.FindTarget(target)
	if node == nil {
		logz.Error("Target not found:", target)
		os.Exit(1)
	}
	err = executor.RunTarget(node)
	if err != nil {
		logz.Error("Execution failed:", err)
		os.Exit(1)
	}
}
