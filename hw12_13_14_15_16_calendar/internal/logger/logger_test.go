package logger

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	logg := NewLogger("WARN")
	logg.output = &buf

	t.Run("logger levels", func(t *testing.T) {
		logg.Debug("do not show")
		require.Empty(t, buf)
		logg.Info("do not show")
		require.Empty(t, buf)
		logg.Warn("show this")
		require.Contains(t, buf.String(), "[WARN]: show this\n")
		buf.Reset()
		logg.Error("and this")
		require.Contains(t, buf.String(), "[ERROR]: and this\n")
	})
}

func BenchmarkLogger(b *testing.B) {
	logg := NewLogger("INFO")
	logg.output = io.Discard

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logg.Info("benchmark message Value: %d", b.N)
		}
	})
}
