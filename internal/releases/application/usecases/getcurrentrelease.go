package usecases

import (
	"github.com/Masterminds/semver"

	"gomander/internal/releases"
)

type GetCurrentRelease struct {
}

func NewGetCurrentRelease() *GetCurrentRelease {
	return &GetCurrentRelease{}
}

func (uc *GetCurrentRelease) Execute() string {
	return semver.MustParse(releases.CurrentRelease).String()
}
