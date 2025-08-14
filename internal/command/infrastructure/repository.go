package commandinfrastructure

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

func (g GormCommandRepository) GetCommand(commandId string) (*domain.Command, error) {
	cmd, err := gorm.G[CommandModel](g.db).Where("id = ?", commandId).First(g.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	command := ToDomainCommand(cmd)

	return &command, nil
}

func (g GormCommandRepository) GetCommands(projectId string) ([]domain.Command, error) {
	cmds, err := gorm.G[CommandModel](g.db).Where("project_id = ?", projectId).Order("position").Find(g.ctx)
	if err != nil {
		return nil, err
	}

	return array.Map(cmds, ToDomainCommand), nil
}

func (g GormCommandRepository) SaveCommand(command *domain.Command) error {
	commandModel := ToCommandModel(command)

	err := gorm.G[CommandModel](g.db).Create(g.ctx, &commandModel)
	if err != nil {
		return err
	}

	return nil
}

func (g GormCommandRepository) EditCommand(command *domain.Command) error {
	commandModel := ToCommandModel(command)

	_, err := gorm.G[CommandModel](g.db).Where("id = ?", commandModel.Id).Updates(g.ctx, commandModel)
	if err != nil {
		return err
	}

	return nil
}

func (g GormCommandRepository) DeleteCommand(commandId string) error {
	_, err := gorm.G[CommandModel](g.db).Where("id = ?", commandId).Delete(g.ctx)
	if err != nil {
		return err
	}

	return nil
}
