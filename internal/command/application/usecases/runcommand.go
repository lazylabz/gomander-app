package usecases

import (
	"gomander/internal/command/domain"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
	"gomander/internal/runner"
)

type RunCommand interface {
	Execute(commandId string) error
}

type DefaultRunCommand struct {
	configRepository  configdomain.Repository
	commandRepository domain.Repository
	projectRepository projectdomain.Repository
	commandRunner     runner.Runner
}

func NewRunCommand(
	configRepo configdomain.Repository,
	commandRepo domain.Repository,
	projectRepo projectdomain.Repository,
	runner runner.Runner,
) *DefaultRunCommand {
	return &DefaultRunCommand{
		configRepository:  configRepo,
		commandRepository: commandRepo,
		projectRepository: projectRepo,
		commandRunner:     runner,
	}
}

func (uc *DefaultRunCommand) Execute(commandId string) error {
	cmd, err := uc.commandRepository.Get(commandId)
	if err != nil {
		return err
	}

	userConfig, err := uc.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	currentProject, err := uc.projectRepository.Get(cmd.ProjectId)
	if err != nil {
		return err
	}

	environmentPathsStrings := array.Map(userConfig.EnvironmentPaths, func(ep configdomain.EnvironmentPath) string {
		return ep.Path
	})

	err = uc.commandRunner.RunCommand(cmd, environmentPathsStrings, currentProject.WorkingDirectory)
	if err != nil {
		return err
	}

	return nil
}
