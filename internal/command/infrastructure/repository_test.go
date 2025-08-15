package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"gomander/internal/testutils"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gomander/internal/command/domain"
	_ "gomander/migrations"
)

func commandDataToDomain(data testutils.CommandData) *domain.Command {
	return &domain.Command{
		Id:               data.Id,
		ProjectId:        data.ProjectId,
		Name:             data.Name,
		Command:          data.Command,
		WorkingDirectory: data.WorkingDirectory,
		Position:         data.Position,
	}
}

func commandDataToModel(data testutils.CommandData) *CommandModel {
	return &CommandModel{
		Id:               data.Id,
		ProjectId:        data.ProjectId,
		Name:             data.Name,
		Command:          data.Command,
		WorkingDirectory: data.WorkingDirectory,
		Position:         data.Position,
	}
}

type testHelper struct {
	t      *testing.T
	repo   *GormCommandRepository
	dbPath string
}

func newTestHelper(t *testing.T, preloaded []*CommandModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo, dbPath := arrange(preloaded)

	helper := &testHelper{
		t:      t,
		repo:   repo,
		dbPath: dbPath,
	}

	// Automatic cleanup when test finishes
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(helper.dbPath), "Failed to cleanup test database")
	})

	return helper
}

func TestGormCommandRepository_Get(t *testing.T) {
	t.Parallel()
	t.Run("Should return command when it exists", func(t *testing.T) {
		cmdData := testutils.NewCommand().Data()

		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmdData),
		}

		expectedCommand := commandDataToDomain(cmdData)

		h := newTestHelper(t, preloadedCommandModels)

		cmd, err := h.repo.Get(cmdData.Id)
		assert.NoError(t, err)

		assert.Equal(t, expectedCommand, cmd)
	})
	t.Run("Should return nil when it doesn't exists", func(t *testing.T) {
		var preloadedCommandModels []*CommandModel

		h := newTestHelper(t, preloadedCommandModels)

		cmd, err := h.repo.Get("nonexistent")
		assert.NoError(t, err)

		assert.Nil(t, cmd)
	})
}

func TestGormCommandRepository_GetAll(t *testing.T) {
	t.Run("Should return all commands for a project, sorted", func(t *testing.T) {
		projectId := "proj1"
		cmd1 := testutils.NewCommand().
			WithProjectId(projectId).
			WithPosition(1).
			Data()
		cmd2 := testutils.NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()

		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmd2),
			commandDataToModel(cmd1),
		}

		expectedCommands := []*domain.Command{
			commandDataToDomain(cmd2),
			commandDataToDomain(cmd1),
		}

		h := newTestHelper(t, preloadedCommandModels)

		cmds, err := h.repo.GetAll(projectId)
		assert.NoError(t, err)

		for i, cmd := range cmds {
			assert.Equal(t, expectedCommands[i], &cmd)
		}
	})
}

func TestGormCommandRepository_Save(t *testing.T) {
	t.Parallel()
	t.Run("Should save a new command", func(t *testing.T) {
		var preloadedCommandModels []*CommandModel

		h := newTestHelper(t, preloadedCommandModels)

		cmdData := testutils.NewCommand().
			WithId("cmd3").
			WithProjectId("proj1").
			WithName("New Command").
			WithCommand("echo 'New Command'").
			WithWorkingDirectory("/tmp").
			WithPosition(2).
			Data()
		createdCommand := commandDataToDomain(cmdData)

		err := h.repo.Create(createdCommand)
		assert.NoError(t, err)

		actual, err := h.repo.Get("cmd3")
		assert.NoError(t, err)

		assert.Equal(t, createdCommand, actual)
	})
}

func TestGormCommandRepository_Edit(t *testing.T) {
	t.Parallel()
	t.Run("Should edit an existing command", func(t *testing.T) {
		existingCommand := testutils.NewCommand().
			WithProjectId("proj1").
			WithName("Old Command").
			WithCommand("echo 'Old Command'").
			WithWorkingDirectory("/tmp").
			WithPosition(0)
		preloadedCommandModels := []*CommandModel{
			commandDataToModel(existingCommand.Data()),
		}

		h := newTestHelper(t, preloadedCommandModels)

		editedCommand := existingCommand.
			WithName("Edited Command").
			WithCommand("echo 'Edited Command'").
			Data()
		domainEditedCommand := commandDataToDomain(editedCommand)

		err := h.repo.Update(domainEditedCommand)
		assert.NoError(t, err)

		actual, err := h.repo.Get(existingCommand.Data().Id)
		assert.NoError(t, err)

		assert.Equal(t, domainEditedCommand, actual)
	})
}

func TestGormCommandRepository_Delete(t *testing.T) {
	t.Parallel()
	t.Run("Should delete an existing command", func(t *testing.T) {
		cmd := testutils.NewCommand().Data()
		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmd),
		}

		h := newTestHelper(t, preloadedCommandModels)

		err := h.repo.Delete(cmd.Id)
		assert.NoError(t, err)

		deletedCommand, err := h.repo.Get(cmd.Id)
		assert.NoError(t, err)

		assert.Nil(t, deletedCommand)
	})
	t.Run("Should delete an existing command and adjust other commands positions", func(t *testing.T) {
		projectId := "proj1"
		cmd1 := testutils.NewCommand().
			WithProjectId(projectId).
			WithPosition(0).
			Data()
		cmd2 := testutils.NewCommand().
			WithProjectId(projectId).
			WithPosition(1).
			Data()
		cmd3 := testutils.NewCommand().
			WithProjectId(projectId).
			WithPosition(2).
			Data()

		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmd1),
			commandDataToModel(cmd2),
			commandDataToModel(cmd3),
		}

		h := newTestHelper(t, preloadedCommandModels)

		err := h.repo.Delete(cmd2.Id)
		assert.NoError(t, err)

		cmd1AfterDelete, err := h.repo.Get(cmd1.Id)
		assert.NoError(t, err)

		cmd3AfterDelete, err := h.repo.Get(cmd3.Id)
		assert.NoError(t, err)

		assert.Equal(t, cmd1.Position, cmd1AfterDelete.Position, "Expected cmd1 position to remain unchanged")
		assert.Equal(t, cmd3.Position-1, cmd3AfterDelete.Position, "Expected cmd3 position to be adjusted after deletion of cmd2")
	})
}

func arrange(preloadedCommandModels []*CommandModel) (repo *GormCommandRepository, tmpDbFilePath string) {
	// Initialize the database
	ctx := context.Background()

	tmp := os.TempDir()
	id := uuid.New().String()
	tmpDbFilePath = filepath.Join(tmp, id+".db")

	gormDb, err := gorm.Open(sqlite.Open(tmpDbFilePath), &gorm.Config{})
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
