package infrastructure_test

import (
	"context"
	"slices"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	commanddomain "gomander/internal/command/domain"
	commandinfrastructure "gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/commandgroup/infrastructure"
	"gomander/internal/helpers/array"
	"gomander/internal/testutils"

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
		projectId := "project1"

		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()
		cmd2 := testutils.NewCommand().WithName("Command 2").WithProjectId(projectId).Data()
		cmd3 := testutils.NewCommand().WithName("Command 3").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3).Data()
		cmdGroup2 := testutils.NewCommandGroup().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd1, cmd3, cmd2).Data()

		groupModel, commandToCommandGroupModels, commandModels := commandGroupDataToModel(cmdGroup1)
		groupModel2, commandToCommandGroupModels2, _ := commandGroupDataToModel(cmdGroup2)

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel, groupModel2},
			slices.Concat(commandToCommandGroupModels, commandToCommandGroupModels2),
		)

		result, err := helper.repo.GetAll(projectId)

		expectedCommandGroups := []domain.CommandGroup{
			commandGroupDataToDomain(cmdGroup1),
			commandGroupDataToDomain(cmdGroup2),
		}

		assert.Nil(t, err)
		for i, group := range result {
			assert.Equal(t, expectedCommandGroups[i], group)
		}
	})
}

func TestGormCommandGroupRepository_Get(t *testing.T) {
	t.Run("Should return a command group by id with sorted commands", func(t *testing.T) {
		projectId := "project1"

		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()
		cmd2 := testutils.NewCommand().WithName("Command 2").WithProjectId(projectId).Data()
		cmd3 := testutils.NewCommand().WithName("Command 3").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3).Data()
		cmdGroup2 := testutils.NewCommandGroup().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd1, cmd3, cmd2).Data()

		groupModel, commandToCommandGroupModels, commandModels := commandGroupDataToModel(cmdGroup1)
		groupModel2, commandToCommandGroupModels2, _ := commandGroupDataToModel(cmdGroup2)

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel, groupModel2},
			slices.Concat(commandToCommandGroupModels, commandToCommandGroupModels2),
		)

		result, err := helper.repo.Get(cmdGroup1.Id)
		assert.Nil(t, err)

		expectedCommandGroup := commandGroupDataToDomain(cmdGroup1)

		assert.Equal(t, &expectedCommandGroup, result)
	})
	t.Run("Should return nil if command group does not exist", func(t *testing.T) {
		helper := newTestHelper(t, nil, nil, nil)

		result, err := helper.repo.Get("non-existent-id")

		assert.Nil(t, err)
		assert.Nil(t, result)
	})
}

func TestGormCommandGroupRepository_Create(t *testing.T) {
	t.Run("Should create a new command group and its associations", func(t *testing.T) {
		projectId := "project1"

		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()
		cmd2 := testutils.NewCommand().WithName("Command 2").WithProjectId(projectId).Data()
		cmd3 := testutils.NewCommand().WithName("Command 3").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3).Data()

		_, _, commandModels := commandGroupDataToModel(cmdGroup1)

		helper := newTestHelper(
			t,
			commandModels,
			nil,
			nil,
		)

		newGroup := commandGroupDataToDomain(cmdGroup1)

		err := helper.repo.Create(&newGroup)
		assert.Nil(t, err)

		result, err := helper.repo.Get(newGroup.Id)
		assert.Nil(t, err)
		assert.Equal(t, &newGroup, result)
	})
}

func TestGormCommandGroupRepository_Update(t *testing.T) {
	t.Run("Should update an existing command group and its associations", func(t *testing.T) {
		projectId := "project1"

		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()
		cmd2 := testutils.NewCommand().WithName("Command 2").WithProjectId(projectId).Data()
		cmd3 := testutils.NewCommand().WithName("Command 3").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3)

		groupModel, commandToCommandGroupModels, commandModels := commandGroupDataToModel(cmdGroup1.Data())

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel},
			commandToCommandGroupModels,
		)

		updatedGroup := cmdGroup1.WithName("Updated Group 1").WithCommands(cmd1, cmd2, cmd3).Data()

		groupToUpdate := commandGroupDataToDomain(updatedGroup)

		err := helper.repo.Update(&groupToUpdate)
		assert.Nil(t, err)

		result, err := helper.repo.Get(groupToUpdate.Id)
		assert.Nil(t, err)
		assert.Equal(t, &groupToUpdate, result)
	})
}

