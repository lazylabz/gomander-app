package apptest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/execution"
	"gomander/internal/openedproject"
	projecttest "gomander/internal/project/domain/test"
)

func TestRunningACommand(t *testing.T) {
	t.Run("Should start the command in the project's working directory, with the configured environment paths", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().WithWorkingDirectory("/work").Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)
		h.GivenEnvironmentPaths("/usr/local/bin", "/opt/bin")

		command := commandtest.NewCommandBuilder().
			WithProjectId(project.Id).
			WithCommand("pnpm dev").
			Build()
		h.GivenCommands(command)

		// Act
		err := h.UseCases.RunCommand.Execute(command.Id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []apptest.StartedProcess{{
			Command: command,
			Environment: execution.Environment{
				Paths:                []string{"/usr/local/bin", "/opt/bin"},
				BaseWorkingDirectory: project.WorkingDirectory,
			},
		}}, h.StartedProcesses())
		assert.Equal(t, []string{command.Id}, h.UseCases.GetRunningCommandIds.Execute())
	})

	t.Run("Should not start a command that is already running", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		assert.NoError(t, h.UseCases.RunCommand.Execute(command.Id))

		// Act
		err := h.UseCases.RunCommand.Execute(command.Id)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, h.StartedProcesses(), 1)
	})

	t.Run("Should refuse to run while no project is open, instead of dying", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		command := commandtest.NewCommandBuilder().WithProjectId("deleted-project").Build()
		h.GivenCommands(command)

		// Act
		err := h.UseCases.RunCommand.Execute(command.Id)

		// Assert
		assert.ErrorIs(t, err, openedproject.ErrNoneOpen)
		assert.Empty(t, h.StartedProcesses())
	})
}

func TestStoppingACommand(t *testing.T) {
	t.Run("Should stop the running command", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		assert.NoError(t, h.UseCases.RunCommand.Execute(command.Id))

		// Act
		err := h.UseCases.StopCommand.Execute(command.Id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{command.Id}, h.StoppedProcessIds())
		assert.Empty(t, h.UseCases.GetRunningCommandIds.Execute())
	})

	t.Run("Should do nothing when the command is not running", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
		h.GivenCommands(command)

		// Act
		err := h.UseCases.StopCommand.Execute(command.Id)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, h.StoppedProcessIds())
	})
}

func TestRunningACommandGroup(t *testing.T) {
	t.Run("Should start every command of the group, in order, in the project's working directory", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().WithWorkingDirectory("/work").Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)
		h.GivenEnvironmentPaths("/usr/local/bin")

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
		err := h.UseCases.RunCommandGroup.Execute(group.Id)

		// Assert
		assert.NoError(t, err)
		environment := execution.Environment{
			Paths:                []string{"/usr/local/bin"},
			BaseWorkingDirectory: project.WorkingDirectory,
		}
		assert.Equal(t, []apptest.StartedProcess{
			{Command: first, Environment: environment},
			{Command: second, Environment: environment},
			{Command: third, Environment: environment},
		}, h.StartedProcesses())
	})

	t.Run("Should refuse to run while no project is open, instead of dying", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		command := commandtest.NewCommandBuilder().WithProjectId("deleted-project").Build()
		h.GivenCommands(command)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId("deleted-project").
			WithCommands(command).
			Build()
		h.GivenCommandGroups(group)

		// Act
		err := h.UseCases.RunCommandGroup.Execute(group.Id)

		// Assert
		assert.ErrorIs(t, err, openedproject.ErrNoneOpen)
		assert.Empty(t, h.StartedProcesses())
	})

	t.Run("Should stop at the first command that will not start, leaving the earlier ones running", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		third := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(2).Build()
		h.GivenCommands(first, second, third)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(first, second, third).
			Build()
		h.GivenCommandGroups(group)

		failure := errors.New("no such working directory")
		h.GivenProcessThatCannotStart(second.Id, failure)

		// Act
		err := h.UseCases.RunCommandGroup.Execute(group.Id)

		// Assert
		assert.ErrorIs(t, err, failure)
		assert.Equal(t, []string{first.Id}, h.UseCases.GetRunningCommandIds.Execute())
	})
}

func TestStoppingACommandGroup(t *testing.T) {
	t.Run("Should stop every running command of the group", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(first, second)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(first, second).
			Build()
		h.GivenCommandGroups(group)

		assert.NoError(t, h.UseCases.RunCommandGroup.Execute(group.Id))

		// Act
		err := h.UseCases.StopCommandGroup.Execute(group.Id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{first.Id, second.Id}, h.StoppedProcessIds())
		assert.Empty(t, h.UseCases.GetRunningCommandIds.Execute())
	})

	t.Run("Should stop at the first command that will not stop", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(first, second)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(first, second).
			Build()
		h.GivenCommandGroups(group)

		assert.NoError(t, h.UseCases.RunCommandGroup.Execute(group.Id))

		failure := errors.New("no such process")
		h.GivenProcessThatCannotStop(first.Id, failure)

		// Act
		err := h.UseCases.StopCommandGroup.Execute(group.Id)

		// Assert
		assert.ErrorIs(t, err, failure)
		assert.Empty(t, h.StoppedProcessIds())
		assert.Equal(t, []string{first.Id, second.Id}, h.UseCases.GetRunningCommandIds.Execute())
	})
}

func TestRunningACommandEitherWay(t *testing.T) {
	t.Run("Should run it in the same environment alone as it does inside a command group", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)

		project := projecttest.NewProjectBuilder().WithWorkingDirectory("/work").Build()
		h.GivenProjects(project)
		h.GivenOpenedProject(project.Id)
		h.GivenEnvironmentPaths("/usr/local/bin", "/opt/bin")

		alone := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		grouped := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(alone, grouped)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(grouped).
			Build()
		h.GivenCommandGroups(group)

		// Act
		assert.NoError(t, h.UseCases.RunCommand.Execute(alone.Id))
		assert.NoError(t, h.UseCases.RunCommandGroup.Execute(group.Id))

		// Assert
		environment := execution.Environment{
			Paths:                []string{"/usr/local/bin", "/opt/bin"},
			BaseWorkingDirectory: project.WorkingDirectory,
		}
		started := h.StartedProcesses()
		assert.Len(t, started, 2)
		assert.Equal(t, environment, started[0].Environment)
		assert.Equal(t, environment, started[1].Environment)
	})
}
