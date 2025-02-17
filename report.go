package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/m-manu/go-find-duplicates/entity"
)

const bytesPerLineGuess = 500

func reportDuplicates(duplicates *entity.DigestToFiles, outputMode string, allFiles entity.FilePathToMeta,
	runID string, reportFile io.Writer) error {
	var err error
	if outputMode == entity.OutputModeStdOut {
		printReportToStdOut(runID, duplicates)
	} else if outputMode == entity.OutputModeTextFile {
		err = createTextFileReport(duplicates, reportFile)
	} else if outputMode == entity.OutputModeCsvFile {
		err = createCsvReport(duplicates, allFiles, reportFile)
	} else if outputMode == entity.OutputModeJSON {
		err = createJSONReport(duplicates, reportFile)
	}
	return err
}

func createTextFileReport(duplicates *entity.DigestToFiles, reportFile io.Writer) error {
	reportBB := getReportAsText(duplicates)
	_, rcErr := reportFile.Write(reportBB.Bytes())
	return rcErr
}

func getReportAsText(duplicates *entity.DigestToFiles) bytes.Buffer {
	var bb bytes.Buffer
	bb.Grow(duplicates.Size() * bytesPerLineGuess)
	for iter := duplicates.Iterator(); iter.HasNext(); {
		digest, paths := iter.Next()
		sort.Strings(paths)
		bb.WriteString(fmt.Sprintf("%s: %d duplicate(s)\n", digest, len(paths)-1))
		for _, path := range paths {
			bb.WriteString(fmt.Sprintf("\t%s\n", path))
		}
	}
	return bb
}

func printReportToStdOut(runID string, duplicates *entity.DigestToFiles) {
	reportBB := getReportAsText(duplicates)
	fmt.Printf(`
==========================
Report (run id %s)
==========================
`, runID)
	fmt.Println(reportBB.String())
}

func createCsvReport(duplicates *entity.DigestToFiles, allFiles entity.FilePathToMeta, reportFile io.Writer) error {
	var bb bytes.Buffer
	bb.Grow(duplicates.Size() * bytesPerLineGuess)
	cf := csv.NewWriter(&bb)
	_ = cf.Write([]string{"file hash", "file size", "last modified", "file path"})
	for iter := duplicates.Iterator(); iter.HasNext(); {
		digest, paths := iter.Next()
		for _, path := range paths {
			_ = cf.Write([]string{
				digest.FileHash,
				strconv.FormatInt(digest.FileSize, 10),
				time.Unix(allFiles[path].ModifiedTimestamp, 0).Format("02-Jan-2006 03:04:05 PM"),
				path,
			})
		}
	}
	cf.Flush()
	_, err := reportFile.Write(bb.Bytes())
	return err
}

func createJSONReport(duplicates *entity.DigestToFiles, reportFile io.Writer) error {
	type duplicateFile struct {
		entity.FileDigest
		Paths []string `json:"paths"`
	}
	var duplicatesToMarshall []duplicateFile
	for iter := duplicates.Iterator(); iter.HasNext(); {
		digest, paths := iter.Next()
		duplicatesToMarshall = append(duplicatesToMarshall, duplicateFile{
			*digest,
			paths,
		})
	}
	jsonBytes, err := json.Marshal(duplicatesToMarshall)
	if err != nil {
		return err
	}
	_, err = reportFile.Write(jsonBytes)
	return err
}
