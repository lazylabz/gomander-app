package facade

import "github.com/wailsapp/wails/v3/pkg/application"

type WindowFacade struct {
	window *application.WebviewWindow
}

func NewWindowFacade(window *application.WebviewWindow) WindowFacade {
	return WindowFacade{window: window}
}

func (w WindowFacade) IsFocused() bool {
	return w.window.IsFocused()
}

func (w WindowFacade) SetZoom(level float64) {
	w.window.SetZoom(level)
}
