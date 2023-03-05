package utils

import (
	set "github.com/deckarep/golang-set/v2"
	"os"
	"path/filepath"
	"strings"
)

// IsReadableDirectory checks whether argument is a readable directory
func IsReadableDirectory(path string) bool {
	fileInfo, statErr := os.Stat(path)
	if statErr != nil {
		return false
	}
	return fileInfo.IsDir()
}

// IsReadableFile checks whether argument is a readable file
func IsReadableFile(path string) bool {
	fileInfo, statErr := os.Stat(path)
	if statErr != nil {
		return false
	}
	return fileInfo.Mode().IsRegular()
}

// LineSeparatedStrToMap converts a line-separated string to a map with keys and empty values
func LineSeparatedStrToMap(lineSeparatedString string) (entries set.Set[string], firstFew []string) {
	entries = set.NewThreadUnsafeSet[string]()
	firstFew = []string{}
	for _, e := range strings.Split(lineSeparatedString, "\n") {
		entries.Add(strings.TrimSpace(e))
		firstFew = append(firstFew, e)
	}
	if len(firstFew) > 3 {
		firstFew = firstFew[0:3]
	}
	entries.Each(func(e string) bool {
		if e == "" {
			entries.Remove(e)
		}
		return false
	})
	return
}

// GetFileExt gets extension of file, in lower case
func GetFileExt(path string) string {
	ext := filepath.Ext(path)
	return strings.ToLower(ext)
}
