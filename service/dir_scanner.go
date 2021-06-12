package service

import (
	"errors"
	"fmt"
	"github.com/m-manu/go-find-duplicates/entity"
	"github.com/m-manu/go-find-duplicates/fmte"
	"io/fs"
	"path/filepath"
	"strings"
)

func populateFilesFromDirectory(dirPathToScan string, exclusions map[string]struct{}, fileSizeThreshold int64,
	allFiles entity.FilePathToMeta) (
	sizeOfScannedFiles int64,
	err error,
) {
	wErr := filepath.WalkDir(dirPathToScan, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmte.PrintfErr("skipping \"%s\": %+v\n", path, errors.Unwrap(err))
			return nil
		}
		// If the file/directory is in excluded allFiles list, ignore it
		if _, exists := exclusions[d.Name()]; exists {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if _, exists := allFiles[path]; exists {
			return nil
		}
		// Ignore dot allFiles (Mac)
		if strings.HasPrefix(d.Name(), "._") {
			return nil
		}
		if d.Type().IsRegular() {
			info, infoErr := d.Info()
			if infoErr != nil {
				fmte.PrintfErr("couldn't get metadata of \"%s\": %+v\n", path, infoErr)
				return nil
			}
			if info.Size() < fileSizeThreshold {
				return nil
			}
			allFiles[path] = entity.FileMeta{Size: info.Size(), ModifiedTimestamp: info.ModTime().Unix()}
			sizeOfScannedFiles += info.Size()
		}
		return nil
	})
	if wErr != nil {
		return -1, fmt.Errorf("couldn't scan directory %s: %v", dirPathToScan, wErr)
	}
	return sizeOfScannedFiles, nil
}
