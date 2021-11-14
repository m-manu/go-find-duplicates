package entity

// Different output modes
const (
	OutputModeCsvFile  = "csv"
	OutputModeTextFile = "text"
	OutputModeStdOut   = "print"
	OutputModeJSON     = "json"
)

// OutputModes and their brief descriptions
var OutputModes = map[string]string{
	OutputModeTextFile: "creates a text file in current directory with basic information",
	OutputModeCsvFile:  "creates a csv file in current directory with detailed information",
	OutputModeStdOut:   "just prints the report without creating any file",
	OutputModeJSON:     "creates a JSON file in the current directory with detailed information",
}
