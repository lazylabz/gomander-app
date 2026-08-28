package infrastructure

import (
	"context"
)

// DesktopShell is the app's own quit, as the desktop shell offers it: bound to
// the context the shell only produces once the app has started.
type DesktopShell interface {
	CloseApp(ctx context.Context)
}

// ShellAppControl holds that context, so nothing in the application layer has
// to carry the shell's request-scoped plumbing to ask the app to quit.
type ShellAppControl struct {
	ctx   context.Context
	shell DesktopShell
}

func NewShellAppControl(ctx context.Context, shell DesktopShell) *ShellAppControl {
	return &ShellAppControl{
		ctx:   ctx,
		shell: shell,
	}
}

func (c *ShellAppControl) CloseApp() {
	c.shell.CloseApp(c.ctx)
}
