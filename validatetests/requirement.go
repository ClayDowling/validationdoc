package main

import (
	"bufio"
	"container/list"
	"encoding/csv"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type RequirementItem struct {
	Id          string
	Description string
}

// Requirements contains the list of formal requirements and associated tests
var Requirements = list.New()

var IdColumn = 0
var DescriptionColumn = 1
var RequirementsSkipFirstLine = false

// ParseRequirements reads a single CSV file with requirements into memory
// and returns associated RequirementItems
func ParseRequirements(src io.Reader) ([]RequirementItem, error) {

	reader := csv.NewReader(src)
	contents, err := reader.ReadAll()
	if err != nil {
		return []RequirementItem{}, err
	}

	var items = []RequirementItem{}

	for i := 0; i < len(contents); i++ {
		if i == 0 && RequirementsSkipFirstLine {
			continue
		}
		item := RequirementItem{
			Id:          strings.TrimSpace(contents[i][IdColumn]),
			Description: strings.TrimSpace(contents[i][DescriptionColumn]),
		}
		items = append(items, item)
	}

	return items, nil
}

func LoadRequirementsFolder(folder string) error {

	err := filepath.WalkDir(folder, LoadRequirementsWalkDirFunc)
	if err != nil {
		return err
	}

	return nil
}

func LoadRequirementsWalkDirFunc(path string, d fs.DirEntry, err error) error {

	if d.Type().IsRegular() && filepath.Ext(path) == ".csv" {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		items, err := ParseRequirements(reader)
		if err != nil {
			return err
		}
		for _, item := range items {
			Requirements.PushBack(item)
		}
	}

	return nil
}
