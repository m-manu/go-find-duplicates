package entity

import (
	"sync"
)

// DigestToFiles is a multi-map with FileDigest keys and string values.
// Writes to this is goroutine-safe.
type DigestToFiles struct {
	mx   *sync.Mutex
	data map[FileDigest][]string
}

// NewDigestToFiles creates new DigestToFiles
func NewDigestToFiles(size int) (m *DigestToFiles) {
	return &DigestToFiles{
		data: make(map[FileDigest][]string, size),
		mx:   &sync.Mutex{},
	}
}

// Set sets a value for the key
func (m *DigestToFiles) Set(key FileDigest, value string) {
	m.mx.Lock()
	m.data[key] = append(m.data[key], value)
	m.mx.Unlock()
}

// Remove removes entry in the map
func (m *DigestToFiles) Remove(fd FileDigest) {
	delete(m.data, fd)
}

// Map returns internal map to iterate over
func (m *DigestToFiles) Map() map[FileDigest][]string {
	return m.data
}

// Size returns size of map
func (m *DigestToFiles) Size() int {
	return len(m.data)
}
