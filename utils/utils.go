package utils

import (
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
func LineSeparatedStrToMap(lineSeparatedString string) (set Set[string], firstFew []string) {
	set = NewSet[string]()
	firstFew = []string{}
	for _, e := range strings.Split(lineSeparatedString, "\n") {
		set.Add(e)
		firstFew = append(firstFew, e)
	}
	if len(firstFew) > 3 {
		firstFew = firstFew[0:3]
	}
	for e := range set {
		if strings.TrimSpace(e) == "" {
			set.Delete(e)
		}
	}
	return
}

// GetFileExt gets extension of file, in lower case
func GetFileExt(path string) string {
	ext := filepath.Ext(path)
	return strings.ToLower(ext)
}
