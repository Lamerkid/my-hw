package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	bar := ProgressBar{}
	bar.Start()
	defer bar.Cleanup()

	// check for regular file
	stat, _ := os.Stat(fromPath)
	if regular := stat.Mode().IsRegular(); !regular {
		return ErrUnsupportedFile
	}

	bar.Increment(20)

	// read from file
	inFile, err := os.ReadFile(fromPath)
	if err != nil {
		log.Panic("[ERROR] error reading file: ", err)
	}

	bar.Increment(10)

	// check for offset
	if int64(len(inFile)) < offset {
		return ErrOffsetExceedsFileSize
	}

	bar.Increment(10)

	// make default 0 limit as infinite
	if limit == 0 {
		limit = -1
	}

	bar.Increment(10)

	// create reader from byte[] file
	inMem := bytes.NewReader(inFile)

	bar.Increment(10)

	// sectionReader from bytes reader
	reader := io.NewSectionReader(inMem, offset, limit)

	bar.Increment(10)

	// make output dir
	if err := os.MkdirAll(filepath.Dir(toPath), 0o777); err != nil {
		log.Panic("[ERROR] error creating output dir: ", err)
	}

	bar.Increment(10)

	// make output file
	outFile, err := os.Create(toPath)
	if err != nil {
		log.Panic("[ERROR] error creating output file: ", err)
	}
	defer outFile.Close()

	bar.Increment(10)

	// write to output file
	if _, err = io.Copy(outFile, reader); err != nil {
		log.Panic("[ERROR] error writing to output file: ", err)
	}

	bar.Increment(10)
	bar.Finish()
	return nil
}
