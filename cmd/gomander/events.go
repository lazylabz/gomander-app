package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"

	"gomander/internal/event"
)

// Registering every Event lets Wails type its payload in the generated
// frontend bindings and check it on emit.
func init() {
	application.RegisterEvent[string](string(event.ProcessStarted))
	application.RegisterEvent[string](string(event.ProcessFinished))
	application.RegisterEvent[map[string]string](string(event.NewLogEntry))
	application.RegisterEvent[string](string(event.CommandGroupDeleted))
	application.RegisterEvent[string](string(event.CommandErrorDetected))
}
