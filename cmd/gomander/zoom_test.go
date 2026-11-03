package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type zoomSpy struct {
	pressed []string
	failure error
}

func (s *zoomSpy) operations() zoomOperations {
	record := func(name string) func() error {
		return func() error {
			s.pressed = append(s.pressed, name)
			return s.failure
		}
	}

	return zoomOperations{in: record("in"), out: record("out"), reset: record("reset")}
}

func TestMacOSZoomKeyBindings(t *testing.T) {
	t.Run("Should zoom in, out and back to 100% from their shortcuts", func(t *testing.T) {
		// Arrange
		spy := &zoomSpy{}
		bindings := macOSZoomKeyBindings("darwin", spy.operations(), func(string) {})

		// Act
		for _, accelerator := range []string{"CmdOrCtrl+=", "CmdOrCtrl+Shift+=", "CmdOrCtrl+plus", "CmdOrCtrl+-", "CmdOrCtrl+0"} {
			bindings[accelerator](nil)
		}

		// Assert
		assert.Equal(t, []string{"in", "in", "in", "out", "reset"}, spy.pressed)
	})

	t.Run("Should log a zoom that fails", func(t *testing.T) {
		// Arrange
		spy := &zoomSpy{failure: errors.New("storage refused")}
		logged := make([]string, 0)
		bindings := macOSZoomKeyBindings("darwin", spy.operations(), func(message string) { logged = append(logged, message) })

		// Act
		bindings["CmdOrCtrl+-"](nil)

		// Assert
		assert.Equal(t, []string{"storage refused"}, logged)
	})
	t.Run("Should leave the shortcuts to the frontend outside macOS", func(t *testing.T) {
		for _, goos := range []string{"windows", "linux"} {
			assert.Empty(t, macOSZoomKeyBindings(goos, (&zoomSpy{}).operations(), func(string) {}))
		}
	})
}
