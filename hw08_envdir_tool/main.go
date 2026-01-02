package main

import (
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("No directory path to read environment variables is specified.")
		return
	}

	if len(os.Args) < 3 {
		slog.Error("No command to execute is specified.")
		return
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		slog.Error("error parsing envronment variables", "error", err)
	}

	returnCode := RunCmd(os.Args[2:], env)

	slog.Info("Program exited with", "code", returnCode)
}
