package domain

import (
	configdomain "gomander/internal/config/domain"
)

type EnvironmentPath struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type Config struct {
	LastOpenedProjectId string            `json:"lastOpenedProjectId"`
	EnvironmentPaths    []EnvironmentPath `json:"environmentPaths"`
	Locale              string            `json:"locale"`
	LogsFontSize        int               `json:"logsFontSize"`
	ZoomLevel           float64           `json:"zoomLevel"`
}

func FromConfig(config configdomain.Config) Config {
	return Config{
		LastOpenedProjectId: config.LastOpenedProjectId,
		EnvironmentPaths:    mapSlice(config.EnvironmentPaths, fromEnvironmentPath),
		Locale:              config.Locale,
		LogsFontSize:        config.LogsFontSize,
		ZoomLevel:           float64(config.ZoomLevel),
	}
}

func (c Config) ToDomain() configdomain.Config {
	return configdomain.Config{
		LastOpenedProjectId: c.LastOpenedProjectId,
		EnvironmentPaths:    mapSlice(c.EnvironmentPaths, EnvironmentPath.ToDomain),
		Locale:              c.Locale,
		LogsFontSize:        c.LogsFontSize,
		ZoomLevel:           configdomain.ZoomLevel(c.ZoomLevel),
	}
}

func fromEnvironmentPath(environmentPath configdomain.EnvironmentPath) EnvironmentPath {
	return EnvironmentPath{
		Id:   environmentPath.Id,
		Path: environmentPath.Path,
	}
}

func (e EnvironmentPath) ToDomain() configdomain.EnvironmentPath {
	return configdomain.EnvironmentPath{
		Id:   e.Id,
		Path: e.Path,
	}
}
