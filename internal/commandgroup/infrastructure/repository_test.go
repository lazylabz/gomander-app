package infrastructure_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/command/domain/test"
	commandinfrastructure "gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	test2 "gomander/internal/commandgroup/domain/test"
	"gomander/internal/commandgroup/infrastructure"
	_ "gomander/migrations" // Import migrations to ensure they are executed
)

type testHelper struct {
	t      *testing.T
	repo   *infrastructure.GormCommandGroupRepository
	gormDb *gorm.DB
}

func newTestHelper(t *testing.T,
	preloadedCommandModels []commandinfrastructure.CommandModel,
	preloadedCommandGroupModels []infrastructure.CommandGroupModel,
	preloadedCommandToCommandGroupModels []infrastructure.CommandToCommandGroupModel) *testHelper {
	t.Helper() // IMPORTANT: This marks the function as a helper, so error traces will point to the test instead of here

	repo, gormDb := arrange(
		preloadedCommandModels,
		preloadedCommandGroupModels,
		preloadedCommandToCommandGroupModels,
	)

	helper := &testHelper{
		t:      t,
		repo:   repo,
		gormDb: gormDb,
	}

	return helper
}

func TestGormCommandGroupRepository_GetAll(t *testing.T) {
	t.Run("Should return all command groups sorted by position with their commands sorted by position", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(projectId).Build()

		cmd1Model := commandinfrastructure.ToCommandModel(&cmd1)
		cmd2Model := commandinfrastructure.ToCommandModel(&cmd2)
		cmd3Model := commandinfrastructure.ToCommandModel(&cmd3)

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd1, cmd3, cmd2).Build()

		cmdGroup1Model := infrastructure.ToCommandGroupModel(&cmdGroup1)
		cmdGroup2Model := infrastructure.ToCommandGroupModel(&cmdGroup2)

		cmdToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			// CommandGroup 1 associations
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd2.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       1,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd3.Id,
				Position:       2,
			},
			// CommandGroup 2 associations
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd1.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd3.Id,
				Position:       1,
			},
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd2.Id,
				Position:       2,
			},
		}

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{cmd1Model, cmd2Model, cmd3Model},
			[]infrastructure.CommandGroupModel{cmdGroup1Model, cmdGroup2Model},
			cmdToCommandGroupModels,
		)

		// Act
		result, err := helper.repo.GetAll(projectId)

		// Assert
		expectedCommandGroups := []domain.CommandGroup{
			cmdGroup1,
			cmdGroup2,
		}

		assert.Nil(t, err)
		for i, group := range result {
			assert.Equal(t, expectedCommandGroups[i], group)
		}
	})
}

func TestGormCommandGroupRepository_Get(t *testing.T) {
	t.Run("Should return a command group by id with sorted commands", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(projectId).Build()

		cmd1Model := commandinfrastructure.ToCommandModel(&cmd1)
		cmd2Model := commandinfrastructure.ToCommandModel(&cmd2)
		cmd3Model := commandinfrastructure.ToCommandModel(&cmd3)

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd1, cmd3, cmd2).Build()

		cmdGroup1Model := infrastructure.ToCommandGroupModel(&cmdGroup1)
		cmdGroup2Model := infrastructure.ToCommandGroupModel(&cmdGroup2)

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			// CommandGroup 1 associations
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd2.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       1,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd3.Id,
				Position:       2,
			},
			// CommandGroup 2 associations
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd1.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd3.Id,
				Position:       1,
			},
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd2.Id,
				Position:       2,
			},
		}

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{cmd1Model, cmd2Model, cmd3Model},
			[]infrastructure.CommandGroupModel{cmdGroup1Model, cmdGroup2Model},
			commandToCommandGroupModels,
		)

		// Act
		result, err := helper.repo.Get(cmdGroup1.Id)

		// Assert
		assert.Nil(t, err)

		assert.Equal(t, &cmdGroup1, result)
	})
	t.Run("Should return nil if command group does not exist", func(t *testing.T) {
		// Arrange
		helper := newTestHelper(t, nil, nil, nil)

		// Act
		result, err := helper.repo.Get("non-existent-id")

		// Assert
		assert.Nil(t, err)
		assert.Nil(t, result)
	})
}

func TestGormCommandGroupRepository_Create(t *testing.T) {
	t.Run("Should create a new command group and its associations", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(projectId).Build()

		cmd1Model := commandinfrastructure.ToCommandModel(&cmd1)
		cmd2Model := commandinfrastructure.ToCommandModel(&cmd2)
		cmd3Model := commandinfrastructure.ToCommandModel(&cmd3)

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3).Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{cmd1Model, cmd2Model, cmd3Model},
			nil,
			nil,
		)

		// Act
		err := helper.repo.Create(&cmdGroup1)

		// Assert
		assert.Nil(t, err)

		// Verify the group was created correctly
		result, err := helper.repo.Get(cmdGroup1.Id)
		assert.Nil(t, err)
		assert.Equal(t, &cmdGroup1, result)
	})
}

