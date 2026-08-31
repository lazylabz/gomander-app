package main

import (
	transport "gomander/cmd/gomander/transport/domain"
	"gomander/internal/localization"
	projectusecases "gomander/internal/project/application/usecases"
	"gomander/internal/usecases"
)

type WailsControllers struct {
	useCases usecases.Registry
}

func NewWailsControllers() *WailsControllers {
	return &WailsControllers{}
}

func (wc *WailsControllers) loadUseCases(useCases usecases.Registry) {
	wc.useCases = useCases
}

// User config controllers

func (wc *WailsControllers) GetUserConfigController() (*transport.Config, error) {
	config, err := wc.useCases.GetUserConfig.Execute()
	return transport.Optional(config, transport.FromConfig), err
}

func (wc *WailsControllers) SaveUserConfigController(newConfig transport.Config) error {
	return wc.useCases.SaveUserConfig.Execute(newConfig.ToDomain())
}

// Project controllers

func (wc *WailsControllers) GetCurrentProjectController() (*transport.Project, error) {
	project, err := wc.useCases.GetCurrentProject.Execute()
	return transport.Optional(project, transport.FromProject), err
}

func (wc *WailsControllers) GetAvailableProjectsController() ([]transport.Project, error) {
	projects, err := wc.useCases.GetAvailableProjects.Execute()
	return transport.FromProjects(projects), err
}

func (wc *WailsControllers) OpenProjectController(projectId string) error {
	return wc.useCases.OpenProject.Execute(projectId)
}

func (wc *WailsControllers) CreateProjectController(project transport.Project) error {
	return wc.useCases.CreateProject.Execute(project.ToDomain())
}

func (wc *WailsControllers) EditProjectController(project transport.Project) error {
	return wc.useCases.EditProject.Execute(project.ToDomain())
}

func (wc *WailsControllers) CloseProjectController() error {
	return wc.useCases.CloseProject.Execute()
}

func (wc *WailsControllers) DeleteProjectController(projectId string) error {
	return wc.useCases.DeleteProject.Execute(projectId)
}

func (wc *WailsControllers) ExportProjectController(projectId string) (string, error) {
	return wc.useCases.ExportProject.Execute(projectId)
}

func (wc *WailsControllers) ImportProjectController(blueprint transport.ProjectBlueprint, name, workingDirectory string) error {
	return wc.useCases.ImportProject.Execute(blueprint.ToDomain(), name, workingDirectory)
}

func (wc *WailsControllers) GetProjectToImportController() (*transport.ProjectBlueprint, error) {
	return wc.projectToImport(projectusecases.FileTypeGomander)
}

func (wc *WailsControllers) GetProjectToImportFromPackageJsonController() (*transport.ProjectBlueprint, error) {
	return wc.projectToImport(projectusecases.FileTypePackageJSON)
}

func (wc *WailsControllers) projectToImport(fileType projectusecases.FileType) (*transport.ProjectBlueprint, error) {
	blueprint, err := wc.useCases.GetProjectToImport.Execute(fileType)
	if err != nil || blueprint == nil {
		return nil, err
	}

	dto := transport.FromBlueprint(*blueprint)

	return &dto, nil
}

// CommandGroup controllers

func (wc *WailsControllers) GetCommandGroupsController() ([]transport.CommandGroup, error) {
	commandGroups, err := wc.useCases.GetCommandGroups.Execute()
	return transport.FromCommandGroups(commandGroups), err
}

func (wc *WailsControllers) CreateCommandGroupController(commandGroup transport.CommandGroup) error {
	newCommandGroup := commandGroup.ToDomain()
	return wc.useCases.CreateCommandGroup.Execute(&newCommandGroup)
}

func (wc *WailsControllers) UpdateCommandGroupController(commandGroup transport.CommandGroup) error {
	updatedCommandGroup := commandGroup.ToDomain()
	return wc.useCases.UpdateCommandGroup.Execute(&updatedCommandGroup)
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

func (wc *WailsControllers) RunCommandGroupController(commandGroupId string) error {
	return wc.useCases.RunCommandGroup.Execute(commandGroupId)
}

func (wc *WailsControllers) StopCommandGroupController(commandGroupId string) error {
	return wc.useCases.StopCommandGroup.Execute(commandGroupId)
}

// Command controllers

func (wc *WailsControllers) GetCommandsController() ([]transport.Command, error) {
	commands, err := wc.useCases.GetCommands.Execute()
	return transport.FromCommands(commands), err
}

func (wc *WailsControllers) AddCommandController(command transport.Command) error {
	return wc.useCases.AddCommand.Execute(command.ToDomain())
}

func (wc *WailsControllers) DuplicateCommandController(commandId string, targetGroupId string) error {
	return wc.useCases.DuplicateCommand.Execute(commandId, targetGroupId)
}

func (wc *WailsControllers) RemoveCommandController(commandId string) error {
	return wc.useCases.RemoveCommand.Execute(commandId)
}

func (wc *WailsControllers) EditCommandController(command transport.Command) error {
	return wc.useCases.EditCommand.Execute(command.ToDomain())
}

func (wc *WailsControllers) ReorderCommandsController(orderedIds []string) error {
	return wc.useCases.ReorderCommands.Execute(orderedIds)
}

func (wc *WailsControllers) RunCommandController(commandId string) error {
	return wc.useCases.RunCommand.Execute(commandId)
}

func (wc *WailsControllers) StopCommandController(commandId string) error {
	return wc.useCases.StopCommand.Execute(commandId)
}

// Release controllers

func (wc *WailsControllers) GetCurrentReleaseController() string {
	return wc.useCases.GetCurrentRelease.Execute()
}

func (wc *WailsControllers) CheckForNewReleaseController() (string, error) {
	return wc.useCases.CheckForNewRelease.Execute()
}

func (wc *WailsControllers) DownloadReleaseController(version string) (string, error) {
	return wc.useCases.DownloadRelease.Execute(version)
}

func (wc *WailsControllers) InstallReleaseAndQuitController(binaryPath string) error {
	return wc.useCases.InstallReleaseAndQuit.Execute(binaryPath)
}

// Localization controllers

func (wc *WailsControllers) GetTranslationController(locale string) (*localization.Localization, error) {
	return wc.useCases.GetTranslation.Execute(locale)
}

func (wc *WailsControllers) GetSupportedLanguagesController() ([]string, error) {
	return wc.useCases.GetSupportedLanguages.Execute()
}
