package facade

import (
	"context"

	"github.com/skratchdot/open-golang/open"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// DefaultRuntimeFacade is the desktop shell seen from the inside: it satisfies
// the log sink, the event sink, the app control and the file manager ports its
// consumers declare, and is the only place any of them reaches Wails.
type DefaultRuntimeFacade struct{}

func (d DefaultRuntimeFacade) EventsEmit(ctx context.Context, eventName string, optionalData interface{}) {
	runtime.EventsEmit(ctx, eventName, optionalData)
}

func (d DefaultRuntimeFacade) LogInfo(ctx context.Context, message string) {
	runtime.LogInfo(ctx, message)
}

func (d DefaultRuntimeFacade) LogDebug(ctx context.Context, message string) {
	runtime.LogDebug(ctx, message)
}

func (d DefaultRuntimeFacade) LogError(ctx context.Context, message string) {
	runtime.LogError(ctx, message)
}

func (d DefaultRuntimeFacade) OpenFolder(path string) error {
	return open.Run(path)
}

func (d DefaultRuntimeFacade) CloseApp(ctx context.Context) {
	runtime.Quit(ctx)
}
