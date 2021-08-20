package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type RequirementResult struct {
	Id          string
	Description string
	Validated   bool
	Tests       []RequirementReference
}

var sourcefolder = flag.String("source", ".", "Source folder to scan for tests")
var requirementsfolder = flag.String("requirements", ".", "Folder containing CSV files with requirements")
var testresultsfolder = flag.String("testresults", ".", "Folder containing test result files")
var RequirementsIdColumn = flag.Int("id", 0, "Column (starting from 0) for ID in requirements")
var RequirementsDescriptionColumn = flag.Int("description", 1, "Column (starting from 0) for Description in requirements")
var RequirementsFirstLineIsHeader = flag.Bool("header", false, "Treat the first line of requirements CSVs as a header")
var TemplateFileName = flag.String("template", "report.txt", "Name of the template file to generate from")
var PrintJson = flag.Bool("json", false, "Print json output instead of text")

func main() {

	flag.Parse()

	IdColumn = *RequirementsIdColumn
	DescriptionColumn = *RequirementsDescriptionColumn
	RequirementsSkipFirstLine = *RequirementsFirstLineIsHeader

	TokenizeFolder(*sourcefolder)
	LoadRequirementsFolder(*requirementsfolder)
	LoadTestResults(*testresultsfolder)

	var data = []RequirementResult{}

	for e := Requirements.Front(); e != nil; e = e.Next() {
		requirement := e.Value.(RequirementItem)
		refs := ReferencesForRequirement(requirement.Id)
		d := RequirementResult{
			Id:          requirement.Id,
			Description: requirement.Description,
			Tests:       refs,
			Validated:   false,
		}

		if len(refs) > 0 {
			d.Validated = true
			for _, ref := range refs {
				if !TestResults[ref.FunctionName()] {
					d.Validated = false
				}
			}
		}

		data = append(data, d)
	}

	if *PrintJson {
		json, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(json))
	} else {
		templatepath, err := filepath.Abs(*TemplateFileName)
		if err != nil {
			log.Fatal(err)
		}

		contents, err := os.ReadFile(templatepath)
		t := template.New("report")
		t.Parse(string(contents))
		if err != nil {
			log.Fatal(err)
		}

		err = t.Execute(os.Stdout, data)
		if err != nil {
			log.Fatal(err)
		}
	}

}
