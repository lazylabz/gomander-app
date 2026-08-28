package fs

import (
	"path/filepath"

	"gomander/internal/dialog"
	"gomander/internal/facade"
)

type UIFsHelper struct {
	dialogs dialog.Dialogs
	runtime facade.RuntimeFacade
}

func NewUIFsHelper(dialogs dialog.Dialogs, runtime facade.RuntimeFacade) *UIFsHelper {
	return &UIFsHelper{
		dialogs: dialogs,
		runtime: runtime,
	}
}

func (h *UIFsHelper) AskForDirPath() (string, error) {
	return h.dialogs.AskForDirectory(dialog.PickDirectoryRequest{})
}

func (h *UIFsHelper) OpenFileFolder(filePath string) error {
	cleanPath := filepath.Clean(filePath)
	folderPath := filepath.Dir(cleanPath)

	return h.runtime.OpenFolderInFileManager(folderPath)
}
