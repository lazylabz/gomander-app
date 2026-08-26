package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
)

func TestACommandGroupLosingItsLastCommand(t *testing.T) {
	t.Run("Should cease to exist when its only command is deleted, and tell the frontend", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		onlyCommand := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		h.GivenCommands(onlyCommand)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(onlyCommand).
			Build()
		h.GivenCommandGroups(group)

		// Act
		err := h.UseCases.RemoveCommand.Execute(onlyCommand.Id)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, commandGroupsOf(t, h))
		assert.Equal(t, []apptest.EmittedEvent{
			{Name: event.CommandGroupDeleted, Payload: group.Id},
		}, h.EmittedEvents())
	})

	t.Run("Should keep its remaining commands, in order, when one of several is deleted", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		third := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(2).Build()
		h.GivenCommands(first, second, third)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(first, second, third).
			Build()
		h.GivenCommandGroups(group)

		// Act
		err := h.UseCases.RemoveCommand.Execute(second.Id)

		// Assert
		assert.NoError(t, err)
		groups := commandGroupsOf(t, h)
		assert.Len(t, groups, 1)
		assert.Equal(t, []string{first.Id, third.Id}, array.Map(groups[0].Commands, commandId))
		assert.Empty(t, h.EmittedEvents())
	})

	t.Run("Should refuse to let the app remove its last command from it", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		onlyCommand := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(onlyCommand)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(onlyCommand).
			Build()
		h.GivenCommandGroups(group)

		// Act
		err := h.UseCases.RemoveCommandFromCommandGroup.Execute(onlyCommand.Id, group.Id)

		// Assert
		assert.Error(t, err)
		groups := commandGroupsOf(t, h)
		assert.Len(t, groups, 1)
		assert.Equal(t, []string{onlyCommand.Id}, array.Map(groups[0].Commands, commandId))
	})
}
