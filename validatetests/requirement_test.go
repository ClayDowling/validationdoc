package main

import (
	"strings"
	"testing"

	"github.com/assertgo/assert"
)

func TestParseRequirements_GivenIdAndDescriptionColumn_ReturnsIdAndDescription(t *testing.T) {
	var content = `ID.1,  First Requirement
 ID.2,"Second ""Requirement"""
`

	actual, err := ParseRequirements(strings.NewReader(content), 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.ThatInt(len(actual)).IsEqualTo(2)
	assert.ThatString(actual[0].Id).IsEqualTo("ID.1")
	assert.ThatString(actual[0].Description).IsEqualTo("First Requirement")
	assert.ThatString(actual[1].Id).IsEqualTo("ID.2")
	assert.ThatString(actual[1].Description).IsEqualTo(`Second "Requirement"`)
}
