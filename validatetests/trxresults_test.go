package main

import (
	"testing"
)

func TestTrxResults_GivenTrxFile_ReturnsListOfTests(t *testing.T) {
	actual, err := TrxResults("../c-sharp/gameoflife-tests/TestResults/testresults.trx")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range actual {
		t.Logf("UnitTestResult %s: %v", k, v)
	}

	if len(actual) != 3 {
		t.Errorf("Expected 3 entries, found %d", len(actual))
	}
}
