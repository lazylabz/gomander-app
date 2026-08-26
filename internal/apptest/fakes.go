package apptest

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	commanddomain "gomander/internal/command/domain"
	"gomander/internal/event"
)

type StartedProcess struct {
	Command              commanddomain.Command
	EnvironmentPaths     []string
	BaseWorkingDirectory string
}

type EmittedEvent struct {
	Name    event.Event
	Payload interface{}
}

// processRunnerFake stands in for spawning processes, the lowest edge the
// backend reaches. It answers like the real runner does - starting an already
// running command is a no-op, stopping one that is not running is not an error
// - but records what it was asked to do instead of touching the machine.
//
// It replaces the whole Runner rather than exec alone, because the Runner has
// no seam below itself: what it does with a request - computing the working
// directory, injecting the environment paths, and emitting the process and log
// events while the process lives - stays covered by its own tests against real
// processes, not by this.
type processRunnerFake struct {
	mutex             sync.Mutex
	startedProcesses  []StartedProcess
	stoppedProcessIds []string
	runningIds        []string
}

func (r *processRunnerFake) RunCommand(command *commanddomain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isRunning(command.Id) {
		return nil
	}

	r.startedProcesses = append(r.startedProcesses, StartedProcess{
		Command:              *command,
		EnvironmentPaths:     environmentPaths,
		BaseWorkingDirectory: baseWorkingDirectory,
	})
	r.runningIds = append(r.runningIds, command.Id)

	return nil
}

func (r *processRunnerFake) RunCommands(commands []commanddomain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	for i := range commands {
		if err := r.RunCommand(&commands[i], environmentPaths, baseWorkingDirectory); err != nil {
			return err
		}
	}

	return nil
}

func (r *processRunnerFake) StopRunningCommand(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isRunning(id) {
		return nil
	}

	r.stoppedProcessIds = append(r.stoppedProcessIds, id)
	r.runningIds = slices.DeleteFunc(r.runningIds, func(runningId string) bool { return runningId == id })

	return nil
}

func (r *processRunnerFake) StopRunningCommands(commands []commanddomain.Command) error {
	for _, command := range commands {
		if err := r.StopRunningCommand(command.Id); err != nil {
			return err
		}
	}

	return nil
}

func (r *processRunnerFake) StopAllRunningCommands() []error {
	r.mutex.Lock()
	running := append([]string(nil), r.runningIds...)
	r.mutex.Unlock()

	errs := make([]error, 0)
	for _, id := range running {
		if err := r.StopRunningCommand(id); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (r *processRunnerFake) GetRunningCommandIds() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return append([]string(nil), r.runningIds...)
}

func (r *processRunnerFake) isRunning(id string) bool {
	return slices.Contains(r.runningIds, id)
}

func (r *processRunnerFake) started() []StartedProcess {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return append([]StartedProcess(nil), r.startedProcesses...)
}

func (r *processRunnerFake) stopped() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return append([]string(nil), r.stoppedProcessIds...)
}

// runtimeFake stands in for the desktop runtime: the event emitter the frontend
// listens to, the log sink, and the file dialogs the user answers.
type runtimeFake struct {
	mutex              sync.Mutex
	emittedEvents      []EmittedEvent
	saveFileDialogPath string
	openFileDialogPath string
}

func (f *runtimeFake) EventsEmit(_ context.Context, eventName string, payload interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.emittedEvents = append(f.emittedEvents, EmittedEvent{Name: event.Event(eventName), Payload: payload})
}

func (f *runtimeFake) emitted() []EmittedEvent {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return append([]EmittedEvent(nil), f.emittedEvents...)
}

func (f *runtimeFake) SaveFileDialog(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
	return f.saveFileDialogPath, nil
}

func (f *runtimeFake) OpenFileDialog(_ context.Context, _ runtime.OpenDialogOptions) (string, error) {
	return f.openFileDialogPath, nil
}

func (f *runtimeFake) OpenDirectoryDialog(_ context.Context, _ runtime.OpenDialogOptions) (string, error) {
	return "", nil
}

func (f *runtimeFake) LogInfo(_ context.Context, _ string)  {}
func (f *runtimeFake) LogDebug(_ context.Context, _ string) {}
func (f *runtimeFake) LogError(_ context.Context, _ string) {}

func (f *runtimeFake) OpenFolderInFileManager(_ string) error { return nil }

func (f *runtimeFake) CloseApp(_ context.Context) {}

// fsFacadeFake keeps the files the backend reads and writes in memory.
type fsFacadeFake struct {
	mutex sync.Mutex
	files map[string][]byte
}

func (f *fsFacadeFake) WriteFile(path string, data []byte, _ os.FileMode) error {
	f.put(path, data)

	return nil
}

func (f *fsFacadeFake) put(path string, data []byte) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.files[path] = data
}

func (f *fsFacadeFake) ReadFile(path string) ([]byte, error) {
	data, exists := f.file(path)
	if !exists {
		return nil, fmt.Errorf("read %s: %w", path, os.ErrNotExist)
	}

	return data, nil
}

func (f *fsFacadeFake) file(path string) ([]byte, bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	data, exists := f.files[path]

	return data, exists
}
