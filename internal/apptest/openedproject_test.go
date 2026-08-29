package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/domainerrors"
	"gomander/internal/helpers/array"
	"gomander/internal/openedproject"
	projecttest "gomander/internal/project/domain/test"
)

func TestResolvingTheOpenedProject(t *testing.T) {
	t.Run("Should be the project that was opened", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		other := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project, other)

		assert.NoError(t, h.UseCases.OpenProject.Execute(project.Id))

		// Act
		current, err := h.UseCases.GetCurrentProject.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &project, current)
	})

	t.Run("Should be nothing before any project is opened", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		h.GivenProjects(projecttest.NewProjectBuilder().Build())

		// Act
		current, err := h.UseCases.GetCurrentProject.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, current)
	})

	t.Run("Should be nothing once the project is closed", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		givenAnOpenedProject(h)

		assert.NoError(t, h.UseCases.CloseProject.Execute())

		// Act
		current, err := h.UseCases.GetCurrentProject.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, current)
	})

	t.Run("Should refuse to open a project that no longer exists, and leave the open one alone", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		opened := givenAnOpenedProject(h)

		// Act
		err := h.UseCases.OpenProject.Execute("deleted-project")

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)

		current, err := h.UseCases.GetCurrentProject.Execute()
		assert.NoError(t, err)
		assert.Equal(t, &opened, current)
	})

	t.Run("Should be nothing once the opened project has been deleted", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		assert.NoError(t, h.UseCases.DeleteProject.Execute(project.Id))

		// Act
		current, err := h.UseCases.GetCurrentProject.Execute()

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, current)
	})

	t.Run("Should scope the commands the app shows to the opened project", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		opened := projecttest.NewProjectBuilder().Build()
		other := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(opened, other)
		h.GivenOpenedProject(opened.Id)

		ownCommand := commandtest.NewCommandBuilder().WithProjectId(opened.Id).Build()
		otherCommand := commandtest.NewCommandBuilder().WithProjectId(other.Id).Build()
		h.GivenCommands(ownCommand, otherCommand)

		// Act
		commands := commandsOf(t, h)

		// Assert
		assert.Equal(t, []string{ownCommand.Id}, array.Map(commands, commandId))
	})

	t.Run("Should scope the command groups the app shows to the opened project", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		opened := projecttest.NewProjectBuilder().Build()
		other := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(opened, other)
		h.GivenOpenedProject(opened.Id)

		ownCommand := commandtest.NewCommandBuilder().WithProjectId(opened.Id).Build()
		otherCommand := commandtest.NewCommandBuilder().WithProjectId(other.Id).Build()
		h.GivenCommands(ownCommand, otherCommand)

		ownGroup := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(opened.Id).WithCommands(ownCommand).Build()
		otherGroup := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(other.Id).WithCommands(otherCommand).Build()
		h.GivenCommandGroups(ownGroup, otherGroup)

		// Act
		groups := commandGroupsOf(t, h)

		// Assert
		assert.Equal(t, []string{ownGroup.Id}, array.Map(groups, commandGroupId))
	})

	t.Run("Should show no commands and no command groups while no project is open", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)
		h.GivenCommandGroups(commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).WithCommands(command).Build())

		// Act & Assert
		assert.Empty(t, commandsOf(t, h))
		assert.Empty(t, commandGroupsOf(t, h))
	})
}

// The operations below place something in the opened project, so there is
// nothing they can mean while the user has none open.
func TestOperatingWithoutAnOpenedProject(t *testing.T) {
	operations := []struct {
		name    string
		execute func(h *apptest.Harness) error
	}{
		{
			name: "add a command",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.AddCommand.Execute(commandtest.NewCommandBuilder().Build())
			},
		},
		{
			name: "duplicate a command",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.DuplicateCommand.Execute("any-command", "")
			},
		},
		{
			name: "reorder the commands",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.ReorderCommands.Execute([]string{"any-command"})
			},
		},
		{
			name: "create a command group",
			execute: func(h *apptest.Harness) error {
				group := commandgrouptest.NewCommandGroupBuilder().Build()
				return h.UseCases.CreateCommandGroup.Execute(&group)
			},
		},
		{
			name: "reorder the command groups",
			execute: func(h *apptest.Harness) error {
				return h.UseCases.ReorderCommandGroups.Execute([]string{"any-group"})
			},
		},
	}

	for _, operation := range operations {
		t.Run("Should refuse to "+operation.name, func(t *testing.T) {
			// Arrange
			h := apptest.New(t)

			// Act
			err := operation.execute(h)

			// Assert
			assert.ErrorIs(t, err, openedproject.ErrNoneOpen)
		})
	}
}
