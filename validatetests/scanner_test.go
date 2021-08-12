package main

import (
	"strings"
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

func TestGetToken_GivenRandomLines_ReturnsIgnoredLine(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken(" private void SetUp() {  ", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(IgnoredLine)
	assert.ThatString(actual.Value).IsEqualTo("")
}

func TestTokenizeStream_GivenTestClass_ReturnsTokensWithLineNumbers(t *testing.T) {

	assert := assert.New(t)
	r := strings.NewReader(
		`
using XUnit;
using FluentAssertions;
using Requirements; // Totally bogus

namespace Life {
	
	public class TestClass {
		
		[Fact]
		// Requirement Cell-5
		public void DeadCellsStayDead() {
			Cell c = new Cell(false);
			c.Lives(2).Should().BeFalse();
		}

		/// Requirement Cell-1
		[Theory]
		[InlineData(0)]
		[InlineData(1)]
		public void LiveCellsDieFromLoneliness(int n)
		{
			Cell c = new Cell(true);
			c.Lives(n).Should().BeFalse();
		}
	}
}`)

	tokens := TokenizeStream(r, Patterns)
	if len(tokens) != 5 {
		t.Errorf("Expected 5 tokens, got %d", len(tokens))
	}

	assert.ThatInt(int(tokens[0].Type)).IsEqualTo(ClassName)
	assert.ThatString(tokens[0].Value).IsEqualTo("TestClass")
	assert.ThatInt(tokens[0].Line).IsEqualTo(8)
	assert.ThatString(tokens[0].Filename).IsEmpty()

	assert.ThatInt(int(tokens[1].Type)).IsEqualTo(Requirement)
	assert.ThatString(tokens[1].Value).IsEqualTo("Cell-5")
	assert.ThatInt(tokens[1].Line).IsEqualTo(11)
	assert.ThatString(tokens[1].Filename).IsEmpty()

	assert.ThatInt(int(tokens[2].Type)).IsEqualTo(MethodName)
	assert.ThatString(tokens[2].Value).IsEqualTo("DeadCellsStayDead")
	assert.ThatInt(tokens[2].Line).IsEqualTo(12)
	assert.ThatString(tokens[2].Filename).IsEmpty()

	assert.ThatInt(int(tokens[3].Type)).IsEqualTo(Requirement)
	assert.ThatString(tokens[3].Value).IsEqualTo("Cell-1")
	assert.ThatInt(tokens[3].Line).IsEqualTo(17)
	assert.ThatString(tokens[3].Filename).IsEmpty()

	assert.ThatInt(int(tokens[4].Type)).IsEqualTo(MethodName)
	assert.ThatString(tokens[4].Value).IsEqualTo("LiveCellsDieFromLoneliness")
	assert.ThatInt(tokens[4].Line).IsEqualTo(21)
	assert.ThatString(tokens[4].Filename).IsEmpty()
}