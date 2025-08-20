package app

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
)

type CommandGroupJSONv1 struct {
	Name       string   `json:"name"`
	CommandIds []string `json:"commandIds"`
}

type CommandJSONv1 struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

type ProjectExportJSONv1 struct {
	Version       int                  `json:"version"`
	Name          string               `json:"name"`
	Commands      []CommandJSONv1      `json:"commands"`
	CommandGroups []CommandGroupJSONv1 `json:"commandGroups"`
}

func (a *App) ExportProject(projectConfigId string) (err error) {
	project, err := a.projectRepository.Get(projectConfigId)
	if err != nil {
		return err
	}

	filePath, err := a.runtimeFacade.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{Title: "Select a destination", CanCreateDirectories: true, DefaultFilename: project.Name + ".json"})
	if err != nil {
		return err
	}

	if filePath == "" {
		return nil
	}

	commands, err := a.commandRepository.GetAll(projectConfigId)
	if err != nil {
		return err
	}
	commandGroups, err := a.commandGroupRepository.GetAll(projectConfigId)
	if err != nil {
		return err
	}

	exportData := ProjectExportJSONv1{
		Version: 1,
	}

	// Load project data
	exportData.Name = project.Name

	// Load commands
	for _, cmd := range commands {
		exportData.Commands = append(exportData.Commands, CommandJSONv1{
			Id:               cmd.Id,
			Name:             cmd.Name,
			Command:          cmd.Command,
			WorkingDirectory: cmd.WorkingDirectory,
		})
	}

	// Load command groups
	for _, group := range commandGroups {
		exportData.CommandGroups = append(exportData.CommandGroups, CommandGroupJSONv1{
			Name:       group.Name,
			CommandIds: array.Map(group.Commands, func(cmd domain.Command) string { return cmd.Id }),
		})
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to encode project data: "+err.Error())
		return err
	}

	// Write directly to file
	err = a.fsFacade.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to write file: "+err.Error())
		return err
	}

	return nil
}

func (a *App) ImportProject(projectJSON ProjectExportJSONv1, name, workingDirectory string) error {
	project := projectdomain.Project{
		Id:               uuid.New().String(),
		Name:             name,
		WorkingDirectory: workingDirectory,
	}

	commands := make([]domain.Command, 0, len(projectJSON.Commands))
	commandGroups := make([]commandgroupdomain.CommandGroup, 0, len(projectJSON.CommandGroups))

	commandIdsToNewRandomIds := make(map[string]string)
	newIdsToCommand := make(map[string]domain.Command)

	for i, cmd := range projectJSON.Commands {
		newCommand := domain.Command{
			Id:               uuid.New().String(),
			Name:             cmd.Name,
			Command:          cmd.Command,
			WorkingDirectory: cmd.WorkingDirectory,
			ProjectId:        project.Id,
			Position:         i,
		}

		commands = append(commands, newCommand)
		commandIdsToNewRandomIds[cmd.Id] = newCommand.Id
		newIdsToCommand[newCommand.Id] = newCommand
	}

	for i, group := range projectJSON.CommandGroups {
		newGroup := commandgroupdomain.CommandGroup{
			Id:        uuid.New().String(),
			Name:      group.Name,
			ProjectId: project.Id,
			Position:  i,
		}

		for _, cmdId := range group.CommandIds {
			if newCmdId, exists := commandIdsToNewRandomIds[cmdId]; exists {
				newGroup.Commands = append(newGroup.Commands, newIdsToCommand[newCmdId])
			}
		}

		commandGroups = append(commandGroups, newGroup)
	}

	// Persist everything
	err := a.projectRepository.Create(project)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to import project: "+err.Error())
		return err
	}

	for _, command := range commands {
		err = a.commandRepository.Create(&command)
		if err != nil {
			a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to import command: "+err.Error())
			return err
		}
	}

	for _, group := range commandGroups {
		err = a.commandGroupRepository.Create(&group)
		if err != nil {
			a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to import command group: "+err.Error())
			return err
		}
	}

	return nil
}

func (a *App) GetProjectToImport() (projectJSON *ProjectExportJSONv1, err error) {
	filePath, err := a.runtimeFacade.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select a project file",
		Filters: []runtime.FileFilter{{DisplayName: "JSON Files", Pattern: "*.json"}},
	})
	if err != nil {
		return nil, err
	}

	if filePath == "" {
		return nil, nil // User canceled
	}

	// Read entire file into memory
	fileData, err := a.fsFacade.ReadFile(filePath)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to read file: "+err.Error())
		return nil, err
	}

	// Unmarshal JSON data
	err = json.Unmarshal(fileData, &projectJSON)
	if err != nil {
		a.eventEmitter.EmitEvent(event.ErrorNotification, "Failed to decode project data: "+err.Error())
		return nil, err
	}

	return
}
