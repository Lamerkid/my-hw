package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrNotADirectory      = errors.New("not a directory")
	ErrNoFilesInDirectory = errors.New("no files inside directory")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	stat, _ := os.Stat(dir)
	if isDirectory := stat.Mode().IsDir(); !isDirectory {
		return nil, ErrNotADirectory
	}

	env := make(Environment)
	if len(entries) == 0 {
		return nil, ErrNoFilesInDirectory
	}

	for _, e := range entries {
		stat, _ := e.Info()

		if stat.Size() == 0 {
			env[e.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		if regular := stat.Mode().IsRegular(); !regular {
			continue
		}

		fileBytes, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}

		byteReader := bytes.NewReader(fileBytes)
		var word string
		for {
			ch, _, err := byteReader.ReadRune()
			if err != nil || ch == '\n' {
				word = strings.TrimRight(word, ` \s`)
				break
			}
			if ch == 0x00 {
				ch = '\n'
			}
			word += string(ch)
		}
		if word == "" {
			env[e.Name()] = EnvValue{word, true}
		} else {
			env[e.Name()] = EnvValue{word, false}
		}
	}
	return env, nil
}
