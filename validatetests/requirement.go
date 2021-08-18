package main

import "container/list"

type RequirementItem struct {
	Id string
	Description string
}

// Requirements contains the list of formal requirements and associated tests
var Requirements = list.New()

// LoadRequirements reads a single CSV file with requirements into memory
// and returns associated RequirementItems
func LoadRequirements(contents []byte, idcolumn int, descriptioncolumn int) ([]RequirementItem, error) {

	return []RequirementItem{}, nil
}
