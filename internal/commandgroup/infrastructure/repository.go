package infrastructure

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
)

type GormCommandGroupRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormCommandGroupRepository(db *gorm.DB, ctx context.Context) *GormCommandGroupRepository {
	err := db.SetupJoinTable(&CommandGroupModel{}, "Commands", &CommandToCommandGroupModel{})
	if err != nil {
		panic(err)
	}
	return &GormCommandGroupRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r GormCommandGroupRepository) GetCommandGroups(projectId string) ([]domain.CommandGroup, error) {
	var cgModels []CommandGroupModel
	err := r.db.Where("project_id = ?", projectId).
		Order("position ASC").
		Preload("Commands", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN command_group_command ON command_group_command.command_id = command.id").
				Order("command_group_command.position")
		}).
		Find(&cgModels).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return array.Map(cgModels, func(cgModel CommandGroupModel) domain.CommandGroup {
		return *ToDomainCommandGroup(cgModel)
	}), err
}

func (r GormCommandGroupRepository) GetCommandGroupById(id string) (*domain.CommandGroup, error) {
	var cgModel CommandGroupModel
	err := r.db.Where("id = ?", id).
		Preload("Commands", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN command_group_command ON command_group_command.command_id = command.id").
				Order("command_group_command.position")
		}).
		First(&cgModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return ToDomainCommandGroup(cgModel), err
}

func (r GormCommandGroupRepository) CreateCommandGroup(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := gorm.G[CommandGroupModel](tx).Create(r.ctx, &commandGroupModel)
		if err != nil {
			return err
		}

		if len(commandGroup.Commands) > 0 {
			// Create command associations
			for i, cmd := range commandGroup.Commands {
				cmdToGroup := CommandToCommandGroupModel{
					CommandId:      cmd.Id,
					CommandGroupId: commandGroupModel.Id,
					Position:       i,
				}
				err = gorm.G[CommandToCommandGroupModel](tx).Create(r.ctx, &cmdToGroup)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandGroupRepository) UpdateCommandGroup(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Update the command group data
		_, err := gorm.G[CommandGroupModel](tx).Where("id = ?", commandGroupModel.Id).Updates(r.ctx, commandGroupModel)
		if err != nil {
			return err
		}

		// Delete existing command associations
		_, err = gorm.G[CommandToCommandGroupModel](tx).
			Where("command_group_id = ?", commandGroupModel.Id).
			Delete(r.ctx)
		if err != nil {
			return err
		}

		// Create new command associations
		for i, cmd := range commandGroup.Commands {
			cmdToGroup := CommandToCommandGroupModel{
				CommandId:      cmd.Id,
				CommandGroupId: commandGroupModel.Id,
				Position:       i,
			}
			err = gorm.G[CommandToCommandGroupModel](tx).Create(r.ctx, &cmdToGroup)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandGroupRepository) DeleteCommandGroup(commandGroupId string) error {
	existingGroup, err := gorm.G[CommandGroupModel](r.db).Where("id = ?", commandGroupId).First(r.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // If the command group does not exist, nothing to delete
		}
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Delete the command group
		_, err = gorm.G[CommandGroupModel](tx).Where("id = ?", commandGroupId).Delete(r.ctx)
		if err != nil {
			return err
		}

		// Delete all command associations for this command group
		_, err = gorm.G[CommandToCommandGroupModel](tx).
			Where("command_group_id = ?", commandGroupId).
			Delete(r.ctx)
		if err != nil {
			return err
		}

		// Decrease the position of all command groups with a higher position
		_, err = gorm.G[CommandGroupModel](tx).
			Where("project_id = ? AND position > ?", existingGroup.ProjectId, existingGroup.Position).
			Update(r.ctx, "position", gorm.Expr("position - 1"))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
