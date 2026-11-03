package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type windowHooks interface {
	RegisterHook(eventType events.WindowEventType, callback func(event *application.WindowEvent)) func()
}

type quitter interface {
	Context() context.Context
	Quit()
}

// The close button goes through Quit so ShouldQuit decides for it as for Cmd+Q;
// the close that quitting itself triggers comes after the app context is cancelled.
func closeWindowByQuitting(window windowHooks, app quitter) {
	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if app.Context().Err() != nil {
			return
		}
		event.Cancel()
		app.Quit()
	})
}
