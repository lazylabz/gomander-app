package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commandtest "gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/domain"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
)

func TestCommandPlacements(t *testing.T) {
	t.Run("Should place no command when the group holds none", func(t *testing.T) {
		// Arrange
		commandGroup := commandgrouptest.NewCommandGroupBuilder().Build()

		// Act & Assert
		assert.Empty(t, commandGroup.CommandPlacements())
	})

	t.Run("Should place the commands densely, in the order the group holds them", func(t *testing.T) {
		// Arrange
		first := commandtest.NewCommandBuilder().Build()
		second := commandtest.NewCommandBuilder().Build()
		third := commandtest.NewCommandBuilder().Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithCommands(first, second, third).Build()

		// Act & Assert
		assert.Equal(t, []domain.CommandPlacement{
			{CommandId: first.Id, Position: 0},
			{CommandId: second.Id, Position: 1},
			{CommandId: third.Id, Position: 2},
		}, commandGroup.CommandPlacements())
	})

	t.Run("Should ignore where the commands sit among their project's commands", func(t *testing.T) {
		// Arrange
		last := commandtest.NewCommandBuilder().WithPosition(7).Build()
		first := commandtest.NewCommandBuilder().WithPosition(2).Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithCommands(last, first).Build()

		// Act & Assert
		assert.Equal(t, []domain.CommandPlacement{
			{CommandId: last.Id, Position: 0},
			{CommandId: first.Id, Position: 1},
		}, commandGroup.CommandPlacements())
	})
}
