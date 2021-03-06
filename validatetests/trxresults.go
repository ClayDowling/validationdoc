package main

import (
	"encoding/xml"
	"os"
	"strings"
)

// UnitTestResult holds the results from an actual unit test.
type UnitTestResult struct {
	TestName string `xml:"testName,attr"` // TestName is the concatenated class, test, and argument name of a particular test
	Outcome  string `xml:"outcome,attr"`  // Outcome, which is Passed, Failed, or Error
}

type ResultsType struct {
	Entries []UnitTestResult `xml:"UnitTestResult"`
}

// TestRun holds the results of a test run for a dotnet test file
type TestRun struct {
	Results ResultsType // Results of the individual test runs
}

// LoadTrxResults scans a Microsoft test results file and returns a map of all tests and whether they pass or fail.
// Test names are unique by class and method name.  For parameterized tests, the collection of all invocations is
// considered the same test.  Any result other than "Passed" across all invocations of a test method is considered
// a failing test.
func LoadTrxResults(filename string) (map[string]bool, error) {

	tests := map[string]bool{}

	content, err := os.ReadFile(filename)
	if err != nil {
		return tests, err
	}

	var run = TestRun{}
	err = xml.Unmarshal(content, &run)
	if err != nil {
		return tests, err
	}

	for _, t := range run.Results.Entries {
		name := t.TestName
		parenindex := strings.IndexRune(t.TestName, '(')
		if parenindex != -1 {
			name = name[:parenindex]
		}

		nameparts := strings.Split(name, ".")
		if len(nameparts) == 3 {
			name = nameparts[1] + "." + nameparts[2]
		}

		existingvalue, ok := tests[name]
		if ok && !existingvalue {
			continue
		}
		tests[name] = t.Outcome == "Passed"
	}

	return tests, nil
}
