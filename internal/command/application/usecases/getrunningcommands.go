package usecases

import "gomander/internal/runner"

type GetRunningCommandIds struct {
	runner runner.Runner
}

func NewGetRunningCommandIds(runner runner.Runner) *GetRunningCommandIds {
	return &GetRunningCommandIds{
		runner: runner,
	}
}

func (uc *GetRunningCommandIds) Execute() []string {
	return uc.runner.GetRunningCommandIds()
}
