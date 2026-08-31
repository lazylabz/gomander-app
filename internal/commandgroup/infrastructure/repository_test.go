package infrastructure_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"gomander/internal/command/domain/test"
	commandinfrastructure "gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	test2 "gomander/internal/commandgroup/domain/test"
	"gomander/internal/commandgroup/infrastructure"
	"gomander/internal/domainerrors"
	"gomander/internal/testdb"
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
		t,
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
		assert.Equal(t, expectedCommandGroups, result)
	})

	t.Run("Should return a command group that has no commands", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		empty := test2.NewCommandGroupBuilder().WithName("Empty").WithProjectId(projectId).WithPosition(0).Build()

		helper := newTestHelper(t, nil, []infrastructure.CommandGroupModel{infrastructure.ToCommandGroupModel(&empty)}, nil)

		// Act
		result, err := helper.repo.GetAll(projectId)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, []domain.CommandGroup{empty}, result)
	})

	t.Run("Should read the command groups and their commands in a single query", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()

		cmdGroup1 := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1).Build()
		cmdGroup2 := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd2).Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{
				commandinfrastructure.ToCommandModel(&cmd1),
				commandinfrastructure.ToCommandModel(&cmd2),
			},
			[]infrastructure.CommandGroupModel{
				infrastructure.ToCommandGroupModel(&cmdGroup1),
				infrastructure.ToCommandGroupModel(&cmdGroup2),
			},
			[]infrastructure.CommandToCommandGroupModel{
				{CommandGroupId: cmdGroup1.Id, CommandId: cmd1.Id, Position: 0},
				{CommandGroupId: cmdGroup2.Id, CommandId: cmd2.Id, Position: 0},
			},
		)

		recorder := &sqlRecorder{Interface: gormlogger.Discard}
		repo := infrastructure.NewGormCommandGroupRepository(
			helper.gormDb.Session(&gorm.Session{Logger: recorder}),
			context.Background(),
		)

		// Act
		result, err := repo.GetAll(projectId)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, []domain.CommandGroup{cmdGroup1, cmdGroup2}, result)
		assert.Len(t, recorder.statements, 1)
	})
}

// sqlRecorder keeps every statement GORM executes, so a test can say how many a
// read costs.
type sqlRecorder struct {
	gormlogger.Interface
	statements []string
}

func (r *sqlRecorder) Trace(_ context.Context, _ time.Time, statement func() (string, int64), _ error) {
	sql, _ := statement()
	r.statements = append(r.statements, sql)
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

		assert.Equal(t, cmdGroup1, result)
	})
	t.Run("Should report a command group that does not exist as not found", func(t *testing.T) {
		// Arrange
		helper := newTestHelper(t, nil, nil, nil)

		// Act
		_, err := helper.repo.Get("non-existent-id")

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	})
}

func TestGormCommandGroupRepository_GetAllWithCommandIds(t *testing.T) {
	t.Run("Should return all command groups sorted by position, each naming its commands in position order", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(projectId).Build()

		cmdGroup1Builder := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1, cmd3)
		cmdGroup2Builder := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd1, cmd3, cmd2)

		cmdGroup1 := cmdGroup1Builder.Build()
		cmdGroup2 := cmdGroup2Builder.Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{
				commandinfrastructure.ToCommandModel(&cmd1),
				commandinfrastructure.ToCommandModel(&cmd2),
				commandinfrastructure.ToCommandModel(&cmd3),
			},
			[]infrastructure.CommandGroupModel{
				infrastructure.ToCommandGroupModel(&cmdGroup1),
				infrastructure.ToCommandGroupModel(&cmdGroup2),
			},
			[]infrastructure.CommandToCommandGroupModel{
				{CommandGroupId: cmdGroup1.Id, CommandId: cmd2.Id, Position: 0},
				{CommandGroupId: cmdGroup1.Id, CommandId: cmd1.Id, Position: 1},
				{CommandGroupId: cmdGroup1.Id, CommandId: cmd3.Id, Position: 2},
				{CommandGroupId: cmdGroup2.Id, CommandId: cmd1.Id, Position: 0},
				{CommandGroupId: cmdGroup2.Id, CommandId: cmd3.Id, Position: 1},
				{CommandGroupId: cmdGroup2.Id, CommandId: cmd2.Id, Position: 2},
			},
		)

		// Act
		result, err := helper.repo.GetAllWithCommandIds(projectId)

		// Assert
		expected := []domain.CommandGroupWithCommandIds{
			cmdGroup1Builder.BuildWithCommandIds(),
			cmdGroup2Builder.BuildWithCommandIds(),
		}

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Should return a command group that names no commands", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		emptyBuilder := test2.NewCommandGroupBuilder().WithName("Empty").WithProjectId(projectId).WithPosition(0)
		empty := emptyBuilder.Build()

		helper := newTestHelper(t, nil, []infrastructure.CommandGroupModel{infrastructure.ToCommandGroupModel(&empty)}, nil)

		// Act
		result, err := helper.repo.GetAllWithCommandIds(projectId)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, []domain.CommandGroupWithCommandIds{emptyBuilder.BuildWithCommandIds()}, result)
	})

	t.Run("Should read the command groups and the commands they name in a single query", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()

		cmdGroup1Builder := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1)
		cmdGroup1 := cmdGroup1Builder.Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{commandinfrastructure.ToCommandModel(&cmd1)},
			[]infrastructure.CommandGroupModel{infrastructure.ToCommandGroupModel(&cmdGroup1)},
			[]infrastructure.CommandToCommandGroupModel{{CommandGroupId: cmdGroup1.Id, CommandId: cmd1.Id, Position: 0}},
		)

		recorder := &sqlRecorder{Interface: gormlogger.Discard}
		repo := infrastructure.NewGormCommandGroupRepository(
			helper.gormDb.Session(&gorm.Session{Logger: recorder}),
			context.Background(),
		)

		// Act
		result, err := repo.GetAllWithCommandIds(projectId)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, []domain.CommandGroupWithCommandIds{cmdGroup1Builder.BuildWithCommandIds()}, result)
		assert.Len(t, recorder.statements, 1)
	})
}

