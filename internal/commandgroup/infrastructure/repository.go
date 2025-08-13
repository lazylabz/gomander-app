package infrastructure

import (
	"context"
	"errors"
	"slices"
	"sort"

	"gorm.io/gorm"

	"gomander/internal/command/infrastructure"
	"gomander/internal/commandgroup/domain"
	"gomander/internal/helpers/array"
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

func (r GormCommandGroupRepository) GetCommandGroups(projectId string) ([]*domain.CommandGroup, error) {
	commandGroupModels, err := gorm.G[CommandGroupModel](r.db).
		Where("project_id = ?", projectId).
		Order("position ASC").
		Find(r.ctx)
	if err != nil {
		return nil, err
	}

	cgs := make([]*domain.CommandGroup, len(commandGroupModels))

	for i, model := range commandGroupModels {
		cg, err := r.GetCommandGroupById(model.Id)
		if err != nil {
			return nil, err
		}
		if cg == nil {
			return nil, errors.New("command group not found")
		}
		cgs[i] = cg
	}

	return cgs, nil
}

func (r GormCommandGroupRepository) GetCommandGroupById(id string) (*domain.CommandGroup, error) {
	commandGroupModel, err := gorm.G[CommandGroupModel](r.db).
		Where("id = ?", id).
		First(r.ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	relations, err := gorm.G[CommandToCommandGroupModel](r.db).
		Where("command_group_id = ?", id).
		Order("position ASC").
		Find(r.ctx)

	commandIds := array.Map(relations, func(r CommandToCommandGroupModel) string { return r.CommandId })

	commandModels, err := gorm.G[infrastructure.CommandModel](r.db).
		Where("id IN ?", commandIds).
		Find(r.ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(commandModels, func(i, j int) bool {
		return slices.Index(commandIds, commandModels[i].Id) < slices.Index(commandIds, commandModels[j].Id)
	})

	commandGroup := ToDomainCommandGroup(commandGroupModel)
	commandGroup.Commands = array.Map(commandModels, infrastructure.ToDomainCommand)

	return commandGroup, nil
}

func (r GormCommandGroupRepository) CreateCommandGroup(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	err := gorm.G[CommandGroupModel](r.db).Create(r.ctx, &commandGroupModel)

	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandGroupRepository) UpdateCommandGroup(commandGroup *domain.CommandGroup) error {
	commandGroupModel := ToCommandGroupModel(commandGroup)

	_, err := gorm.G[CommandGroupModel](r.db).Where("id = ?", commandGroupModel.Id).Updates(r.ctx, commandGroupModel)
	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandGroupRepository) DeleteCommandGroup(commandGroupId string) error {
	_, err := gorm.G[CommandGroupModel](r.db).Where("id = ?", commandGroupId).Delete(r.ctx)
	if err != nil {
		return err
	}

	return nil
}
