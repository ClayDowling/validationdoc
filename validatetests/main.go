package main

import (
	"flag"
)

var sourcefolder = flag.String("source", ".", "Source folder to scan for tests")
var requirementsfolder = flag.String("requirements", ".", "Folder containing CSV files with requirements")

func main() {

	TokenizeFolder(*sourcefolder)

}
