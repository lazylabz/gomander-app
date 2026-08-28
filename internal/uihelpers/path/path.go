package path

import (
	commanddomain "gomander/internal/command/domain"
)

type UiPathHelper struct {
}

func NewUiPathHelper() *UiPathHelper {
	return &UiPathHelper{}
}

// GetComputedPath previews the directory a Command with this working directory
// would run in, asking the same domain rule the Runner asks so the preview and
// the Process cannot drift apart.
func (ph UiPathHelper) GetComputedPath(basePath, path string) string {
	return commanddomain.Command{WorkingDirectory: path}.ResolveWorkingDirectory(basePath)
}
