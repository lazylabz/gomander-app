package app

import (
	"gomander/internal/event"
)

func (a *App) StopCommand(id string) {
	// Check if the command exists before trying to stop it
	_, err := a.commandRepository.Get(id)

	if err != nil {
		a.logger.Error(err.Error())
		return
	}

	err = a.commandRunner.StopRunningCommand(id)

	if err != nil {
		a.logger.Error(err.Error())
		return
	}

	a.logger.Info("Command stopped: " + id)

	a.eventEmitter.EmitEvent(event.ProcessFinished, id)
}
