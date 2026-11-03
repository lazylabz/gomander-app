package facade

import (
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// RuntimeFacade is the desktop shell seen from the inside: it satisfies the log
// sink, the event sink, the app control and the file manager ports its
// consumers declare, and is the only place any of them reaches Wails.
type RuntimeFacade struct {
	app *application.App
}

func NewRuntimeFacade(app *application.App) RuntimeFacade {
	return RuntimeFacade{app: app}
}

func (r RuntimeFacade) EventsEmit(eventName string, payload interface{}) {
	r.app.Event.Emit(eventName, payload)
}

func (r RuntimeFacade) LogInfo(message string) {
	r.app.Logger.Info(message)
}

func (r RuntimeFacade) LogDebug(message string) {
	r.app.Logger.Debug(message)
}

func (r RuntimeFacade) LogError(message string) {
	r.app.Logger.Error(message)
}

func (r RuntimeFacade) OpenFolder(path string) error {
	return open.Run(path)
}

func (r RuntimeFacade) CloseApp() {
	r.app.Quit()
}
