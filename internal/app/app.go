package app

import (
	"context"

	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/logger"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/runner"
)

// App struct
type App struct {
	ctx context.Context

	openedProjectId string

	logger        logger.Logger
	eventEmitter  event.EventEmitter
	commandRunner runner.Runner

	commandRepository      commanddomain.Repository
	commandGroupRepository commandgroupdomain.Repository
	projectRepository      projectdomain.Repository
	userConfigRepository   configdomain.Repository
}

func (a *App) LoadDependencies(l logger.Logger,
	ee event.EventEmitter,
	r runner.Runner,
	commandRepository commanddomain.Repository,
	commandGroupRepository commandgroupdomain.Repository,
	projectRepository projectdomain.Repository,
	configRepository configdomain.Repository,
) {
	a.logger = l
	a.eventEmitter = ee
	a.commandRunner = r

	a.commandRepository = commandRepository
	a.commandGroupRepository = commandGroupRepository
	a.projectRepository = projectRepository
	a.userConfigRepository = configRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}
