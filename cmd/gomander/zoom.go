package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"

	"gomander/internal/usecases"
)

const (
	zoomInAccelerator    = "CmdOrCtrl+plus"
	zoomOutAccelerator   = "CmdOrCtrl+-"
	resetZoomAccelerator = "CmdOrCtrl+0"
)

type zoomOperations struct {
	in    func() error
	out   func() error
	reset func() error
}

// The registry is read when a shortcut fires, not now: the window the
// shortcuts are bound to has to exist before the use cases that zoom it.
func zoomOperationsOf(useCases *usecases.Registry) zoomOperations {
	return zoomOperations{
		in:    func() error { return useCases.ZoomIn.Execute() },
		out:   func() error { return useCases.ZoomOut.Execute() },
		reset: func() error { return useCases.ResetZoom.Execute() },
	}
}

// Only macOS: Wails names the zoom keys "oem_plus"/"oem_minus" on Windows and
// not at all on Linux, so the frontend catches them there. These run before
// the menu sees the press, and macOS reports zoom-in as "=", shifted or not.
func macOSZoomKeyBindings(goos string, zoom zoomOperations, logError func(string)) map[string]func(application.Window) {
	if goos != "darwin" {
		return nil
	}

	zoomIn := onKey(zoom.in, logError)

	return map[string]func(application.Window){
		"CmdOrCtrl+=":        zoomIn,
		"CmdOrCtrl+Shift+=":  zoomIn,
		zoomInAccelerator:    zoomIn,
		zoomOutAccelerator:   onKey(zoom.out, logError),
		resetZoomAccelerator: onKey(zoom.reset, logError),
	}
}

// macOSApplicationMenu replaces the default one Wails gives macOS, whose zoom
// items zoom the webview without storing the Zoom Level. The app and edit menus
// stay, or Quit and copy and paste would stop answering their shortcuts.
func macOSApplicationMenu(zoom zoomOperations, logError func(string)) *application.Menu {
	menu := application.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)

	view := menu.AddSubmenu("View")
	view.Add("Actual Size").SetAccelerator(resetZoomAccelerator).OnClick(onClick(zoom.reset, logError))
	view.Add("Zoom In").SetAccelerator(zoomInAccelerator).OnClick(onClick(zoom.in, logError))
	view.Add("Zoom Out").SetAccelerator(zoomOutAccelerator).OnClick(onClick(zoom.out, logError))
	view.AddSeparator()
	view.AddRole(application.ToggleFullscreen)

	menu.AddRole(application.WindowMenu)

	return menu
}

func onKey(operation func() error, logError func(string)) func(application.Window) {
	return func(application.Window) {
		run(operation, logError)
	}
}

func onClick(operation func() error, logError func(string)) func(*application.Context) {
	return func(*application.Context) {
		run(operation, logError)
	}
}

func run(operation func() error, logError func(string)) {
	if err := operation(); err != nil {
		logError(err.Error())
	}
}
