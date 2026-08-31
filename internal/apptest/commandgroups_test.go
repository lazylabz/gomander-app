package apptest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/apptest"
	commandtest "gomander/internal/command/domain/test"
	commandgrouptest "gomander/internal/commandgroup/domain/test"
	"gomander/internal/domainerrors"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
	projecttest "gomander/internal/project/domain/test"
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
		assert.Equal(t, []string{first.Id, third.Id}, groups[0].CommandIds)
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
		assert.Equal(t, []string{onlyCommand.Id}, groups[0].CommandIds)
	})
}

const missingCommandGroupId = "deleted-command-group"

func TestRemovingACommandFromACommandGroup(t *testing.T) {
	t.Run("Should leave the group with its other commands, and the command itself untouched", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		removed := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(0).Build()
		kept := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithPosition(1).Build()
		h.GivenCommands(removed, kept)

		group := commandgrouptest.NewCommandGroupBuilder().
			WithProjectId(project.Id).
			WithCommands(removed, kept).
			Build()
		h.GivenCommandGroups(group)

		// Act
		err := h.UseCases.RemoveCommandFromCommandGroup.Execute(removed.Id, group.Id)

		// Assert
		assert.NoError(t, err)

		groups := commandGroupsOf(t, h)
		assert.Len(t, groups, 1)
		assert.Equal(t, []string{kept.Id}, groups[0].CommandIds)

		// Leaving a group is not leaving the project.
		assert.Equal(t, []string{removed.Id, kept.Id}, array.Map(commandsOf(t, h), commandId))
	})
}

// A Command Group id reaches the backend from UI state that can be stale: the
// Group it names may be gone by the time the operation runs.
func TestOperatingOnACommandGroupThatIsNotThere(t *testing.T) {
	// The two operations that name a Command as well as a Group are handed one
	// that is really there, so the absence under test is only ever the Group's.
	operations := []struct {
		name    string
		execute func(h *apptest.Harness, commandId string) error
	}{
		{
			name: "run it",
			execute: func(h *apptest.Harness, _ string) error {
				return h.UseCases.RunCommandGroup.Execute(missingCommandGroupId)
			},
		},
		{
			name: "stop it",
			execute: func(h *apptest.Harness, _ string) error {
				return h.UseCases.StopCommandGroup.Execute(missingCommandGroupId)
			},
		},
		{
			name: "delete it",
			execute: func(h *apptest.Harness, _ string) error {
				return h.UseCases.DeleteCommandGroup.Execute(missingCommandGroupId)
			},
		},
		{
			name: "remove a command from it",
			execute: func(h *apptest.Harness, commandId string) error {
				return h.UseCases.RemoveCommandFromCommandGroup.Execute(commandId, missingCommandGroupId)
			},
		},
		{
			name: "duplicate a command into it",
			execute: func(h *apptest.Harness, commandId string) error {
				return h.UseCases.DuplicateCommand.Execute(commandId, missingCommandGroupId)
			},
		},
	}

	for _, operation := range operations {
		t.Run("Should refuse to "+operation.name, func(t *testing.T) {
			// Arrange
			h := apptest.New(t)
			project := givenAnOpenedProject(h)

			command := commandtest.NewCommandBuilder().WithProjectId(project.Id).Build()
			h.GivenCommands(command)

			// Act
			err := operation.execute(h, command.Id)

			// Assert: the Command Group is what is missing, not the Project or
			// the Command the operations resolve on the way to it.
			assert.ErrorIs(t, err, domainerrors.ErrNotFound)
			assert.ErrorContains(t, err, missingCommandGroupId)
			assert.Empty(t, h.StartedProcesses())
			assert.Empty(t, h.StoppedProcessIds())
			assert.Empty(t, h.EmittedEvents())
		})
	}
}

// Duplicating creates the copy before the event that files it into a Group is
// published, so a Group that has gone leaves the copy behind, loose in the
// project. Pinned rather than fixed here: making the two atomic is the cascade
// work in #241.
func TestDuplicatingACommandIntoACommandGroupThatIsNotThere(t *testing.T) {
	t.Run("Should leave the copy in the project it was made in", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		project := givenAnOpenedProject(h)

		original := commandtest.NewCommandBuilder().WithProjectId(project.Id).WithName("build").Build()
		h.GivenCommands(original)

		// Act
		err := h.UseCases.DuplicateCommand.Execute(original.Id, missingCommandGroupId)

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
		assert.Equal(t, []string{"build", "build (copy)"}, array.Map(commandsOf(t, h), commandName))
		assert.Empty(t, commandGroupsOf(t, h))
	})
}

func TestDeletingAProject(t *testing.T) {
	t.Run("Should take its command groups with it, and tell the frontend about each", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		deleted := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(deleted)
		h.GivenOpenedProject(deleted.Id)

		first := commandtest.NewCommandBuilder().WithProjectId(deleted.Id).WithPosition(0).Build()
		second := commandtest.NewCommandBuilder().WithProjectId(deleted.Id).WithPosition(1).Build()
		h.GivenCommands(first, second)

		firstGroup := commandgrouptest.NewCommandGroupBuilder().WithProjectId(deleted.Id).WithPosition(0).WithCommands(first).Build()
		secondGroup := commandgrouptest.NewCommandGroupBuilder().WithProjectId(deleted.Id).WithPosition(1).WithCommands(second).Build()
		h.GivenCommandGroups(firstGroup, secondGroup)

		// Act
		err := h.UseCases.DeleteProject.Execute(deleted.Id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []apptest.EmittedEvent{
			{Name: event.CommandGroupDeleted, Payload: firstGroup.Id},
			{Name: event.CommandGroupDeleted, Payload: secondGroup.Id},
		}, h.EmittedEvents())
	})

	t.Run("Should leave another project's commands and command groups alone", func(t *testing.T) {
		// Arrange
		h := apptest.New(t)
		kept := projecttest.NewProjectBuilder().Build()
		deleted := projecttest.NewProjectBuilder().Build()
		h.GivenProjects(kept, deleted)

		keptCommand := commandtest.NewCommandBuilder().WithProjectId(kept.Id).WithPosition(0).Build()
		deletedCommand := commandtest.NewCommandBuilder().WithProjectId(deleted.Id).WithPosition(0).Build()
		h.GivenCommands(keptCommand, deletedCommand)

		keptGroup := commandgrouptest.NewCommandGroupBuilder().WithProjectId(kept.Id).WithPosition(0).WithCommands(keptCommand).Build()
		deletedGroup := commandgrouptest.NewCommandGroupBuilder().WithProjectId(deleted.Id).WithPosition(0).WithCommands(deletedCommand).Build()
		h.GivenCommandGroups(keptGroup, deletedGroup)

		// Act
		err := h.UseCases.DeleteProject.Execute(deleted.Id)

		// Assert
		assert.NoError(t, err)

		h.GivenOpenedProject(kept.Id)
		assert.Equal(t, []string{keptCommand.Id}, array.Map(commandsOf(t, h), commandId))
		assert.Equal(t, []string{keptGroup.Id}, array.Map(commandGroupsOf(t, h), commandGroupId))
	})
}
