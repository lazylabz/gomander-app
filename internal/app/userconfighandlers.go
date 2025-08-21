package app

import (
	"gomander/internal/config/domain"
)

func (a *App) GetUserConfig() (*domain.Config, error) {
	return a.userConfigRepository.GetOrCreate()
}

func (a *App) SaveUserConfig(newUserConfig domain.Config) error {
	err := a.userConfigRepository.Update(&newUserConfig)
	if err != nil {
		a.logger.Error("Failed to save user configuration: " + err.Error())
		return err
	}

	a.logger.Info("User configuration saved successfully")
	return nil
}
