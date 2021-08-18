package main

import (
	"container/list"
	"encoding/csv"
	"io"
	"strings"
)

type RequirementItem struct {
	Id          string
	Description string
}

// Requirements contains the list of formal requirements and associated tests
var Requirements = list.New()

// ParseRequirements reads a single CSV file with requirements into memory
// and returns associated RequirementItems
func ParseRequirements(src io.Reader, idcolumn int, descriptioncolumn int) ([]RequirementItem, error) {

	reader := csv.NewReader(src)
	contents, err := reader.ReadAll()
	if err != nil {
		return []RequirementItem{}, err
	}

	var items = []RequirementItem{}

	for i := 0; i < len(contents); i++ {
		item := RequirementItem{
			Id:          strings.TrimSpace(contents[i][idcolumn]),
			Description: strings.TrimSpace(contents[i][descriptioncolumn]),
		}
		items = append(items, item)
	}

	return items, nil
}
