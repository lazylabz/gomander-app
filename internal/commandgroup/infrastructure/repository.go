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

func (r GormCommandGroupRepository) GetAll(projectId string) ([]domain.CommandGroup, error) {
	var ids []string
	err := r.db.Model(&CommandGroupModel{}).
		Where("project_id = ?", projectId).
		Order("position ASC").
		Pluck("id", &ids).Error

	if err != nil {
		return nil, err
	}

	commandGroups := make([]domain.CommandGroup, 0)
	for _, id := range ids {
		cg, err := r.Get(id)
		if err != nil {
			return nil, err
		}
		if cg != nil {
			commandGroups = append(commandGroups, *cg)
		}
	}
	return commandGroups, nil
}

func (r GormCommandGroupRepository) Get(id string) (*domain.CommandGroup, error) {
	var cgModel CommandGroupModel
	err := r.db.Where("id = ?", id).
		Preload("Commands", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN command_group_command ON command_group_command.command_id = command.id AND command_group_command.command_group_id = ?", id).
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

func (r GormCommandGroupRepository) Create(commandGroup *domain.CommandGroup) error {
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

func (r GormCommandGroupRepository) Update(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Update the command group data
		_, err := gorm.G[CommandGroupModel](tx).Where("id = ?", commandGroupModel.Id).Select("*").Updates(r.ctx, commandGroupModel)
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

func (r GormCommandGroupRepository) Delete(commandGroupId string) error {
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

func (r GormCommandGroupRepository) RemoveCommandFromCommandGroups(commandId string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Find all command group associations for the command
		relations, err := gorm.G[CommandToCommandGroupModel](tx).
			Where("command_id = ?", commandId).
			Find(r.ctx)

		if err != nil {
			return err
		}

		// Update positions of command groups after the removed command
		for _, relation := range relations {
			_, err = gorm.G[CommandToCommandGroupModel](tx).
				Where("command_group_id = ? AND position > ?", relation.CommandGroupId, relation.Position).
				Update(r.ctx, "position", gorm.Expr("position - 1"))
			if err != nil {
				return err
			}
		}

		// Delete the command from all command groups
		_, err = gorm.G[CommandToCommandGroupModel](tx).
			Where("command_id = ?", commandId).
			Delete(r.ctx)

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

func (r GormCommandGroupRepository) DeleteEmpty() ([]string, error) {
	query := "id NOT IN (SELECT DISTINCT command_group_id FROM command_group_command)"

	entriesToDelete, err := gorm.G[CommandGroupModel](r.db).
		Where(query).
		Find(r.ctx)

	if err != nil {
		return nil, err
	}

	_, err = gorm.G[CommandGroupModel](r.db).
		Where(query).
		Delete(r.ctx)

	if err != nil {
		return nil, err
	}

	return array.Map(entriesToDelete, func(entry CommandGroupModel) string { return entry.Id }), nil
}

func (r GormCommandGroupRepository) DeleteAll(projectId string) ([]string, error) {
	var commandGroupIds []string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		commandGroups, err := gorm.G[CommandGroupModel](tx).Where("project_id = ?", projectId).Find(r.ctx)
		if err != nil {
			return err
		}

		commandGroupIds = array.Map(commandGroups, func(commandGroup CommandGroupModel) string {
			return commandGroup.Id
		})

		_, err = gorm.G[CommandGroupModel](tx).Where("project_id = ?", projectId).Delete(r.ctx)
		if err != nil {
			return err
		}

		_, err = gorm.G[CommandToCommandGroupModel](tx).Where("command_group_id IN ?", commandGroupIds).Delete(r.ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return commandGroupIds, nil
}
