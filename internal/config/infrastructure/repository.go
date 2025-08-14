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

func (g GormConfigRepository) GetOrCreateConfig() (*domain.Config, error) {
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

func (g GormConfigRepository) SaveConfig(config *domain.Config) error {
	configModel, pathModels := ToModelConfig(config)
	if configModel == nil {
		return errors.New("config cannot be nil")
	}

	_, err := gorm.G[ConfigModel](g.db).Where("id = ?", 1).Updates(g.ctx, *configModel)
	if err != nil {
		return err
	}

	_, err = gorm.G[EnvironmentPathModel](g.db).Where("id NOT NULL").Delete(g.ctx)
	if err != nil {
		return err
	}

	for _, pathModel := range pathModels {
		if err := gorm.G[EnvironmentPathModel](g.db).Create(g.ctx, &pathModel); err != nil {
			return err
		}
	}

	return nil
}
