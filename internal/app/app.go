package app

import (
	"context"

	commandhandlers "gomander/internal/command/application/handlers"
	commanddomain "gomander/internal/command/domain"
	commandgrouphandlers "gomander/internal/commandgroup/application/handlers"
	commandgroupusecases "gomander/internal/commandgroup/application/usecases"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configusecases "gomander/internal/config/application/usecases"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/event"
	"gomander/internal/eventbus"
	"gomander/internal/facade"
	"gomander/internal/logger"
	projectusecases "gomander/internal/project/application/usecases"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/runner"
)

type EventHandlers struct {
	CleanCommandGroupsOnCommandDeleted   commandgrouphandlers.CleanCommandGroupsOnCommandDeleted
	CleanCommandGroupsOnProjectDeleted   commandgrouphandlers.CleanCommandGroupsOnProjectDeleted
	CleanCommandsOnProjectDeleted        commandhandlers.CleanCommandsOnProjectDeleted
	AddCommandToGroupOnCommandDuplicated commandgrouphandlers.AddCommandToGroupOnCommandDuplicated
}

type UseCases struct {
	GetUserConfig                 configusecases.GetUserConfig
	SaveUserConfig                configusecases.SaveUserConfig
	GetCurrentProject             projectusecases.GetCurrentProject
	GetAvailableProjects          projectusecases.GetAvailableProjects
	OpenProject                   projectusecases.OpenProject
	CreateProject                 projectusecases.CreateProject
	EditProject                   projectusecases.EditProject
	CloseProject                  projectusecases.CloseProject
	DeleteProject                 projectusecases.DeleteProject
	GetCommandGroups              commandgroupusecases.GetCommandGroups
	CreateCommandGroup            commandgroupusecases.CreateCommandGroup
	UpdateCommandGroup            commandgroupusecases.UpdateCommandGroup
	DeleteCommandGroup            commandgroupusecases.DeleteCommandGroup
	RemoveCommandFromCommandGroup commandgroupusecases.RemoveCommandFromCommandGroup
}

// App struct
type App struct {
	ctx context.Context

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
	UseCases      UseCases
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
	UseCases      UseCases
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
	a.UseCases = d.UseCases
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
