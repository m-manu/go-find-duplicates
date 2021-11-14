# Go Find Duplicates

[![build-and-test](https://github.com/m-manu/go-find-duplicates/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/m-manu/go-find-duplicates/actions/workflows/build-and-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/m-manu/go-find-duplicates)](https://goreportcard.com/report/github.com/m-manu/go-find-duplicates)
[![Go Reference](https://pkg.go.dev/badge/github.com/m-manu/go-find-duplicates.svg)](https://pkg.go.dev/github.com/m-manu/go-find-duplicates)
[![License](https://img.shields.io/badge/License-Apache%202-blue.svg)](./LICENSE)

## Introduction

A blazingly-fast simple-to-use tool to find duplicate files (photos, videos, music, documents etc.) on your computer,
portable hard drives etc.

## How to install?

1. Install Go version at least **1.16**
    * On Ubuntu: `snap install go`
    * On Mac: `brew install go`
    * For any other OS: [Go downloads page](https://golang.org/dl/)
2. Run command:
   ```bash
   go install github.com/m-manu/go-find-duplicates
   ```
3. Ensure `$HOME/go/bin` is part of `$PATH`

## How to use?

```bash
go-find-duplicates {dir-1} {dir-2} ... {dir-n}
```

Above command just creates a *duplicates report*. Note that this tool just reads your files. It does *not* delete or
otherwise modify your files in any way.

## Command line options

Running `go-find-duplicates -help` displays following:

```
go-find-duplicates is a tool to find duplicate files and directories

Usage:
  go-find-duplicates [flags] <dir-1> <dir-2> ... <dir-n>

where,
  arguments are readable directories that need to be scanned for duplicates

Flags (all optional):
  -exclusions string
    	path to file containing newline separated list of file/directory names to be excluded
    	(if this is not set, by default these will be ignored:
    	.DS_Store, System Volume Information, $RECYCLE.BIN etc.)
  -help
    	display help
  -minsize uint
    	minimum size of file in KiB to consider (default 4)
  -output string
    	following modes are accepted:
    	 text = creates a text file in current directory with basic information
    	  csv = creates a csv file in current directory with detailed information
    	print = just prints the report without creating any file
       json = creates a JSON file in the current directory with detailed information
    	 (default "text")
  -parallelism uint
    	extent of parallelism (defaults to number of cores minus 1)
  -thorough
    	apply thorough check of uniqueness of files
    	(caution: this makes the scan very slow!)

For more details: https://github.com/m-manu/go-find-duplicates
```

## Running this through a Docker container

```bash
docker run --rm -v /Volumes/PortableHD:/mnt/PortableHD manumk/go-find-duplicates:latest go-find-duplicates -output=print /mnt/PortableHD
```

In above command:

* option `--rm` removes the container when it exits
* option `-v` is mounts host directory `/Volumes/PortableHD` as `/mnt/PortableHD` inside the container

## How does this identify duplicates?

**By default**, this tool identifies duplicates if _all_ of the following conditions match:

1. file extension is same
2. file size is same
3. CRC32 hash of "crucial bytes" is same

If above default isn't enough for your requirements, you could use the command line option `-thorough` to switch to
SHA-256 hash of *entire file contents*. But remember, with this, scan becomes much slower!

When tested on my portable hard drive containing >172k files (videos, audio files, images and documents), with and
without `-thorough` option, the results were same!
