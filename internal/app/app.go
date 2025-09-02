package app

import (
	"context"

	commandhandlers "gomander/internal/command/application/handlers"
	commanddomain "gomander/internal/command/domain"
	commandgrouphandlers "gomander/internal/commandgroup/application/handlers"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/facade"
	"gomander/internal/logger"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/runner"
)

type EventHandlers struct {
	CleanCommandGroupsOnCommandDeleted   commandgrouphandlers.CleanCommandGroupsOnCommandDeleted
	CleanCommandGroupsOnProjectDeleted   commandgrouphandlers.CleanCommandGroupsOnProjectDeleted
	CleanCommandsOnProjectDeleted        commandhandlers.CleanCommandsOnProjectDeleted
	AddCommandToGroupOnCommandDuplicated commandgrouphandlers.AddCommandToGroupOnCommandDuplicated
}

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

	eventHandlers EventHandlers
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

	EventHandlers EventHandlers
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

	a.eventHandlers = d.EventHandlers
}

func (a *App) RegisterHandlers() {
	// Register event handlers
	a.eventBus.RegisterHandler(a.eventHandlers.CleanCommandGroupsOnCommandDeleted)
	a.eventBus.RegisterHandler(a.eventHandlers.CleanCommandGroupsOnProjectDeleted)
	a.eventBus.RegisterHandler(a.eventHandlers.CleanCommandsOnProjectDeleted)
	a.eventBus.RegisterHandler(a.eventHandlers.AddCommandToGroupOnCommandDuplicated)
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}
