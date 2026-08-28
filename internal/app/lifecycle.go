package app

import (
	"context"
)

func (a *App) Startup(_ context.Context) error {
	a.logger.Info("Loading configuration...")

	if _, err := a.userConfigRepository.GetOrCreate(); err != nil {
		return err
	}

	a.logger.Info("Configuration loaded successfully")

	return nil
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
