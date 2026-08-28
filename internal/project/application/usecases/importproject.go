package usecases

import (
	"github.com/google/uuid"

	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/unitofwork"
)

type ImportProject struct {
	unitOfWork unitofwork.UnitOfWork
}

func NewImportProject(unitOfWork unitofwork.UnitOfWork) *ImportProject {
	return &ImportProject{
		unitOfWork: unitOfWork,
	}
}

func (uc *ImportProject) Execute(projectJSON projectdomain.ProjectExportJSONv1, name, workingDirectory string) error {
	project := projectdomain.Project{
		Id:               uuid.New().String(),
		Name:             name,
		WorkingDirectory: workingDirectory,
	}

	commands := make([]domain.Command, 0, len(projectJSON.Commands))
	commandGroups := make([]commandgroupdomain.CommandGroup, 0, len(projectJSON.CommandGroups))

	commandIdsToNewRandomIds := make(map[string]string)
	newIdsToCommand := make(map[string]domain.Command)

	for _, cmd := range projectJSON.Commands {
		newCommand := domain.Command{
			Id:               uuid.New().String(),
			Name:             cmd.Name,
			Command:          cmd.Command,
			WorkingDirectory: cmd.WorkingDirectory,
			ProjectId:        project.Id,
			Position:         domain.Order.End(commands),
		}

		commands = append(commands, newCommand)
		commandIdsToNewRandomIds[cmd.Id] = newCommand.Id
		newIdsToCommand[newCommand.Id] = newCommand
	}

	for _, group := range projectJSON.CommandGroups {
		newGroup := commandgroupdomain.CommandGroup{
			Id:        uuid.New().String(),
			Name:      group.Name,
			ProjectId: project.Id,
			Position:  commandgroupdomain.Order.End(commandGroups),
		}

		for _, cmdId := range group.CommandIds {
			if newCmdId, exists := commandIdsToNewRandomIds[cmdId]; exists {
				newGroup.Commands = append(newGroup.Commands, newIdsToCommand[newCmdId])
			}
		}

		commandGroups = append(commandGroups, newGroup)
	}

	return uc.unitOfWork.Do(func(repositories unitofwork.Repositories) error {
		if err := repositories.Projects.Create(project); err != nil {
			return err
		}

		for _, command := range commands {
			if err := repositories.Commands.Create(&command); err != nil {
				return err
			}
		}

		for _, group := range commandGroups {
			if err := repositories.CommandGroups.Create(&group); err != nil {
				return err
			}
		}

		return nil
	})
}
