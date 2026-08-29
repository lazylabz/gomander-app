package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	"gomander/internal/domainerrors"
)

// A Command id reaches the backend from UI state that can be stale: the
// Command it names may be gone by the time the operation runs.
func TestOperatingOnACommandThatIsNotThere(t *testing.T) {
	const missingCommandId = "deleted-command"

	operations := []struct {
		name    string
		execute func(h *apptest.Harness) error
	}{
		{
			name: "run it",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.RunCommand.Execute(missingCommandId)
			},
		},
		{
			name: "stop it",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.StopCommand.Execute(missingCommandId)
			},
		},
		{
			name: "remove it",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.RemoveCommand.Execute(missingCommandId)
			},
		},
		{
			name: "duplicate it",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.DuplicateCommand.Execute(missingCommandId, "")
			},
		},
	}

	for _, operation := range operations {
		t.Run("Should refuse to "+operation.name, func(t *testing.T) {
			// Arrange
			h := apptest.New(t)
			givenAnOpenedProject(h)

			// Act
			err := operation.execute(h)

			// Assert: the Command is what is missing, not the Project the
			// operations resolve on the way to it.
			assert.ErrorIs(t, err, domainerrors.ErrNotFound)
			assert.ErrorContains(t, err, missingCommandId)
			assert.Empty(t, h.StartedProcesses())
			assert.Empty(t, h.StoppedProcessIds())
		})
	}
}
