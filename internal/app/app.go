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

	fsFacade      FsFacade
	runtimeFacade RuntimeFacade
}

type Dependencies struct {
	Logger                 logger.Logger
	EventEmitter           event.EventEmitter
	Runner                 runner.Runner
	CommandRepository      commanddomain.Repository
	CommandGroupRepository commandgroupdomain.Repository
	ProjectRepository      projectdomain.Repository
	ConfigRepository       configdomain.Repository
	FsFacade               FsFacade
	RuntimeFacade          RuntimeFacade
}

func (a *App) LoadDependencies(d Dependencies) {
	a.logger = d.Logger
	a.eventEmitter = d.EventEmitter
	a.commandRunner = d.Runner

	a.commandRepository = d.CommandRepository
	a.commandGroupRepository = d.CommandGroupRepository
	a.projectRepository = d.ProjectRepository
	a.userConfigRepository = d.ConfigRepository
	a.fsFacade = d.FsFacade
	a.runtimeFacade = d.RuntimeFacade

	a.ctx = context.Background()
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}
