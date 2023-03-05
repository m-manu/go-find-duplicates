package entity

import (
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

// DigestToFiles is a multi-map with FileDigest keys and string values.
// Writes to this is goroutine-safe.
type DigestToFiles struct {
	mx   *sync.Mutex
	data *treemap.Map
}

// FileDigestComparator is a comparator for FileDigest that compares FileSize, FileExtension and FileHash in that order
func FileDigestComparator(a, b any) int {
	fa := a.(FileDigest)
	fb := b.(FileDigest)
	if fa.FileSize < fb.FileSize {
		return 1
	} else if fa.FileSize > fb.FileSize {
		return -1
	} else {
		if fa.FileExtension < fb.FileExtension {
			return 1
		} else if fa.FileExtension > fb.FileExtension {
			return -1
		} else {
			if fa.FileHash < fb.FileHash {
				return 1
			} else if fa.FileHash > fb.FileHash {
				return -1
			} else {
				return 0
			}
		}
	}
}

// NewDigestToFiles creates new DigestToFiles
func NewDigestToFiles() (m *DigestToFiles) {
	return &DigestToFiles{
		data: treemap.NewWith(FileDigestComparator),
		mx:   &sync.Mutex{},
	}
}

// Set sets a value for the key
func (m *DigestToFiles) Set(key FileDigest, value string) {
	m.mx.Lock()
	valuesRaw, found := m.data.Get(key)
	var values []string
	if found {
		values = valuesRaw.([]string)
		values = append(values, value)
	} else {
		values = []string{value}
	}
	m.data.Put(key, values)
	m.mx.Unlock()
}

// Remove removes entry in the map
func (m *DigestToFiles) Remove(fd FileDigest) {
	m.data.Remove(fd)
}

// Size returns size of map
func (m *DigestToFiles) Size() int {
	return m.data.Size()
}

type digestToFilesIterator struct {
	iter treemap.Iterator
}

// Iterator returns an iterator for a DigestToFiles	map
func (m *DigestToFiles) Iterator() *digestToFilesIterator {
	return &digestToFilesIterator{m.data.Iterator()}
}

// HasNext returns true if there are more elements in the iterator
func (m *digestToFilesIterator) HasNext() bool {
	return m.iter.Next()
}

// Next returns the next element in the iterator
func (m *digestToFilesIterator) Next() (digest *FileDigest, paths []string) {
	fd := m.iter.Key().(FileDigest)
	filePaths := m.iter.Value().([]string)
	return &fd, filePaths
}
