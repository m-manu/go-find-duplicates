package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
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
	"strconv"
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
)

const helpStr = "Run \"go-find-duplicates -h\" for usage\n"

//go:embed default_exclusions.txt
var defaultExclusionsStr string

func setupExclusionsOpt() func() map[string]struct{} {
	const exclusionsFlag = "exclusions"
	const exclusionsDefaultValue = ""
	_, defaultExclusionsExamples := utils.LineSeparatedStrToMap(defaultExclusionsStr)
	excludesListFilePathPtr := flag.String(exclusionsFlag, exclusionsDefaultValue,
		fmt.Sprintf("path to file containing newline separated list of file/directory names to be excluded "+
			"(if this is not set, by default these will be ignored: %s etc.)",
			strings.Join(defaultExclusionsExamples, ", ")))
	return func() map[string]struct{} {
		excludesListFilePath := *excludesListFilePathPtr
		var exclusions map[string]struct{}
		if excludesListFilePath == exclusionsDefaultValue {
			defaultExclusions, _ := utils.LineSeparatedStrToMap(defaultExclusionsStr)
			exclusions = defaultExclusions
		} else {
			if !utils.IsReadableFile(excludesListFilePath) {
				fmte.PrintfErr("error: argument to flag -%s should be a file\n"+helpStr, exclusionsFlag)
				os.Exit(exitCodeInvalidExclusions)
			}
			rawContents, err := os.ReadFile(excludesListFilePath)
			if err != nil {
				fmte.PrintfErr("error: argument to flag -%s isn't readable: %+v\n"+helpStr, exclusionsFlag, err)
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
			fmt.Printf("error: invalid output mode '%s'\n"+helpStr, outputModeStr)
			os.Exit(exitCodeInvalidOutputMode)
		}
		return outputModeStr
	}
}

func setupUsage() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage:
  go-find-duplicates [flags] <dir-1> <dir-2> ... <dir-n>

where,
	arguments are readable directories that need to be scanned for duplicates

Flags (all optional):
`)
		flag.PrintDefaults()
	}
}

func readDirectories() (directories []string) {
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(exitCodeInvalidNumArgs)
	}
	for i, p := range flag.Args() {
		if !utils.IsReadableDirectory(p) {
			fmte.PrintfErr("error: input #%d \"%v\" isn't a readable directory\n"+helpStr, i+1, p)
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
	fmte.Printf(`go-find-duplicates is a tool to find duplicate files and directories
For more details: https://github.com/m-manu/go-find-duplicates
`)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage()
	os.Exit(exitCodeSuccess)
}

func reportDuplicates(duplicates *entity.DigestToFiles, outputMode string, allFiles entity.FilePathToMeta) {
	runID := time.Now().Format("150405")
	if outputMode == entity.OutputModeStdOut || outputMode == entity.OutputModeTextFile {
		var bb bytes.Buffer
		bb.Grow(duplicates.Size() * 400)
		for digest, paths := range duplicates.Map() {
			bb.WriteString(fmt.Sprintf("%s: %d duplicate(s)\n", digest, len(paths)-1))
			for _, path := range paths {
				bb.WriteString(fmt.Sprintf("\t%s\n", path))
			}
		}
		if outputMode == entity.OutputModeTextFile {
			reportFileName := fmt.Sprintf("./duplicates_%s.txt", runID)
			rcErr := os.WriteFile(reportFileName, bb.Bytes(), 0644)
			if rcErr != nil {
				fmte.PrintfErr("error while creating report file %s: %+v", reportFileName, rcErr)
				os.Exit(exitCodeErrorCreatingReport)
			}
			fmte.Printf("View duplicates report here: %s\n", reportFileName)
		} else if outputMode == entity.OutputModeStdOut {
			fmte.Printf(`
==========================
Report (run id %s)
==========================
`, runID)
		}
		fmte.Printf(bb.String())
	} else if outputMode == entity.OutputModeCsvFile {
		var bb bytes.Buffer
		bb.Grow(duplicates.Size() * 500)
		cf := csv.NewWriter(&bb)
		cf.Write([]string{"file hash", "file size", "last modified", "file path"})
		for digest, paths := range duplicates.Map() {
			for _, path := range paths {
				cf.Write([]string{
					digest.FileFuzzyHash,
					strconv.FormatInt(digest.FileSize, 10),
					time.Unix(allFiles[path].ModifiedTimestamp, 0).Format("02-Jan-2006 03:04:05 PM"),
					path,
				})
			}
		}
		cf.Flush()
		reportFileName := fmt.Sprintf("./duplicates_%s.csv", runID)
		os.WriteFile(reportFileName, bb.Bytes(), 0644)
		fmte.Printf("View duplicates report here: %s\n", reportFileName)
	}
}

func setupCmdOptions() (
	isHelp func() bool, getExcludedFiles func() map[string]struct{}, getMinSize func() int64,
	getOutputMode func() string, getParallelism func() int,
) {
	getExcludedFiles = setupExclusionsOpt()
	isHelp = setupHelpOpt()
	getMinSize = setupMinSizeOpt()
	getOutputMode = setupOutputModeOpt()
	getParallelism = setupParallelismOpt()
	setupUsage()
	return
}

func main() {
	defer handlePanic()
	isHelp, getExcludedFiles, getMinSize, getOutputMode, getParallelism := setupCmdOptions()
	flag.Parse()
	if isHelp() {
		showHelpAndExit()
	}
	outputMode := getOutputMode()
	directories := readDirectories()
	duplicates, duplicateTotalCount, savingsSize, allFiles, fdErr :=
		service.FindDuplicates(directories, getExcludedFiles(), getMinSize(), getParallelism())
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
	reportDuplicates(duplicates, outputMode, allFiles)
}
