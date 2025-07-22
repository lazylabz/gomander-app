package main

func (a *App) SaveUserConfig(userConfig UserConfig) {
	a.savedConfig.ExtraPaths = userConfig.ExtraPaths

	a.persistSavedConfig()

	a.logger.info("Extra paths saved successfully")
	a.eventEmitter.emitEvent(SuccessNotification, "Extra paths saved successfully")
	a.eventEmitter.emitEvent(GetUserConfig, nil)
}

func (a *App) GetUserConfig() UserConfig {
	return UserConfig{
		ExtraPaths: a.savedConfig.ExtraPaths,
	}
}