func TestGormCommandGroupRepository_GetWithCommandIds(t *testing.T) {
	t.Run("Should return a command group by id naming its commands in position order", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()

		cmdGroup1Builder := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1)
		cmdGroup2Builder := test2.NewCommandGroupBuilder().WithName("Group 2").WithProjectId(projectId).WithPosition(1).WithCommands(cmd1)

		cmdGroup1 := cmdGroup1Builder.Build()
		cmdGroup2 := cmdGroup2Builder.Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{
				commandinfrastructure.ToCommandModel(&cmd1),
				commandinfrastructure.ToCommandModel(&cmd2),
			},
			[]infrastructure.CommandGroupModel{
				infrastructure.ToCommandGroupModel(&cmdGroup1),
				infrastructure.ToCommandGroupModel(&cmdGroup2),
			},
			[]infrastructure.CommandToCommandGroupModel{
				{CommandGroupId: cmdGroup1.Id, CommandId: cmd2.Id, Position: 0},
				{CommandGroupId: cmdGroup1.Id, CommandId: cmd1.Id, Position: 1},
				{CommandGroupId: cmdGroup2.Id, CommandId: cmd1.Id, Position: 0},
			},
		)

		// Act
		result, err := helper.repo.GetWithCommandIds(cmdGroup1.Id)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, cmdGroup1Builder.BuildWithCommandIds(), result)
	})

	t.Run("Should report a command group that does not exist as not found", func(t *testing.T) {
		// Arrange
		helper := newTestHelper(t, nil, nil, nil)

		// Act
		_, err := helper.repo.GetWithCommandIds("non-existent-id")

		// Assert
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)
	})
}

