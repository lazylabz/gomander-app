package fs

import (
	"context"

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

func (h *UIFsHelper) GetDirPath() (string, error) {
	return h.runtime.OpenDirectoryDialog(h.ctx, runtime.OpenDialogOptions{})
}
