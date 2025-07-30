package app

import (
	"gomander/internal/config"
	"gomander/internal/event"
)

func (a *App) GetUserConfig() *config.UserConfig {
	return a.userConfig
}

func (a *App) SaveExtraPaths(extraPaths []string) error {
	a.userConfig.ExtraPaths = extraPaths

	err := a.persistUserConfig()
	if err != nil {
		return err
	}

	a.logger.Info("Extra paths saved successfully")
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Extra paths saved successfully")
	a.eventEmitter.EmitEvent(event.GetUserConfig, nil)

	return nil
}
