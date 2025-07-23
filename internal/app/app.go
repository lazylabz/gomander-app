package app

import (
	"context"
	"gomander/internal/command"
	"gomander/internal/config"
	"gomander/internal/event"
	"gomander/internal/logger"
)

// App struct
type App struct {
	ctx                context.Context
	logger             *logger.Logger
	eventEmitter       *event.EventEmitter
	commandRunner      *command.CommandRunner
	commandsRepository *command.CommandRepository
	savedConfig        *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	logger := logger.NewLogger(ctx)
	eventEmitter := event.NewEventEmitter(ctx)
	commandRunner := command.NewCommandRunner(logger, eventEmitter)

	a.logger = logger
	a.eventEmitter = eventEmitter
	a.commandRunner = commandRunner

	a.logger.Info("Loading configuration...")

	a.savedConfig = config.LoadConfigOrPanic()

	a.logger.Info("Configuration loaded successfully")

	a.commandsRepository = command.NewCommandRepository(a.savedConfig.Commands)
}

func (a *App) persistConfig() {
	config.SaveConfigOrPanic(&config.Config{
		Commands:      a.commandsRepository.GetCommands(),
		ExtraPaths:    a.savedConfig.ExtraPaths,
		CommandGroups: a.savedConfig.CommandGroups,
	})
}
