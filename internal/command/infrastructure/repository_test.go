package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gomander/internal/command/domain/test"

	"gomander/internal/command/domain"
	"gomander/internal/domainerrors"
	"gomander/internal/testdb"
)

type testHelper struct {
	t    *testing.T
	repo *GormCommandRepository
}

func newTestHelper(t *testing.T, preloaded []*CommandModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo := arrange(t, preloaded)

	helper := &testHelper{
		t:    t,
		repo: repo,
	}

	return helper
}

func TestGormCommandRepository_Get(t *testing.T) {
	t.Parallel()
	t.Run("Should return command when it exists", func(t *testing.T) {
		// Arrange
		cmd := test.NewCommandBuilder().Build()
		model := ToCommandModel(&cmd)

		preloadedCommandModels := []*CommandModel{&model}

		expectedCommand := cmd

		h := newTestHelper(t, preloadedCommandModels)

		// Act
		got, err := h.repo.Get(cmd.Id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCommand, got)
	})
	t.Run("Should report a command that does not exist as not found", func(t *testing.T) {
		// Arrange
		var preloadedCommandModels []*CommandModel
		h := newTestHelper(t, preloadedCommandModels)

		// Act
		_, err := h.repo.Get("nonexistent")

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	})
}

func TestGormCommandRepository_GetAll(t *testing.T) {
	t.Run("Should return all commands for a project, sorted", func(t *testing.T) {
		// Arrange
		projectId := "proj1"
		cmd1 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(1).
			Build()
		cmd2 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			Build()
		cmd1model := ToCommandModel(&cmd1)
		cmd2model := ToCommandModel(&cmd2)

		preloadedCommandModels := []*CommandModel{
			&cmd2model,
			&cmd1model,
		}

		expectedCommands := []*domain.Command{
			&cmd2,
			&cmd1,
		}

		h := newTestHelper(t, preloadedCommandModels)

		// Act
		cmds, err := h.repo.GetAll(projectId)

		// Assert
		assert.NoError(t, err)
		for i, cmd := range cmds {
			assert.Equal(t, expectedCommands[i], &cmd)
		}
	})
}

func TestGormCommandRepository_Save(t *testing.T) {
	t.Parallel()
	t.Run("Should save a new command", func(t *testing.T) {
		// Arrange
		var preloadedCommandModels []*CommandModel
		h := newTestHelper(t, preloadedCommandModels)

		cmd := test.NewCommandBuilder().
			WithId("cmd3").
			WithProjectId("proj1").
			WithName("New Command").
			WithCommand("echo 'New Command'").
			WithWorkingDirectory("/tmp").
			WithPosition(2).
			Build()

		// Act
		err := h.repo.Create(&cmd)

		// Assert
		assert.NoError(t, err)

		// Verify the command was saved
		actual, err := h.repo.Get("cmd3")
		assert.NoError(t, err)
		assert.Equal(t, cmd, actual)
	})
}

func TestGormCommandRepository_Edit(t *testing.T) {
	t.Parallel()
	t.Run("Should edit an existing command", func(t *testing.T) {
		// Arrange
		existingCommandBuilder := test.NewCommandBuilder().
			WithProjectId("proj1").
			WithName("Old Command").
			WithCommand("echo 'Old Command'").
			WithWorkingDirectory("/tmp").
			WithPosition(0)
		existingCommand := existingCommandBuilder.Build()
		existingCommandModel := ToCommandModel(&existingCommand)

		preloadedCommandModels := []*CommandModel{
			&existingCommandModel,
		}

		h := newTestHelper(t, preloadedCommandModels)

		editedCommand := existingCommandBuilder.
			WithName("Edited Command").
			WithCommand("echo 'Edited Command'").
			Build()

		// Act
		err := h.repo.Update(&editedCommand)

		// Assert
		assert.NoError(t, err)

		// Verify the command was updated
		actual, err := h.repo.Get(existingCommandBuilder.Build().Id)
		assert.NoError(t, err)
		assert.Equal(t, editedCommand, actual)
	})
}

