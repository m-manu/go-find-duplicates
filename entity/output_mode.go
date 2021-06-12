package entity

const (
	OutputModeCsvFile  = "csv"
	OutputModeTextFile = "text"
	OutputModeStdOut   = "print"
)

var OutputModes = map[string]string{
	OutputModeTextFile: "creates a text file in current directory with basic information",
	OutputModeCsvFile:  "creates a csv file in current directory with detailed information",
	OutputModeStdOut:   "just prints the report without creating any file",
}
