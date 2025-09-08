package usecases

import (
	"context"
	"encoding/json"

	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/facade"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ExportProject interface {
	Execute(projectId string) error
}

type DefaultExportProject struct {
	ctx                    context.Context
	projectRepository      projectdomain.Repository
	commandRepository      domain.Repository
	commandGroupRepository commandgroupdomain.Repository
	runtimeFacade          facade.RuntimeFacade
	fsFacade               facade.FsFacade
}

func NewExportProject(ctx context.Context, projectRepo projectdomain.Repository, commandRepo domain.Repository, commandGroupRepo commandgroupdomain.Repository, runtimeFacade facade.RuntimeFacade, fsFacade facade.FsFacade) *DefaultExportProject {
	return &DefaultExportProject{
		ctx:                    ctx,
		projectRepository:      projectRepo,
		commandRepository:      commandRepo,
		commandGroupRepository: commandGroupRepo,
		runtimeFacade:          runtimeFacade,
		fsFacade:               fsFacade,
	}
}

func (uc *DefaultExportProject) Execute(projectId string) error {
	project, err := uc.projectRepository.Get(projectId)
	if err != nil {
		return err
	}

	filePath, err := uc.runtimeFacade.SaveFileDialog(uc.ctx, runtime.SaveDialogOptions{
		Title:                "Select a destination",
		CanCreateDirectories: true,
		DefaultFilename:      project.Name + ".json",
	})
	if err != nil {
		return err
	}

	if filePath == "" {
		// User canceled the dialog
		return nil
	}

	commands, err := uc.commandRepository.GetAll(projectId)
	if err != nil {
		return err
	}

	commandGroups, err := uc.commandGroupRepository.GetAll(projectId)
	if err != nil {
		return err
	}

	exportData := projectdomain.ProjectExportJSONv1{
		Version: 1,
		Name:    project.Name,
	}

	// Prepare commands for export
	for _, cmd := range commands {
		exportData.Commands = append(exportData.Commands, projectdomain.CommandJSONv1{
			Id:               cmd.Id,
			Name:             cmd.Name,
			Command:          cmd.Command,
			WorkingDirectory: cmd.WorkingDirectory,
		})
	}

	// Prepare command groups for export
	for _, group := range commandGroups {
		exportData.CommandGroups = append(exportData.CommandGroups, projectdomain.CommandGroupJSONv1{
			Name:       group.Name,
			CommandIds: array.Map(group.Commands, func(cmd domain.Command) string { return cmd.Id }),
		})
	}

	// Marshal to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	err = uc.fsFacade.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
