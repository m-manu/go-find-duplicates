package entity

import (
	"fmt"
	"github.com/m-manu/go-find-duplicates/bytesutil"
)

// FileDigest contains properties of a file that makes the file unique to a very high degree of confidence
type FileDigest struct {
	FileExtension string
	FileSize      int64
	FileHash      string
}

func (f FileDigest) String() string {
	return fmt.Sprintf("%v/%v/%v", f.FileExtension, f.FileHash, bytesutil.BinaryFormat(f.FileSize))
}
