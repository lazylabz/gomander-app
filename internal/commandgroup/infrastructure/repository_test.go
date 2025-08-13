package infrastructure_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	command_domain "gomander/internal/command/domain"
	command_infrastructure "gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/commandgroup/infrastructure"

	_ "gomander/migrations" // Import migrations to ensure they are executed
)

var TestGetCommandGroupsTestCases = []struct {
	name                                 string
	preloadedCommandModels               []*command_infrastructure.CommandModel
	preloadedCommandGroupModels          []*infrastructure.CommandGroupModel
	preloadedCommandToCommandGroupModels []*infrastructure.CommandToCommandGroupModel
	expectedCommandGroups                []domain.CommandGroup
}{
	{
		name: "Should return all command groups with their commands",
		preloadedCommandModels: []*command_infrastructure.CommandModel{
			{
				Id:               "1",
				ProjectId:        "test-project",
				Name:             "Test Command",
				Command:          "echo 'Hello, World!'",
				WorkingDirectory: ".",
				Position:         0,
			},
		},
		preloadedCommandGroupModels: []*infrastructure.CommandGroupModel{
			{
				Id:        "1",
				ProjectId: "test-project",
				Name:      "Test Command Group",
				Position:  0,
			},
		},
		preloadedCommandToCommandGroupModels: []*infrastructure.CommandToCommandGroupModel{
			{
				CommandId:      "1",
				CommandGroupId: "1",
				Position:       0,
			},
		},
		expectedCommandGroups: []domain.CommandGroup{
			{
				Id:        "1",
				ProjectId: "test-project",
				Name:      "Test Command Group",
				Position:  0,
				Commands: []*command_domain.Command{
					{
						Id:               "1",
						ProjectId:        "test-project",
						Name:             "Test Command",
						Command:          "echo 'Hello, World!'",
						WorkingDirectory: ".",
						Position:         0,
					},
				},
			},
		},
	},
	{
		name: "Should return all command groups sorted by position with their commands sorted by position",
		preloadedCommandModels: []*command_infrastructure.CommandModel{
			{
				Id:               "1",
				ProjectId:        "test-project",
				Name:             "Command A",
				Command:          "echo 'A'",
				WorkingDirectory: ".",
				Position:         1,
			},
			{
				Id:               "2",
				ProjectId:        "test-project",
				Name:             "Command B",
				Command:          "echo 'B'",
				WorkingDirectory: ".",
				Position:         0,
			},
		},
		preloadedCommandGroupModels: []*infrastructure.CommandGroupModel{
			{
				Id:        "1",
				ProjectId: "test-project",
				Name:      "Group B",
				Position:  1,
			},
			{
				Id:        "2",
				ProjectId: "test-project",
				Name:      "Group A",
				Position:  0,
			},
		},
		preloadedCommandToCommandGroupModels: []*infrastructure.CommandToCommandGroupModel{
			{
				CommandId:      "1",
				CommandGroupId: "1",
				Position:       1,
			},
			{
				CommandId:      "2",
				CommandGroupId: "1",
				Position:       0,
			},
		},
		expectedCommandGroups: []domain.CommandGroup{
			{
				Id:        "2",
				ProjectId: "test-project",
				Name:      "Group A",
				Position:  0,
				Commands:  []*command_domain.Command{},
			},
			{
				Id:        "1",
				ProjectId: "test-project",
				Name:      "Group B",
				Position:  1,
				Commands: []*command_domain.Command{
					{
						Id:               "2",
						ProjectId:        "test-project",
						Name:             "Command B",
						Command:          "echo 'B'",
						WorkingDirectory: ".",
						Position:         0,
					},
					{
						Id:               "1",
						ProjectId:        "test-project",
						Name:             "Command A",
						Command:          "echo 'A'",
						WorkingDirectory: ".",
						Position:         1,
					},
				},
			},
		},
	},
}

func TestGormCommandGroupRepository_GetCommandGroups(t *testing.T) {
	t.Parallel()
	for _, testCase := range TestGetCommandGroupsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo, tempDbFilePath := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

			result, err := repo.GetCommandGroups("test-project")

			if err != nil {
				t.Fatalf("Failed to get command groups: %v", err)
			}
			for i, group := range result {
				if !group.Equals(&testCase.expectedCommandGroups[i]) {
					t.Errorf("Expected command group %v, got %v", testCase.expectedCommandGroups[i], group)
				}
			}

			if err := os.Remove(tempDbFilePath); err != nil {
				t.Fatalf("Failed to remove temporary database file: %v", err)
			}
		})
	}
}

