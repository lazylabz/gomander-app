package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gomander/internal/releases/infrastructure"
)

type desktopShellSpy struct {
	closedWith context.Context
}

func (s *desktopShellSpy) CloseApp(ctx context.Context) {
	s.closedWith = ctx
}

func TestShellAppControl_CloseApp(t *testing.T) {
	t.Run("Should close the app with the context it was built with", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		shell := &desktopShellSpy{}

		appControl := infrastructure.NewShellAppControl(ctx, shell)

		// Act
		appControl.CloseApp()

		// Assert
		assert.Equal(t, ctx, shell.closedWith)
	})
}