func TestGormCommandGroupRepository_GetAllContainingWithCommandIds(t *testing.T) {
	t.Run("Should return every command group holding the command, naming all of the commands they hold", func(t *testing.T) {
		// Arrange
		projectId := "project1"
		otherProjectId := "project2"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()
		cmd3 := test.NewCommandBuilder().WithName("Command 3").WithProjectId(otherProjectId).Build()

		holdingBuilder := test2.NewCommandGroupBuilder().WithName("Holding").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1)
		notHoldingBuilder := test2.NewCommandGroupBuilder().WithName("Not holding").WithProjectId(projectId).WithPosition(1).WithCommands(cmd2)
		holdingElsewhereBuilder := test2.NewCommandGroupBuilder().WithName("Holding elsewhere").WithProjectId(otherProjectId).WithPosition(0).WithCommands(cmd1, cmd3)

		holding := holdingBuilder.Build()
		notHolding := notHoldingBuilder.Build()
		holdingElsewhere := holdingElsewhereBuilder.Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{
				commandinfrastructure.ToCommandModel(&cmd1),
				commandinfrastructure.ToCommandModel(&cmd2),
				commandinfrastructure.ToCommandModel(&cmd3),
			},
			[]infrastructure.CommandGroupModel{
				infrastructure.ToCommandGroupModel(&holding),
				infrastructure.ToCommandGroupModel(&notHolding),
				infrastructure.ToCommandGroupModel(&holdingElsewhere),
			},
			[]infrastructure.CommandToCommandGroupModel{
				{CommandGroupId: holding.Id, CommandId: cmd2.Id, Position: 0},
				{CommandGroupId: holding.Id, CommandId: cmd1.Id, Position: 1},
				{CommandGroupId: notHolding.Id, CommandId: cmd2.Id, Position: 0},
				{CommandGroupId: holdingElsewhere.Id, CommandId: cmd1.Id, Position: 0},
				{CommandGroupId: holdingElsewhere.Id, CommandId: cmd3.Id, Position: 1},
			},
		)

		// Act
		result, err := helper.repo.GetAllContainingWithCommandIds(cmd1.Id)

		// Assert
		assert.Nil(t, err)
		assert.ElementsMatch(t, []domain.CommandGroupWithCommandIds{
			holdingBuilder.BuildWithCommandIds(),
			holdingElsewhereBuilder.BuildWithCommandIds(),
		}, result)
	})

	t.Run("Should still name a command whose own record is gone", func(t *testing.T) {
		// Arrange
		projectId := "project1"

		deleted := test.NewCommandBuilder().WithName("Deleted").WithProjectId(projectId).Build()
		survivor := test.NewCommandBuilder().WithName("Survivor").WithProjectId(projectId).Build()

		holdingBuilder := test2.NewCommandGroupBuilder().WithName("Holding").WithProjectId(projectId).WithPosition(0).WithCommands(deleted, survivor)
		holding := holdingBuilder.Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{commandinfrastructure.ToCommandModel(&survivor)},
			[]infrastructure.CommandGroupModel{infrastructure.ToCommandGroupModel(&holding)},
			[]infrastructure.CommandToCommandGroupModel{
				{CommandGroupId: holding.Id, CommandId: deleted.Id, Position: 0},
				{CommandGroupId: holding.Id, CommandId: survivor.Id, Position: 1},
			},
		)

		// Act
		result, err := helper.repo.GetAllContainingWithCommandIds(deleted.Id)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, []domain.CommandGroupWithCommandIds{holdingBuilder.BuildWithCommandIds()}, result)
	})

	t.Run("Should return nothing when no command group holds the command", func(t *testing.T) {
		// Arrange
		helper := newTestHelper(t, nil, nil, nil)

		// Act
		result, err := helper.repo.GetAllContainingWithCommandIds("unheld-command")

		// Assert
		assert.Nil(t, err)
		assert.Empty(t, result)
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
		assert.Equal(t, cmdGroup1, result)
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
		assert.Equal(t, updatedGroup, result)
	})
}

func TestGormCommandGroupRepository_UpdateWithCommandIds(t *testing.T) {
	t.Run("Should update an existing command group and the commands it names", func(t *testing.T) {
		projectId := "project1"

		cmd1 := test.NewCommandBuilder().WithName("Command 1").WithProjectId(projectId).Build()
		cmd2 := test.NewCommandBuilder().WithName("Command 2").WithProjectId(projectId).Build()

		commandModels := []commandinfrastructure.CommandModel{
			commandinfrastructure.ToCommandModel(&cmd1),
			commandinfrastructure.ToCommandModel(&cmd2),
		}

		cmdGroup1Builder := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(cmd1, cmd2)
		cmdGroup1 := cmdGroup1Builder.Build()

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
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
			[]infrastructure.CommandGroupModel{infrastructure.ToCommandGroupModel(&cmdGroup1)},
			commandToCommandGroupModels,
		)

		updatedGroup := cmdGroup1Builder.WithName("Updated Group 1").WithCommands(cmd2).BuildWithCommandIds()

		err := helper.repo.UpdateWithCommandIds(&updatedGroup)
		assert.Nil(t, err)

		result, err := helper.repo.GetWithCommandIds(updatedGroup.Id)
		assert.Nil(t, err)
		assert.Equal(t, updatedGroup, result)
	})

	t.Run("Should write a command group that names a command whose own record is gone", func(t *testing.T) {
		projectId := "project1"

		kept := test.NewCommandBuilder().WithName("Kept").WithProjectId(projectId).Build()

		cmdGroup1Builder := test2.NewCommandGroupBuilder().WithName("Group 1").WithProjectId(projectId).WithPosition(0).WithCommands(kept)
		cmdGroup1 := cmdGroup1Builder.Build()

		helper := newTestHelper(
			t,
			[]commandinfrastructure.CommandModel{commandinfrastructure.ToCommandModel(&kept)},
			[]infrastructure.CommandGroupModel{infrastructure.ToCommandGroupModel(&cmdGroup1)},
			[]infrastructure.CommandToCommandGroupModel{
				{
					CommandGroupId: cmdGroup1.Id,
					CommandId:      kept.Id,
					Position:       0,
				},
			},
		)

		namingAGoneCommand := domain.CommandGroupWithCommandIds{
			Id:         cmdGroup1.Id,
			ProjectId:  projectId,
			Name:       "Group 1",
			CommandIds: []string{kept.Id, "deleted-command"},
			Position:   0,
		}

		err := helper.repo.UpdateWithCommandIds(&namingAGoneCommand)
		assert.Nil(t, err)

		result, err := helper.repo.GetWithCommandIds(cmdGroup1.Id)
		assert.Nil(t, err)
		assert.Equal(t, namingAGoneCommand, result)
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

		_, err = helper.repo.Get(cmdGroup1.Id)
		assert.ErrorIs(t, err, domainerrors.ErrNotFound)

		existingRelations, err := gorm.G[infrastructure.CommandToCommandGroupModel](helper.gormDb).Where("command_group_id = ?", cmdGroup1.Id).Find(context.Background())
		assert.Nil(t, err)
		assert.Len(t, existingRelations, 0)

		existingCommands, err := gorm.G[commandinfrastructure.CommandModel](helper.gormDb).Where("id = ?", cmd1.Id).Find(context.Background())
		assert.Nil(t, err)
		assert.Len(t, existingCommands, 1)
	})
	t.Run("Should leave the positions of the remaining command groups untouched", func(t *testing.T) {
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

		// Closing the gap is the ordering module's job, not the repository's
		assert.Equal(t, group1Model.Position, resultGroup1.Position)
		assert.Equal(t, group3Model.Position, resultGroup3.Position)
	})
}

