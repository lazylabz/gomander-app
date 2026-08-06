package runner

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"gomander/internal/command/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/path"
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

type RunningCommand struct {
	process commandProcess
	done    chan struct{}
}

type DefaultRunner struct {
	runningCommands map[string]RunningCommand
	eventEmitter    event.EventEmitter
	logger          logger.Logger
	mutex           sync.Mutex
}

type Runner interface {
	RunCommand(command *domain.Command, environmentPaths []string, baseWorkingDirectory string) error
	RunCommands(commands []domain.Command, environmentPaths []string, baseWorkingDirectory string) error
	StopRunningCommand(id string) error
	StopAllRunningCommands() []error
	StopRunningCommands(commands []domain.Command) error
	GetRunningCommandIds() []string
}

func NewDefaultRunner(logger logger.Logger, emitter event.EventEmitter) *DefaultRunner {
	return &DefaultRunner{
		runningCommands: make(map[string]RunningCommand),
		eventEmitter:    emitter,
		logger:          logger,
	}
}

func (c *DefaultRunner) RunCommands(commands []domain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	for _, command := range commands {
		err := c.RunCommand(&command, environmentPaths, baseWorkingDirectory)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DefaultRunner) StopRunningCommands(commands []domain.Command) error {
	for _, command := range commands {
		err := c.StopRunningCommand(command.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunCommand executes a command and streams its output.
func (c *DefaultRunner) RunCommand(command *domain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	c.mutex.Lock()

	if _, exists := c.runningCommands[command.Id]; exists {
		// Command is already running, skip it
		c.mutex.Unlock()
		return nil
	}

	done := make(chan struct{})
	process := newCommandProcess(
		command.Command,
		path.GetComputedPath(baseWorkingDirectory, command.WorkingDirectory),
		commandEnvironment(environmentPaths),
	)
	runningCommand := RunningCommand{
		process: process,
		done:    done,
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
	c.runningCommands[command.Id] = runningCommand
	c.mutex.Unlock()

	var scanWg sync.WaitGroup
	scanWg.Add(len(readers))

	for _, reader := range readers {
		go func(pipe io.ReadCloser) {
			defer scanWg.Done()
			defer pipe.Close()
			c.streamOutput(command, pipe)
		}(reader)
	}

	// Wait in background until the command finishes, because it ends naturally or because it is stopped.
	go func() {
		defer close(done)

		// Notify the event emitter that the command has finished and remove it from the runningCommands map
		defer func() {
			c.mutex.Lock()
			delete(c.runningCommands, command.Id)
			c.mutex.Unlock()
			c.logger.Info("Command execution ended: " + command.Id)
			c.eventEmitter.EmitEvent(event.ProcessFinished, command.Id)
		}()

		err := process.Wait()
		scanWg.Wait()

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
	c.eventEmitter.EmitEvent(event.NewLogEntry, map[string]string{
		"id":   command.Id,
		"line": "\033[1;36m" + command.Command + "\033[0m",
	})
}

func (c *DefaultRunner) StopRunningCommand(id string) error {
	c.mutex.Lock()
	runningCommand, exists := c.runningCommands[id]
	c.mutex.Unlock()

	if !exists {
		return nil
	}

	return stopProcessGracefully(runningCommand.process, runningCommand.done)
}

func (c *DefaultRunner) StopAllRunningCommands() []error {
	errs := make([]error, 0)

	// Create a slice to hold commands to stop
	// this is necessary because we should not modify the map while iterating over it
	c.mutex.Lock()
	commandsToStop := make([]RunningCommand, 0, len(c.runningCommands))

	for _, command := range c.runningCommands {
		commandsToStop = append(commandsToStop, command)
	}
	c.mutex.Unlock()

	for _, command := range commandsToStop {
		err := stopProcessGracefully(command.process, command.done)

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
		if shouldSkipProcessOutputLine(line) {
			continue
		}
		c.logger.Debug(line)
		c.sendStreamLine(command, line)
	}
}

func (c *DefaultRunner) sendStreamLine(command *domain.Command, line string) {
	c.processStreamLine(command, line)

	c.eventEmitter.EmitEvent(event.NewLogEntry, map[string]string{
		"id":   command.Id,
		"line": line,
	})
}

func (c *DefaultRunner) processStreamLine(command *domain.Command, line string) {
	c.checkLineForErrors(command, line)
}

func (c *DefaultRunner) checkLineForErrors(command *domain.Command, line string) {
	errorPatterns := command.ErrorPatterns

	for _, pattern := range errorPatterns {
		matchString := strings.Contains(line, pattern)

		if matchString {
			c.eventEmitter.EmitEvent(event.CommandErrorDetected, command.Id)
			break
		}
	}
}

func (c *DefaultRunner) GetRunningCommands() map[string]RunningCommand {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	runningCommands := make(map[string]RunningCommand, len(c.runningCommands))
	for id, command := range c.runningCommands {
		runningCommands[id] = command
	}
	return runningCommands
}

func (c *DefaultRunner) WaitForCommand(commandId string) {
	c.mutex.Lock()
	runningCommand, exists := c.runningCommands[commandId]
	c.mutex.Unlock()

	if exists {
		<-runningCommand.done
	} else {
		return
	}
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
