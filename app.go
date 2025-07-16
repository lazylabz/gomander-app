package main

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ExecCommand(arg string) error {
	chunks := strings.Fields(arg)

	cmd := exec.Command(chunks[0], chunks[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	}

	if err := cmd.Start(); err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	}

	// Stream stdout
	go a.streamOutput(stdout, "stdout")
	// Stream stderr
	go a.streamOutput(stderr, "stderr")

	// Optional: Wait in background
	go func() {
		err := cmd.Wait()
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			return
		}
		runtime.EventsEmit(a.ctx, "processFinished", nil)
	}()

	return nil
}

func (a *App) streamOutput(pipeReader io.ReadCloser, streamType string) {
	scanner := bufio.NewScanner(pipeReader)
	for scanner.Scan() {
		line := scanner.Text()
		runtime.LogDebug(a.ctx, line)
		runtime.EventsEmit(a.ctx, "newLog", map[string]string{
			"type": streamType,
			"line": line,
		})
	}
}
