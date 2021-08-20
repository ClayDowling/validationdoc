package main

import (
	"os"
	"strings"
	"testing"

	"github.com/assertgo/assert"
)

func TestTrxResults_GivenTrxFile_ReturnsListOfTests(t *testing.T) {
	actual, err := LoadTrxResults("testresults.trx")
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 5 {
		t.Errorf("Expected 3 entries, found %d", len(actual))
	}

	for k := range actual {
		if strings.ContainsRune(k, '(') {
			t.Errorf("Found name with parens '%s'", k)
		}
	}
}

func TestTrxResults_GivenTestWithOnePassAndOneFail_MarksTestAsFailing(t *testing.T) {

	blob := `
<TestRun>
  <Results>
    <UnitTestResult testName="Namespace.TestClass.FirstUnitTest(arg: 1)" outcome="Passed" />
	<UnitTestResult testName="Namespace.TestClass.FirstUnitTest(arg: 2)" outcome="Failed" />
	<UnitTestResult testName="Namespace.TestClass.SecondUnitTest" outcome="Passed" />
	<UnitTestResult testName="Namespace.TestClass.ThirdUnitTest(arg: 3)" outcome="Passed" />
	<UnitTestResult testName="Namespace.TestClass.ThirdUnitTest(arg: 4)" outcome="Error" />
  </Results>
</TestRun>
`
	os.WriteFile("sample.trx", []byte(blob), 0644)

	actual, err := LoadTrxResults("sample.trx")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.ThatBool(actual["TestClass.FirstUnitTest"]).IsFalse()
	assert.ThatBool(actual["TestClass.SecondUnitTest"]).IsTrue()
	assert.ThatBool(actual["TestClass.ThirdUnitTest"]).IsFalse()

	os.Remove("sample.trx")

}
