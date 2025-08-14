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

	commanddomain "gomander/internal/command/domain"
	commandinfrastructure "gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/commandgroup/infrastructure"

	_ "gomander/migrations" // Import migrations to ensure they are executed
)

var (
	command1 = commandinfrastructure.CommandModel{
		Id:               "1",
		ProjectId:        "test-project",
		Name:             "Test Command",
		Command:          "echo 'Hello, World!'",
		WorkingDirectory: ".",
		Position:         0,
	}
	command2 = commandinfrastructure.CommandModel{
		Id:               "2",
		ProjectId:        "test-project",
		Name:             "Command A",
		Command:          "echo 'A'",
		WorkingDirectory: ".",
		Position:         1,
	}
	command3 = commandinfrastructure.CommandModel{
		Id:               "3",
		ProjectId:        "test-project",
		Name:             "Command B",
		Command:          "echo 'B'",
		WorkingDirectory: ".",
		Position:         2,
	}
	group1Model = infrastructure.CommandGroupModel{
		Id:        "1",
		ProjectId: "test-project",
		Name:      "Test Command Group 1",
		Position:  0,
	}
	group2Model = infrastructure.CommandGroupModel{
		Id:        "2",
		ProjectId: "test-project",
		Name:      "Test Command Group 2",
		Position:  1,
	}
	group3Model = infrastructure.CommandGroupModel{
		Id:        "3",
		ProjectId: "test-project",
		Name:      "Test Command Group 3",
		Position:  2,
	}
	group1command1relation = infrastructure.CommandToCommandGroupModel{
		CommandId:      "1",
		CommandGroupId: "1",
		Position:       0,
	}
	group2command2relation = infrastructure.CommandToCommandGroupModel{
		CommandId:      "2",
		CommandGroupId: "2",
		Position:       1,
	}
	group2command3relation = infrastructure.CommandToCommandGroupModel{
		CommandId:      "3",
		CommandGroupId: "2",
		Position:       0,
	}
	domainCommand1 = commanddomain.Command{
		Id:               "1",
		ProjectId:        "test-project",
		Name:             "Test Command",
		Command:          "echo 'Hello, World!'",
		WorkingDirectory: ".",
		Position:         0,
	}
	domainCommand2 = commanddomain.Command{
		Id:               "2",
		ProjectId:        "test-project",
		Name:             "Command A",
		Command:          "echo 'A'",
		WorkingDirectory: ".",
		Position:         1,
	}
	domainCommand3 = commanddomain.Command{
		Id:               "3",
		ProjectId:        "test-project",
		Name:             "Command B",
		Command:          "echo 'B'",
		WorkingDirectory: ".",
		Position:         2,
	}
	group1Domain = domain.CommandGroup{
		Id:        "1",
		ProjectId: "test-project",
		Name:      "Test Command Group 1",
		Position:  0,
		Commands:  []commanddomain.Command{domainCommand1},
	}
	group2domain = domain.CommandGroup{
		Id:        "2",
		ProjectId: "test-project",
		Name:      "Test Command Group 2",
		Position:  1,
		Commands:  []commanddomain.Command{domainCommand3, domainCommand2},
	}
	group3domain = domain.CommandGroup{
		Id:        "3",
		ProjectId: "test-project",
		Name:      "Test Command Group 3",
		Position:  2,
		Commands:  []commanddomain.Command{},
	}
)

var TestGetCommandGroupsTestCases = []struct {
	name                                 string
	preloadedCommandModels               []commandinfrastructure.CommandModel
	preloadedCommandGroupModels          []infrastructure.CommandGroupModel
	preloadedCommandToCommandGroupModels []infrastructure.CommandToCommandGroupModel
	expectedCommandGroups                []domain.CommandGroup
}{
	{
		name: "Should return all command groups with their commands",
		preloadedCommandModels: []commandinfrastructure.CommandModel{
			command1,
		},
		preloadedCommandGroupModels: []infrastructure.CommandGroupModel{
			group1Model,
		},
		preloadedCommandToCommandGroupModels: []infrastructure.CommandToCommandGroupModel{
			group1command1relation,
		},
		expectedCommandGroups: []domain.CommandGroup{
			group1Domain,
		},
	},
	{
		name: "Should return all command groups sorted by position with their commands sorted by position",
		preloadedCommandModels: []commandinfrastructure.CommandModel{
			command2,
			command3,
		},
		preloadedCommandGroupModels: []infrastructure.CommandGroupModel{
			group2Model,
			group3Model,
		},
		preloadedCommandToCommandGroupModels: []infrastructure.CommandToCommandGroupModel{
			group2command2relation,
			group2command3relation,
		},
		expectedCommandGroups: []domain.CommandGroup{
			group2domain,
			group3domain,
		},
	},
}

