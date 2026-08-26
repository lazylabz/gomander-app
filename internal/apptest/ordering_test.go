package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/helpers/array"
)

func TestOrderingCommands(t *testing.T) {
	t.Run("Should add a command at the end of the list", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		h.GivenCommands(first)

		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()

		// Act
		err := h.UseCases.AddCommand.Execute(second)

		// Assert
		assert.NoError(t, err)
		commands := commandsOf(t, h)
		assert.Equal(t, []string{first.Id, second.Id}, array.Map(commands, commandId))
		assert.Equal(t, []int{0, 1}, array.Map(commands, commandPosition))
	})

	t.Run("Should add a duplicated command at the end of the list", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		original := commandtest.NewCommandBuilder().
			WithProjectId(project.Id).
			WithName("Dev server").
			WithPosition(0).
			Build()
		other := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(original, other)

		// Act
		err := h.UseCases.DuplicateCommand.Execute(original.Id, "")

		// Assert
		assert.NoError(t, err)
		commands := commandsOf(t, h)
		assert.Equal(t, []int{0, 1, 2}, array.Map(commands, commandPosition))
		assert.Equal(t, "Dev server (copy)", commands[2].Name)
		assert.Equal(t, original.Command, commands[2].Command)
		assert.NotEqual(t, original.Id, commands[2].Id)
	})

	t.Run("Should renumber the list when commands are reordered", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		third := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(2).Build()
		h.GivenCommands(first, second, third)

		// Act
		err := h.UseCases.ReorderCommands.Execute([]string{third.Id, first.Id, second.Id})

		// Assert
		assert.NoError(t, err)
		commands := commandsOf(t, h)
		assert.Equal(t, []string{third.Id, first.Id, second.Id}, array.Map(commands, commandId))
		assert.Equal(t, []int{0, 1, 2}, array.Map(commands, commandPosition))
	})

	t.Run("Should close the gap when a command is removed", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		third := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(2).Build()
		h.GivenCommands(first, second, third)

		// Act
		err := h.UseCases.RemoveCommand.Execute(second.Id)

		// Assert
		assert.NoError(t, err)
		commands := commandsOf(t, h)
		assert.Equal(t, []string{first.Id, third.Id}, array.Map(commands, commandId))
		assert.Equal(t, []int{0, 1}, array.Map(commands, commandPosition))
	})
}

func TestOrderingCommandGroups(t *testing.T) {
	t.Run("Should add a command group at the end of the list", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		first := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(command).
			WithPosition(0).
			Build()
		h.GivenCommandGroups(first)

		second := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(command).
			Build()

		// Act
		err := h.UseCases.CreateCommandGroup.Execute(&second)

		// Assert
		assert.NoError(t, err)
		groups := commandGroupsOf(t, h)
		assert.Equal(t, []string{first.Id, second.Id}, array.Map(groups, commandGroupId))
		assert.Equal(t, []int{0, 1}, array.Map(groups, commandGroupPosition))
	})

	t.Run("Should renumber the list when command groups are reordered", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		first := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).WithCommands(command).WithPosition(0).Build()
		second := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).WithCommands(command).WithPosition(1).Build()
		h.GivenCommandGroups(first, second)

		// Act
		err := h.UseCases.ReorderCommandGroups.Execute([]string{second.Id, first.Id})

		// Assert
		assert.NoError(t, err)
		groups := commandGroupsOf(t, h)
		assert.Equal(t, []string{second.Id, first.Id}, array.Map(groups, commandGroupId))
		assert.Equal(t, []int{0, 1}, array.Map(groups, commandGroupPosition))
	})

	t.Run("Should close the gap when a command group is deleted", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		first := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).WithCommands(command).WithPosition(0).Build()
		second := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).WithCommands(command).WithPosition(1).Build()
		third := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).WithCommands(command).WithPosition(2).Build()
		h.GivenCommandGroups(first, second, third)

		// Act
		err := h.UseCases.DeleteCommandGroup.Execute(second.Id)

		// Assert
		assert.NoError(t, err)
		groups := commandGroupsOf(t, h)
		assert.Equal(t, []string{first.Id, third.Id}, array.Map(groups, commandGroupId))
		assert.Equal(t, []int{0, 1}, array.Map(groups, commandGroupPosition))
	})
}
