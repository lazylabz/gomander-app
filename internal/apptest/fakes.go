package apptest

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	commandhandlers "gomander/internal/command/application/handlers"
	commanddomain "gomander/internal/command/domain"
	commandevent "gomander/internal/command/domain/event"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	"gomander/internal/dialog"
	"gomander/internal/event"
	"gomander/internal/eventbus"
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
	domainEvents      eventbus.EventBus
	// errorReported holds the running Commands whose Process already reported
	// an Error Pattern; a new run starts without one.
	errorReported map[string]bool
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
	delete(r.errorReported, command.Id)

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

// print reports the first line of a Process that matches its Command's Error
// Patterns, as the real runner does.
func (r *processRunnerFake) print(id string, line string) {
	r.mutex.Lock()
	command, running := r.runningCommand(id)
	report := running && !r.errorReported[id] && command.MatchesErrorPattern(line)
	if report {
		r.errorReported[id] = true
	}
	r.mutex.Unlock()

	if !report {
		return
	}

	// The real runner logs what the handlers report and carries on
	_ = r.domainEvents.PublishSync(commandevent.ErrorPatternMatchedEvent{CommandId: id, Line: line})
}

func (r *processRunnerFake) exit(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.runningIds = slices.DeleteFunc(r.runningIds, func(runningId string) bool { return runningId == id })
}

func (r *processRunnerFake) runningCommand(id string) (commanddomain.Command, bool) {
	if !r.isRunning(id) {
		return commanddomain.Command{}, false
	}

	for i := len(r.startedProcesses) - 1; i >= 0; i-- {
		if r.startedProcesses[i].Command.Id == id {
			return r.startedProcesses[i].Command, true
		}
	}

	return commanddomain.Command{}, false
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

// eventSinkFake stands in for the desktop runtime's event sink: what the
// frontend would be listening to, recorded so a test can assert on it.
type eventSinkFake struct {
	mutex         sync.Mutex
	emittedEvents []EmittedEvent
}

func (f *eventSinkFake) EventsEmit(eventName string, payload interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.emittedEvents = append(f.emittedEvents, EmittedEvent{Name: event.Event(eventName), Payload: payload})
}

func (f *eventSinkFake) emitted() []EmittedEvent {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return append([]EmittedEvent(nil), f.emittedEvents...)
}

// logSinkFake swallows the log lines the desktop runtime would have written:
// nothing a test asserts on goes through the logger.
type logSinkFake struct{}

func (f *logSinkFake) LogInfo(_ string)  {}
func (f *logSinkFake) LogDebug(_ string) {}
func (f *logSinkFake) LogError(_ string) {}

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

// releaseFeedFake stands in for the feed the app's own releases are published
// to: the version a test has published, or the failure reaching it reports.
type releaseFeedFake struct {
	latestRelease string
	failure       error
}

func (f *releaseFeedFake) GetLatestRelease() (string, error) {
	if f.failure != nil {
		return "", f.failure
	}

	return f.latestRelease, nil
}

// releaseDownloaderFake downloads nothing: it records the versions it was asked
// for and answers with the path the binary would have landed at.
type releaseDownloaderFake struct {
	downloadedVersions []string
}

func (f *releaseDownloaderFake) Download(version string) (string, error) {
	f.downloadedVersions = append(f.downloadedVersions, version)

	return filepath.Join(os.TempDir(), "gomander-"+version), nil
}

// releaseInstallerFake stands in for handing a downloaded binary to the
// operating system, which is the last thing that happens before the app quits.
type releaseInstallerFake struct {
	installedBinaries []string
	failure           error
}

func (f *releaseInstallerFake) Install(binaryPath string) error {
	if f.failure != nil {
		return f.failure
	}

	f.installedBinaries = append(f.installedBinaries, binaryPath)

	return nil
}

// appControlFake records the quit an install asks for instead of closing the
// window the test is not running in.
type appControlFake struct {
	closed bool
}

func (f *appControlFake) CloseApp() {
	f.closed = true
}

// windowFake stands in for the desktop window, keeping the zoom it was last
// set to.
type windowFake struct {
	mutex     sync.Mutex
	zoomLevel float64
	focused   bool
}

func (f *windowFake) IsFocused() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.focused
}

func (f *windowFake) focus(focused bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.focused = focused
}

func (f *windowFake) SetZoom(level float64) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.zoomLevel = level
}

func (f *windowFake) zoom() float64 {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.zoomLevel
}

// storageFailures is what a test has made storage refuse. The repositories
// share it by pointer, so arranging a failure reaches every path a Command
// Group can be written through - directly, or through a Unit of Work.
type storageFailures struct {
	commandGroupWrite error
	zoomLevelWrite    error
}

// unitOfWorkWithFailingWrites runs the real Unit of Work and lets a test make
// a write inside it fail, the way storage refuses one partway through an
// operation. Nothing below it is faked: what a test arranges this for is the
// rollback the real transaction performs.
type unitOfWorkWithFailingWrites struct {
	unitOfWork unitofwork.UnitOfWork
	failures   *storageFailures
}

func (u *unitOfWorkWithFailingWrites) Do(change func(unitofwork.Repositories) error) error {
	return u.unitOfWork.Do(func(repositories unitofwork.Repositories) error {
		repositories.CommandGroups = commandGroupRepositoryThatCanFail{
			Repository: repositories.CommandGroups,
			failures:   u.failures,
		}

		return change(repositories)
	})
}

// commandGroupRepositoryThatCanFail refuses every write a test has arranged a
// failure for, and answers like the repository it wraps otherwise.
type commandGroupRepositoryThatCanFail struct {
	commandgroupdomain.Repository
	failures *storageFailures
}

func (r commandGroupRepositoryThatCanFail) Create(commandGroup *commandgroupdomain.CommandGroup) error {
	if r.failures.commandGroupWrite != nil {
		return r.failures.commandGroupWrite
	}

	return r.Repository.Create(commandGroup)
}

func (r commandGroupRepositoryThatCanFail) Update(commandGroup *commandgroupdomain.CommandGroup) error {
	if r.failures.commandGroupWrite != nil {
		return r.failures.commandGroupWrite
	}

	return r.Repository.Update(commandGroup)
}

func (r commandGroupRepositoryThatCanFail) Delete(commandGroupId string) error {
	if r.failures.commandGroupWrite != nil {
		return r.failures.commandGroupWrite
	}

	return r.Repository.Delete(commandGroupId)
}

type SentNotification struct {
	Title string
	Body  string
}

// notifierFake stands in for the operating system's notification centre.
type notifierFake struct {
	mutex   sync.Mutex
	sent    []SentNotification
	failure error
}

func (f *notifierFake) Notify(notification commandhandlers.Notification) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.failure != nil {
		return f.failure
	}

	f.sent = append(f.sent, SentNotification{Title: notification.Title, Body: notification.Body})

	return nil
}

func (f *notifierFake) fail(err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.failure = err
}

func (f *notifierFake) delivered() []SentNotification {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return append([]SentNotification(nil), f.sent...)
}
