package main

import (
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("No directory path to read environment variables is specified.")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		slog.Error("No command to execute is specified.")
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		slog.Error("error parsing envronment variables", "error", err)
		os.Exit(1)
	}

	returnCode := RunCmd(os.Args[2:], env)

	os.Exit(returnCode)
}
