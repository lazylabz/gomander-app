package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/command/domain"
	"gomander/internal/command/domain/test"
)

func TestRunningCommands(t *testing.T) {
	t.Run("Should hold a command the runner has a process for", func(t *testing.T) {
		// Arrange
		command := test.NewCommandBuilder().Build()
		running := domain.NewRunningCommands([]string{command.Id})

		// Act & Assert
		assert.True(t, running.IsRunning(command.Id))
		assert.Equal(t, domain.Running, running.StatusOf(command.Id))
	})

	t.Run("Should not hold a command the runner has no process for", func(t *testing.T) {
		// Arrange
		stopped := test.NewCommandBuilder().Build()
		running := domain.NewRunningCommands([]string{test.NewCommandBuilder().Build().Id})

		// Act & Assert
		assert.False(t, running.IsRunning(stopped.Id))
		assert.Equal(t, domain.Stopped, running.StatusOf(stopped.Id))
	})

	t.Run("Should report every command as stopped when nothing is running", func(t *testing.T) {
		// Arrange
		command := test.NewCommandBuilder().Build()
		running := domain.NewRunningCommands(nil)

		// Act & Assert
		assert.False(t, running.IsRunning(command.Id))
		assert.Equal(t, domain.Stopped, running.StatusOf(command.Id))
	})

	t.Run("Should count only the running ones of the commands given", func(t *testing.T) {
		// Arrange
		first := test.NewCommandBuilder().Build()
		second := test.NewCommandBuilder().Build()
		third := test.NewCommandBuilder().Build()
		running := domain.NewRunningCommands([]string{first.Id, third.Id})

		// Act & Assert
		assert.Equal(t, 2, running.CountIn([]string{first.Id, second.Id, third.Id}))
	})

	t.Run("Should not count a running command that is not one of the ones given", func(t *testing.T) {
		// Arrange
		outsider := test.NewCommandBuilder().Build()
		command := test.NewCommandBuilder().Build()
		running := domain.NewRunningCommands([]string{outsider.Id})

		// Act & Assert
		assert.Equal(t, 0, running.CountIn([]string{command.Id}))
	})

	t.Run("Should count nothing in a command group with no commands", func(t *testing.T) {
		// Arrange
		running := domain.NewRunningCommands([]string{test.NewCommandBuilder().Build().Id})

		// Act & Assert
		assert.Equal(t, 0, running.CountIn(nil))
	})
}
