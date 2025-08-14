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
