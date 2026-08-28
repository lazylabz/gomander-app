package apptest

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"

	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/dialog"
	"gomander/internal/event"
	"gomander/internal/execution"
	"gomander/internal/unitofwork"
)

type StartedProcess struct {
	Command     commanddomain.Command
	Environment execution.Environment
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
	startFailures     map[string]error
	stopFailures      map[string]error
}

func (r *processRunnerFake) RunCommand(command *commanddomain.Command, environment execution.Environment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.isRunning(command.Id) {
		return nil
	}

	if err := r.startFailures[command.Id]; err != nil {
		return err
	}

	r.startedProcesses = append(r.startedProcesses, StartedProcess{
		Command:     *command,
		Environment: environment,
	})
	r.runningIds = append(r.runningIds, command.Id)

	return nil
}

func (r *processRunnerFake) StopRunningCommand(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isRunning(id) {
		return nil
	}

	if err := r.stopFailures[id]; err != nil {
		return err
	}

	r.stoppedProcessIds = append(r.stoppedProcessIds, id)
	r.runningIds = slices.DeleteFunc(r.runningIds, func(runningId string) bool { return runningId == id })

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

func (r *processRunnerFake) failStart(commandId string, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.startFailures[commandId] = err
}

func (r *processRunnerFake) failStop(commandId string, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.stopFailures[commandId] = err
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
// listens to and the log sink.
type runtimeFake struct {
	mutex         sync.Mutex
	emittedEvents []EmittedEvent
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

func (f *runtimeFake) LogInfo(_ context.Context, _ string)  {}
func (f *runtimeFake) LogDebug(_ context.Context, _ string) {}
func (f *runtimeFake) LogError(_ context.Context, _ string) {}

func (f *runtimeFake) OpenFolderInFileManager(_ string) error { return nil }

func (f *runtimeFake) CloseApp(_ context.Context) {}

// dialogsFake answers the file dialogs in the user's place: with the path they
// would have picked, with an empty path when they cancel, or with the failure
// the desktop toolkit reports when it cannot put a dialog on screen at all.
type dialogsFake struct {
	saveFilePath string
	openFilePath string
	failure      error
}

func (f *dialogsFake) AskForFileToOpen(_ dialog.OpenFileRequest) (string, error) {
	if f.failure != nil {
		return "", f.failure
	}

	return f.openFilePath, nil
}

func (f *dialogsFake) AskWhereToSaveFile(_ dialog.SaveFileRequest) (string, error) {
	if f.failure != nil {
		return "", f.failure
	}

	return f.saveFilePath, nil
}

func (f *dialogsFake) AskForDirectory(_ dialog.PickDirectoryRequest) (string, error) {
	if f.failure != nil {
		return "", f.failure
	}

	return "", nil
}

// fsFacadeFake keeps the files the backend reads and writes in memory.
type fsFacadeFake struct {
	mutex        sync.Mutex
	files        map[string][]byte
	writeFailure error
}

func (f *fsFacadeFake) WriteFile(path string, data []byte, _ os.FileMode) error {
	if err := f.failedWrite(); err != nil {
		return err
	}

	f.put(path, data)

	return nil
}

func (f *fsFacadeFake) failedWrite() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.writeFailure
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

// unitOfWorkWithFailingWrites runs the real Unit of Work and lets a test make
// one write inside it fail, the way storage refuses one partway through an
// operation. Nothing below it is faked: what a test arranges this for is the
// rollback the real transaction performs.
type unitOfWorkWithFailingWrites struct {
	unitOfWork               unitofwork.UnitOfWork
	commandGroupWriteFailure error
}

func (u *unitOfWorkWithFailingWrites) Do(change func(unitofwork.Repositories) error) error {
	return u.unitOfWork.Do(func(repositories unitofwork.Repositories) error {
		if u.commandGroupWriteFailure != nil {
			repositories.CommandGroups = commandGroupRepositoryThatCannotWrite{
				Repository: repositories.CommandGroups,
				failure:    u.commandGroupWriteFailure,
			}
		}

		return change(repositories)
	})
}

type commandGroupRepositoryThatCannotWrite struct {
	commandgroupdomain.Repository
	failure error
}

func (r commandGroupRepositoryThatCannotWrite) Create(*commandgroupdomain.CommandGroup) error {
	return r.failure
}
