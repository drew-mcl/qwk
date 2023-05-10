package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"qwk/pkg/qwkparser"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

type TestTarget struct {
	Output string `yaml:"output"`
}

type TestTargets struct {
	All map[string]TestTarget `yaml:"all"`
}

// Capture the standard output of the function
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)

	return strings.TrimSpace(buf.String())
}

func TestQwk(t *testing.T) {
	// Read the expected_output.yaml
	expectedOutputData, err := ioutil.ReadFile("../../tests/expected/qwk.yml")
	if err != nil {
		t.Fatalf("Error reading expected_output.yaml: %v", err)
	}

	var expectedOutput TestTargets
	err = yaml.Unmarshal(expectedOutputData, &expectedOutput)
	if err != nil {
		t.Fatalf("Error unmarshaling expected_output.yaml: %v", err)
	}

	// Read the sample qwkfile
	qwkfileContent, err := ioutil.ReadFile("../../tests/qwkfiles/qwk")
	if err != nil {
		t.Fatalf("Error reading qwkfile: %v", err)
	}

	mf, err := qwkparser.ParseQwkfile(string(qwkfileContent))
	if err != nil {
		t.Fatalf("Error parsing qwkfile: %v", err)
	}

	args := []string{"test"}

	for targetName, testTarget := range expectedOutput.All {
		target, ok := mf.Targets[targetName]
		if !ok {
			t.Fatalf("Target not found in qwkfile: %s", targetName)
		}

		// Execute the target and capture its output
		output := captureStdout(func() {
			err := ExecuteTargets([]*qwkparser.Target{target}, mf, args)
			if err != nil {
				t.Fatalf("Error executing target %s: %v", targetName, err)
			}
		})

		if output != testTarget.Output {
			t.Errorf("Unexpected output for target %s. Expected: %q, Got: %q", targetName, testTarget.Output, output)
		}
	}
}
