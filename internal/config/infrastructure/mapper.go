package infrastructure

import "gomander/internal/config/domain"

func ToDomainConfig(model *ConfigModel, paths []EnvironmentPathModel) *domain.Config {
	if model == nil {
		return nil
	}

	config := &domain.Config{
		LastOpenedProjectId:  model.LastOpenedProjectId,
		EnvironmentPaths:     make([]domain.EnvironmentPath, 0),
		Locale:               model.Locale,
		LogsFontSize:         model.LogsFontSize,
		ZoomLevel:            domain.ZoomLevel(model.ZoomLevel),
		NotificationsEnabled: model.NotificationsEnabled,
	}

	for _, pathModel := range paths {
		config.EnvironmentPaths = append(config.EnvironmentPaths, domain.EnvironmentPath{
			Id:   pathModel.Id,
			Path: pathModel.Path,
		})
	}

	return config
}

func ToModelConfig(config *domain.Config) (*ConfigModel, []EnvironmentPathModel) {
	if config == nil {
		return nil, nil
	}

	model := &ConfigModel{
		Id:                   1,
		LastOpenedProjectId:  config.LastOpenedProjectId,
		Locale:               config.Locale,
		LogsFontSize:         config.LogsFontSize,
		ZoomLevel:            float64(config.ZoomLevel),
		NotificationsEnabled: config.NotificationsEnabled,
	}

	var pathModels []EnvironmentPathModel
	for _, path := range config.EnvironmentPaths {
		pathModels = append(pathModels, EnvironmentPathModel{
			Id:   path.Id,
			Path: path.Path,
		})
	}

	return model, pathModels
}
