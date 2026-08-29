// Package openedproject owns which Project the user has open. It is the one
// place that turns the user configuration into a Project, so no operation has
// to know that "the open Project" is stored as a configuration field, nor
// decide for itself what an absent one means.
package openedproject

import (
	"errors"

	configdomain "gomander/internal/config/domain"
	"gomander/internal/execution"
	"gomander/internal/helpers/array"
	projectdomain "gomander/internal/project/domain"
)

// ErrNoneOpen reports that an operation needed the open Project while the user
// had none open: that is this error, and anything else is a storage failure.
var ErrNoneOpen = errors.New("no project is open")

type OpenedProject interface {
	// Get reports having no Project open as ErrNoneOpen, so a nil error
	// guarantees a usable Project; Find is for the callers where having none
	// open is a legitimate outcome.
	Get() (projectdomain.Project, error)
	Find() (project projectdomain.Project, open bool, err error)
	// ExecutionEnvironment is resolved here so that neither the Runner nor
	// the operations that reach it assemble one of their own.
	ExecutionEnvironment() (execution.Environment, error)
	Open(projectId string) error
	Close() error
}

type DefaultOpenedProject struct {
	configRepository  configdomain.Repository
	projectRepository projectdomain.Repository
}

func NewOpenedProject(configRepo configdomain.Repository, projectRepo projectdomain.Repository) *DefaultOpenedProject {
	return &DefaultOpenedProject{
		configRepository:  configRepo,
		projectRepository: projectRepo,
	}
}

func (op *DefaultOpenedProject) Get() (projectdomain.Project, error) {
	config, err := op.configRepository.GetOrCreate()
	if err != nil {
		return projectdomain.Project{}, err
	}

	return op.openedIn(config)
}

func (op *DefaultOpenedProject) Find() (projectdomain.Project, bool, error) {
	config, err := op.configRepository.GetOrCreate()
	if err != nil {
		return projectdomain.Project{}, false, err
	}

	return op.projectRepository.Find(config.LastOpenedProjectId)
}

func (op *DefaultOpenedProject) ExecutionEnvironment() (execution.Environment, error) {
	config, err := op.configRepository.GetOrCreate()
	if err != nil {
		return execution.Environment{}, err
	}

	project, err := op.openedIn(config)
	if err != nil {
		return execution.Environment{}, err
	}

	return execution.Environment{
		Paths: array.Map(config.EnvironmentPaths, func(ep configdomain.EnvironmentPath) string {
			return ep.Path
		}),
		BaseWorkingDirectory: project.WorkingDirectory,
	}, nil
}

func (op *DefaultOpenedProject) openedIn(config *configdomain.Config) (projectdomain.Project, error) {
	project, open, err := op.projectRepository.Find(config.LastOpenedProjectId)
	if err != nil {
		return projectdomain.Project{}, err
	}
	if !open {
		return projectdomain.Project{}, ErrNoneOpen
	}

	return project, nil
}

func (op *DefaultOpenedProject) Open(projectId string) error {
	if _, err := op.projectRepository.Get(projectId); err != nil {
		return err
	}

	return op.remember(projectId)
}

func (op *DefaultOpenedProject) Close() error {
	return op.remember("")
}

func (op *DefaultOpenedProject) remember(projectId string) error {
	config, err := op.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	config.LastOpenedProjectId = projectId

	return op.configRepository.Update(config)
}
