package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type TestCase struct {
	Name      string `xml:"name,attr"`
	ClassName string `xml:"classname,attr"`
	Failure   string `xml:"failure"`
}

type TestSuite struct {
	Entries []TestCase `xml:"testcase"`
}

func LoadJUnitResults(filename string) (map[string]bool, error) {

	var results = map[string]bool{}
	content, err := os.ReadFile(filename)
	if err != nil {
		return results, err
	}

	var run = TestSuite{}
	err = xml.Unmarshal(content, &run)
	if err != nil {
		return results, err
	}

	for _, t := range run.Entries {
		name := t.Name
		parenindex := strings.IndexRune(name, '(')
		if parenindex != -1 {
			name = name[:parenindex]
		}
		fullname := fmt.Sprintf("%s.%s", t.ClassName, name)

		existingvalue, ok := results[fullname]
		if ok && !existingvalue {
			continue
		}
		results[fullname] = t.Failure == ""
	}

	return results, nil
}
