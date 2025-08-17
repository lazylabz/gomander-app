package app

import (
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type FsFacade interface {
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) ([]byte, error)
}

type DefaultFsFacade struct{}

func (d DefaultFsFacade) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (d DefaultFsFacade) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

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
