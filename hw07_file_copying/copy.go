package main

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// check file stats
	stat, _ := os.Stat(fromPath)
	if regular := stat.Mode().IsRegular(); !regular {
		return ErrUnsupportedFile
	}
	if size := stat.Size(); size < offset {
		return ErrOffsetExceedsFileSize
	}

	// read from source file
	inFile, err := os.ReadFile(fromPath)
	if err != nil {
		slog.Error("error reading from file", "error", err)
		os.Exit(1)
	}

	// respect the offset
	inFile = inFile[offset:]
	inMem := bytes.NewReader(inFile)

	// create output file
	outFile, err := os.Create(toPath)
	if err != nil {
		slog.Error("error creating output file", "error", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// init progress bar
	bar := ProgressBar{}
	if limit == 0 || limit > int64(len(inFile)) {
		bar.Start(len(inFile))
	} else {
		bar.Start(int(limit))
	}
	defer bar.Cleanup()

	// copy function
	var bufSize int64 = 16
	var bytesRead int64
	// while limit is not exceeded
	for limit == 0 || limit != bytesRead {
		bytesToCopy := bufSize
		// if limit is less then buffer
		if limit != 0 && bytesToCopy > limit-bytesRead {
			bytesToCopy = limit - bytesRead
		}
		bytesCopied, err := io.CopyN(outFile, inMem, bytesToCopy)
		// EOF - write remaining
		if err == io.EOF {
			bytesRead += bytesCopied
			bar.Increment(int(bytesCopied))
			bar.Finish()
			return nil
		}
		if err != nil {
			slog.Error("error copying file", "error", err)
			os.Exit(1)
		}
		bytesRead += bytesCopied
		bar.Increment(int(bytesCopied))
	}

	bar.Finish()
	return nil
}
