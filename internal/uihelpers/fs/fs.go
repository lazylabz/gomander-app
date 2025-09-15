package fs

import (
	"context"
	"net/url"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/facade"
)

type UIFsHelper struct {
	runtime facade.RuntimeFacade
	ctx     context.Context
}

func NewUIFsHelper(runtime facade.RuntimeFacade) *UIFsHelper {
	return &UIFsHelper{
		runtime: runtime,
	}
}

func (h *UIFsHelper) SetContext(ctx context.Context) {
	h.ctx = ctx
}

func (h *UIFsHelper) AskForDirPath() (string, error) {
	return h.runtime.OpenDirectoryDialog(h.ctx, runtime.OpenDialogOptions{})
}

func (h *UIFsHelper) OpenFileFolder(filePath string) {
	cleanPath := filepath.Clean(filePath)
	folderPath := filepath.Dir(cleanPath)

	folderUrl := url.URL{Scheme: "file", Path: folderPath}

	h.runtime.BrowserOpenURL(h.ctx, folderUrl.String())
}
