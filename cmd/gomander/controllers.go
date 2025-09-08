package main

import (
	"gomander/internal/app"
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	projectdomain "gomander/internal/project/domain"
)

type WailsControllers struct {
	useCases app.UseCases
}

func NewWailsControllers() *WailsControllers {
	return &WailsControllers{}
}

func (wc *WailsControllers) loadUseCases(useCases app.UseCases) {
	wc.useCases = useCases
}

// User config controllers

func (wc *WailsControllers) GetUserConfigController() (*configdomain.Config, error) {
	return wc.useCases.GetUserConfig.Execute()
}

func (wc *WailsControllers) SaveUserConfigController(newConfig configdomain.Config) error {
	return wc.useCases.SaveUserConfig.Execute(newConfig)
}

// Project controllers

func (wc *WailsControllers) GetCurrentProjectController() (*projectdomain.Project, error) {
	return wc.useCases.GetCurrentProject.Execute()
}

func (wc *WailsControllers) GetAvailableProjectsController() ([]projectdomain.Project, error) {
	return wc.useCases.GetAvailableProjects.Execute()
}

func (wc *WailsControllers) OpenProjectController(projectId string) error {
	return wc.useCases.OpenProject.Execute(projectId)
}

func (wc *WailsControllers) CreateProjectController(project projectdomain.Project) error {
	return wc.useCases.CreateProject.Execute(project)
}

func (wc *WailsControllers) EditProjectController(project projectdomain.Project) error {
	return wc.useCases.EditProject.Execute(project)
}

func (wc *WailsControllers) CloseProjectController() error {
	return wc.useCases.CloseProject.Execute()
}

func (wc *WailsControllers) DeleteProjectController(projectId string) error {
	return wc.useCases.DeleteProject.Execute(projectId)
}

// CommandGroup controllers

func (wc *WailsControllers) GetCommandGroupsController() ([]commandgroupdomain.CommandGroup, error) {
	return wc.useCases.GetCommandGroups.Execute()
}

func (wc *WailsControllers) CreateCommandGroupController(commandGroup commandgroupdomain.CommandGroup) error {
	return wc.useCases.CreateCommandGroup.Execute(&commandGroup)
}

func (wc *WailsControllers) UpdateCommandGroupController(commandGroup commandgroupdomain.CommandGroup) error {
	return wc.useCases.UpdateCommandGroup.Execute(&commandGroup)
}

func (wc *WailsControllers) DeleteCommandGroupController(commandGroupId string) error {
	return wc.useCases.DeleteCommandGroup.Execute(commandGroupId)
}

func (wc *WailsControllers) RemoveCommandFromCommandGroupController(commandId string, commandGroupId string) error {
	return wc.useCases.RemoveCommandFromCommandGroup.Execute(commandId, commandGroupId)
}

func (wc *WailsControllers) ReorderCommandGroupsController(newOrderedIds []string) error {
	return wc.useCases.ReorderCommandGroups.Execute(newOrderedIds)
}

// Command controllers

func (wc *WailsControllers) GetCommandsController() ([]commanddomain.Command, error) {
	return wc.useCases.GetCommands.Execute()
}

func (wc *WailsControllers) AddCommandController(command commanddomain.Command) error {
	return wc.useCases.AddCommand.Execute(command)
}

func (wc *WailsControllers) DuplicateCommandController(commandId string, targetGroupId string) error {
	return wc.useCases.DuplicateCommand.Execute(commandId, targetGroupId)
}

func (wc *WailsControllers) RemoveCommandController(commandId string) error {
	return wc.useCases.RemoveCommand.Execute(commandId)
}

func (wc *WailsControllers) EditCommandController(command commanddomain.Command) error {
	return wc.useCases.EditCommand.Execute(command)
}
