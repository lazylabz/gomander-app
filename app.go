package main

import (
	"context"
)

// App struct
type App struct {
	ctx                context.Context
	logger             *Logger
	eventEmitter       *EventEmitter
	commandRunner      *CommandRunner
	commandsRepository *CommandRepository
	savedConfig        *SavedConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger := NewLogger(ctx)
	eventEmitter := NewEventEmitter(ctx)
	commandRunner := NewCommandRunner(logger, eventEmitter)

	a.logger = logger
	a.eventEmitter = eventEmitter
	a.commandRunner = commandRunner

	a.logger.info("Loading configuration...")

	a.savedConfig = loadConfigOrPanic()

	a.logger.info("Configuration loaded successfully")

	a.commandsRepository = NewCommandRepository(a.savedConfig.Commands)
}

func (a *App) persistSavedConfig() {
	saveConfigOrPanic(&SavedConfig{
		Commands:      a.commandsRepository.commands,
		ExtraPaths:    a.savedConfig.ExtraPaths,
		CommandGroups: a.savedConfig.CommandGroups,
	})
}
