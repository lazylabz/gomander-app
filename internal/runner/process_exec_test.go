package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecProcessOutputRemainsReadableAfterWait(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestExecProcessHelper$")
	cmd.Env = append(os.Environ(), "GO_WANT_EXEC_PROCESS_HELPER=1")
	process := &execProcess{cmd: cmd}

	readers, err := process.Start()
	require.NoError(t, err)
	require.Len(t, readers, 2)
	t.Cleanup(func() { closeReaders(readers[0], readers[1]) })

	require.NoError(t, process.Wait())
	stdout, err := io.ReadAll(readers[0])
	require.NoError(t, err)
	assert.Equal(t, "tail output\n", string(stdout))
}

func TestExecProcessHelper(t *testing.T) {
	if os.Getenv("GO_WANT_EXEC_PROCESS_HELPER") != "1" {
		return
	}

	_, _ = fmt.Fprintln(os.Stdout, "tail output")
	os.Exit(0)
}
