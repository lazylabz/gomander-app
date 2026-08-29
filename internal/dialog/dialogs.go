// Package dialog is how the backend asks the user for a path: where to save a
// file, which file to open, which directory to pick. The requests are written
// in the app's own words, so the desktop toolkit that puts the window on screen
// is named by the adapter in this package and nowhere else.
package dialog

// FileFilter narrows what the user can pick, as the plain strings every desktop
// toolkit understands - a label and a glob such as "*.json".
type FileFilter struct {
	DisplayName string
	Pattern     string
}

type OpenFileRequest struct {
	Title   string
	Filters []FileFilter
}

type SaveFileRequest struct {
	Title                string
	DefaultFilename      string
	CanCreateDirectories bool
}

type PickDirectoryRequest struct {
	Title string
}

// Dialogs answers with the path the user chose, or with an empty path when they
// cancelled. Cancelling is an outcome, not an error.
type Dialogs interface {
	AskForFileToOpen(request OpenFileRequest) (string, error)
	AskWhereToSaveFile(request SaveFileRequest) (string, error)
	AskForDirectory(request PickDirectoryRequest) (string, error)
}
