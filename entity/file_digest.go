package entity

import (
	"fmt"
	"github.com/m-manu/go-find-duplicates/bytesutil"
)

// FileDigest contains properties of a file that makes the file unique to a very high degree of confidence
type FileDigest struct {
	FileExtension string `json:"ext"`
	FileSize      int64  `json:"size"`
	FileHash      string `json:"hash"`
}

// String returns a string representation of FileDigest
func (f FileDigest) String() string {
	return fmt.Sprintf("%v/%v/%v", f.FileExtension, f.FileHash, bytesutil.BinaryFormat(f.FileSize))
}
