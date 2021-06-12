package service

import (
	"github.com/m-manu/go-find-duplicates/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const exclusionsStr = `.DS_Store
VERSION
`

func TestFindDuplicates(t *testing.T) {
	goRoot, ok := os.LookupEnv("GOROOT")
	if !ok {
		assert.FailNow(t, "Can't run test as GOROOT is not set")
	}
	directories := []string{filepath.Join(goRoot, "pkg"), filepath.Join(goRoot, "src"), filepath.Join(goRoot, "test")}
	exclusions, _ := utils.LineSeparatedStrToMap(exclusionsStr)
	duplicates, duplicateCount, savingsSize, _, err := FindDuplicates(directories, exclusions,
		4_196, 2)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, duplicates.Size(), 0)
	assert.GreaterOrEqual(t, duplicateCount, int64(0))
	assert.GreaterOrEqual(t, savingsSize, int64(0))
}
