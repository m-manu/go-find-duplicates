package service

import (
	"github.com/m-manu/go-find-duplicates/bytesutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, int64(0), thresholdFileSize%(4*bytesutil.KIBI))
}

func TestGetDigest(t *testing.T) {
	goRoot, ok := os.LookupEnv("GOROOT")
	if !ok {
		assert.FailNow(t, "Can't run test as GOROOT is not set")
	}
	var paths = []string{
		goRoot + "/src/io/io.go",
		goRoot + "/src/io/pipe.go",
	}
	for _, path := range paths {
		digest, err := GetDigest(path)
		assert.Equal(t, nil, err)
		assert.Greater(t, digest.FileSize, int64(0))
		assert.Equal(t, 9, len(digest.FileFuzzyHash))
		assert.Greater(t, len(digest.FileExtension), 0)
	}
}
