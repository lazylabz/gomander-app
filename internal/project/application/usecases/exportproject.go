package usecases

import (
	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/dialog"
	projectdomain "gomander/internal/project/domain"
)

// BlueprintWriter writes a Project to the file the user picked, in whatever
// format that file is.
type BlueprintWriter interface {
	Write(filePath string, blueprint projectdomain.Blueprint) error
}

type ExportProject struct {
	projectRepository      projectdomain.Repository
	commandRepository      domain.Repository
	commandGroupRepository commandgroupdomain.Repository
	dialogs                dialog.Dialogs
	files                  BlueprintWriter
}

func NewExportProject(projectRepo projectdomain.Repository, commandRepo domain.Repository, commandGroupRepo commandgroupdomain.Repository, dialogs dialog.Dialogs, files BlueprintWriter) *ExportProject {
	return &ExportProject{
		projectRepository:      projectRepo,
		commandRepository:      commandRepo,
		commandGroupRepository: commandGroupRepo,
		dialogs:                dialogs,
		files:                  files,
	}
}

func (uc *ExportProject) Execute(projectId string) (string, error) {
	project, err := uc.projectRepository.Get(projectId)
	if err != nil {
		return "", err
	}

	filePath, err := uc.dialogs.AskWhereToSaveFile(dialog.SaveFileRequest{
		Title:                "Select a destination",
		CanCreateDirectories: true,
		DefaultFilename:      project.Name + ".json",
	})
	if err != nil {
		return "", err
	}

	if filePath == "" {
		// User canceled the dialog
		return "", nil
	}

	commands, err := uc.commandRepository.GetAll(projectId)
	if err != nil {
		return "", err
	}

	commandGroups, err := uc.commandGroupRepository.GetAllWithCommandIds(projectId)
	if err != nil {
		return "", err
	}

	// An export carries no working directory: it is the one thing the user is
	// asked for again when the file is imported somewhere else.
	blueprint := projectdomain.Blueprint{
		Name: project.Name,
	}

	for _, cmd := range commands {
		blueprint.Commands = append(blueprint.Commands, projectdomain.BlueprintCommand{
			Id:               cmd.Id,
			Name:             cmd.Name,
			Command:          cmd.Command,
			WorkingDirectory: cmd.WorkingDirectory,
		})
	}

	for _, group := range commandGroups {
		blueprint.CommandGroups = append(blueprint.CommandGroups, projectdomain.BlueprintCommandGroup{
			Id:         group.Id,
			Name:       group.Name,
			CommandIds: group.CommandIds,
		})
	}

	if err := uc.files.Write(filePath, blueprint); err != nil {
		return "", err
	}

	return filePath, nil
}
