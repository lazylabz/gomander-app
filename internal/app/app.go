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

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	l := logger.NewLogger(ctx)
	ee := event.NewEventEmitter(ctx)
	cr := command.NewCommandRunner(l, ee)

	a.logger = l
	a.eventEmitter = ee
	a.commandRunner = cr

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
