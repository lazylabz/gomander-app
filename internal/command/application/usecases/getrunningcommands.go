package usecases

import "gomander/internal/execution"

type GetRunningCommandIds struct {
	runner execution.Runner
}

func NewGetRunningCommandIds(runner execution.Runner) *GetRunningCommandIds {
	return &GetRunningCommandIds{
		runner: runner,
	}
}

func (uc *GetRunningCommandIds) Execute() []string {
	return uc.runner.GetRunningCommandIds()
}
