package uihelpers

import (
	"gomander/internal/helpers"
)

type UiPathHelper struct {
}

func NewUiPathHelper() *UiPathHelper {
	return &UiPathHelper{}
}

func (ph UiPathHelper) GetComputedPath(basePath, path string) string {
	return helpers.GetComputedPath(basePath, path)
}
