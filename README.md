# Go Find Duplicates

[![Build Status](https://www.travis-ci.com/m-manu/go-find-duplicates.svg?branch=main)](https://www.travis-ci.com/m-manu/go-find-duplicates)
[![Go Report Card](https://goreportcard.com/badge/github.com/m-manu/go-find-duplicates)](https://goreportcard.com/report/github.com/m-manu/go-find-duplicates)
[![Go Reference](https://pkg.go.dev/badge/github.com/m-manu/go-find-duplicates.svg)](https://pkg.go.dev/github.com/m-manu/go-find-duplicates)
[![License](https://img.shields.io/badge/License-Apache%202-blue.svg)](./LICENSE)

## Introduction

A blazingly-fast simple-to-use tool to find duplicate files (photos, videos, music, documents etc.) on your computer,
portable hard drives etc.

## How to install and use?

Two ways: (one direct, one through docker)

### Direct

To install:

1. Install Go version at least **1.16**
    * On Ubuntu: `snap install go`
    * On Mac: `brew install go`
    * For any other OS: [Go downloads page](https://golang.org/dl/)
2. Run command:
   ```bash
   go get github.com/m-manu/go-find-duplicates
   ```
3. Ensure `$HOME/go/bin` is part of `$PATH`

To use:

```bash
go-find-duplicates {dir-1} {dir-2} ... {dir-n}
```

For more options and help, run:

```bash
go-find-duplicates -help
```

### Through Docker

```bash
docker run --rm -v /Volumes/PortableHD:/mnt/PortableHD manumk/go-find-duplicates:latest go-find-duplicates -output=print /mnt/PortableHD
```

In above command:

* option `--rm` removes the container when it exits
* option `-v` is mounts host directory `/Volumes/PortableHD` as `/mnt/PortableHD` inside the container
