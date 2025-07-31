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

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	l := logger.NewLogger(ctx)
	ee := event.NewEventEmitter(ctx)

	a.logger = l
	a.eventEmitter = ee
	a.commandRunner = project.NewCommandRunner(l, ee)

	uc, err := config.LoadUserConfig()
	if uc == nil || err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to load user config")
		return
	}

	a.userConfig = uc

	var p *project.Project

	if a.userConfig.LastOpenedProjectId != "" {
		p, err = project.LoadProject(a.userConfig.LastOpenedProjectId)
		if err != nil {
			a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to load last opened project config")
			a.userConfig.LastOpenedProjectId = ""
			err = a.persistUserConfig()
			if err != nil {
				panic(err)
			}
			return
		}
	}

	a.selectedProject = p

	a.logger.Info("Loading configuration...")
	a.logger.Info("Configuration loaded successfully")
}

func (a *App) persistSelectedProjectConfig() error {
	err := project.SaveProject(&project.Project{
		Id:            a.selectedProject.Id,
		Name:          a.selectedProject.Name,
		Commands:      a.selectedProject.Commands,
		CommandGroups: a.selectedProject.CommandGroups,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *App) persistUserConfig() error {
	err := config.SaveUserConfig(&config.UserConfig{
		ExtraPaths:          a.userConfig.ExtraPaths,
		LastOpenedProjectId: a.userConfig.LastOpenedProjectId,
	})
	if err != nil {
		return err
	}

	return nil
}
