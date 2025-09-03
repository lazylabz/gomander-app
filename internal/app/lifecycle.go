package app

import (
	"context"
)

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.logger.Info("Loading configuration...")

	_, err := a.userConfigRepository.GetOrCreate()
	if err != nil {
		panic(err)
	}

	a.logger.Info("Configuration loaded successfully")
}

func (a *App) OnBeforeClose(_ context.Context) (prevent bool) {
	errs := a.commandRunner.StopAllRunningCommands()

	if len(errs) > 0 {
		for _, err := range errs {
			a.logger.Error(err.Error())
		}
		return true // Prevent the application from closing
	}

	return false // Allow the application to close
}
