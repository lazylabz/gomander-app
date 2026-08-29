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

func TestClassifyWindowsVersion(t *testing.T) {
	tests := []struct {
		name        string
		major       uint32
		build       uint32
		productType byte
		expected    HostEnvironment
	}{
		{name: "Windows 10 before ConPTY", major: 10, build: 17134, productType: 1},
		{name: "Windows 10 1809", major: 10, build: 17763, productType: 1, expected: HostEnvironmentWindows10},
		{name: "Windows 10 22H2", major: 10, build: 19045, productType: 1, expected: HostEnvironmentWindows10},
		{name: "Windows Server", major: 10, build: 17763, productType: 3},
		{name: "Windows 11", major: 10, build: 22000, productType: 1},
		{name: "future major version", major: 11, build: 100, productType: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, classifyWindowsVersion(test.major, test.build, test.productType))
		})
	}
}

func TestShouldUseConPTYRequiresAllowlistAndAPI(t *testing.T) {
	enabled := Config{ConPTYEnvironments: []HostEnvironment{HostEnvironmentWindows10}}

	assert.True(t, shouldUseConPTY(enabled, HostEnvironmentWindows10, true))
	assert.False(t, shouldUseConPTY(Config{}, HostEnvironmentWindows10, true))
	assert.False(t, shouldUseConPTY(enabled, "", true))
	assert.False(t, shouldUseConPTY(enabled, HostEnvironmentWindows10, false))
}

func TestWindowsProcessUsesTerminalOnlyWhenConfigured(t *testing.T) {
	if currentWindowsHostEnvironment() != HostEnvironmentWindows10 || !isConPTYAvailable() {
		t.Skip("the current host is not a supported Windows 10 ConPTY environment")
	}

	terminalProcess := newCommandProcess(
		`powershell -NoProfile -Command Write-Output ([Console]::IsOutputRedirected)`,
		`C:\`,
		os.Environ(),
		Config{ConPTYEnvironments: []HostEnvironment{HostEnvironmentWindows10}},
	)
	terminalOutput, err := runAndReadProcess(terminalProcess)
	require.NoError(t, err)
	assert.Contains(t, terminalOutput, "False")
	assert.NotContains(t, terminalOutput, string(conPTYClearScreen))

	pipeProcess := newCommandProcess(
		`powershell -NoProfile -Command Write-Output ([Console]::IsOutputRedirected)`,
		`C:\`,
		os.Environ(),
		Config{},
	)
	pipeOutput, err := runAndReadProcess(pipeProcess)
	require.NoError(t, err)
	assert.Contains(t, pipeOutput, "True")
}

func TestWindowsConPTYProcessReturnsCommandFailure(t *testing.T) {
	if currentWindowsHostEnvironment() != HostEnvironmentWindows10 || !isConPTYAvailable() {
		t.Skip("the current host is not a supported Windows 10 ConPTY environment")
	}

	process := newCommandProcess(
		"definitely-not-a-real-command-12345",
		`C:\`,
		os.Environ(),
		Config{ConPTYEnvironments: []HostEnvironment{HostEnvironmentWindows10}},
	)
	_, err := runAndReadProcess(process)
	assert.Error(t, err)
}

func TestCommandInterpreterDefaultsToCmd(t *testing.T) {
	t.Setenv("COMSPEC", "")
	assert.Equal(t, "cmd.exe", commandInterpreter())
}

func runAndReadProcess(process commandProcess) (string, error) {
	readers, err := process.Start()
	if err != nil {
		return "", err
	}

	outputs := make(chan string, len(readers))
	for _, reader := range readers {
		go func(reader io.ReadCloser) {
			output, _ := io.ReadAll(reader)
			_ = reader.Close()
			outputs <- string(output)
		}(reader)
	}

	waitErr := process.Wait()
	var output strings.Builder
	for range readers {
		output.WriteString(<-outputs)
	}
	return output.String(), waitErr
}
