package infrastructure

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gomander/internal/command/domain/test"

	"gomander/internal/command/domain"
	_ "gomander/migrations"
)

type testHelper struct {
	t    *testing.T
	repo *GormCommandRepository
}

func newTestHelper(t *testing.T, preloaded []*CommandModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo := arrange(preloaded)

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
		assert.Equal(t, &expectedCommand, got)
	})
	t.Run("Should return nil when it doesn't exists", func(t *testing.T) {
		// Arrange
		var preloadedCommandModels []*CommandModel
		h := newTestHelper(t, preloadedCommandModels)

		// Act
		cmd, err := h.repo.Get("nonexistent")

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, cmd)
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
		assert.Equal(t, &cmd, actual)
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
		assert.Equal(t, &editedCommand, actual)
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
		deletedCommand, err := h.repo.Get(cmd.Id)
		assert.NoError(t, err)
		assert.Nil(t, deletedCommand)
	})
	t.Run("Should delete an existing command and adjust other commands positions", func(t *testing.T) {
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

		// Verify positions were adjusted correctly
		cmd1AfterDelete, err := h.repo.Get(cmd1.Id)
		assert.NoError(t, err)

		cmd3AfterDelete, err := h.repo.Get(cmd3.Id)
		assert.NoError(t, err)

		assert.Equal(t, cmd1.Position, cmd1AfterDelete.Position, "Expected cmd1 position to remain unchanged")
		assert.Equal(t, cmd3.Position-1, cmd3AfterDelete.Position, "Expected cmd3 position to be adjusted after deletion of cmd2")
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
		cmd1After, _ := h.repo.Get(cmd1.Id)
		cmd2After, _ := h.repo.Get(cmd2.Id)
		assert.Nil(t, cmd1After)
		assert.Nil(t, cmd2After)

		// Verify other project command remains
		cmdOtherAfter, _ := h.repo.Get(cmdOther.Id)
		assert.NotNil(t, cmdOtherAfter)
	})
}

func arrange(preloadedCommandModels []*CommandModel) (repo *GormCommandRepository) {
	// Initialize the database
	ctx := context.Background()
	gormDb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db, err := gormDb.DB()
	if err != nil {
		panic(err)
	}

	// Execute migrations
	err = goose.SetDialect("sqlite3")
	if err != nil {
		panic(err)
	}

	err = goose.UpContext(ctx, db, ".")
	if err != nil {
		panic(err)
	}

	for _, m := range preloadedCommandModels {
		err = gorm.G[CommandModel](gormDb).Create(ctx, m)
		if err != nil {
			panic(err)
		}
	}

	repo = NewGormCommandRepository(gormDb, ctx)

	return
}
