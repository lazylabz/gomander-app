package infrastructure

import (
	"context"

	"gorm.io/gorm"

	commandinfrastructure "gomander/internal/command/infrastructure"
	commandgroupinfrastructure "gomander/internal/commandgroup/infrastructure"
	projectinfrastructure "gomander/internal/project/infrastructure"
	"gomander/internal/unitofwork"
)

type GormUnitOfWork struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormUnitOfWork(db *gorm.DB, ctx context.Context) *GormUnitOfWork {
	return &GormUnitOfWork{
		db:  db,
		ctx: ctx,
	}
}

func (u GormUnitOfWork) Do(change func(unitofwork.Repositories) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return change(unitofwork.Repositories{
			Projects:      projectinfrastructure.NewGormProjectRepository(tx, u.ctx),
			Commands:      commandinfrastructure.NewGormCommandRepository(tx, u.ctx),
			CommandGroups: commandgroupinfrastructure.NewGormCommandGroupRepository(tx, u.ctx),
		})
	})
}
