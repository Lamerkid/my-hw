package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)

	t.Run("empty command", func(t *testing.T) {
		code := RunCmd([]string{}, env)
		require.Equal(t, 1, code)
	})

	t.Run("invalid command", func(t *testing.T) {
		code := RunCmd([]string{"1j2k3"}, env)
		require.NotEqual(t, 0, code)
	})

	t.Run("success", func(t *testing.T) {
		code := RunCmd([]string{"echo"}, env)
		require.Equal(t, 0, code)
	})

	env["TEST"] = EnvValue{"test_value", false}

	t.Run("success with env", func(t *testing.T) {
		code := RunCmd([]string{"ls", "-la"}, env)
		require.Equal(t, 0, code)
		require.Contains(t, os.Environ(), "TEST=test_value")
	})
}