func TestGormCommandGroupRepository_Update(t *testing.T) {
	t.Run("Should update an existing command group and its associations", func(t *testing.T) {
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(projectId).Build()

		commandModels := []commandinfrastructure.CommandModel{
			commandinfrastructure.ToCommandModel(&cmd1),
			commandinfrastructure.ToCommandModel(&cmd2),
			commandinfrastructure.ToCommandModel(&cmd3),
		}

		cmdGroup1Builder := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3)
		cmdGroup1 := cmdGroup1Builder.Build()

		groupModel := infrastructure.ToCommandGroupModel(&cmdGroup1)

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd2.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       1,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd3.Id,
				Position:       2,
			},
		}

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel},
			commandToCommandGroupModels,
		)

		updatedGroup := cmdGroup1Builder.WithName("Updated Group 1").WithCommands(cmd1, cmd2, cmd3).Build()

		err := helper.repo.Update(&updatedGroup)
		assert.Nil(t, err)

		result, err := helper.repo.Get(updatedGroup.Id)
		assert.Nil(t, err)
		assert.Equal(t, &updatedGroup, result)
	})
}

func TestGormCommandGroupRepository_Delete(t *testing.T) {
	t.Run("Should delete an existing command group and its associations", func(t *testing.T) {
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()

		commandModels := []commandinfrastructure.CommandModel{
			commandinfrastructure.ToCommandModel(&cmd1),
		}

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1).Build()

		groupModel := infrastructure.ToCommandGroupModel(&cmdGroup1)

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       0,
			},
		}

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel},
			commandToCommandGroupModels,
		)

		err := helper.repo.Delete(cmdGroup1.Id)
		assert.Nil(t, err)

		result, err := helper.repo.Get(cmdGroup1.Id)
		assert.Nil(t, err)
		assert.Nil(t, result)

		existingRelations, err := gorm.G[infrastructure.CommandToCommandGroupModel](helper.gormDb).Where("command_group_id = ?", cmdGroup1.Id).Find(context.Background())
		assert.Nil(t, err)
		assert.Len(t, existingRelations, 0)

		existingCommands, err := gorm.G[commandinfrastructure.CommandModel](helper.gormDb).Where("id = ?", cmd1.Id).Find(context.Background())
		assert.Nil(t, err)
		assert.Len(t, existingCommands, 1)
	})
	t.Run("Should delete an existing command groups and correctly update positions of other command groups", func(t *testing.T) {
		projectId := "project1"

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(1).Build()
		cmdGroup3 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(2).Build()

		group1Model := infrastructure.ToCommandGroupModel(&cmdGroup1)
		group2Model := infrastructure.ToCommandGroupModel(&cmdGroup2)
		group3Model := infrastructure.ToCommandGroupModel(&cmdGroup3)

		helper := newTestHelper(
			t,
			nil,
			[]infrastructure.CommandGroupModel{group1Model, group2Model, group3Model},
			nil,
		)

		err := helper.repo.Delete(group2Model.Id)
		assert.Nil(t, err)

		resultGroup1, err := helper.repo.Get(group1Model.Id)
		assert.Nil(t, err)

		resultGroup3, err := helper.repo.Get(group3Model.Id)
		assert.Nil(t, err)

		// Check if the positions of the remaining command groups are updated correctly
		assert.Equal(t, resultGroup1.Position, group1Model.Position)
		assert.Equal(t, resultGroup3.Position, group3Model.Position-1)
	})
}

func TestGormCommandGroupRepository_RemoveCommandFromCommandGroups(t *testing.T) {
	t.Run("Should remove a command from all groups and update positions", func(t *testing.T) {
		projectId := "project1"
		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(projectId).Build()

		commandModels := []commandinfrastructure.CommandModel{
			commandinfrastructure.ToCommandModel(&cmd1),
			commandinfrastructure.ToCommandModel(&cmd2),
			commandinfrastructure.ToCommandModel(&cmd3),
		}

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1, cmd2, cmd3).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd2, cmd1).Build()

		groupModel1 := infrastructure.ToCommandGroupModel(&cmdGroup1)
		groupModel2 := infrastructure.ToCommandGroupModel(&cmdGroup2)

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			// CommandGroup 1 associations
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd2.Id,
				Position:       1,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd3.Id,
				Position:       2,
			},
			// CommandGroup 2 associations
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd2.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd1.Id,
				Position:       1,
			},
		}

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel1, groupModel2},
			commandToCommandGroupModels,
		)

		err := helper.repo.RemoveCommandFromCommandGroups(cmd1.Id)
		assert.Nil(t, err)

		group1, _ := helper.repo.Get(cmdGroup1.Id)
		group2, _ := helper.repo.Get(cmdGroup2.Id)

		expectedGroup1Commands := []commanddomain.Command{cmd2, cmd3}
		expectedGroup2Commands := []commanddomain.Command{cmd2}

		assert.Equal(t, expectedGroup1Commands, group1.Commands)
		assert.Equal(t, expectedGroup2Commands, group2.Commands)
	})
}