func TestGormCommandGroupRepository_GetCommandGroupById(t *testing.T) {
	t.Run("Should return a command group by id", func(t *testing.T) {
		testCase := TestGetCommandGroupsTestCases[0]

		repo, tempDbFilePath := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

		result, err := repo.GetCommandGroupById("1")

		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result == nil {
			t.Fatal("Expected command group, got nil")
		}
		if !result.Equals(&testCase.expectedCommandGroups[0]) {
			t.Errorf("Expected command group %v, got %v", testCase.expectedCommandGroups[0], result)
		}
		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
	t.Run("Should return a command group by id with sorted commands", func(t *testing.T) {
		testCase := TestGetCommandGroupsTestCases[1]

		repo, tempDbFilePath := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

		result, err := repo.GetCommandGroupById("1")
		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result == nil {
			t.Fatal("Expected command group, got nil")
		}
		if !testCase.expectedCommandGroups[1].Equals(result) {
			t.Errorf("Expected command group %v, got %v", testCase.expectedCommandGroups[1], result)
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
	t.Run("Should return nil if command group does not exist", func(t *testing.T) {
		repo, tempDbFilePath := arrange(nil, nil, nil)

		result, err := repo.GetCommandGroupById("non-existent-id")

		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result != nil {
			t.Fatal("Expected nil command group, got non-nil")
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandGroupRepository_CreateCommandGroup(t *testing.T) {
	t.Run("Should create a new command group", func(t *testing.T) {
		repo, tempDbFilePath := arrange(
			make([]*command_infrastructure.CommandModel, 0),
			make([]*infrastructure.CommandGroupModel, 0),
			make([]*infrastructure.CommandToCommandGroupModel, 0),
		)

		newGroup := &domain.CommandGroup{
			Id:        "1",
			ProjectId: "test-project",
			Name:      "New Command Group",
			Position:  0,
		}

		err := repo.CreateCommandGroup(newGroup)
		if err != nil {
			t.Fatalf("Failed to create command group: %v", err)
		}

		result, err := repo.GetCommandGroupById("1")
		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result == nil {
			t.Fatal("Expected command group, got nil")
		}
		if !result.Equals(newGroup) {
			t.Errorf("Expected command group %v, got %v", newGroup, result)
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandGroupRepository_UpdateCommandGroup(t *testing.T) {
	t.Run("Should update an existing command group", func(t *testing.T) {
		testCase := TestGetCommandGroupsTestCases[0]

		repo, tempDbFilePath := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

		groupToUpdate := &domain.CommandGroup{
			Id:        "1",
			ProjectId: "test-project",
			Name:      "Updated Command Group",
			Position:  1,
		}

		err := repo.UpdateCommandGroup(groupToUpdate)
		if err != nil {
			t.Fatalf("Failed to update command group: %v", err)
		}

		result, err := repo.GetCommandGroupById("1")
		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result == nil {
			t.Fatal("Expected command group, got nil")
		}
		if result.Name != groupToUpdate.Name {
			t.Errorf("Expected command group name %s, got %s", groupToUpdate.Name, result.Name)
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandGroupRepository_DeleteCommandGroup(t *testing.T) {
	t.Run("Should delete an existing command group", func(t *testing.T) {
		testCase := TestGetCommandGroupsTestCases[0]

		repo, tempDbFilePath := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

		err := repo.DeleteCommandGroup("1")
		if err != nil {
			t.Fatalf("Failed to delete command group: %v", err)
		}

		result, err := repo.GetCommandGroupById("1")
		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result != nil {
			t.Fatal("Expected nil command group, got non-nil")
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func arrange(preloadedCommandModels []*command_infrastructure.CommandModel, preloadedCommandGroupModels []*infrastructure.CommandGroupModel, preloadedCommandToCommandGroupModels []*infrastructure.CommandToCommandGroupModel) (repo *infrastructure.GormCommandGroupRepository, tmpDbFilePath string) {
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
		err = gorm.G[command_infrastructure.CommandModel](gormDb).Create(ctx, m)
		if err != nil {
			panic(err)
		}
	}

	for _, m := range preloadedCommandGroupModels {
		err = gorm.G[infrastructure.CommandGroupModel](gormDb).Create(ctx, m)
		if err != nil {
			panic(err)
		}
	}

	for _, m := range preloadedCommandToCommandGroupModels {
		err = gorm.G[infrastructure.CommandToCommandGroupModel](gormDb).Create(ctx, m)
		if err != nil {
			panic(err)
		}
	}

	repo = infrastructure.NewGormCommandGroupRepository(
		gormDb,
		ctx,
	)

	return
}
