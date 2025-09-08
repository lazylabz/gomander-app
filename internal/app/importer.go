package app

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
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
		return err
	}

	for _, command := range commands {
		err = a.commandRepository.Create(&command)
		if err != nil {
			return err
		}
	}

	for _, group := range commandGroups {
		err = a.commandGroupRepository.Create(&group)
		if err != nil {
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
		return nil, err
	}

	// Unmarshal JSON data
	err = json.Unmarshal(fileData, &projectJSON)
	if err != nil {
		return nil, err
	}

	return
}
