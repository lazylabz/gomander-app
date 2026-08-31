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

// commandGroupIdentityQuery names the Commands a Command Group holds instead
// of hydrating them. It joins the membership table alone, so a Command Group
// still names a Command whose row is gone. Ordering by the join table's
// position is what a preload cannot do, and is why this is written out rather
// than built by GORM.
const commandGroupIdentityQuery = `
SELECT
	command_group.id,
	command_group.project_id,
	command_group.name,
	command_group.position,
	command_group_command.command_id
FROM command_group
LEFT JOIN command_group_command ON command_group_command.command_group_id = command_group.id
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
	var rows []commandGroupIdentityRow

	err := r.db.WithContext(r.ctx).Raw(fmt.Sprintf(commandGroupIdentityQuery, condition), args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return ToDomainCommandGroups(rows), nil
}

func (r GormCommandGroupRepository) Create(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	return r.db.Transaction(func(tx *gorm.DB) error {
		err := gorm.G[CommandGroupModel](tx).Create(r.ctx, &commandGroupModel)
		if err != nil {
			return err
		}

		return r.writeCommandPlacements(tx, ToCommandToCommandGroupModels(commandGroup))
	})
}

func (r GormCommandGroupRepository) Update(commandGroup *domain.CommandGroup) error {
	return r.rewrite(ToCommandGroupModel(commandGroup), ToCommandToCommandGroupModels(commandGroup))
}

// rewrite replaces the Command Group and the membership rows it owns, so a
// Command the Group no longer holds leaves no row behind.
func (r GormCommandGroupRepository) rewrite(
	commandGroupModel CommandGroupModel,
	placements []CommandToCommandGroupModel,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[CommandGroupModel](tx).Where("id = ?", commandGroupModel.Id).Select("*").Updates(r.ctx, commandGroupModel)
		if err != nil {
			return err
		}

		_, err = gorm.G[CommandToCommandGroupModel](tx).
			Where("command_group_id = ?", commandGroupModel.Id).
			Delete(r.ctx)
		if err != nil {
			return err
		}

		return r.writeCommandPlacements(tx, placements)
	})
}

func (r GormCommandGroupRepository) writeCommandPlacements(tx *gorm.DB, placements []CommandToCommandGroupModel) error {
	for _, model := range placements {
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
