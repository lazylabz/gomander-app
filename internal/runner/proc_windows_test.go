//go:build windows

package runner

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowsProcessUsesTerminalOutput(t *testing.T) {
	if !isConPTYAvailable() {
		t.Skip("ConPTY is unavailable on this Windows version")
	}

	process := newCommandProcess(
		`powershell -NoProfile -Command Write-Output ([Console]::IsOutputRedirected)`,
		`C:\`,
		os.Environ(),
	)
	readers, err := process.Start()
	require.NoError(t, err)

	outputs := make(chan string, len(readers))
	for _, reader := range readers {
		go func() {
			output, _ := io.ReadAll(reader)
			outputs <- string(output)
		}()
	}

	require.NoError(t, process.Wait())

	var output strings.Builder
	for range readers {
		output.WriteString(<-outputs)
	}
	assert.Contains(t, output.String(), "False")
	assert.NotContains(t, output.String(), string(conPTYClearScreen))
}

func TestWindowsProcessReturnsCommandFailure(t *testing.T) {
	if !isConPTYAvailable() {
		t.Skip("ConPTY is unavailable on this Windows version")
	}

	process := newCommandProcess(
		"definitely-not-a-real-command-12345",
		`C:\`,
		os.Environ(),
	)
	readers, err := process.Start()
	require.NoError(t, err)
	drained := make(chan struct{}, len(readers))
	for _, reader := range readers {
		go func() {
			_, _ = io.Copy(io.Discard, reader)
			_ = reader.Close()
			drained <- struct{}{}
		}()
	}

	assert.Error(t, process.Wait())
	for range readers {
		<-drained
	}
}
