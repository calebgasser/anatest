package anatest

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// KeyValue represents a generic key-value pair used in 'deps' and 'compare'.
type KeyValue struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Test represents a single test case block.
type Test struct {
	URL     string     `yaml:"url"`
	Name    string     `yaml:"name"`
	Deps    []KeyValue `yaml:"deps"`
	Compare []KeyValue `yaml:"compare"`
}

// Config is the root structure that holds the list of tests.
type TestConfig struct {
	Tests []Test `yaml:"tests"`
}

func parseTests(testFile string) TestConfig {
	var testConfig TestConfig
	// Read yaml test config file

	yamlData, err := os.ReadFile(testFile)
	if err != nil {
		log.Fatalf("Failed to read test config: %s", err)
	}
	// Unmarshal the YAML byte slice into the 'config' struct.
	// The library uses the `yaml:` tags to map YAML keys to struct fields.
	err = yaml.Unmarshal([]byte(yamlData), &testConfig)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}

	// Print the parsed data to verify it works.
	fmt.Printf("Successfully parsed %d test(s).\n\n", len(testConfig.Tests))

	return testConfig
}
