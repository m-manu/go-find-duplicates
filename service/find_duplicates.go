package service

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	set "github.com/deckarep/golang-set/v2"
	"github.com/m-manu/go-find-duplicates/bytesutil"
	"github.com/m-manu/go-find-duplicates/entity"
	"github.com/m-manu/go-find-duplicates/fmte"
	"github.com/m-manu/go-find-duplicates/utils"
)

// FindDuplicates finds duplicate files in a given set of directories and matching criteria
func FindDuplicates(directories []string, excludedFiles set.Set[string], fileSizeThreshold int64, parallelism int,
	isThorough bool) (
	duplicates *entity.DigestToFiles, duplicateTotalCount int64, savingsSize int64,
	allFiles entity.FilePathToMeta, err error,
) {
	fmte.Printf("Scanning %d directories...\n", len(directories))
	allFiles = make(entity.FilePathToMeta, 10_000)
	var totalSize int64
	for _, dirPath := range directories {
		size, pErr := populateFilesFromDirectory(dirPath, excludedFiles, fileSizeThreshold, allFiles)
		if pErr != nil {
			err = fmt.Errorf("error while scaning directory %s: %+v", dirPath, pErr)
			return
		}
		totalSize += size
	}
	fmte.Printf("Done. Found %d files of total size %s.\n", len(allFiles), bytesutil.BinaryFormat(totalSize))
	if len(allFiles) == 0 {
		return
	}
	fmte.Printf("Finding potential duplicates... \n")
	shortlist := identifyShortList(allFiles)
	if len(shortlist) == 0 {
		return
	}
	fmte.Printf("Completed. Found %d files that may have one or more duplicates!\n", len(shortlist))
	if isThorough {
		fmte.Printf("Thoroughly scanning for duplicates... \n")
	} else {
		fmte.Printf("Scanning for duplicates... \n")
	}
	var processedCount int32
	var wg sync.WaitGroup
	wg.Add(2)
	go func(pc *int32, fc int32) {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		for atomic.LoadInt32(pc) < fc {
			time.Sleep(2 * time.Second)
			progress := float64(atomic.LoadInt32(pc)) / float64(fc)
			fmte.Printf("%2.0f%% processed so far\n", progress*100.0)
		}
	}(&processedCount, int32(len(shortlist)))
	go func(p *int32) {
		defer wg.Done()
		duplicates = entity.NewDigestToFiles()
		computeDigestsAndGroupThem(shortlist, parallelism, p, duplicates, isThorough)
		for iter := duplicates.Iterator(); iter.HasNext(); {
			digest, files := iter.Next()
			numDuplicates := int64(len(files)) - 1
			duplicateTotalCount += numDuplicates
			savingsSize += numDuplicates * digest.FileSize
		}
	}(&processedCount)
	wg.Wait()
	fmte.Printf("Scan completed.\n")
	return
}

func computeDigestsAndGroupThem(shortlist entity.FileExtAndSizeToFiles, parallelism int,
	processedCount *int32, duplicates *entity.DigestToFiles, isThorough bool,
) {
	// Find potential duplicates:
	slKeys := make([]entity.FileExtAndSize, 0, len(shortlist))
	for extAndSize := range shortlist {
		slKeys = append(slKeys, extAndSize)
	}
	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func(shard int, wg *sync.WaitGroup, count *int32) {
			defer wg.Done()
			low := shard * len(slKeys) / parallelism
			high := (shard + 1) * len(slKeys) / parallelism
			for _, fileExtAndSize := range slKeys[low:high] {
				for _, path := range shortlist[fileExtAndSize] {
					digest, err := GetDigest(path, isThorough)
					if err != nil {
						fmte.Printf("error while scanning %s: %+v\n", path, err)
						continue
					}
					duplicates.Set(digest, path)
				}
				atomic.AddInt32(count, 1)
			}
		}(i, &wg, processedCount)
	}
	wg.Wait()
	// Remove non-duplicates
	var duplicateKeys []entity.FileDigest
	for iter := duplicates.Iterator(); iter.HasNext(); {
		digest, files := iter.Next()
		if len(files) <= 1 {
			duplicateKeys = append(duplicateKeys, *digest)
		}
	}
	for _, key := range duplicateKeys {
		duplicates.Remove(key)
	}
	return
}

// identifyShortList identifies the files that may have duplicates
func identifyShortList(filesAndMeta entity.FilePathToMeta) (shortlist entity.FileExtAndSizeToFiles) {
	shortlist = make(entity.FileExtAndSizeToFiles, len(filesAndMeta))
	// Group the files that have same extension and same size
	for path, meta := range filesAndMeta {
		fileExtAndSize := entity.FileExtAndSize{FileExtension: utils.GetFileExt(path), FileSize: meta.Size}
		shortlist[fileExtAndSize] = append(shortlist[fileExtAndSize], path)
	}
	// Remove non-duplicates
	for fileExtAndSize, paths := range shortlist {
		if len(paths) <= 1 {
			delete(shortlist, fileExtAndSize)
		}
	}
	return shortlist
}
