package app

import (
	configdomain "gomander/internal/config/domain"
	"gomander/internal/eventbus"
	"gomander/internal/execution"
)

// Logger is where the app reports the startup and shutdown it goes through.
type Logger interface {
	Info(message string)
	Error(message string)
}

type EventHandlers struct {
	CleanCommandGroupsOnCommandDeleted   eventbus.EventHandler
	CleanCommandGroupsOnProjectDeleted   eventbus.EventHandler
	CleanCommandsOnProjectDeleted        eventbus.EventHandler
	AddCommandToGroupOnCommandDuplicated eventbus.EventHandler
}

type App struct {
	logger               Logger
	commandRunner        execution.Runner
	userConfigRepository configdomain.Repository
	eventBus             eventbus.EventBus
	eventHandlers        EventHandlers
}

func NewApp(
	logger Logger,
	commandRunner execution.Runner,
	userConfigRepository configdomain.Repository,
	eventBus eventbus.EventBus,
	eventHandlers EventHandlers,
) *App {
	return &App{
		logger:               logger,
		commandRunner:        commandRunner,
		userConfigRepository: userConfigRepository,
		eventBus:             eventBus,
		eventHandlers:        eventHandlers,
	}
}

func (a *App) RegisterHandlers() {
	a.eventBus.RegisterHandler(a.eventHandlers.CleanCommandGroupsOnCommandDeleted)
	a.eventBus.RegisterHandler(a.eventHandlers.CleanCommandGroupsOnProjectDeleted)
	a.eventBus.RegisterHandler(a.eventHandlers.CleanCommandsOnProjectDeleted)
	a.eventBus.RegisterHandler(a.eventHandlers.AddCommandToGroupOnCommandDuplicated)
}
