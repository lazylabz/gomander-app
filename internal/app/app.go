package app

import (
	"context"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/commandgroup/application/handlers"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/facade"
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

	fsFacade      facade.FsFacade
	runtimeFacade facade.RuntimeFacade

	eventBus eventbus.EventBus

	cleanCommandGroupsOnCommandDeletedHandler handlers.CleanCommandGroupsOnCommandDeleted
}

type Dependencies struct {
	Logger       logger.Logger
	EventEmitter event.EventEmitter
	Runner       runner.Runner

	CommandRepository      commanddomain.Repository
	CommandGroupRepository commandgroupdomain.Repository
	ProjectRepository      projectdomain.Repository
	ConfigRepository       configdomain.Repository

	FsFacade      facade.FsFacade
	RuntimeFacade facade.RuntimeFacade

	EventBus eventbus.EventBus

	CleanCommandGroupsOnCommandDeletedHandler handlers.CleanCommandGroupsOnCommandDeleted
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

	a.eventBus = d.EventBus

	a.cleanCommandGroupsOnCommandDeletedHandler = d.CleanCommandGroupsOnCommandDeletedHandler
}

func (a *App) RegisterHandlers() {
	// Register event handlers
	a.eventBus.RegisterHandler(a.cleanCommandGroupsOnCommandDeletedHandler)
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}
