package infrastructure

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gomander/internal/domainerrors"
	"gomander/internal/helpers/array"
	"gomander/internal/project/domain"
)

type GormProjectRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormProjectRepository(db *gorm.DB, ctx context.Context) *GormProjectRepository {
	return &GormProjectRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r GormProjectRepository) GetAll() ([]domain.Project, error) {
	projectModels, err := gorm.G[ProjectModel](r.db).Find(r.ctx)
	if err != nil {
		return nil, err
	}
	return array.Map(projectModels, ToDomainProject), nil
}

func (r GormProjectRepository) Get(id string) (domain.Project, error) {
	project, found, err := r.Find(id)
	if err != nil {
		return domain.Project{}, err
	}
	if !found {
		return domain.Project{}, domainerrors.NotFound("project", id)
	}

	return project, nil
}

func (r GormProjectRepository) Find(id string) (domain.Project, bool, error) {
	project, err := gorm.G[ProjectModel](r.db).Where("id = ?", id).First(r.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Project{}, false, nil
		}
		return domain.Project{}, false, err
	}

	return ToDomainProject(project), true, nil
}

func (r GormProjectRepository) Create(project domain.Project) error {
	projectModel := ToProjectModel(project)
	err := gorm.G[ProjectModel](r.db).Create(r.ctx, &projectModel)
	if err != nil {
		return err
	}
	return nil
}

func (r GormProjectRepository) Update(project domain.Project) error {
	projectModel := ToProjectModel(project)
	_, err := gorm.G[ProjectModel](r.db).Where("id = ?", project.Id).Select("*").Updates(r.ctx, projectModel)
	if err != nil {
		return err
	}
	return nil
}

func (r GormProjectRepository) Delete(id string) error {
	_, err := gorm.G[ProjectModel](r.db).Where("id = ?", id).Delete(r.ctx)
	if err != nil {
		return err
	}
	return nil
}
