package app

import (
	"gomander/internal/config/domain"
	"gomander/internal/event"
)

func (a *App) GetUserConfig() (*domain.Config, error) {
	return a.userConfigRepository.GetOrCreate()
}

func (a *App) SaveUserConfig(newUserConfig domain.Config) error {
	err := a.userConfigRepository.Update(&newUserConfig)
	if err != nil {
		a.logger.Error("Failed to save user configuration: " + err.Error())
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to save user configuration")
		return err
	}

	a.logger.Info("User configuration saved successfully")
	a.eventEmitter.EmitEvent(event.SuccessNotification, "User configuration saved successfully")
	a.eventEmitter.EmitEvent(event.GetUserConfig, nil)

	return nil
}
