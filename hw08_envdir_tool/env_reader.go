package main

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
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
	env := make(Environment)
	reg := regexp.MustCompile(`[\s||\p{P}$+<=>^|~]`)
	emptySpaces := regexp.MustCompile(`\s`)
	stat, _ := os.Stat(dir)
	if isDirectory := stat.Mode().IsDir(); !isDirectory {
		return nil, ErrNotADirectory
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("Error reading directory", "error", err)
		return nil, err
	}

	for _, e := range entries {
		stat, _ := e.Info()
		if regular := stat.Mode().IsRegular(); !regular {
			continue
		}

		fileBytes, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			slog.Error("Error opening file", "error", err)
			return nil, err
		}

		byteReader := bytes.NewReader(fileBytes)
		var word string
		for {
			ch, _, err := byteReader.ReadRune()
			if err != nil || ch == '\n' || ch == 0x00 {
				break
			}
			word += string(ch)
		}

		regWord := reg.ReplaceAllString(word, "")
		if emptySpaces.MatchString(regWord) || regWord == "" {
			continue
		}

		env[e.Name()] = EnvValue{regWord, false}
	}
	return env, nil
}
