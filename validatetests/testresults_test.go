package main

import (
	"testing"

	"github.com/assertgo/assert"
)

func TestLoadTestResults_GivenFolderWithTestResults_PopulatesTestResults(t *testing.T) {
	LoadTestResults(".")

	assert := assert.New(t)
	assert.ThatInt(len(TestResults)).IsGreaterThan(0)
}

func TestLoadTestResults_GivenMultipleFiles_PopulatesTestResultsFromAllFiles(t *testing.T) {
	LoadTestResults(".")

	assert := assert.New(t)

	value, present := TestResults["CellTest.LiveCell_WithTwoOrThreeNeighbors_Lives"]
	assert.ThatBool(value).IsTrue()
	assert.ThatBool(present).IsTrue()

	value, present = TestResults["CellTests.Lives_GivenZeroNeighbors_ReturnsFalse"]
	assert.ThatBool(value).IsTrue()
	assert.ThatBool(present).IsTrue()

}
