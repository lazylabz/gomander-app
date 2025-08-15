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

	_, err := gorm.G[CommandModel](r.db).Where("id = ?", commandModel.Id).Updates(r.ctx, commandModel)
	if err != nil {
		return err
	}

	return nil
}

func (r GormCommandRepository) Delete(commandId string) error {
	originalCommand, err := r.Get(commandId)
	if err != nil {
		return err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[CommandModel](r.db).Where("id = ?", commandId).Delete(r.ctx)
		if err != nil {
			return err
		}

		// Decrease position of all commands with position greater than the deleted command's position
		if originalCommand != nil {
			_, err = gorm.G[CommandModel](tx).
				Where("project_id = ? AND position > ?", originalCommand.ProjectId, originalCommand.Position).
				Update(r.ctx, "position", gorm.Expr("position - 1"))

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
