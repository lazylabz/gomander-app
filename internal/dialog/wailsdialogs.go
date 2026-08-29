package dialog

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gomander/internal/helpers/array"
)

type WailsDialogs struct {
	ctx context.Context
}

func NewWailsDialogs() *WailsDialogs {
	return &WailsDialogs{}
}

// SetWailsDialogsContext hands the adapter the context Wails only produces once
// the app starts, after the objects bound to the frontend have been built. It
// is a package function rather than a method so that the adapter's method set
// stays exactly the Dialogs port.
func SetWailsDialogsContext(d *WailsDialogs, ctx context.Context) {
	d.ctx = ctx
}

func (d *WailsDialogs) AskForFileToOpen(request OpenFileRequest) (string, error) {
	return runtime.OpenFileDialog(d.ctx, runtime.OpenDialogOptions{
		Title: request.Title,
		Filters: array.Map(request.Filters, func(filter FileFilter) runtime.FileFilter {
			return runtime.FileFilter{DisplayName: filter.DisplayName, Pattern: filter.Pattern}
		}),
	})
}

func (d *WailsDialogs) AskWhereToSaveFile(request SaveFileRequest) (string, error) {
	return runtime.SaveFileDialog(d.ctx, runtime.SaveDialogOptions{
		Title:                request.Title,
		DefaultFilename:      request.DefaultFilename,
		CanCreateDirectories: request.CanCreateDirectories,
	})
}

func (d *WailsDialogs) AskForDirectory(request PickDirectoryRequest) (string, error) {
	return runtime.OpenDirectoryDialog(d.ctx, runtime.OpenDialogOptions{
		Title: request.Title,
	})
}
