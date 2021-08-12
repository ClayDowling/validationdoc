package main

import (
	"testing"

	"github.com/assertgo/assert"
)

func TestGetToken_GivenAClassName_ReturnsClassNameToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken("public class RumpleStiltskin : fairytale", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(ClassName)
	assert.ThatString(actual.Value).IsEqualTo("RumpleStiltskin")
}

func TestGetToken_GivenARequirement_ReturnsRequirementToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken(" /// Requirement CELL-1 ", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(Requirement)
	assert.ThatString(actual.Value).IsEqualTo("CELL-1")
}

func TestGetToken_GivenAMethodName_ReturnsMethodNameToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken(" public void ShouldReturnAMethodName  ( int number )  ", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(MethodName)
	assert.ThatString(actual.Value).IsEqualTo("ShouldReturnAMethodName")
}
