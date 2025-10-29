package runner

import (
	"bufio"
	"io"
	"os"
	"os/exec"
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
	cmd *exec.Cmd
	wg  *sync.WaitGroup
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

	// Get the command object based on the project string and OS
	cmd := GetCommand(command.Command)

	// Enable color output and set terminal type
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")
	cmd.Dir = path.GetComputedPath(baseWorkingDirectory, command.WorkingDirectory)

	// Set project attributes based on OS
	SetProcAttributes(cmd)
	SetProcEnv(cmd, environmentPaths)

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
	runningCommand := RunningCommand{
		cmd: cmd,
		wg:  &wg,
	}

	c.sendStartingLine(command)
	if err := cmd.Start(); err != nil {
		c.sendStreamLine(command, err.Error())
		c.mutex.Unlock()
		return err
	}

	c.eventEmitter.EmitEvent(event.ProcessStarted, command.Id)

	// Save the command in the runningCommands map
	c.runningCommands[command.Id] = runningCommand
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

		// If the command is still running (StopRunningCommand has not been called, this is a self-ended command), wait for it to finish
		if cmd.ProcessState == nil {
			err := cmd.Wait()
			if err != nil {
				c.sendStreamLine(command, err.Error())

				if !isExpectedError(err) {
					c.logger.Error("[ERROR - Waiting for project]: " + err.Error())
				}
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

	return StopProcessGracefully(runningCommand.cmd)
}

func (c *DefaultRunner) StopAllRunningCommands() []error {
	errs := make([]error, 0)

	// Create a slice to hold commands to stop
	// this is necessary because we should not modify the map while iterating over it
	commandsToStop := make([]*exec.Cmd, 0, len(c.runningCommands))

	for _, cmd := range c.runningCommands {
		commandsToStop = append(commandsToStop, cmd.cmd)
	}

	for _, cmd := range commandsToStop {
		err := StopProcessGracefully(cmd)

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
	go c.asyncProcessStreamLine(command, line)

	c.eventEmitter.EmitEvent(event.NewLogEntry, map[string]string{
		"id":   command.Id,
		"line": line,
	})
}

func (c *DefaultRunner) asyncProcessStreamLine(command *domain.Command, line string) {
	checkLineForErrors(command, line, c)
}

func checkLineForErrors(command *domain.Command, line string, c *DefaultRunner) {
	errorPatterns := command.GetErrorPatterns()

	for _, pattern := range errorPatterns {
		matchString := strings.Contains(line, pattern)

		if matchString {
			c.eventEmitter.EmitEvent(event.CommandErrorDetected, command.Id)
			break
		}
	}
}

func (c *DefaultRunner) GetRunningCommands() map[string]RunningCommand {
	return c.runningCommands
}

func (c *DefaultRunner) WaitForCommand(commandId string) {
	c.mutex.Lock()
	runningCommand, exists := c.runningCommands[commandId]
	c.mutex.Unlock()

	if exists {
		runningCommand.wg.Wait()
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
