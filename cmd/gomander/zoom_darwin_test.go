//go:build darwin

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestMacOSApplicationMenu(t *testing.T) {
	// The app menu names itself after the running application, which Wails only
	// knows once an App exists.
	application.New(application.Options{Name: "Gomander"})

	t.Run("Should keep the app and edit menus, or Quit and copy and paste stop answering their shortcuts", func(t *testing.T) {
		// Act
		menu := macOSApplicationMenu((&zoomSpy{}).operations(), func(string) {})

		// Assert
		assert.NotNil(t, menu.FindByRole(application.AppMenu))
		assert.NotNil(t, menu.FindByRole(application.EditMenu))
	})

	t.Run("Should show each zoom item under the shortcut that zooms that way", func(t *testing.T) {
		// Arrange
		spy := &zoomSpy{}
		bindings := macOSZoomKeyBindings("darwin", spy.operations(), func(string) {})
		bindingByAccelerator := make(map[string]func(application.Window), len(bindings))
		for accelerator, binding := range bindings {
			bindingByAccelerator[normalised(accelerator)] = binding
		}

		// Act
		view := macOSApplicationMenu(spy.operations(), func(string) {}).FindByLabel("View").GetSubmenu()

		// Assert
		for label, operation := range map[string]string{"Actual Size": "reset", "Zoom In": "in", "Zoom Out": "out"} {
			item := view.FindByLabel(label)
			require.NotNil(t, item, label)

			binding, bound := bindingByAccelerator[item.GetAccelerator()]
			require.True(t, bound, "%s shows %s, which no key binding answers", label, item.GetAccelerator())

			spy.pressed = nil
			binding(nil)
			assert.Equal(t, []string{operation}, spy.pressed, label)
		}
	})
}

// normalised spells an accelerator the way Wails stores it once parsed, so a
// binding written "CmdOrCtrl+plus" compares equal to the menu's.
func normalised(accelerator string) string {
	return application.NewMenuItem("").SetAccelerator(accelerator).GetAccelerator()
}