func TestGormCommandGroupRepository_GetCommandGroups(t *testing.T) {
	t.Parallel()
	for _, testCase := range TestGetCommandGroupsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo, tempDbFilePath, _ := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

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

		repo, tempDbFilePath, _ := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

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

		repo, tempDbFilePath, _ := arrange(testCase.preloadedCommandModels, testCase.preloadedCommandGroupModels, testCase.preloadedCommandToCommandGroupModels)

		result, err := repo.GetCommandGroupById("3")
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
		repo, tempDbFilePath, _ := arrange(nil, nil, nil)

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
	t.Run("Should create a new command group and its associations", func(t *testing.T) {
		preloadedCommandModels := []commandinfrastructure.CommandModel{command2, command3}
		repo, tempDbFilePath, _ := arrange(preloadedCommandModels, nil, nil)

		newGroup := &domain.CommandGroup{
			Id:        "4",
			ProjectId: "test-project",
			Name:      "New Command Group",
			Position:  0,
			Commands: []commanddomain.Command{
				domainCommand2,
				domainCommand3,
			},
		}

		err := repo.CreateCommandGroup(newGroup)
		if err != nil {
			t.Fatalf("Failed to create command group: %v", err)
		}

		result, err := repo.GetCommandGroupById("4")
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
	t.Run("Should update an existing command group and its associations", func(t *testing.T) {
		preloadedCommandModels := []commandinfrastructure.CommandModel{command1, command2}
		preloadedCommandGroupModels := []infrastructure.CommandGroupModel{group1Model}
		preloadedCommandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{group1command1relation}

		repo, tempDbFilePath, _ := arrange(preloadedCommandModels, preloadedCommandGroupModels, preloadedCommandToCommandGroupModels)

		groupToUpdate := &domain.CommandGroup{
			Id:        group1Model.Id,
			ProjectId: group1Model.ProjectId,
			Name:      "Updated Command Group",
			Position:  group1Model.Position,
			Commands: []commanddomain.Command{
				domainCommand1, // Assuming command1 is already preloaded
				domainCommand2,
			},
		}

		err := repo.UpdateCommandGroup(groupToUpdate)
		if err != nil {
			t.Fatalf("Failed to update command group: %v", err)
		}

		result, err := repo.GetCommandGroupById(groupToUpdate.Id)
		if err != nil {
			t.Fatalf("Failed to get command group by id: %v", err)
		}
		if result == nil {
			t.Fatal("Expected command group, got nil")
		}
		if !result.Equals(groupToUpdate) {
			t.Errorf("Expected command group %v, got %v", groupToUpdate, result)
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func TestGormCommandGroupRepository_DeleteCommandGroup(t *testing.T) {
	t.Run("Should delete an existing command group", func(t *testing.T) {
		preloadedCommandModels := []commandinfrastructure.CommandModel{command1}
		preloadedCommandGroupModels := []infrastructure.CommandGroupModel{group1Model}
		preloadedCommandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{group1command1relation}

		repo, tempDbFilePath, db := arrange(preloadedCommandModels, preloadedCommandGroupModels, preloadedCommandToCommandGroupModels)

		err := repo.DeleteCommandGroup(group1Model.Id)
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

		existingRelations, err := gorm.G[infrastructure.CommandToCommandGroupModel](db).Where("command_group_id = ?", group1Model.Id).Find(context.Background())
		if err != nil {
			t.Fatalf("Failed to check command group relations: %v", err)
		}
		if len(existingRelations) > 0 {
			t.Fatal("Expected no command group relations, found some")
		}

		existingCommands, err := gorm.G[commandinfrastructure.CommandModel](db).Where("id = ?", command1.Id).Find(context.Background())
		if err != nil {
			t.Fatalf("Failed to check command existence: %v", err)
		}
		if len(existingCommands) == 0 {
			t.Fatal("Expected command to still exist, it was deleted")
		}

		if err := os.Remove(tempDbFilePath); err != nil {
			t.Fatalf("Failed to remove temporary database file: %v", err)
		}
	})
}

func arrange(
	preloadedCommandModels []commandinfrastructure.CommandModel,
	preloadedCommandGroupModels []infrastructure.CommandGroupModel,
	preloadedCommandToCommandGroupModels []infrastructure.CommandToCommandGroupModel,
) (repo *infrastructure.GormCommandGroupRepository, tmpDbFilePath string, gormDb *gorm.DB) {
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
		err = gorm.G[commandinfrastructure.CommandModel](gormDb).Create(ctx, &m)
		if err != nil {
			panic(err)
		}
	}

	for _, m := range preloadedCommandGroupModels {
		err = gorm.G[infrastructure.CommandGroupModel](gormDb).Create(ctx, &m)
		if err != nil {
			panic(err)
		}
	}

	for _, m := range preloadedCommandToCommandGroupModels {
		err = gorm.G[infrastructure.CommandToCommandGroupModel](gormDb).Create(ctx, &m)
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
