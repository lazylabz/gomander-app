package main

func (a *App) GetCommandGroups() []CommandGroup {
	return a.savedConfig.CommandGroups
}

func (a *App) SaveCommandGroups(commandGroups []CommandGroup) {
	a.savedConfig.CommandGroups = commandGroups
	a.persistSavedConfig()
	a.eventEmitter.emitEvent(SuccessNotification, "Command groups saved successfully")
	a.eventEmitter.emitEvent(GetCommandGroups, nil)
}
