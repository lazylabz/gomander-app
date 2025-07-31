package app

import (
	"context"

	"gomander/internal/config"
	"gomander/internal/event"
	"gomander/internal/logger"
	"gomander/internal/project"
)

// App struct
type App struct {
	ctx context.Context

	logger        *logger.Logger
	eventEmitter  *event.EventEmitter
	commandRunner *project.Runner

	userConfig      *config.UserConfig
	selectedProject *project.Project
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) persistSelectedProjectConfig() error {
	err := project.SaveProject(&project.Project{
		Id:                   a.selectedProject.Id,
		Name:                 a.selectedProject.Name,
		BaseWorkingDirectory: a.selectedProject.BaseWorkingDirectory,
		Commands:             a.selectedProject.Commands,
		CommandGroups:        a.selectedProject.CommandGroups,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *App) persistUserConfig() error {
	err := config.SaveUserConfig(&config.UserConfig{
		EnvironmentPaths:    a.userConfig.EnvironmentPaths,
		LastOpenedProjectId: a.userConfig.LastOpenedProjectId,
	})
	if err != nil {
		return err
	}

	return nil
}
