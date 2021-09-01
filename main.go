package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/m-manu/go-find-duplicates/bytesutil"
	"github.com/m-manu/go-find-duplicates/entity"
	"github.com/m-manu/go-find-duplicates/fmte"
	"github.com/m-manu/go-find-duplicates/service"
	"github.com/m-manu/go-find-duplicates/utils"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// Exit codes for this program
const (
	exitCodeSuccess = iota
	exitCodeInvalidNumArgs
	exitCodeInvalidExclusions
	exitCodeInputDirectoryNotReadable
	exitCodeExclusionFilesError
	exitCodeErrorFindingDuplicates
	exitCodeErrorCreatingReport
	exitCodeInvalidOutputMode
	exitCodeReportFileCreationFailed
)

//go:embed default_exclusions.txt
var defaultExclusionsStr string

func setupExclusionsOpt() func() map[string]struct{} {
	const exclusionsFlag = "exclusions"
	const exclusionsDefaultValue = ""
	_, defaultExclusionsExamples := utils.LineSeparatedStrToMap(defaultExclusionsStr)
	excludesListFilePathPtr := flag.String(exclusionsFlag, exclusionsDefaultValue,
		fmt.Sprintf("path to file containing newline separated list of file/directory names to be excluded\n"+
			"(if this is not set, by default these will be ignored:\n%s etc.)",
			strings.Join(defaultExclusionsExamples, ", ")))
	return func() map[string]struct{} {
		excludesListFilePath := *excludesListFilePathPtr
		var exclusions map[string]struct{}
		if excludesListFilePath == exclusionsDefaultValue {
			defaultExclusions, _ := utils.LineSeparatedStrToMap(defaultExclusionsStr)
			exclusions = defaultExclusions
		} else {
			if !utils.IsReadableFile(excludesListFilePath) {
				fmte.PrintfErr("error: argument to flag -%s should be a file\n", exclusionsFlag)
				flag.Usage()
				os.Exit(exitCodeInvalidExclusions)
			}
			rawContents, err := os.ReadFile(excludesListFilePath)
			if err != nil {
				fmte.PrintfErr("error: argument to flag -%s isn't a readable file: %+v\n", exclusionsFlag, err)
				flag.Usage()
				os.Exit(exitCodeExclusionFilesError)
			}
			contents := strings.ReplaceAll(string(rawContents), "\r\n", "\n") // Windows
			exclusions, _ = utils.LineSeparatedStrToMap(contents)
		}
		return exclusions
	}
}

func setupHelpOpt() func() bool {
	helpPtr := flag.Bool("help", false, "display help")
	return func() bool {
		return *helpPtr
	}
}

func setupThoroughOpt() func() bool {
	thoroughPtr := flag.Bool("thorough", false, "apply thorough check of uniqueness of files\n"+
		"(caution: this makes the scan very slow!)")
	return func() bool {
		return *thoroughPtr
	}
}

func setupMinSizeOpt() func() int64 {
	fileSizeThresholdPtr := flag.Uint64("minsize", 4, "minimum size of file in KiB to consider")
	return func() int64 {
		return int64(*fileSizeThresholdPtr) * bytesutil.KIBI
	}
}

func setupParallelismOpt() func() int {
	const defaultParallelismValue = 0
	parallelismPtr := flag.Uint("parallelism", defaultParallelismValue,
		"extent of parallelism (defaults to number of cores minus 1)")
	return func() int {
		if *parallelismPtr == defaultParallelismValue {
			n := runtime.NumCPU()
			if n > 1 {
				return n - 1
			}
			return 1
		}
		return int(*parallelismPtr)
	}
}

