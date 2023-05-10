package qwkparser

import (
	"io/ioutil"
	"qwk/pkg/logger"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

type ExpectedResults struct {
	Variables map[string]string `yaml:"variables"`
	Targets   []string          `yaml:"targets"`
}

func TestParseQwkfile(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	content, err := ioutil.ReadFile("../../tests/qwkfiles/qwkparser")
	if err != nil {
		t.Fatalf("Error reading qwkfile: %v", err)
	}

	mf, err := ParseQwkfile(string(content))
	if err != nil {
		t.Fatalf("Error parsing qwkfile: %v", err)
	}

	expectedYAML, err := ioutil.ReadFile("../../tests/expected/qwkparser.yml")
	if err != nil {
		t.Fatalf("Error reading expected results YAML: %v", err)
	}

	var expectedResults ExpectedResults
	err = yaml.Unmarshal(expectedYAML, &expectedResults)
	if err != nil {
		t.Fatalf("Error unmarshaling expected results YAML: %v", err)
	}

	if !reflect.DeepEqual(mf.Variables, expectedResults.Variables) {
		t.Errorf("Expected variables: %v, got: %v", expectedResults.Variables, mf.Variables)
	}

	for _, targetName := range expectedResults.Targets {
		if _, ok := mf.Targets[targetName]; !ok {
			t.Errorf("Expected target not found: %s", targetName)
		}
	}
}