func TestGormCommandRepository_Delete(t *testing.T) {
	t.Parallel()
	t.Run("Should delete an existing command", func(t *testing.T) {
		// Arrange
		cmd := test.NewCommandBuilder().Build()
		cmdModel := ToCommandModel(&cmd)

		preloadedCommandModels := []*CommandModel{
			&cmdModel,
		}

		h := newTestHelper(t, preloadedCommandModels)

		// Act
		err := h.repo.Delete(cmd.Id)

		// Assert
		assert.NoError(t, err)

		// Verify the command was deleted
		_, err = h.repo.Get(cmd.Id)
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	})
	t.Run("Should leave the positions of the remaining commands untouched", func(t *testing.T) {
		// Arrange
		projectId := "proj1"
		cmd1 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(0).
			Build()
		cmd2 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(1).
			Build()
		cmd3 := test.NewCommandBuilder().
			WithProjectId(projectId).
			WithPosition(2).
			Build()

		cmd1Model := ToCommandModel(&cmd1)
		cmd2Model := ToCommandModel(&cmd2)
		cmd3Model := ToCommandModel(&cmd3)

		preloadedCommandModels := []*CommandModel{
			&cmd1Model,
			&cmd2Model,
			&cmd3Model,
		}

		h := newTestHelper(t, preloadedCommandModels)

		// Act
		err := h.repo.Delete(cmd2.Id)

		// Assert
		assert.NoError(t, err)

		// Closing the gap is the ordering module's job, not the repository's
		cmd1AfterDelete, err := h.repo.Get(cmd1.Id)
		assert.NoError(t, err)

		cmd3AfterDelete, err := h.repo.Get(cmd3.Id)
		assert.NoError(t, err)

		assert.Equal(t, cmd1.Position, cmd1AfterDelete.Position)
		assert.Equal(t, cmd3.Position, cmd3AfterDelete.Position)
	})
}

func TestGormCommandRepository_DeleteAll(t *testing.T) {
	t.Run("Should delete all commands for a project and not affect others", func(t *testing.T) {
		// Arrange
		projectId := "proj1"
		otherProjectId := "proj2"

		cmd1 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(0).Build()
		cmd2 := test.NewCommandBuilder().WithProjectId(projectId).WithPosition(1).Build()
		cmdOther := test.NewCommandBuilder().WithProjectId(otherProjectId).WithPosition(0).Build()

		cmd1Model := ToCommandModel(&cmd1)
		cmd2Model := ToCommandModel(&cmd2)
		cmdOtherModel := ToCommandModel(&cmdOther)

		preloadedCommandModels := []*CommandModel{
			&cmd1Model,
			&cmd2Model,
			&cmdOtherModel,
		}

		h := newTestHelper(t, preloadedCommandModels)

		// Act
		err := h.repo.DeleteAll(projectId)

		// Assert
		assert.NoError(t, err)

		// Verify project commands were deleted
		_, cmd1Err := h.repo.Get(cmd1.Id)
		_, cmd2Err := h.repo.Get(cmd2.Id)
		assert.ErrorIs(t, cmd1Err, domainerrors.ErrNotFound)
		assert.ErrorIs(t, cmd2Err, domainerrors.ErrNotFound)

		// Verify other project command remains
		cmdOtherAfter, err := h.repo.Get(cmdOther.Id)
		assert.NoError(t, err)
		assert.Equal(t, cmdOther, cmdOtherAfter)
	})
}

func arrange(t *testing.T, preloadedCommandModels []*CommandModel) (repo *GormCommandRepository) {
	t.Helper()

	ctx := context.Background()
	gormDb := testdb.New(t)

	for _, m := range preloadedCommandModels {
		err := gorm.G[CommandModel](gormDb).Create(ctx, m)
		if err != nil {
			t.Fatalf("failed to preload command: %v", err)
		}
	}

	repo = NewGormCommandRepository(gormDb, ctx)

	return
}
