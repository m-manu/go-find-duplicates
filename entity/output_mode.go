package entity

// Different output modes
const (
	OutputModeTextFile = "text"
	OutputModeCsvFile  = "csv"
	OutputModeStdOut   = "print"
	OutputModeJSON     = "json"
)

// OutputModes and their brief descriptions
var OutputModes = map[string]string{
	OutputModeStdOut:   "just prints the report without creating any file",
	OutputModeTextFile: "creates a text file in the output directory with basic information",
	OutputModeCsvFile:  "creates a csv file in the output directory with detailed information",
	OutputModeJSON:     "creates a JSON file in the output directory with basic information",
}
