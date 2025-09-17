package usecases

import "gomander/internal/runner"

type GetRunningCommandIds interface {
	Execute() []string
}

type DefaultGetRunningCommandIds struct {
	runner runner.Runner
}

func NewGetRunningCommandIds(runner runner.Runner) *DefaultGetRunningCommandIds {
	return &DefaultGetRunningCommandIds{
		runner: runner,
	}
}

func (uc *DefaultGetRunningCommandIds) Execute() []string {
	return uc.runner.GetRunningCommandIds()
}
