package main

import (
	"io/fs"
	"path/filepath"
)

var TestResults = map[string]bool{}

type TestResult struct {
	Name   string
	Passed bool
}

// LoadTestResults loads test results and puts them into
func LoadTestResults(folder string) error {
	return filepath.WalkDir(folder, TestResultsWalkDirFunc)
}

func TestResultsWalkDirFunc(path string, d fs.DirEntry, err error) error {

	if !d.Type().IsRegular() {
		return nil
	}

	var results map[string]bool

	switch filepath.Ext(path) {
	case ".trx":
		results, err = LoadTrxResults(path)
	case ".xml":
		results, err = LoadJUnitResults(path)
	}

	if err != nil {
		return err
	}
	for k, v := range results {
		TestResults[k] = v
	}

	return nil
}
