package facade

import (
	"context"
	"github.com/skratchdot/open-golang/open"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type RuntimeFacade interface {
	SaveFileDialog(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
	OpenFileDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	OpenDirectoryDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	EventsEmit(ctx context.Context, eventName string, payload interface{})
	LogInfo(ctx context.Context, message string)
	LogDebug(ctx context.Context, message string)
	LogError(ctx context.Context, message string)
	OpenFolderInFileManager(path string) error
}

type DefaultRuntimeFacade struct{}

func (d DefaultRuntimeFacade) SaveFileDialog(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
	return runtime.SaveFileDialog(ctx, dialogOptions)
}

func (d DefaultRuntimeFacade) OpenFileDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenFileDialog(ctx, dialogOptions)
}

func (d DefaultRuntimeFacade) EventsEmit(ctx context.Context, eventName string, optionalData interface{}) {
	runtime.EventsEmit(ctx, eventName, optionalData)
}

func (d DefaultRuntimeFacade) LogInfo(ctx context.Context, message string) {
	runtime.LogInfo(ctx, message)
}

func (d DefaultRuntimeFacade) LogDebug(ctx context.Context, message string) {
	runtime.LogDebug(ctx, message)
}

func (d DefaultRuntimeFacade) LogError(ctx context.Context, message string) {
	runtime.LogError(ctx, message)
}

func (d DefaultRuntimeFacade) OpenDirectoryDialog(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error) {
	return runtime.OpenDirectoryDialog(ctx, dialogOptions)
}

func (d DefaultRuntimeFacade) OpenFolderInFileManager(path string) error {
	return open.Run(path)
}
