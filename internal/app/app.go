package app

import (
	"context"

	"gomander/internal/command"
	"gomander/internal/commandgroup"
	"gomander/internal/config"
	"gomander/internal/event"
	"gomander/internal/extrapath"
	"gomander/internal/logger"
)

// App struct
type App struct {
	ctx           context.Context
	logger        *logger.Logger
	eventEmitter  *event.EventEmitter
	commandRunner *command.Runner

	commandsRepository      *command.Repository
	commandGroupsRepository *commandgroup.Repository
	extraPathRepository     *extrapath.Repository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	cfg := config.LoadConfigOrPanic()

	a.ctx = ctx
	l := logger.NewLogger(ctx)
	ee := event.NewEventEmitter(ctx)

	a.logger = l
	a.eventEmitter = ee
	a.commandRunner = command.NewCommandRunner(l, ee)

	a.commandsRepository = command.NewCommandRepository(cfg.Commands)
	a.commandGroupsRepository = commandgroup.NewCommandGroupRepository(cfg.CommandGroups)
	a.extraPathRepository = extrapath.NewExtraPathRepository(cfg.ExtraPaths)

	a.logger.Info("Loading configuration...")
	a.logger.Info("Configuration loaded successfully")
}

func (a *App) persistRepositoryInformation() {
	config.SaveConfigOrPanic(&config.Config{
		Commands:      a.commandsRepository.GetCommands(),
		ExtraPaths:    a.extraPathRepository.GetExtraPaths(),
		CommandGroups: a.commandGroupsRepository.GetCommandGroups(),
	})
}
