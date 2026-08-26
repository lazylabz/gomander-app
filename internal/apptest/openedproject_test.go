package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	"gomander/internal/helpers/array"
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
}
