package usecases

import (
	"gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/runner"
)

type RunCommandGroup interface {
	Execute(commandGroupId string) error
}

type DefaultRunCommandGroup struct {
	configRepository       configdomain.Repository
	commandRepository      domain.Repository
	commandGroupRepository commandgroupdomain.Repository
	projectRepository      projectdomain.Repository
	commandRunner          runner.Runner
}

func NewRunCommandGroup(
	configRepo configdomain.Repository,
	commandRepo domain.Repository,
	commandGroupRepo commandgroupdomain.Repository,
	projectRepo projectdomain.Repository,
	runner runner.Runner,
) *DefaultRunCommandGroup {
	return &DefaultRunCommandGroup{
		configRepository:       configRepo,
		commandRepository:      commandRepo,
		commandGroupRepository: commandGroupRepo,
		projectRepository:      projectRepo,
		commandRunner:          runner,
	}
}

func (uc *DefaultRunCommandGroup) Execute(commandGroupId string) error {
	cmdGroup, err := uc.commandGroupRepository.Get(commandGroupId)
	if err != nil {
		return err
	}

	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	currentProject, err := uc.projectRepository.Get(cmdGroup.ProjectId)
	if err != nil {
		return err
	}

	environmentPathsStrings := array.Map(userConfig.EnvironmentPaths, func(ep configdomain.EnvironmentPath) string {
		return ep.Path
	})

	err = uc.commandRunner.RunCommands(cmdGroup.Commands, environmentPathsStrings, currentProject.WorkingDirectory)
	if err != nil {
		return err
	}

	return nil
}
