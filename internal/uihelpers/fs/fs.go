package fs

import (
	"path/filepath"

	"gomander/internal/dialog"
)

// FileManager reveals a path in whatever the operating system uses to browse
// files.
type FileManager interface {
	OpenFolder(path string) error
}

type UIFsHelper struct {
	dialogs     dialog.Dialogs
	fileManager FileManager
}

func NewUIFsHelper(dialogs dialog.Dialogs, fileManager FileManager) *UIFsHelper {
	return &UIFsHelper{
		dialogs:     dialogs,
		fileManager: fileManager,
	}
}

func (h *UIFsHelper) AskForDirPath() (string, error) {
	return h.dialogs.AskForDirectory(dialog.PickDirectoryRequest{})
}

func (h *UIFsHelper) OpenFileFolder(filePath string) error {
	cleanPath := filepath.Clean(filePath)
	folderPath := filepath.Dir(cleanPath)

	return h.fileManager.OpenFolder(folderPath)
}
