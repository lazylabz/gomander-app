package app

import (
	"gomander/internal/config"
	"gomander/internal/event"
)

func (a *App) GetUserConfig() *config.UserConfig {
	return a.userConfig
}

func (a *App) SaveUserConfig(newUserConfig config.UserConfig) error {
	a.userConfig = &newUserConfig

	err := a.persistUserConfig()
	if err != nil {
		return err
	}

	a.logger.Info("User configuration saved successfully")
	a.eventEmitter.EmitEvent(event.SuccessNotification, "User configuration saved successfully")
	a.eventEmitter.EmitEvent(event.GetUserConfig, nil)

	return nil
}
