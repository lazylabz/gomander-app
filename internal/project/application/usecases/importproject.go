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

func (uc *ImportProject) Execute(blueprint projectdomain.Blueprint, name, workingDirectory string) error {
	project := projectdomain.Project{
		Id:               uuid.New().String(),
		Name:             name,
		WorkingDirectory: workingDirectory,
	}

	commands := make([]domain.Command, 0, len(blueprint.Commands))
	commandGroups := make([]commandgroupdomain.CommandGroupWithCommandIds, 0, len(blueprint.CommandGroups))

	// The Ids a Blueprint carries only name its own Commands; the ones stored
	// are new, so a Group has to be told which is which.
	commandsByBlueprintId := make(map[string]domain.Command, len(blueprint.Commands))

	for _, cmd := range blueprint.Commands {
		newCommand := domain.Command{
			Id:               uuid.New().String(),
			Name:             cmd.Name,
			Command:          cmd.Command,
			WorkingDirectory: cmd.WorkingDirectory,
			ProjectId:        project.Id,
			Position:         domain.Order.End(commands),
		}

		commands = append(commands, newCommand)
		commandsByBlueprintId[cmd.Id] = newCommand
	}

	for _, group := range blueprint.CommandGroups {
		newGroup := commandgroupdomain.CommandGroupWithCommandIds{
			Id:        uuid.New().String(),
			Name:      group.Name,
			ProjectId: project.Id,
			Position:  commandgroupdomain.Order.End(commandGroups),
		}

		for _, cmdId := range group.CommandIds {
			if command, exists := commandsByBlueprintId[cmdId]; exists {
				newGroup.CommandIds = append(newGroup.CommandIds, command.Id)
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
			if err := repositories.CommandGroups.CreateWithCommandIds(&group); err != nil {
				return err
			}
		}

		return nil
	})
}
