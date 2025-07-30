package app

import (
	"gomander/internal/event"
	"gomander/internal/project"
)

func (a *App) GetCommandGroups() []project.CommandGroup {
	return a.selectedProject.GetCommandGroups()
}

func (a *App) SaveCommandGroups(commandGroups []project.CommandGroup) error {
	a.selectedProject.SetCommandGroups(commandGroups)
	err := a.persistSelectedProjectConfig()
	if err != nil {
		return err
	}
	a.eventEmitter.EmitEvent(event.SuccessNotification, "Command groups saved successfully")
	a.eventEmitter.EmitEvent(event.GetCommandGroups, nil)

	return nil
}
