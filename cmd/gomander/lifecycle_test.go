package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type windowHooksSpy struct {
	hooks map[events.WindowEventType]func(*application.WindowEvent)
}

func (w *windowHooksSpy) RegisterHook(eventType events.WindowEventType, callback func(*application.WindowEvent)) func() {
	w.hooks[eventType] = callback
	return func() {}
}

func (w *windowHooksSpy) close() *application.WindowEvent {
	closing := application.NewWindowEvent()
	w.hooks[events.Common.WindowClosing](closing)
	return closing
}

type appSpy struct {
	ctx       context.Context
	quitAsked int
}

func (a *appSpy) Context() context.Context {
	return a.ctx
}

func (a *appSpy) Quit() {
	a.quitAsked++
}

func TestCloseWindowByQuitting(t *testing.T) {
	t.Run("Should keep the window open and ask the app to quit, so quitting decides whether it closes", func(t *testing.T) {
		// Arrange
		window := &windowHooksSpy{hooks: map[events.WindowEventType]func(*application.WindowEvent){}}
		app := &appSpy{ctx: context.Background()}

		closeWindowByQuitting(window, app)

		// Act
		closing := window.close()

		// Assert
		assert.True(t, closing.IsCancelled())
		assert.Equal(t, 1, app.quitAsked)
	})

	t.Run("Should let the window close once the app is shutting down", func(t *testing.T) {
		// Arrange
		window := &windowHooksSpy{hooks: map[events.WindowEventType]func(*application.WindowEvent){}}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		app := &appSpy{ctx: ctx}

		closeWindowByQuitting(window, app)

		// Act
		closing := window.close()

		// Assert
		assert.False(t, closing.IsCancelled())
		assert.Equal(t, 0, app.quitAsked)
	})
}
