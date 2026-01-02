package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var offset int64
	var limit int64
	testDir := "testdata"
	tmpDir := t.TempDir()
	entries, _ := os.ReadDir(testDir)

	t.Run("not regular file test", func(t *testing.T) {
		offset = 10
		limit = 10
		if runtime.GOOS == "windows" {
			t.Skip("Skipping /dev/urandom test on windows")
		}
		tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
		err := Copy("/dev/urandom", tmpFile.Name(), offset, limit)
		require.Error(t, err, "must return ErrUnsupportedFile")
	})

	t.Run("offset more than file length", func(t *testing.T) {
		offset = 1000000
		limit = 1
		for _, e := range entries {
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)
			require.Error(t, err, "must return ErrOffsetExceedsFileSize")

			tmpFile.Close()
		}
	})

	t.Run("copy all test files", func(t *testing.T) {
		offset = 0
		limit = 0
		for _, e := range entries {
			inFileStat, _ := e.Info()
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)

			outFileStat, _ := tmpFile.Stat()
			require.NoError(t, err, "must return nil")
			require.Equal(t, inFileStat.Size(), outFileStat.Size(), "files must be the same size")

			tmpFile.Close()
		}
	})

	t.Run("test limit", func(t *testing.T) {
		offset = 0
		limit = 100
		for _, e := range entries {
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)

			outFileStat, _ := tmpFile.Stat()
			require.NoError(t, err, "must return nil")
			require.LessOrEqual(t, outFileStat.Size(), limit, "files must be less or equal to limit")

			tmpFile.Close()
		}
	})

	t.Run("test offset", func(t *testing.T) {
		offset = 100
		limit = 0
		for _, e := range entries {
			inFileStat, _ := e.Info()
			if inFileStat.Size() <= offset {
				continue
			}
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)

			outFileStat, _ := tmpFile.Stat()
			require.NoError(t, err, "must return nil")
			require.Equal(t, inFileStat.Size()-offset, outFileStat.Size(), "copied files must be lesser by offset")

			tmpFile.Close()
		}
	})
}
