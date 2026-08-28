package runner

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"

	"gomander/internal/command/domain"
	"gomander/internal/event"
	"gomander/internal/execution"
	"gomander/internal/logger"
)

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
	cmd *exec.Cmd
	wg  *sync.WaitGroup
	// exited is closed by the goroutine that owns cmd.Wait, so that stopping a
	// Command can await the exit without calling Wait a second time: os/exec
	// forbids concurrent Wait on one Cmd.
	exited chan struct{}
}

type DefaultRunner struct {
	runningCommands map[string]runningCommand
	eventEmitter    event.EventEmitter
	logger          logger.Logger
	mutex           sync.Mutex
}

type Runner interface {
	RunCommand(command *domain.Command, environment execution.Environment) error
	StopRunningCommand(id string) error
	StopAllRunningCommands() []error
	GetRunningCommandIds() []string
}

func NewDefaultRunner(logger logger.Logger, emitter event.EventEmitter) *DefaultRunner {
	return &DefaultRunner{
		runningCommands: make(map[string]runningCommand),
		eventEmitter:    emitter,
		logger:          logger,
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

	// Get the command object based on the project string and OS
	cmd := GetCommand(command.Command)

	// Enable color output and set terminal type
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")
	cmd.Dir = command.ResolveWorkingDirectory(environment.BaseWorkingDirectory)

	// Set project attributes based on OS
	SetProcAttributes(cmd)
	SetProcEnv(cmd, environment.Paths)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.sendStreamLine(command, err.Error())
		c.mutex.Unlock()
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.sendStreamLine(command, err.Error())
		c.mutex.Unlock()
		return err
	}

	var wg sync.WaitGroup
	exited := make(chan struct{})
	running := runningCommand{
		cmd:    cmd,
		wg:     &wg,
		exited: exited,
	}

	c.sendStartingLine(command)
	if err := cmd.Start(); err != nil {
		c.sendStreamLine(command, err.Error())
		c.mutex.Unlock()
		return err
	}

	c.eventEmitter.EmitEvent(event.ProcessStarted, command.Id)

	// Save the command in the runningCommands map
	c.runningCommands[command.Id] = running
	c.mutex.Unlock()

	// Add to WaitGroup before starting goroutines to avoid race conditions
	wg.Add(3) // stdout, stderr, and wait goroutines

	var scanWg sync.WaitGroup

	scanWg.Add(2) // For stdout and stderr streaming

	// Stream stdout
	go func() {
		defer scanWg.Done()
		defer wg.Done()
		c.streamOutput(command, stdout)
	}()
	// Stream stderr
	go func() {
		defer scanWg.Done()
		defer wg.Done()
		c.streamOutput(command, stderr)
	}()

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

		// Wait for all pipes to finish
		scanWg.Wait()

		err := cmd.Wait()
		close(exited)

		if err != nil {
			c.sendStreamLine(command, err.Error())

			if !isExpectedError(err) {
				c.logger.Error("[ERROR - Waiting for project]: " + err.Error())
			}
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

	return StopProcessGracefully(running.cmd, running.exited)
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
		err := StopProcessGracefully(running.cmd, running.exited)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
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

func (c *DefaultRunner) streamOutput(command *domain.Command, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)
	scanner.Buffer(make([]byte, 1024), 1024*1024) // Set buffer size to 1MB

	for scanner.Scan() {
		line := scanner.Text()
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
