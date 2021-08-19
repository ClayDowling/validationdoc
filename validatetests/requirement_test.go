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

	actual, err := ParseRequirements(strings.NewReader(content))
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

func TestLoadRequirementsFolder_GivenRequirementsFolder_GetsAtLeastGameAndCellRequirements(t *testing.T) {

	err := LoadRequirementsFolder("../requirements")
	if err != nil {
		t.Fatal(err)
	}

	var foundGame = false
	var foundCell = false

	for e := Requirements.Front(); e != nil; e = e.Next() {
		r := e.Value.(RequirementItem)
		t.Logf("%s -> %s", r.Id, r.Description)
		if r.Id == "CELL-5" {
			foundCell = true
		}
		if r.Id == "GAME-3" {
			foundGame = true
		}
	}

	assert := assert.New(t)
	assert.ThatBool(foundGame).IsTrue()
	assert.ThatBool(foundCell).IsTrue()

}
