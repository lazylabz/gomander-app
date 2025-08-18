package facade

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RuntimeFacade interface {
	SaveFileDialog(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
	OpenFileDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
}

type DefaultRuntimeFacade struct{}

func (d DefaultRuntimeFacade) SaveFileDialog(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(ctx, dialogOptions)
}

func (d DefaultRuntimeFacade) OpenFileDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(ctx, dialogOptions)
}
