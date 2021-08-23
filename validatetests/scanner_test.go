package main

import (
	"os"
	"testing"

	"github.com/assertgo/assert"
)

const fileContents = `
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
}`

func TestGetToken_GivenAClassName_ReturnsClassNameToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken("public class RumpleStiltskin : fairytale", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(ClassName)
	assert.ThatString(actual.Value).IsEqualTo("RumpleStiltskin")
}

func TestGetToken_GivenARequirement_ReturnsRequirementToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken(" /// Requirement CELL-1 ", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(RequirementLabel)
	assert.ThatString(actual.Value).IsEqualTo("CELL-1")
}

func TestGetToken_GivenAMethodName_ReturnsMethodNameToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken(" public void ShouldReturnAMethodName  ( int number )  ", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(MethodName)
	assert.ThatString(actual.Value).IsEqualTo("ShouldReturnAMethodName")
}

func TestGetToken_GivenDefaultVisibilityTestFunction_ReturnsMethodNameToken(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken("    void DifferentTestMethod ()", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(MethodName)
	assert.ThatString(actual.Value).IsEqualTo("DifferentTestMethod")
}

func TestGetToken_GivenPrivateLines_ReturnsIgnoredLine(t *testing.T) {
	assert := assert.New(t)

	actual := GetToken(" private void SetUp() {  ", Patterns)

	assert.ThatInt(int(actual.Type)).IsEqualTo(IgnoredLine)
	assert.ThatString(actual.Value).IsEqualTo("")
}

func TestTokenizeStream_GivenTestClass_ReturnsTokensWithLineNumbers(t *testing.T) {

	assert := assert.New(t)

	tokens := TokenizeStream([]byte(fileContents), Patterns)
	if len(tokens) != 5 {
		t.Errorf("Expected 5 tokens, got %d", len(tokens))
	}

	assert.ThatInt(int(tokens[0].Type)).IsEqualTo(ClassName)
	assert.ThatString(tokens[0].Value).IsEqualTo("TestClass")
	assert.ThatInt(tokens[0].Line).IsEqualTo(8)
	assert.ThatString(tokens[0].Filename).IsEmpty()

	assert.ThatInt(int(tokens[1].Type)).IsEqualTo(RequirementLabel)
	assert.ThatString(tokens[1].Value).IsEqualTo("Cell-5")
	assert.ThatInt(tokens[1].Line).IsEqualTo(11)
	assert.ThatString(tokens[1].Filename).IsEmpty()

	assert.ThatInt(int(tokens[2].Type)).IsEqualTo(MethodName)
	assert.ThatString(tokens[2].Value).IsEqualTo("DeadCellsStayDead")
	assert.ThatInt(tokens[2].Line).IsEqualTo(12)
	assert.ThatString(tokens[2].Filename).IsEmpty()

	assert.ThatInt(int(tokens[3].Type)).IsEqualTo(RequirementLabel)
	assert.ThatString(tokens[3].Value).IsEqualTo("Cell-1")
	assert.ThatInt(tokens[3].Line).IsEqualTo(17)
	assert.ThatString(tokens[3].Filename).IsEmpty()

	assert.ThatInt(int(tokens[4].Type)).IsEqualTo(MethodName)
	assert.ThatString(tokens[4].Value).IsEqualTo("LiveCellsDieFromLoneliness")
	assert.ThatInt(tokens[4].Line).IsEqualTo(21)
	assert.ThatString(tokens[4].Filename).IsEmpty()
}

func TestTokenizeFile_GivenAFileInThisFolder_ReturnsTokensWithCorrectFileName(t *testing.T) {
	os.WriteFile("sample.cs", []byte(fileContents), 0644)

	assert := assert.New(t)
	tokens, err := TokenizeFile("sample.cs")
	if err != nil {
		t.Fatal(err)
	}

	assert.ThatInt(len(tokens)).IsEqualTo(5)

	for _, t := range tokens {
		assert.ThatString(t.Filename).IsEqualTo("validatetests/sample.cs")
	}

	err = os.Remove("sample.cs")
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseTokens_GivenSampleTokens_ProducesTwoRequirementReferences(t *testing.T) {
	expected := []RequirementReference{
		{
			FileName:   "validatetests/sample.cs",
			ClassName:  "TestClass",
			MethodName: "DeadCellsStayDead",
			Id:         "Cell-5",
			Line:       11,
		},
		{
			FileName:   "validatetests/sample.cs",
			ClassName:  "TestClass",
			MethodName: "LiveCellsDieFromLoneliness",
			Id:         "Cell-1",
			Line:       17,
		},
	}

	os.WriteFile("sample.cs", []byte(fileContents), 0644)
	tokens, err := TokenizeFile("sample.cs")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove("sample.cs")
	if err != nil {
		t.Fatal(err)
	}

	actual := ParseTokens(tokens)

	if len(actual) != 2 {
		t.Fatalf("Expected 2 RequirementReferences, got %d", len(actual))
	}

	for i := 0; i < 2; i++ {
		if actual[i].FileName != expected[i].FileName {
			t.Errorf("Expected Filename %s, got %s", expected[i].FileName, actual[i].FileName)
		}
		if actual[i].ClassName != expected[i].ClassName {
			t.Errorf("Expected ClassName %s, got %s", expected[i].ClassName, actual[i].ClassName)
		}
		if actual[i].MethodName != expected[i].MethodName {
			t.Errorf("Expected MethodName %s, got %s", expected[i].MethodName, actual[i].MethodName)
		}
		if actual[i].Id != expected[i].Id {
			t.Errorf("Expected Id %s, got %s", expected[i].Id, actual[i].Id)
		}
		if actual[i].Line != expected[i].Line {
			t.Errorf("Expected Line %d, got %d", expected[i].Line, actual[i].Line)
		}
	}

}

func TestParseTokens_GivenTwoRequirementsForAMethod_ReturnsTwoEntiresForThatMethod(t *testing.T) {
	tokens := []Token{
		{
			Type:     ClassName,
			Value:    "MyClass",
			Filename: "test1.cs",
		},
		{
			Type:     RequirementLabel,
			Value:    "REQ-1",
			Line:     7,
			Filename: "test1.cs",
		},
		{
			Type:     RequirementLabel,
			Value:    "REQ-2",
			Line:     8,
			Filename: "test1.cs",
		},
		{
			Type:     MethodName,
			Value:    "MyTest",
			Line:     9,
			Filename: "test1.cs",
		},
	}

	assert := assert.New(t)
	actual := ParseTokens(tokens)

	if len(actual) != 2 {
		t.Fatalf("Expected 2 RequirementReferences, found %d", len(actual))
	}

	assert.ThatString(actual[0].Id).IsEqualTo("REQ-1")
	assert.ThatString(actual[1].Id).IsEqualTo("REQ-2")

	assert.ThatString(actual[0].MethodName).IsEqualTo("MyTest")
	assert.ThatString(actual[1].MethodName).IsEqualTo("MyTest")

}

func TestTokenizeFolder_GivenCSharpGameOfLife_ReturnsAtLeastTwoFiles(t *testing.T) {
	err := TokenizeFolder("../c-sharp/gameoflife-tests")
	if err != nil {
		t.Fatal(err)
	}

	files := map[string]int{}
	for ri := References.Front(); ri != nil; ri = ri.Next() {
		r := ri.Value.(RequirementReference)
		files[r.FileName] = files[r.FileName] + 1
	}

	assert := assert.New(t)
	assert.ThatInt(len(files)).IsGreaterOrEqualTo(2)
}

func TestFunctionName_WhenClassAndMethodArePresent_ReturnsClassAndMethodConcatenated(t *testing.T) {

	rr := RequirementReference{
		ClassName:  "MyClass",
		MethodName: "FirstTest",
	}

	assert := assert.New(t)
	assert.ThatString(rr.FunctionName()).IsEqualTo("MyClass.FirstTest")

}

func TestFunctionName_WhenNoClassNamePresent_ReturnsMethodName(t *testing.T) {
	rr := RequirementReference{
		MethodName: "BareFunction",
	}

	assert := assert.New(t)
	assert.ThatString(rr.FunctionName()).IsEqualTo("BareFunction")
}

func TestReferencesForRequirement_WhenTwoReferencesPresent_ReturnsReferences(t *testing.T) {

	var refs = []RequirementReference{
		{
			Id:         "R1",
			FileName:   "file.cs",
			MethodName: "FirstMethod",
			ClassName:  "TestClass",
			Line:       18,
		},
		{
			Id:         "R2",
			FileName:   "otherfile.cs",
			MethodName: "OtherMethod",
			ClassName:  "OtherTestClass",
			Line:       23,
		},
		{
			Id:         "R1",
			FileName:   "file.cs",
			MethodName: "SecondMethod",
			ClassName:  "TestClass",
			Line:       18,
		},
	}

	for _, r := range refs {
		References.PushBack(r)
	}

	actual := ReferencesForRequirement("R1")

	assert := assert.New(t)
	assert.ThatInt(len(actual)).IsEqualTo(2)
	assert.ThatString(actual[0].Id).IsEqualTo("R1")
	assert.ThatString(actual[1].Id).IsEqualTo("R1")
	assert.ThatString(actual[0].MethodName).IsNotEqualTo(actual[1].MethodName)

}

func TestReferencesForRequirement_WhenNoReferencePresent_ReturnsEmptyArray(t *testing.T) {
	References.Init()

	actual := ReferencesForRequirement("R1")

	assert := assert.New(t)
	assert.ThatInt(len(actual)).IsEqualTo(0)
}
