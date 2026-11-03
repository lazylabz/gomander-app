package infrastructure

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gomander/internal/config/domain"
)

type GormConfigRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewGormConfigRepository(db *gorm.DB, ctx context.Context) *GormConfigRepository {
	return &GormConfigRepository{
		db:  db,
		ctx: ctx,
	}
}

func (g GormConfigRepository) GetOrCreate() (*domain.Config, error) {
	var configModel ConfigModel

	configModel, err := gorm.G[ConfigModel](g.db).First(g.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no config exists, create a new one
			configModel = ConfigModel{
				Id:                  1,
				LastOpenedProjectId: "",
			}
			if err := gorm.G[ConfigModel](g.db).Create(g.ctx, &configModel); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	pathModels, err := gorm.G[EnvironmentPathModel](g.db).Find(g.ctx)
	if err != nil {
		return nil, err
	}

	return ToDomainConfig(&configModel, pathModels), nil
}

// One transaction, so nothing reads the environment paths in between the
// delete and the reinsert.
func (g GormConfigRepository) Update(config *domain.Config) error {
	configModel, pathModels := ToModelConfig(config)
	if configModel == nil {
		return errors.New("config cannot be nil")
	}

	return g.db.WithContext(g.ctx).Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[ConfigModel](tx).Where("id = ?", 1).Select("*").Omit("zoom_level").Updates(g.ctx, *configModel)
		if err != nil {
			return err
		}

		_, err = gorm.G[EnvironmentPathModel](tx).Where("id NOT NULL").Delete(g.ctx)
		if err != nil {
			return err
		}

		for _, pathModel := range pathModels {
			if err := gorm.G[EnvironmentPathModel](tx).Create(g.ctx, &pathModel); err != nil {
				return err
			}
		}

		return nil
	})
}

func (g GormConfigRepository) UpdateZoomLevel(level domain.ZoomLevel) error {
	_, err := gorm.G[ConfigModel](g.db).Where("id = ?", 1).Update(g.ctx, "zoom_level", float64(level))

	return err
}
