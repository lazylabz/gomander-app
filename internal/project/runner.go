package project

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"gomander/internal/event"
	"gomander/internal/logger"
	"gomander/internal/platform"
)

type Runner struct {
	runningCommands map[string]*exec.Cmd
	eventEmitter    *event.EventEmitter
	logger          *logger.Logger
}

func NewCommandRunner(logger *logger.Logger, emitter *event.EventEmitter) *Runner {
	return &Runner{
		runningCommands: make(map[string]*exec.Cmd),
		eventEmitter:    emitter,
		logger:          logger,
	}
}

// RunCommand executes a command and streams its output.
func (c *Runner) RunCommand(command Command, environmentPaths []string, baseWorkingDirectory string) error {
	// Get the project object based on the project string and OS
	cmd := platform.GetCommand(command.Command)

	// Enable color output and set terminal type
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")
	cmd.Dir = getCommandExecutionDirectory(command, baseWorkingDirectory)

	// Set project attributes based on OS
	platform.SetProcAttributes(cmd)
	platform.SetProcEnv(cmd, environmentPaths)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.sendStreamError(command, err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.sendStreamError(command, err)
		return err
	}

	if err := cmd.Start(); err != nil {
		c.sendStreamError(command, err)
		return err
	}

	c.eventEmitter.EmitEvent(event.ProcessStarted, command.Id)

	// Save the project in the runningCommands map
	c.runningCommands[command.Id] = cmd

	// Stream stdout
	go c.streamOutput(command.Id, stdout)
	// Stream stderr
	go c.streamOutput(command.Id, stderr)

	// Optional: Wait in background
	go func() {
		err := cmd.Wait()
		if err != nil {
			c.sendStreamError(command, err)
			c.logger.Error("[ERROR - Waiting for project]: " + err.Error())
			return
		}
		c.eventEmitter.EmitEvent(event.ProcessFinished, command.Id)
	}()

	return nil
}

func getCommandExecutionDirectory(command Command, baseWorkingDirectory string) string {
	if command.WorkingDirectory == "" {
		return baseWorkingDirectory
	}
	if filepath.IsAbs(command.WorkingDirectory) {
		return command.WorkingDirectory
	}
	return filepath.Join(baseWorkingDirectory, command.WorkingDirectory)
}

func (c *Runner) StopRunningCommand(id string) error {
	runningCommand, exists := c.runningCommands[id]

	if !exists {
		return errors.New("No running runningCommand for project: " + id)
	}

	return platform.StopProcessGracefully(runningCommand)
}

func (c *Runner) StopAllRunningCommands() []error {
	errs := make([]error, 0)

	for id, cmd := range c.runningCommands {
		err := platform.StopProcessGracefully(cmd)

		if err != nil {
			c.logger.Error("[ERROR - Stopping project]: " + err.Error())
			errs = append(errs, err)
		} else {
			c.eventEmitter.EmitEvent(event.ProcessFinished, id)
		}

		delete(c.runningCommands, id)
	}

	return errs
}

func (c *Runner) streamOutput(commandId string, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)

	for scanner.Scan() {
		line := scanner.Text()
		c.logger.Debug(line)

		c.sendStreamLine(commandId, line)
	}
}

func (c *Runner) sendStreamError(command Command, err error) {
	c.sendStreamLine(command.Id, err.Error())
	c.logger.Error(err.Error())
	c.eventEmitter.EmitEvent(event.ProcessFinished, command.Id)
}

func (c *Runner) sendStreamLine(commandId string, line string) {
	c.eventEmitter.EmitEvent(event.NewLogEntry, map[string]string{
		"id":   commandId,
		"line": line,
	})
}
