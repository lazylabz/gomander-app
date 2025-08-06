package app

import (
	"errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/event"
	"gomander/internal/project"
)

func (a *App) GetCurrentProject() *project.Project {
	return a.selectedProject
}

func (a *App) GetAvailableProjects() ([]*project.Project, error) {
	return project.GetAllProjectsAvailableInProjectsFolder()
}

func (a *App) OpenProject(projectConfigId string) (*project.Project, error) {
	p, err := project.LoadProject(projectConfigId)
	if err != nil {
		return nil, err
	}

	a.userConfig.LastOpenedProjectId = p.Id

	err = a.persistUserConfig()
	if err != nil {
		return nil, err
	}

	a.selectedProject = p

	return p, nil
}

func (a *App) CreateProject(id, name, baseWorkingDirectory string) error {
	err := project.SaveProject(&project.Project{
		Id:                   id,
		Name:                 name,
		BaseWorkingDirectory: baseWorkingDirectory,
		Commands:             make(map[string]project.Command),
		CommandGroups:        make([]project.CommandGroup, 0),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *App) EditProject(p project.Project) error {
	isEditingSelectedProject := a.selectedProject != nil && a.selectedProject.Id == p.Id

	// TODO: Edit project should only be able to update the name and base working directory
	err := project.SaveProject(&p)

	if err != nil {
		return err
	}

	if isEditingSelectedProject {
		a.selectedProject = &p
	}

	a.eventEmitter.EmitEvent(event.SuccessNotification, "Project edited successfully")

	return nil
}

func (a *App) CloseProject() error {
	if a.selectedProject == nil {
		return nil
	}

	a.selectedProject = nil

	a.userConfig.LastOpenedProjectId = ""

	err := a.persistUserConfig()
	if err != nil {
		return err
	}

	errs := a.commandRunner.StopAllRunningCommands()
	if len(errs) > 0 {
		return errors.New("error stopping running commands")
	}

	return nil
}

func (a *App) DeleteProject(projectConfigId string) error {
	err := project.DeleteProject(projectConfigId)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ExportProject(projectConfigId string) error {
	p, err := project.LoadProject(projectConfigId)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
		return err
	}
	if p == nil {
		return errors.New("project not found")
	}

	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{Title: "Select a destination", CanCreateDirectories: true, DefaultFilename: p.Name + ".json"})
	if err != nil {
		return err
	}

	if filePath == "" {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Export cancelled")
		return nil
	}

	err = project.ExportProject(p, filePath)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to export project: "+err.Error())
		return err
	}

	a.eventEmitter.EmitEvent(event.SuccessNotification, "Project exported successfully")
	return nil
}

func (a *App) ImportProject(p project.Project) error {
	err := project.ImportProject(p)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to import project: "+err.Error())
		return err
	}

	a.eventEmitter.EmitEvent(event.SuccessNotification, "Project imported successfully")
	return nil
}

func (a *App) GetProjectToImport() (*project.Project, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Select a project file", Filters: []runtime.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}}})
	if err != nil {
		return nil, err
	}
	if filePath == "" {
		return nil, errors.New("import cancelled")
	}

	return project.LoadProjectFromPath(filePath)
}
