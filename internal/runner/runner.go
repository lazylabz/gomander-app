package runner

import (
	"bufio"
	"io"
	"sync"
	"time"

	"gomander/internal/command/domain"
	"gomander/internal/event"
	"gomander/internal/execution"
)

// pipeDrainGrace is how long the scanners get to drain what the process left
// buffered before the runner stops reading.
const pipeDrainGrace = 2 * time.Second

var ExpectedTerminationLogs = []string{
	"signal: terminated",
	"signal: interrupt",
	"signal: killed",
	"exit status 143",
	"exit status 137",
	"exit status 130",
	"wait: no child processes",
}

type runningCommand struct {
	process commandProcess
	wg      *sync.WaitGroup
	// exited is closed by the goroutine that owns process.Wait, so stopping a
	// Command can await the exit without calling Wait a second time.
	exited chan struct{}
}

// Logger is where the runner reports what a Command's process is doing.
type Logger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
}

// EventEmitter carries a Process's life and output to whoever is watching it.
type EventEmitter interface {
	EmitEvent(event event.Event, payload interface{})
}

type DefaultRunner struct {
	runningCommands map[string]runningCommand
	eventEmitter    EventEmitter
	logger          Logger
	config          Config
	mutex           sync.Mutex
}

func NewDefaultRunner(logger Logger, emitter EventEmitter, config Config) *DefaultRunner {
	return &DefaultRunner{
		runningCommands: make(map[string]runningCommand),
		eventEmitter:    emitter,
		logger:          logger,
		config:          config,
	}
}

// RunCommand executes a command and streams its output.
func (c *DefaultRunner) RunCommand(command *domain.Command, environment execution.Environment) error {
	c.mutex.Lock()

	if _, exists := c.runningCommands[command.Id]; exists {
		// Command is already running, skip it
		c.mutex.Unlock()
		return nil
	}

	process := newCommandProcess(
		command.Command,
		command.ResolveWorkingDirectory(environment.BaseWorkingDirectory),
		commandEnvironment(environment.Paths),
		c.config,
	)

	var wg sync.WaitGroup
	exited := make(chan struct{})
	running := runningCommand{
		process: process,
		wg:      &wg,
		exited:  exited,
	}

	c.sendStartingLine(command)
	readers, err := process.Start()
	if err != nil {
		c.sendStreamLine(command, err.Error())
		c.mutex.Unlock()
		return err
	}

	c.eventEmitter.EmitEvent(event.ProcessStarted, command.Id)

	// Save the command in the runningCommands map
	c.runningCommands[command.Id] = running
	c.mutex.Unlock()

	// Add to WaitGroup before starting goroutines to avoid race conditions
	wg.Add(len(readers) + 1)

	var scanWg sync.WaitGroup

	scanWg.Add(len(readers))
	for _, reader := range readers {
		go func(reader io.ReadCloser) {
			defer reader.Close()
			defer scanWg.Done()
			defer wg.Done()
			c.streamOutput(command, process, reader)
		}(reader)
	}

	// Wait in background until the command finishes, because it ends naturally or because it is stopped.
	go func() {
		defer wg.Done()

		// Notify the event emitter that the command has finished and remove it from the runningCommands map
		defer func() {
			c.mutex.Lock()
			delete(c.runningCommands, command.Id)
			c.mutex.Unlock()
			c.logger.Info("Command execution ended: " + command.Id)
			c.eventEmitter.EmitEvent(event.ProcessFinished, command.Id)
		}()

		err := process.Wait()
		close(exited)

		// A descendant that escaped the kill can hold the write ends open, so the
		// scanners get a grace period to drain what is buffered and are then cut off.
		if !waitFor(&scanWg, pipeDrainGrace) {
			closeProcessReaders(readers)
			scanWg.Wait()
		}

		// An expected termination is what stopping a Command looks like, not
		// something to report on its terminal.
		if err != nil && !isExpectedError(err) {
			c.sendStreamLine(command, err.Error())
			c.logger.Error("[ERROR - Waiting for project]: " + err.Error())
		}
	}()

	return nil
}

func (c *DefaultRunner) sendStartingLine(command *domain.Command) {
	c.emitLogEntry(command, command.Command, event.CommandLogEntry)
}

func (c *DefaultRunner) StopRunningCommand(id string) error {
	c.mutex.Lock()
	running, exists := c.runningCommands[id]
	c.mutex.Unlock()

	if !exists {
		return nil
	}

	return stopProcessGracefully(running.process, running.exited)
}

func (c *DefaultRunner) StopAllRunningCommands() []error {
	// Copy under the lock: stopping is slow and the wait goroutines delete from
	// the map as their Commands end.
	c.mutex.Lock()
	commandsToStop := make([]runningCommand, 0, len(c.runningCommands))
	for _, running := range c.runningCommands {
		commandsToStop = append(commandsToStop, running)
	}
	c.mutex.Unlock()

	errs := make([]error, 0)

	for _, running := range commandsToStop {
		err := stopProcessGracefully(running.process, running.exited)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// alreadyExited reports whether the goroutine owning cmd.Wait has seen the
// process end.
func alreadyExited(exited <-chan struct{}) bool {
	select {
	case <-exited:
		return true
	default:
		return false
	}
}

func waitFor(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// isExpectedError checks if the error is one of the expected termination logs.
func isExpectedError(err error) bool {
	for _, expected := range ExpectedTerminationLogs {
		if err.Error() == expected {
			return true
		}
	}
	return false
}

func (c *DefaultRunner) streamOutput(command *domain.Command, process commandProcess, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)
	scanner.Buffer(make([]byte, 1024), 1024*1024) // Set buffer size to 1MB

	for scanner.Scan() {
		line := scanner.Text()
		if process.shouldSkipOutputLine(line) {
			continue
		}
		c.logger.Debug(line)
		c.sendStreamLine(command, line)
	}
}

func (c *DefaultRunner) sendStreamLine(command *domain.Command, line string) {
	if command.MatchesErrorPattern(line) {
		c.eventEmitter.EmitEvent(event.CommandErrorDetected, command.Id)
	}

	c.emitLogEntry(command, line, event.OutputLogEntry)
}

func (c *DefaultRunner) emitLogEntry(command *domain.Command, line string, kind event.LogEntryKind) {
	c.eventEmitter.EmitEvent(event.NewLogEntry, map[string]string{
		"id":   command.Id,
		"line": line,
		"kind": string(kind),
	})
}

func (c *DefaultRunner) GetRunningCommandIds() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ids := make([]string, 0, len(c.runningCommands))
	for id := range c.runningCommands {
		ids = append(ids, id)
	}
	return ids
}