func TestGormCommandGroupRepository_GetAllContaining(t *testing.T) {
	t.Run("Should return every command group holding the command, with all of their commands", func(t *testing.T) {
		// Arrange
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

		holding := test2.NewCommandGroupBuilder().WithName("Holding").WithProjectId(projectId).WithPosition(0).WithCommands(cmd2, cmd1).Build()
		notHolding := test2.NewCommandGroupBuilder().WithName("Not holding").WithProjectId(projectId).WithPosition(1).WithCommands(cmd2).Build()
		holdingElsewhere := test2.NewCommandGroupBuilder().WithName("Holding elsewhere").WithProjectId(otherProjectId).WithPosition(0).WithCommands(cmd1, cmd3).Build()

		commandToCommandGroupModels := []infrastructure.CommandToCommandGroupModel{
			{CommandGroupId: holding.Id, CommandId: cmd2.Id, Position: 0},
			{CommandGroupId: holding.Id, CommandId: cmd1.Id, Position: 1},
			{CommandGroupId: notHolding.Id, CommandId: cmd2.Id, Position: 0},
			{CommandGroupId: holdingElsewhere.Id, CommandId: cmd1.Id, Position: 0},
			{CommandGroupId: holdingElsewhere.Id, CommandId: cmd3.Id, Position: 1},
		}

		helper := newTestHelper(
			t,
			commandModels,
			[]infrastructure.CommandGroupModel{
				infrastructure.ToCommandGroupModel(&holding),
				infrastructure.ToCommandGroupModel(&notHolding),
				infrastructure.ToCommandGroupModel(&holdingElsewhere),
			},
			commandToCommandGroupModels,
		)

		// Act
		result, err := helper.repo.GetAllContaining(cmd1.Id)

		// Assert
		assert.Nil(t, err)
		assert.ElementsMatch(t, []domain.CommandGroup{holding, holdingElsewhere}, result)
	})

	t.Run("Should return nothing when no command group holds the command", func(t *testing.T) {
		// Arrange
		helper := newTestHelper(t, nil, nil, nil)

		// Act
		result, err := helper.repo.GetAllContaining("unheld-command")

		// Assert
		assert.Nil(t, err)
		assert.Empty(t, result)
	})
}

func arrange(
	t *testing.T,
	preloadedCommandModels []commandinfrastructure.CommandModel,
	preloadedCommandGroupModels []infrastructure.CommandGroupModel,
	preloadedCommandToCommandGroupModels []infrastructure.CommandToCommandGroupModel,
) (repo *infrastructure.GormCommandGroupRepository, gormDb *gorm.DB) {
	t.Helper()

	ctx := context.Background()
	gormDb = testdb.New(t)

	for _, m := range preloadedCommandModels {
		err := gorm.G[commandinfrastructure.CommandModel](gormDb).Create(ctx, &m)
		if err != nil {
			t.Fatalf("failed to preload command: %v", err)
		}
	}

	for _, m := range preloadedCommandGroupModels {
		err := gorm.G[infrastructure.CommandGroupModel](gormDb).Create(ctx, &m)
		if err != nil {
			t.Fatalf("failed to preload command group: %v", err)
		}
	}

	for _, m := range preloadedCommandToCommandGroupModels {
		err := gorm.G[infrastructure.CommandToCommandGroupModel](gormDb).Create(ctx, &m)
		if err != nil {
			t.Fatalf("failed to preload command group relation: %v", err)
		}
	}

	repo = infrastructure.NewGormCommandGroupRepository(
		gormDb,
		ctx,
	)

	return
}
