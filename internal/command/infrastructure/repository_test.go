package commandinfrastructure

import (
	"context"
	"gomander/internal/testutils"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
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

func TestGormCommandRepository_Get(t *testing.T) {
	t.Parallel()
	t.Run("Should return command when it exists", func(t *testing.T) {
		cmdData := testutils.NewCommand().
			WithId("cmd1").
			WithProjectId("proj1").
			WithName("Test Command").
			WithCommand("echo 'Hello, World!'").
			WithWorkingDirectory("/tmp").
			WithPosition(0).
			Data()
		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmdData),
		}

		expectedCommand := commandDataToDomain(cmdData)

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmd, err := repo.Get("cmd1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !cmd.Equals(expectedCommand) {
			t.Errorf("Expected command %v, got %v", expectedCommand, cmd)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
	t.Run("Should return nil when it doesn't exists", func(t *testing.T) {
		var preloadedCommandModels []*CommandModel

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmd, err := repo.Get("nonexistent")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cmd != nil {
			t.Errorf("Expected nil command, got %v", cmd)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandRepository_GetAll(t *testing.T) {
	t.Run("Should return all commands for a project", func(t *testing.T) {
		cmd1 := testutils.NewCommand().
			WithId("cmd1").
			WithProjectId("proj1").
			WithName("Test Command 1").
			WithCommand("echo 'Hello, World!'").
			WithWorkingDirectory("/tmp").
			WithPosition(0).
			Data()
		cmd2 := testutils.NewCommand().
			WithId("cmd2").
			WithProjectId("proj1").
			WithName("Test Command 2").
			WithCommand("echo 'Goodbye, World!'").
			WithWorkingDirectory("/tmp").
			WithPosition(1).
			Data()

		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmd2),
			commandDataToModel(cmd1),
		}

		expectedCommands := []*domain.Command{
			commandDataToDomain(cmd1),
			commandDataToDomain(cmd2),
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmds, err := repo.GetAll("proj1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		for i, cmd := range cmds {
			if !cmd.Equals(expectedCommands[i]) {
				t.Errorf("Expected command %v, got %v", expectedCommands[i], cmd)
			}
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandRepository_Save(t *testing.T) {
	t.Parallel()
	t.Run("Should save a new command", func(t *testing.T) {
		var preloadedCommandModels []*CommandModel

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmdData := testutils.NewCommand().
			WithId("cmd3").
			WithProjectId("proj1").
			WithName("New Command").
			WithCommand("echo 'New Command'").
			WithWorkingDirectory("/tmp").
			WithPosition(2).
			Data()
		newCommand := commandDataToDomain(cmdData)

		err := repo.Create(newCommand)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		savedCommand, err := repo.Get("cmd3")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !savedCommand.Equals(newCommand) {
			t.Errorf("Expected command %v, got %v", newCommand, savedCommand)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandRepository_Edit(t *testing.T) {
	t.Parallel()
	t.Run("Should edit an existing command", func(t *testing.T) {
		existingCommand := testutils.NewCommand().
			WithId("cmd1").
			WithProjectId("proj1").
			WithName("Old Command").
			WithCommand("echo 'Old Command'").
			WithWorkingDirectory("/tmp").
			WithPosition(0)

		repo, tmpDbFilePath := arrange([]*CommandModel{
			commandDataToModel(existingCommand.Data()),
		})

		editedCommand := existingCommand.
			WithName("Edited Command").
			WithCommand("echo 'Edited Command'").
			Data()
		domainEditedCommand := commandDataToDomain(editedCommand)

		err := repo.Update(domainEditedCommand)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		savedCommand, err := repo.Get("cmd1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !savedCommand.Equals(domainEditedCommand) {
			t.Errorf("Expected command %v, got %v", editedCommand, savedCommand)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandRepository_Delete(t *testing.T) {
	t.Parallel()
	t.Run("Should delete an existing command", func(t *testing.T) {
		cmd := testutils.NewCommand().
			WithId("cmd1").
			WithProjectId("proj1").
			WithName("Command to Delete").
			WithCommand("echo 'Delete Me'").
			WithWorkingDirectory("/tmp").
			WithPosition(0).
			Data()
		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmd),
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		err := repo.Delete("cmd1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		deletedCommand, err := repo.Get("cmd1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if deletedCommand != nil {
			t.Errorf("Expected nil command after deletion, got %v", deletedCommand)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
	t.Run("Should delete an existing command and adjust other commands positions", func(t *testing.T) {
		cmd1 := testutils.NewCommand().
			WithId("cmd1").
			WithProjectId("proj1").
			WithName("Command to Delete").
			WithCommand("echo 'Delete Me'").
			WithWorkingDirectory("/tmp").
			WithPosition(0).
			Data()
		cmd2 := testutils.NewCommand().
			WithId("cmd2").
			WithProjectId("proj1").
			WithName("Command 2").
			WithCommand("echo 'Command 2'").
			WithWorkingDirectory("/tmp").
			WithPosition(1).
			Data()
		cmd3 := testutils.NewCommand().
			WithId("cmd3").
			WithProjectId("proj1").
			WithName("Command 3").
			WithCommand("echo 'Command 3'").
			WithWorkingDirectory("/tmp").
			WithPosition(2).
			Data()

		preloadedCommandModels := []*CommandModel{
			commandDataToModel(cmd1),
			commandDataToModel(cmd2),
			commandDataToModel(cmd3),
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		err := repo.Delete("cmd2")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		cmd1AfterDelete, err := repo.Get("cmd1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		cmd3AfterDelete, err := repo.Get("cmd3")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cmd1AfterDelete.Position != cmd1.Position {
			t.Errorf("Expected command cmd1 position to remain %d, got %d", cmd1.Position, cmd1AfterDelete.Position)
		}
		if cmd3AfterDelete.Position != cmd3.Position-1 {
			t.Errorf("Expected command cmd3 position to be adjusted to %d, got %d", cmd3.Position-1, cmd3AfterDelete.Position)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
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
