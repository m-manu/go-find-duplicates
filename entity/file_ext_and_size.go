package entity

import "fmt"

// FileExtAndSize is a struct of file extension and file size
type FileExtAndSize struct {
	FileExtension string
	FileSize      int64
}

// String returns a string representation of FileExtAndSize
func (f FileExtAndSize) String() string {
	return fmt.Sprintf("%v/%v", f.FileExtension, f.FileSize)
}

// FileExtAndSizeToFiles is a multi-map of FileExtAndSize key and string values
type FileExtAndSizeToFiles map[FileExtAndSize][]string
