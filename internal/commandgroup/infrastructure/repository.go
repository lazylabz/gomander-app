package infrastructure

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"gomander/internal/commandgroup/domain"
	"gomander/internal/domainerrors"
)

type GormCommandGroupRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormCommandGroupRepository(db *gorm.DB, ctx context.Context) *GormCommandGroupRepository {
	return &GormCommandGroupRepository{
		db:  db,
		ctx: ctx,
	}
}

// commandGroupQuery reads Command Groups and their Commands in one go: a Group
// with three Commands arrives as three rows, and a Group with none as a single
// row whose command_id is null. Ordering by the join table's position is what a
// preload cannot do, and is why this is written out rather than built by GORM.
const commandGroupQuery = `
SELECT
	command_group.id,
	command_group.project_id,
	command_group.name,
	command_group.position,
	command.id AS command_id,
	COALESCE(command.project_id, '') AS command_project_id,
	COALESCE(command.name, '') AS command_name,
	COALESCE(command.command, '') AS command_command,
	COALESCE(command.working_directory, '') AS command_working_directory,
	COALESCE(command.position, 0) AS command_position,
	COALESCE(command.link, '') AS command_link,
	COALESCE(command.error_patterns, '') AS command_error_patterns
FROM command_group
LEFT JOIN command_group_command ON command_group_command.command_group_id = command_group.id
LEFT JOIN command ON command.id = command_group_command.command_id
WHERE %s
ORDER BY command_group.position ASC, command_group.id ASC, command_group_command.position ASC
`

func (r GormCommandGroupRepository) GetAll(projectId string) ([]domain.CommandGroup, error) {
	return r.find("command_group.project_id = ?", projectId)
}

func (r GormCommandGroupRepository) GetAllContaining(commandId string) ([]domain.CommandGroup, error) {
	return r.find(
		"command_group.id IN (SELECT command_group_id FROM command_group_command WHERE command_id = ?)",
		commandId,
	)
}

func (r GormCommandGroupRepository) Get(id string) (domain.CommandGroup, error) {
	commandGroups, err := r.find("command_group.id = ?", id)
	if err != nil {
		return domain.CommandGroup{}, err
	}

	if len(commandGroups) == 0 {
		return domain.CommandGroup{}, domainerrors.NotFound("command group", id)
	}

	return commandGroups[0], nil
}

func (r GormCommandGroupRepository) find(condition string, args ...any) ([]domain.CommandGroup, error) {
	var rows []commandGroupRow

	err := r.db.WithContext(r.ctx).Raw(fmt.Sprintf(commandGroupQuery, condition), args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return ToDomainCommandGroups(rows), nil
}

func (r GormCommandGroupRepository) Create(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := gorm.G[CommandGroupModel](tx).Create(r.ctx, &commandGroupModel)
		if err != nil {
			return err
		}

		return r.writeCommandPlacements(tx, commandGroup)
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

		return r.writeCommandPlacements(tx, commandGroup)
	})

	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandGroupRepository) writeCommandPlacements(tx *gorm.DB, commandGroup *domain.CommandGroup) error {
	for _, model := range ToCommandToCommandGroupModels(commandGroup) {
		if err := gorm.G[CommandToCommandGroupModel](tx).Create(r.ctx, &model); err != nil {
			return err
		}
	}

	return nil
}

func (r GormCommandGroupRepository) Delete(commandGroupId string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[CommandGroupModel](tx).Where("id = ?", commandGroupId).Delete(r.ctx)
		if err != nil {
			return err
		}

		_, err = gorm.G[CommandToCommandGroupModel](tx).
			Where("command_group_id = ?", commandGroupId).
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

func (r GormCommandGroupRepository) Atomically(change func(domain.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return change(NewGormCommandGroupRepository(tx, r.ctx))
	})
}
