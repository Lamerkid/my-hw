package main

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		slog.Error("no command to execute found")
		return 1
	}
	commandName := cmd[0]
	args := cmd[1:]

	for k, v := range env {
		err := os.Unsetenv(k)
		if err != nil {
			slog.Error("error unseting env", "error", err)
			return 1
		}

		if !v.NeedRemove {
			err = os.Setenv(k, v.Value)
			if err != nil {
				slog.Error("error seting env", "error", err)
				return 1
			}
		}
	}

	command := exec.Command(commandName, args...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return e.ExitCode()
		}
		slog.Error("error executing command", "error", err)
		return 1
	}
	return command.ProcessState.ExitCode()
}
