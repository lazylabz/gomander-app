package usecases

import (
	"github.com/google/uuid"

	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	projectdomain "gomander/internal/project/domain"
)

type ImportProject interface {
	Execute(projectJSON projectdomain.ProjectExportJSONv1, name, workingDirectory string) error
}

type DefaultImportProject struct {
	projectRepository      projectdomain.Repository
	commandRepository      domain.Repository
	commandGroupRepository commandgroupdomain.Repository
}

func NewImportProject(
	projectRepo projectdomain.Repository,
	commandRepo domain.Repository,
	commandGroupRepo commandgroupdomain.Repository,
) *DefaultImportProject {
	return &DefaultImportProject{
		projectRepository:      projectRepo,
		commandRepository:      commandRepo,
		commandGroupRepository: commandGroupRepo,
	}
}

func (uc *DefaultImportProject) Execute(projectJSON projectdomain.ProjectExportJSONv1, name, workingDirectory string) error {
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
	err := uc.projectRepository.Create(project)
	if err != nil {
		return err
	}

	for _, command := range commands {
		err = uc.commandRepository.Create(&command)
		if err != nil {
			return err
		}
	}

	for _, group := range commandGroups {
		err = uc.commandGroupRepository.Create(&group)
		if err != nil {
			return err
		}
	}

	return nil
}
