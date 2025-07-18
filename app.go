package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	commands map[string]Command
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		commands: make(map[string]Command),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.logInfo("Loading configuration...")
	config, err := loadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	a.logInfo("Configuration loaded successfully")

	a.commands = config.Commands
}

func (a *App) logInfo(message string) {
	runtime.LogInfo(a.ctx, message)
}

func (a *App) logDebug(message string) {
	runtime.LogDebug(a.ctx, message)
}

func (a *App) logError(message string) {
	runtime.LogError(a.ctx, message)
}

func (a *App) notifyError(message string) {
	a.emitEvent(ErrorNotification, message)
}
