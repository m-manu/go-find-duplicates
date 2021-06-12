package entity

import (
	"fmt"
	"time"
)

// FileMeta is a combination of file size and its modification timestamp
type FileMeta struct {
	Size              int64
	ModifiedTimestamp int64
}

func (f FileMeta) String() string {
	return fmt.Sprintf("{size: %d, modified: %v}", f.Size, time.Unix(f.ModifiedTimestamp, 0))
}

// FilePathToMeta is a map of file path to its FileMeta
type FilePathToMeta map[string]FileMeta
