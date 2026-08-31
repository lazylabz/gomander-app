package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commandtest "gomander/internal/command/domain/test"
	"gomander/internal/commandgroup/domain"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/helpers/array"
)

func TestRemoveCommandFrom(t *testing.T) {
	t.Run("Should leave a command group that still names other commands standing, without the removed one", func(t *testing.T) {
		// Arrange
		removed := commandtest.NewCommandBuilder().Build()
		kept := commandtest.NewCommandBuilder().Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithCommands(removed, kept).Build()

		// Act
		cascade := domain.RemoveCommandFrom([]domain.CommandGroup{commandGroup}, removed.Id)

		// Assert
		assert.Empty(t, cascade.Deleted)
		assert.Len(t, cascade.Survived, 1)
		assert.Equal(t, []string{kept.Id}, cascade.Survived[0].CommandIds)
	})

	t.Run("Should make a command group that has lost its last command cease to exist", func(t *testing.T) {
		// Arrange
		onlyCommand := commandtest.NewCommandBuilder().Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithCommands(onlyCommand).Build()

		// Act
		cascade := domain.RemoveCommandFrom([]domain.CommandGroup{commandGroup}, onlyCommand.Id)

		// Assert
		assert.Empty(t, cascade.Survived)
		assert.Equal(t, []string{commandGroup.Id}, idsOf(cascade.Deleted))
	})

	t.Run("Should decide each command group on its own", func(t *testing.T) {
		// Arrange
		removed := commandtest.NewCommandBuilder().Build()
		kept := commandtest.NewCommandBuilder().Build()

		emptied := commandgrouptest.NewCommandGroupBuilder().WithCommands(removed).Build()
		surviving := commandgrouptest.NewCommandGroupBuilder().WithCommands(kept, removed).Build()

		// Act
		cascade := domain.RemoveCommandFrom([]domain.CommandGroup{emptied, surviving}, removed.Id)

		// Assert
		assert.Equal(t, []string{surviving.Id}, idsOf(cascade.Survived))
		assert.Equal(t, []string{emptied.Id}, idsOf(cascade.Deleted))
	})

	t.Run("Should leave the command groups it was given untouched", func(t *testing.T) {
		// Arrange
		removed := commandtest.NewCommandBuilder().Build()
		kept := commandtest.NewCommandBuilder().Build()
		commandGroup := commandgrouptest.NewCommandGroupBuilder().WithCommands(removed, kept).Build()

		// Act
		domain.RemoveCommandFrom([]domain.CommandGroup{commandGroup}, removed.Id)

		// Assert
		assert.Equal(t, []string{removed.Id, kept.Id}, commandGroup.CommandIds)
	})

	t.Run("Should decide nothing when there is no command group naming the command", func(t *testing.T) {
		// Act
		cascade := domain.RemoveCommandFrom(nil, "any-command")

		// Assert
		assert.Empty(t, cascade.Survived)
		assert.Empty(t, cascade.Deleted)
	})
}

func idsOf(commandGroups []domain.CommandGroup) []string {
	return array.Map(commandGroups, func(commandGroup domain.CommandGroup) string {
		return commandGroup.Id
	})
}
