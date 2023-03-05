package service

import (
	set "github.com/deckarep/golang-set/v2"
	"github.com/m-manu/go-find-duplicates/entity"
	"github.com/m-manu/go-find-duplicates/fmte"
	"github.com/m-manu/go-find-duplicates/utils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

const exclusionsStr = ` .DS_Store 
vendor 

`

// TestFindDuplicates tests whether FindDuplicates returns a non-nil map of duplicates
func TestFindDuplicates(t *testing.T) {
	goRoot := runtime.GOROOT()
	directories := []string{
		filepath.Join(goRoot, "pkg"),
		filepath.Join(goRoot, "src"),
		filepath.Join(goRoot, "test"),
	}
	exclusions, _ := utils.LineSeparatedStrToMap(exclusionsStr)
	fmte.Off()
	duplicates, duplicateCount, savingsSize, _, err := FindDuplicates(directories, exclusions,
		4_196, 2, false)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, duplicates.Size(), 0)
	assert.GreaterOrEqual(t, duplicateCount, int64(0))
	assert.GreaterOrEqual(t, savingsSize, int64(0))
}

// TestNonThoroughVsNot checks whether FindDuplicates with 'thorough mode' on and off returns the same results
func TestNonThoroughVsNot(t *testing.T) {
	exclusions, _ := utils.LineSeparatedStrToMap(exclusionsStr)
	goRoot := []string{runtime.GOROOT()}
	fmte.Off()
	duplicatesExpected, duplicateCountExpected, savingsSizeExpected, _, tErr := FindDuplicates(goRoot, exclusions,
		4_196, 2, false)
	assert.Nil(t, tErr, "error while scanning for duplicates in GOROOT directory")
	duplicatesActual, duplicateCountActual, savingsSizeActual, _, ntErr := FindDuplicates(goRoot, exclusions,
		4_196, 5, true)
	assert.Nil(t, ntErr, "error while thoroughly scanning for duplicates in GOROOT directory")
	actualDuplicateFilePaths := extractFiles(duplicatesActual)
	expectedDuplicateFilePaths := extractFiles(duplicatesExpected)
	assert.True(t, actualDuplicateFilePaths.Equal(expectedDuplicateFilePaths), "Duplicate files differed between thorough and non-thorough modes")
	assert.Equal(t, duplicateCountExpected, duplicateCountActual, "Number of duplicates differed between thorough and non-thorough modes")
	assert.Equal(t, savingsSizeExpected, savingsSizeActual, "Savings expected differed between thorough and non-thorough modes")
}

func extractFiles(duplicatesExpected *entity.DigestToFiles) set.Set[string] {
	expectedDuplicatesFiles := set.NewThreadUnsafeSet[string]()
	for iter := duplicatesExpected.Iterator(); iter.HasNext(); {
		_, paths := iter.Next()
		for _, path := range paths {
			expectedDuplicatesFiles.Add(path)
		}
	}
	return expectedDuplicatesFiles
}
