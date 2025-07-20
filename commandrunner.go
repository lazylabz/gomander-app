package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
)

type CommandRunner struct {
	runningCommands map[string]*exec.Cmd
	eventEmitter    *EventEmitter
	logger          *Logger
}

func NewCommandRunner(logger *Logger, emitter *EventEmitter) *CommandRunner {
	return &CommandRunner{
		runningCommands: make(map[string]*exec.Cmd),
		eventEmitter:    emitter,
		logger:          logger,
	}
}

// ExecCommand executes a command by its ID and streams its output.
func (c *CommandRunner) RunCommand(command Command) error {
	// Get the command object based on the command string and OS
	cmd := getCommand(command.Command)

	// Enable color output and set terminal type
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1", "TERM=xterm-256color")

	// Set commabd attributes based on OS
	setProcAttributes(cmd)

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

	// Save the command in the runningCommands map
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
			c.logger.error("[ERROR - Waiting for command]: " + err.Error())
			return
		}
		c.eventEmitter.emitEvent(ProcessFinished, command.Id)
	}()

	return nil
}

func (c *CommandRunner) StopRunningCommand(id string) error {
	runningCommand, exists := c.runningCommands[id]

	if !exists {
		return errors.New("No running runningCommand for command: " + id)
	}

	return stopProcessGracefully(runningCommand)
}

func (c *CommandRunner) streamOutput(commandId string, pipeReader io.ReadCloser) {
	scanner := bufio.NewScanner(pipeReader)

	for scanner.Scan() {
		line := scanner.Text()
		c.logger.debug(line)

		c.sendStreamLine(commandId, line)
	}
}

func (c *CommandRunner) sendStreamError(command Command, err error) {
	c.sendStreamLine(command.Id, err.Error())
	c.logger.error(err.Error())
	c.eventEmitter.emitEvent(ProcessFinished, command.Id)
}

func (c *CommandRunner) sendStreamLine(commandId string, line string) {
	c.eventEmitter.emitEvent(NewLogEntry, map[string]string{
		"id":   commandId,
		"line": line,
	})
}
