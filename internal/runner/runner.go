package runner

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"

	"gomander/internal/command/domain"
	"gomander/internal/event"
	"gomander/internal/helpers/path"
	"gomander/internal/logger"
	"gomander/internal/platform"
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

type DefaultRunner struct {
	runningCommands map[string]*exec.Cmd
	eventEmitter    event.EventEmitter
	logger          logger.Logger
}

type Runner interface {
	RunCommand(command *domain.Command, environmentPaths []string, baseWorkingDirectory string) error
	StopRunningCommand(id string) error
	StopAllRunningCommands() []error
}

func NewDefaultRunner(logger logger.Logger, emitter event.EventEmitter) *DefaultRunner {
	return &DefaultRunner{
		runningCommands: make(map[string]*exec.Cmd),
		eventEmitter:    emitter,
		logger:          logger,
	}
}

// RunCommand executes a command and streams its output.
func (c *DefaultRunner) RunCommand(command *domain.Command, environmentPaths []string, baseWorkingDirectory string) error {
	// Get the project object based on the project string and OS
	cmd := platform.GetCommand(command.Command)

	// Enable color output and set terminal type
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")
	cmd.Dir = path.GetComputedPath(baseWorkingDirectory, command.WorkingDirectory)

	// Set project attributes based on OS
	platform.SetProcAttributes(cmd)
	platform.SetProcEnv(cmd, environmentPaths)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.sendStreamErrorWhileStartingCommand(command, err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.sendStreamErrorWhileStartingCommand(command, err)
		return err
	}

	if err := cmd.Start(); err != nil {
		c.sendStreamErrorWhileStartingCommand(command, err)
		return err
	}

	c.eventEmitter.EmitEvent(event.ProcessStarted, command.Id)

	// Save the project in the runningCommands map
	c.runningCommands[command.Id] = cmd

	var wg sync.WaitGroup

	// Stream stdout
	go func() {
		wg.Add(1)
		defer wg.Done()
		c.streamOutput(command.Id, stdout)
	}()

	// Stream stderr
	go func() {
		wg.Add(1)
		defer wg.Done()
		c.streamOutput(command.Id, stderr)
	}()

	// Wait in background until the command finishes, because it ends naturally or because it is stopped.
	go func() {
		err := cmd.Wait()
		wg.Wait() // Wait for both stdout and stderr to finish

		// Notify the event emitter that the command has finished and remove it from the runningCommands map
		defer func() {
			delete(c.runningCommands, command.Id)
			c.logger.Info("Command execution ended: " + command.Id)
			c.eventEmitter.EmitEvent(event.ProcessFinished, command.Id)
		}()

		if err != nil {
			c.sendStreamLine(command.Id, err.Error())

			if !isExpectedError(err) {
				c.logger.Error("[ERROR - Waiting for project]: " + err.Error())
			}
		}
	}()

	return nil
}

func (c *DefaultRunner) StopRunningCommand(id string) error {
	runningCommand, exists := c.runningCommands[id]

	if !exists {
		return errors.New("No running command with id: " + id)
	}

	return platform.StopProcessGracefully(runningCommand)
}

func (c *DefaultRunner) StopAllRunningCommands() []error {
	errs := make([]error, 0)

	// Create a slice to hold commands to stop
	// this is necessary because we should not modify the map while iterating over it
	commandsToStop := make([]*exec.Cmd, 0, len(c.runningCommands))

	for _, cmd := range c.runningCommands {
		commandsToStop = append(commandsToStop, cmd)
	}

	for _, cmd := range commandsToStop {
		err := platform.StopProcessGracefully(cmd)

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

func (c *DefaultRunner) streamOutput(commandId string, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)

	for scanner.Scan() {
		line := scanner.Text()
		c.logger.Debug(line)

		c.sendStreamLine(commandId, line)
	}
}

func (c *DefaultRunner) sendStreamErrorWhileStartingCommand(command *domain.Command, err error) {
	c.sendStreamLine(command.Id, err.Error())
	c.logger.Error(err.Error())
	c.eventEmitter.EmitEvent(event.ProcessFinished, command.Id)
}

func (c *DefaultRunner) sendStreamLine(commandId string, line string) {
	c.eventEmitter.EmitEvent(event.NewLogEntry, map[string]string{
		"id":   commandId,
		"line": line,
	})
}

func (c *DefaultRunner) GetRunningCommands() map[string]*exec.Cmd {
	return c.runningCommands
}
