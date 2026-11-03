package dialog

import (
	"github.com/wailsapp/wails/v3/pkg/application"

	"gomander/internal/helpers/array"
)

type WailsDialogs struct {
	dialogs *application.DialogManager
}

func NewWailsDialogs(dialogs *application.DialogManager) *WailsDialogs {
	return &WailsDialogs{dialogs: dialogs}
}

func (d *WailsDialogs) AskForFileToOpen(request OpenFileRequest) (string, error) {
	return Answer(d.dialogs.OpenFileWithOptions(OpenFileDialogOptions(request)).PromptForSingleSelection())
}

func (d *WailsDialogs) AskWhereToSaveFile(request SaveFileRequest) (string, error) {
	return Answer(d.dialogs.SaveFileWithOptions(SaveFileDialogOptions(request)).PromptForSingleSelection())
}

func (d *WailsDialogs) AskForDirectory(request PickDirectoryRequest) (string, error) {
	return Answer(d.dialogs.OpenFileWithOptions(PickDirectoryDialogOptions(request)).PromptForSingleSelection())
}

// Wails reports a cancelled Windows dialog with an error from an internal
// package, so only its message can be matched. macOS and Linux answer "".
const cancelledOnWindows = "cancelled by user"

func Answer(path string, err error) (string, error) {
	if err != nil && err.Error() == cancelledOnWindows {
		return "", nil
	}
	return path, err
}

func OpenFileDialogOptions(request OpenFileRequest) *application.OpenFileDialogOptions {
	return &application.OpenFileDialogOptions{
		Title:          request.Title,
		CanChooseFiles: true,
		Filters: array.Map(request.Filters, func(filter FileFilter) application.FileFilter {
			return application.FileFilter{DisplayName: filter.DisplayName, Pattern: filter.Pattern}
		}),
	}
}

func SaveFileDialogOptions(request SaveFileRequest) *application.SaveFileDialogOptions {
	return &application.SaveFileDialogOptions{
		Title:                request.Title,
		Filename:             request.DefaultFilename,
		CanCreateDirectories: request.CanCreateDirectories,
	}
}

// PickDirectoryDialogOptions keeps CanCreateDirectories on, as Wails v2's
// directory dialog did: v3 turns it off unless asked.
func PickDirectoryDialogOptions(request PickDirectoryRequest) *application.OpenFileDialogOptions {
	return &application.OpenFileDialogOptions{
		Title:                request.Title,
		CanChooseDirectories: true,
		CanCreateDirectories: true,
	}
}
