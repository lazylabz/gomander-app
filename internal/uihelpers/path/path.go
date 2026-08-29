package path

import (
	commanddomain "gomander/internal/command/domain"
)

type UiPathHelper struct {
}

func NewUiPathHelper() *UiPathHelper {
	return &UiPathHelper{}
}

// Asks the same domain rule the Runner asks, so the path previewed here and the
// one the Process runs in cannot drift apart.
func (ph UiPathHelper) GetComputedPath(basePath, path string) string {
	return commanddomain.Command{WorkingDirectory: path}.ResolveWorkingDirectory(basePath)
}
