package infrastructure

import "gomander/internal/config/domain"

func ToDomainConfig(model *ConfigModel, paths []EnvironmentPathModel) *domain.Config {
	if model == nil {
		return nil
	}

	config := &domain.Config{
		LastOpenedProjectId: model.LastOpenedProjectId,
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
		LastOpenedProjectId: config.LastOpenedProjectId,
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