func TestGormCommandGroupRepository_Delete(t *testing.T) {
	t.Run("Should delete an existing command group and its associations", func(t *testing.T) {
		projectId := "project1"

		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1).Data()

		groupModel, commandToCommandGroupModels, commandModels := commandGroupDataToModel(cmdGroup1)

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

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).Data()
		cmdGroup2 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(1).Data()
		cmdGroup3 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(2).Data()

		group1Model, _, _ := commandGroupDataToModel(cmdGroup1)
		group2Model, _, _ := commandGroupDataToModel(cmdGroup2)
		group3Model, _, _ := commandGroupDataToModel(cmdGroup3)

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
		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()
		cmd2 := testutils.NewCommand().WithName("Command 2").WithProjectId(projectId).Data()
		cmd3 := testutils.NewCommand().WithName("Command 3").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1, cmd2, cmd3).Data()
		cmdGroup2 := testutils.NewCommandGroup().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd2, cmd1).Data()

		groupModel1, commandToCommandGroupModels1, commandModels := commandGroupDataToModel(cmdGroup1)
		groupModel2, commandToCommandGroupModels2, _ := commandGroupDataToModel(cmdGroup2)

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel1, groupModel2},
			slices.Concat(commandToCommandGroupModels1, commandToCommandGroupModels2),
		)

		err := helper.repo.RemoveCommandFromCommandGroups(cmd1.Id)
		assert.Nil(t, err)

		group1, _ := helper.repo.Get(cmdGroup1.Id)
		group2, _ := helper.repo.Get(cmdGroup2.Id)

		expectedGroup1Commands := []commanddomain.Command{commandDataToDomain(cmd2), commandDataToDomain(cmd3)}
		expectedGroup2Commands := []commanddomain.Command{commandDataToDomain(cmd2)}

		assert.Equal(t, expectedGroup1Commands, group1.Commands)
		assert.Equal(t, expectedGroup2Commands, group2.Commands)
	})
}

func TestGormCommandGroupRepository_DeleteEmptyGroups(t *testing.T) {
	t.Run("Should delete only empty command groups", func(t *testing.T) {
		projectId := "project1"
		cmd1 := testutils.NewCommand().WithName("Command 1").WithProjectId(projectId).Data()
		cmd2 := testutils.NewCommand().WithName("Command 2").WithProjectId(projectId).Data()

		cmdGroup1 := testutils.NewCommandGroup().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1, cmd2).Data()
		cmdGroup2 := testutils.NewCommandGroup().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands().Data()

		groupModel1, commandToCommandGroupModels1, commandModels := commandGroupDataToModel(cmdGroup1)
		groupModel2, _, _ := commandGroupDataToModel(cmdGroup2)

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{groupModel1, groupModel2},
			commandToCommandGroupModels1,
		)

		err := helper.repo.DeleteEmptyGroups()
		assert.Nil(t, err)

		group1, _ := helper.repo.Get(cmdGroup1.Id)
		group2, _ := helper.repo.Get(cmdGroup2.Id)
		assert.NotNil(t, group1)
		assert.Nil(t, group2)
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

func commandDataToModel(data testutils.CommandData) commandinfrastructure.CommandModel {
	return commandinfrastructure.CommandModel{
		Id:               data.Id,
		ProjectId:        data.ProjectId,
		Name:             data.Name,
		Command:          data.Command,
		WorkingDirectory: data.WorkingDirectory,
		Position:         data.Position,
	}
}

func commandDataToDomain(data testutils.CommandData) commanddomain.Command {
	return commanddomain.Command{
		Id:               data.Id,
		ProjectId:        data.ProjectId,
		Name:             data.Name,
		Command:          data.Command,
		WorkingDirectory: data.WorkingDirectory,
		Position:         data.Position,
	}
}

func commandGroupDataToModel(data testutils.CommandGroupData) (infrastructure.CommandGroupModel, []infrastructure.CommandToCommandGroupModel, []commandinfrastructure.CommandModel) {
	groupModel := infrastructure.CommandGroupModel{
		Id:        data.Id,
		ProjectId: data.ProjectId,
		Name:      data.Name,
		Position:  data.Position,
	}

	commandRelations := make([]infrastructure.CommandToCommandGroupModel, len(data.Commands))
	commandModels := make([]commandinfrastructure.CommandModel, len(data.Commands))
	for i, cmd := range data.Commands {
		commandModels[i] = commandDataToModel(cmd)
		commandRelations[i] = infrastructure.CommandToCommandGroupModel{
			CommandId:      cmd.Id,
			CommandGroupId: data.Id,
			Position:       i,
		}
	}
	return groupModel, commandRelations, commandModels
}

func commandGroupDataToDomain(data testutils.CommandGroupData) domain.CommandGroup {
	return domain.CommandGroup{
		Id:        data.Id,
		ProjectId: data.ProjectId,
		Name:      data.Name,
		Position:  data.Position,
		Commands:  array.Map(data.Commands, commandDataToDomain),
	}
}