func TestGormCommandGroupRepository_DeleteEmpty(t *testing.T) {
	t.Run("Should delete only empty command groups", func(t *testing.T) {
		projectId := "project1"
		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()

		commandModels := []commandinfrastructure.CommandModel{
			commandinfrastructure.ToCommandModel(&cmd1),
			commandinfrastructure.ToCommandModel(&cmd2),
		}

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1, cmd2).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands().Build()

		groupModel1 := infrastructure.ToCommandGroupModel(&cmdGroup1)
		groupModel2 := infrastructure.ToCommandGroupModel(&cmdGroup2)

		commandToCommandGroupModels1 := []infrastructure.CommandToCommandGroupModel{
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd2.Id,
				Position:       1,
			},
		}

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel1, groupModel2},
			commandToCommandGroupModels1,
		)

		ids, err := helper.repo.DeleteEmpty()
		assert.Nil(t, err)
		assert.Equal(t, []string{cmdGroup2.Id}, ids)

		group1, _ := helper.repo.Get(cmdGroup1.Id)
		group2, _ := helper.repo.Get(cmdGroup2.Id)
		assert.NotNil(t, group1)
		assert.Nil(t, group2)
	})
}

func TestGormCommandGroupRepository_DeleteAll(t *testing.T) {
	t.Run("Should delete all command groups and their associations for a project", func(t *testing.T) {
		projectId := "project1"
		otherProjectId := "project2"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(otherProjectId).Build()

		commandModels := []commandinfrastructure.CommandModel{
			commandinfrastructure.ToCommandModel(&cmd1),
			commandinfrastructure.ToCommandModel(&cmd2),
			commandinfrastructure.ToCommandModel(&cmd3),
		}

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1, cmd2).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd2).Build()
		cmdGroupOther := test2.NewCommandGroupBuilder().WithName("Other Group").WithProjectId(otherProjectId).WithPosition(0).WithCommands(cmd3).Build()

		groupModel1 := infrastructure.ToCommandGroupModel(&cmdGroup1)
		groupModel2 := infrastructure.ToCommandGroupModel(&cmdGroup2)
		groupModelOther := infrastructure.ToCommandGroupModel(&cmdGroupOther)

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			// CommandGroup 1 associations
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd1.Id,
				Position:       0,
			},
			{
				CommandGroupId: cmdGroup1.Id,
				CommandId:      cmd2.Id,
				Position:       1,
			},
			// CommandGroup 2 associations
			{
				CommandGroupId: cmdGroup2.Id,
				CommandId:      cmd2.Id,
				Position:       0,
			},
			// Other CommandGroup associations
			{
				CommandGroupId: cmdGroupOther.Id,
				CommandId:      cmd3.Id,
				Position:       0,
			},
		}

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel1, groupModel2, groupModelOther},
			commandToCommandGroupModels,
		)

		ids, err := helper.repo.DeleteAll(projectId)
		assert.Nil(t, err)
		assert.Equal(t, []string{cmdGroup1.Id, cmdGroup2.Id}, ids)

		// Command groups from the specified project should be deleted
		result1, _ := helper.repo.Get(cmdGroup1.Id)
		result2, _ := helper.repo.Get(cmdGroup2.Id)
		assert.Nil(t, result1)
		assert.Nil(t, result2)

		// Other group from a different project should remain
		resultOther, _ := helper.repo.Get(cmdGroupOther.Id)
		assert.NotNil(t, resultOther)

		// All associations for the deleted groups should be removed
		relations, err := gorm.G[infrastructure.CommandToCommandGroupModel](helper.gormDb).Where("command_group_id IN ?", []string{cmdGroup1.Id, cmdGroup2.Id}).Find(context.Background())
		assert.Nil(t, err)
		assert.Len(t, relations, 0)

		// Associations for the other group should remain
		relationsOther, err := gorm.G[infrastructure.CommandToCommandGroupModel](helper.gormDb).Where("command_group_id = ?", cmdGroupOther.Id).Find(context.Background())
		assert.Nil(t, err)
		assert.NotEmpty(t, relationsOther)
	})
}

func arrange(
	preloadedCommandModels []commandinfrastructure.CommandModel,
	preloadedCommandGroupModels []infrastructure.CommandGroupModel,
	preloadedCommandToCommandGroupModels []infrastructure.CommandToCommandGroupModel,
) (repo *infrastructure.GormCommandGroupRepository, gormDb *gorm.DB) {
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

	// Clean all tables
	_, err = gorm.G[commandinfrastructure.CommandModel](gormDb).Where("true").Delete(ctx)
	if err != nil {
		panic(err)
	}
	_, err = gorm.G[infrastructure.CommandToCommandGroupModel](gormDb).Where("true").Delete(ctx)
	if err != nil {
		panic(err)
	}
	_, err = gorm.G[infrastructure.CommandGroupModel](gormDb).Where("true").Delete(ctx)
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