func setupOutputModeOpt() func() string {
	var sb strings.Builder
	sb.WriteString("following modes are accepted:\n")
	for outputMode, description := range entity.OutputModes {
		sb.WriteString(fmt.Sprintf("%5s = %s\n", outputMode, description))
	}
	outputModeStrPtr := flag.String("output", entity.OutputModeTextFile, sb.String())
	return func() string {
		outputModeStr := strings.ToLower(strings.TrimSpace(*outputModeStrPtr))
		if _, exists := entity.OutputModes[outputModeStr]; !exists {
			fmt.Printf("error: invalid output mode '%s'\n", outputModeStr)
			os.Exit(exitCodeInvalidOutputMode)
		}
		return outputModeStr
	}
}

func setupUsage() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Run \"go-find-duplicates -help\" for usage\n")
	}
}

func readDirectories() (directories []string) {
	if flag.NArg() < 1 {
		fmte.Printf("error: no input directories passed\n")
		flag.Usage()
		os.Exit(exitCodeInvalidNumArgs)
	}
	for i, p := range flag.Args() {
		if !utils.IsReadableDirectory(p) {
			fmte.PrintfErr("error: input #%d \"%v\" isn't a readable directory\n", i+1, p)
			flag.Usage()
			os.Exit(exitCodeInputDirectoryNotReadable)
		}
		abs, _ := filepath.Abs(p)
		directories = append(directories, abs)
	}
	return directories
}

func handlePanic() {
	err := recover()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Program exited unexpectedly. "+
			"Please report the below eror to the author:\n"+
			"%+v\n", err)
		_, _ = fmt.Fprintln(os.Stderr, string(debug.Stack()))
	}
}

func showHelpAndExit() {
	flag.CommandLine.SetOutput(os.Stdout)
	fmt.Printf(`go-find-duplicates is a tool to find duplicate files and directories

Usage:
  go-find-duplicates [flags] <dir-1> <dir-2> ... <dir-n>

where,
  arguments are readable directories that need to be scanned for duplicates

Flags (all optional):
`)
	flag.PrintDefaults()
	fmt.Printf(`
For more details: https://github.com/m-manu/go-find-duplicates
`)
	os.Exit(exitCodeSuccess)
}

func main() {
	defer handlePanic()
	runID := time.Now().Format("150405")
	getExcludedFiles := setupExclusionsOpt()
	isHelp := setupHelpOpt()
	getMinSize := setupMinSizeOpt()
	getOutputMode := setupOutputModeOpt()
	getParallelism := setupParallelismOpt()
	isThorough := setupThoroughOpt()
	setupUsage()
	flag.Parse()
	if isHelp() {
		showHelpAndExit()
	}
	directories := readDirectories()
	outputMode := getOutputMode()
	reportFileName := createReportFileIfApplicable(runID, outputMode)
	duplicates, duplicateTotalCount, savingsSize, allFiles, fdErr :=
		service.FindDuplicates(directories, getExcludedFiles(), getMinSize(), getParallelism(), isThorough())
	if fdErr != nil {
		fmte.PrintfErr("error while finding duplicates: %+v", fdErr)
		os.Exit(exitCodeErrorFindingDuplicates)
	}
	if duplicates == nil || duplicates.Size() == 0 {
		fmte.Printf("No duplicates found!\n")
		return
	}
	fmte.Printf("Found %d duplicates. A total of %s can be saved by removing them.\n",
		duplicateTotalCount, bytesutil.BinaryFormat(savingsSize))
	reportDuplicates(duplicates, outputMode, allFiles, runID, reportFileName)
}

func createReportFileIfApplicable(runID string, outputMode string) (reportFileName string) {
	if outputMode == entity.OutputModeStdOut {
		return
	}
	if outputMode == entity.OutputModeCsvFile {
		reportFileName = fmt.Sprintf("./duplicates_%s.csv", runID)
	} else if outputMode == entity.OutputModeTextFile {
		reportFileName = fmt.Sprintf("./duplicates_%s.txt", runID)
	}
	_, err := os.Create(reportFileName)
	if err != nil {
		fmte.PrintfErr("error: couldn't create report file: %+v\n", err)
		os.Exit(exitCodeReportFileCreationFailed)
	}
	return
}
