package app

import (
	"gomander/internal/event"
	"gomander/internal/project/domain"
)

func (a *App) GetCurrentProject() (*domain.Project, error) {
	return a.projectRepository.GetProjectById(a.openedProjectId)
}

func (a *App) GetAvailableProjects() ([]domain.Project, error) {
	return a.projectRepository.GetAllProjects()
}

func (a *App) OpenProject(projectConfigId string) error {
	config, err := a.userConfigRepository.GetOrCreateConfig()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = projectConfigId

	err = a.userConfigRepository.SaveConfig(config)
	if err != nil {
		return err
	}

	a.openedProjectId = projectConfigId

	return nil
}

func (a *App) CreateProject(project domain.Project) error {
	err := a.projectRepository.CreateProject(project)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) EditProject(project domain.Project) error {
	err := a.projectRepository.UpdateProject(project)
	if err != nil {
		return err
	}

	a.eventEmitter.EmitEvent(event.SuccessNotification, "Project edited successfully")

	return nil
}

func (a *App) CloseProject() error {
	config, err := a.userConfigRepository.GetOrCreateConfig()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = ""

	err = a.userConfigRepository.SaveConfig(config)
	if err != nil {
		return err
	}

	a.openedProjectId = ""

	return nil
}

func (a *App) DeleteProject(projectConfigId string) error {
	err := a.projectRepository.DeleteProject(projectConfigId)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) ExportProject(projectConfigId string) error {
	// TODO: Implement export logic
	//p, err := project.LoadProject(projectConfigId)
	//if err != nil {
	//	a.eventEmitter.EmitEvent(event.ErrorNotification, err.Error())
	//	return err
	//}
	//if p == nil {
	//	return errors.New("project not found")
	//}
	//
	//filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{Title: "Select a destination", CanCreateDirectories: true, DefaultFilename: p.Name + ".json"})
	//if err != nil {
	//	return err
	//}
	//
	//if filePath == "" {
	//	a.eventEmitter.EmitEvent(event.ErrorNotification, "Export cancelled")
	//	return nil
	//}
	//
	//err = project.ExportProject(p, filePath)
	//if err != nil {
	//	a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to export project: "+err.Error())
	//	return err
	//}
	//
	//a.eventEmitter.EmitEvent(event.SuccessNotification, "Project exported successfully")
	return nil
}

func (a *App) ImportProject() error {
	// TODO: Implement import logic
	//err := project.ImportProject(p)
	//if err != nil {
	//	a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to import project: "+err.Error())
	//	return err
	//}
	//
	//a.eventEmitter.EmitEvent(event.SuccessNotification, "Project imported successfully")
	return nil
}

func (a *App) GetProjectToImport() error {
	// TODO: Implement import logic
	//filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Select a project file", Filters: []runtime.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}}})
	//if err != nil {
	//	return nil, err
	//}
	//if filePath == "" {
	//	return nil, errors.New("import cancelled")
	//}
	//
	//return project.LoadProjectFromPath(filePath)

	return nil
}
