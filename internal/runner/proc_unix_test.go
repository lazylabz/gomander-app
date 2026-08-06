//go:build !windows

package runner

import (
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellExecutableDefaultsToBinSh(t *testing.T) {
	t.Setenv("SHELL", "")

	assert.Equal(t, "/bin/sh", shellExecutable())
}

func TestIgnoreNoSuchProcess(t *testing.T) {
	assert.NoError(t, ignoreNoSuchProcess(syscall.ESRCH))

	expected := errors.New("unexpected signal error")
	assert.ErrorIs(t, ignoreNoSuchProcess(expected), expected)
}
