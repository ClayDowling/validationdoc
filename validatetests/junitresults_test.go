package main

import (
	"os"
	"testing"

	"github.com/assertgo/assert"
)

func TestJUnitResults_GivenXmlFile_ReturnsListOfTests(t *testing.T) {
	actual, err := LoadJUnitResults("TEST-CellTest.xml")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.ThatInt(len(actual)).IsEqualTo(6)

	for k := range actual {
		assert.ThatString(k).DoesNotContain("(")
	}
}

func TestJUnitResults_GivenTestWithOnePassAndOneFail_MarksTestAsFailing(t *testing.T) {

	blob := `
	<?xml version="1.0" encoding="UTF-8"?>
	<testsuite name="CellTest" tests="19" skipped="0" failures="5" errors="0" timestamp="2021-08-23T14:26:55" hostname="CPX-V0N9APN9RPX" time="0.108">
	  <properties/>
	  <testcase name="DeadCell_WithExactly3Neighbors_Lives()" classname="CellTest" time="0.001"/>
	  <testcase name="LiveCell_With4OrMoreNeighbors_Dies(6)" classname="CellTest" time="0.009" />
	  <testcase name="LiveCell_With4OrMoreNeighbors_Dies(4)" classname="CellTest" time="0.009">
		<failure message="org.opentest4j.AssertionFailedError: &#13;&#10;Expecting value to be true but was false" type="org.opentest4j.AssertionFailedError">org.opentest4j.AssertionFailedError: 
	Expecting value to be true but was false">
			Drat
	    </failure>
	  </testcase>	
	  <testcase name="ThirdTest" classname="CellTest" time="0.009">
		<failure message="org.opentest4j.AssertionFailedError: &#13;&#10;Expecting value to be true but was false" type="org.opentest4j.AssertionFailedError">org.opentest4j.AssertionFailedError: 	Expecting value to be true but was false" >
		  Oops.
		</failure>
	  </testcase>	
	</testsuite>
`
	os.WriteFile("sample.xml", []byte(blob), 0644)

	actual, err := LoadJUnitResults("sample.xml")
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.ThatBool(actual["CellTest.DeadCell_WithExactly3Neighbors_Lives"]).IsTrue()
	assert.ThatBool(actual["CellTest.LiveCell_With4OrMoreNeighbors_Dies"]).IsFalse()
	assert.ThatBool(actual["CellTest.ThirdTest"]).IsFalse()

	os.Remove("sample.xml")

}
