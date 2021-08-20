package main

import (
	"flag"
)

var sourcefolder = flag.String("source", ".", "Source folder to scan for tests")
var requirementsfolder = flag.String("requirements", ".", "Folder containing CSV files with requirements")
var testresultsfolder = flag.String("testresults", ".", "Folder containing test result files")
var RequirementsIdColumn = flag.Int("id", 0, "Column (starting from 0) for ID in requirements")
var RequirementsDescriptionColumn = flag.Int("description", 1, "Column (starting from 0) for Description in requirements")

func main() {

	flag.Parse()

	IdColumn = *RequirementsIdColumn
	DescriptionColumn = *RequirementsDescriptionColumn

	TokenizeFolder(*sourcefolder)
	LoadRequirementsFolder(*requirementsfolder)
	LoadTestResults(*testresultsfolder)
	

}
