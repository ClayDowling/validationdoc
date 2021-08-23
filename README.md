# Validation Doc

ValidationDoc provides an automated way to document that you have automated tests validating your requirements.  To use it you require three things:

1. Requirements exports as a CSV file.  It is only necessary to include the unique ID of the requirement and the description.  Other fields will be ignored.
2. Automated tests which validate your requirements.  Before each test method, add a comment of the form `// Requirement XXX` where XXX is a unique ID from your requirements document.
3. Test output, usually as an XML file.

There is a working example in the [c-sharp](c-sharp) and [java](java) folders of this repository.

    .\validationdoc.exe --requirements ../requirements --source ..\c-sharp\gameoflife-tests\ --testresults 
      ..\c-sharp\gameoflife-tests\TestResults\ --template ./report.txt --header > validation.md

This will generate `validation.md`, which produces a Markdown file similar to:

### Validation Report

| Id  |  Description | Validated | Tests |
|-----|--------------|-----------|-------|
| CELL-1 | Live cells with fewer than 2 neighbors die, as from loneliness. | true |  c-sharp/gameoflife-tests/CellTests.cs:12  |
| CELL-2 | Live cells with 2 or 3 neighbors live on | true |  c-sharp/gameoflife-tests/CellTests.cs:22  |
| CELL-3 | Dead cells with exactly three neighbors become alive, as from reproduction. | true |  c-sharp/gameoflife-tests/CellTests.cs:30  |
| CELL-4 | Live cells with 4 or more neighbors die, as from overpopulation. | true |  c-sharp/gameoflife-tests/CellTests.cs:42  |
| CELL-5 | Dead cells with more or less than 3 neighbors stay dead | false |  c-sharp/gameoflife-tests/CellTests.cs:48  |
| GAME-1 | The board starts with a randomly determined initial state | false |  c-sharp/gameoflife-tests/BoardTests.cs:17  |
| GAME-2 | On each generation, the cell rules are applied to each cell simultaneously to create the next generation. | false |  |
| GAME-3 | The game ends when all cells have died. | false |  |


## JUnit5 And Parameterized Tests

For reasons known only to the authors of JUnit5, they have changed the way that parameterized tests report their results in the xml results file.  Rather than providing a method name, it provides an index into the list of possible values.  Which is cute if you have more than one parameterized test.  The solution is to use the name parameter to give the results the correct name.  You can see how to manage it in the example tests.


# Building

To build validationdoc you will need a [go develpment environment](https://golang.org).  Version 1.16 is required.  Then type

    cd validatetests
    go build

This should produce an executable for your operating system.  If you need to build for other systems, consult the excellent Go documentation on cross compiling.

# Generating the Report

The report is created from a template in `report.txt`, which you should distribute with your executable.  You can change the generated output using any of the capabilities of the [templating system](https://pkg.go.dev/text/template).

The data structure being templated is an array of [RequirementResult](main.go) objects, with the following properties.

| Name | Type | Source |
|------|------|--------|
| .Id   | string |Id from the requirements CSV files |
| .Description | string | Description from the requirements CSV files |
| .Validated | bool | True if all of the tests which referenced this requirement passed, provided there is at least one test |
| .Tests | array of TestReference | See below |

Each RequirementResult has zero, one, or many [TestReference](scanner.go) objects associated with it.  You can access the list 

    {{ range .Tests }}  Template Stuff Here {{ end }}

| Name | Type | Source |
|------|------|--------|
| .Id  | string | Identifies which requirement this test validates |
| .FileName | string | Path from the top of the repository which this test is found |
| .ClassName | string | Class which holds this test |
| .MethodName | string | Method name for this test |
| .Line | integer | Line number in the source file where the method is implemented |
