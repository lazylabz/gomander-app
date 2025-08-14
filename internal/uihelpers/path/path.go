package path

import pathhelpers "gomander/internal/helpers/path"

type UiPathHelper struct {
}

func NewUiPathHelper() *UiPathHelper {
	return &UiPathHelper{}
}

func (ph UiPathHelper) GetComputedPath(basePath, path string) string {
	return pathhelpers.GetComputedPath(basePath, path)
}
