package infrastructure

import (
	"context"
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

func TestGormCommandRepository_GetCommand(t *testing.T) {
	t.Parallel()
	t.Run("Should return command when it exists", func(t *testing.T) {
		preloadedCommandModels := []*CommandModel{
			{
				Id:               "cmd1",
				ProjectId:        "proj1",
				Name:             "Test Command",
				Command:          "echo 'Hello, World!'",
				WorkingDirectory: "/tmp",
				Position:         0,
			},
		}

		expectedCommand := &domain.Command{
			Id:               "cmd1",
			ProjectId:        "proj1",
			Name:             "Test Command",
			Command:          "echo 'Hello, World!'",
			WorkingDirectory: "/tmp",
			Position:         0,
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmd, err := repo.GetCommand("cmd1")
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
		preloadedCommandModels := []*CommandModel{}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmd, err := repo.GetCommand("nonexistent")
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

func TestGormCommandRepository_GetCommands(t *testing.T) {
	t.Run("Should return all commands for a project", func(t *testing.T) {
		preloadedCommandModels := []*CommandModel{
			{
				Id:               "cmd2",
				ProjectId:        "proj1",
				Name:             "Test Command 2",
				Command:          "echo 'Goodbye, World!'",
				WorkingDirectory: "/tmp",
				Position:         1,
			},
			{
				Id:               "cmd1",
				ProjectId:        "proj1",
				Name:             "Test Command 1",
				Command:          "echo 'Hello, World!'",
				WorkingDirectory: "/tmp",
				Position:         0,
			},
		}

		expectedCommands := []*domain.Command{
			{
				Id:               "cmd1",
				ProjectId:        "proj1",
				Name:             "Test Command 1",
				Command:          "echo 'Hello, World!'",
				WorkingDirectory: "/tmp",
				Position:         0,
			},
			{
				Id:               "cmd2",
				ProjectId:        "proj1",
				Name:             "Test Command 2",
				Command:          "echo 'Goodbye, World!'",
				WorkingDirectory: "/tmp",
				Position:         1,
			},
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		cmds, err := repo.GetCommands("proj1")
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

func TestGormCommandRepository_SaveCommand(t *testing.T) {
	t.Parallel()
	t.Run("Should save a new command", func(t *testing.T) {
		preloadedCommandModels := []*CommandModel{}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		newCommand := &domain.Command{
			Id:               "cmd3",
			ProjectId:        "proj1",
			Name:             "New Command",
			Command:          "echo 'New Command'",
			WorkingDirectory: "/tmp",
			Position:         2,
		}

		err := repo.SaveCommand(newCommand)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		savedCommand, err := repo.GetCommand("cmd3")
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

func TestGormCommandRepository_EditCommand(t *testing.T) {
	t.Parallel()
	t.Run("Should edit an existing command", func(t *testing.T) {
		preloadedCommandModels := []*CommandModel{
			{
				Id:               "cmd1",
				ProjectId:        "proj1",
				Name:             "Old Command",
				Command:          "echo 'Old Command'",
				WorkingDirectory: "/tmp",
				Position:         0,
			},
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		editedCommand := &domain.Command{
			Id:               "cmd1",
			ProjectId:        "proj1",
			Name:             "Edited Command",
			Command:          "echo 'Edited Command'",
			WorkingDirectory: "/tmp",
			Position:         0,
		}

		err := repo.EditCommand(editedCommand)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		savedCommand, err := repo.GetCommand("cmd1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !savedCommand.Equals(editedCommand) {
			t.Errorf("Expected command %v, got %v", editedCommand, savedCommand)
		}

		if err := os.Remove(tmpDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandRepository_DeleteCommand(t *testing.T) {
	t.Parallel()
	t.Run("Should delete an existing command", func(t *testing.T) {
		preloadedCommandModels := []*CommandModel{
			{
				Id:               "cmd1",
				ProjectId:        "proj1",
				Name:             "Command to Delete",
				Command:          "echo 'Delete Me'",
				WorkingDirectory: "/tmp",
				Position:         0,
			},
		}

		repo, tmpDbFilePath := arrange(preloadedCommandModels)

		err := repo.DeleteCommand("")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		deletedCommand, err := repo.GetCommand("cmd1")

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
