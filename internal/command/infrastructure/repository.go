package infrastructure

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gomander/internal/command/domain"
	"gomander/internal/helpers/array"
)

type GormCommandRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormCommandRepository(db *gorm.DB, ctx context.Context) *GormCommandRepository {
	return &GormCommandRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r GormCommandRepository) Get(commandId string) (*domain.Command, error) {
	cmd, err := gorm.G[CommandModel](r.db).Where("id = ?", commandId).First(r.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	command := ToDomainCommand(cmd)

	return &command, nil
}

func (r GormCommandRepository) GetAll(projectId string) ([]domain.Command, error) {
	cmds, err := gorm.G[CommandModel](r.db).Where("project_id = ?", projectId).Order("position").Find(r.ctx)
	if err != nil {
		return nil, err
	}

	return array.Map(cmds, ToDomainCommand), nil
}

func (r GormCommandRepository) Create(command *domain.Command) error {
	commandModel := ToCommandModel(command)

	err := gorm.G[CommandModel](r.db).Create(r.ctx, &commandModel)
	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandRepository) Update(command *domain.Command) error {
	commandModel := ToCommandModel(command)

	_, err := gorm.G[CommandModel](r.db).Where("id = ?", commandModel.Id).Select("*").Updates(r.ctx, commandModel)
	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandRepository) Delete(commandId string) error {
	_, err := gorm.G[CommandModel](r.db).Where("id = ?", commandId).Delete(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandRepository) DeleteAll(projectId string) error {
	_, err := gorm.G[CommandModel](r.db).Where("project_id = ?", projectId).Delete(r.ctx)
	if err != nil {
		return err
	}
	return nil
}
