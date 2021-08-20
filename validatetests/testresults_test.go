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
