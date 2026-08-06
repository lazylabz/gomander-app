//go:build windows

package runner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	pty "github.com/aymanbagabas/go-pty"
	"golang.org/x/sys/windows"
)

const (
	conPTYWidth  = 32766
	conPTYHeight = 1000
)

var (
	conPTYClearScreen = []byte("\x1b[2J")
	conPTYCursorHome  = []byte("\x1b[H")
)

type conPTYOutput struct {
	reader      *bufio.Reader
	pipe        io.Closer
	pending     *bytes.Reader
	pendingErr  error
	initialized bool
}

func (o *conPTYOutput) Read(buffer []byte) (int, error) {
	if !o.initialized {
		o.initialized = true
		firstLine, err := o.reader.ReadBytes('\n')
		// ConPTY clears and homes its screen before the first command output.
		// Gomander already rendered the command header, so keep that host-only
		// initialization from erasing it while preserving application ANSI data.
		if bytes.Contains(firstLine, conPTYClearScreen) && bytes.Contains(firstLine, conPTYCursorHome) {
			firstLine = bytes.Replace(firstLine, conPTYClearScreen, nil, 1)
			firstLine = bytes.Replace(firstLine, conPTYCursorHome, nil, 1)
		}
		o.pending = bytes.NewReader(firstLine)
		o.pendingErr = err
	}

	if o.pending.Len() > 0 {
		return o.pending.Read(buffer)
	}
	if o.pendingErr != nil {
		err := o.pendingErr
		o.pendingErr = nil
		return 0, err
	}

	return o.reader.Read(buffer)
}

func (o *conPTYOutput) Close() error {
	return o.pipe.Close()
}

type windowsProcess struct {
	command          string
	workingDirectory string
	environment      []string
	conPTY           pty.ConPty
	conPTYCommand    *pty.Cmd
	fallback         *execProcess
}

func newCommandProcess(command, workingDirectory string, environment []string) commandProcess {
	return &windowsProcess{
		command:          command,
		workingDirectory: workingDirectory,
		environment:      environment,
	}
}

func shouldSkipProcessOutputLine(line string) bool {
	if !strings.HasPrefix(line, "\x1b]0;") {
		return false
	}

	titleEnd := strings.IndexByte(line, '\a')
	return titleEnd >= 0 && line[titleEnd+1:] == "\x1b[?25h"
}

func (p *windowsProcess) Start() ([]io.ReadCloser, error) {
	if isConPTYAvailable() {
		terminal, err := pty.New()
		if err != nil {
			return nil, err
		}
		pseudoConsole, ok := terminal.(pty.ConPty)
		if !ok {
			_ = terminal.Close()
			return nil, errors.New("Windows PTY does not expose ConPTY pipes")
		}
		if err := pseudoConsole.Resize(conPTYWidth, conPTYHeight); err != nil {
			_ = pseudoConsole.Close()
			return nil, err
		}

		command := pseudoConsole.Command(os.Getenv("COMSPEC"), "/D", "/S", "/C", p.command)
		command.Dir = p.workingDirectory
		command.Env = p.environment
		if err := command.Start(); err != nil {
			_ = pseudoConsole.Close()
			return nil, err
		}

		p.conPTY = pseudoConsole
		p.conPTYCommand = command
		outputPipe := pseudoConsole.OutputPipe()
		output := &conPTYOutput{
			reader: bufio.NewReader(outputPipe),
			pipe:   outputPipe,
		}
		return []io.ReadCloser{output}, nil
	}

	cmd := exec.Command("cmd", "/D", "/S", "/C", p.command)
	cmd.Dir = p.workingDirectory
	cmd.Env = p.environment
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	p.fallback = &execProcess{cmd: cmd}
	return p.fallback.Start()
}

func (p *windowsProcess) Wait() error {
	if p.conPTY == nil {
		return p.fallback.Wait()
	}

	waitErr := p.conPTYCommand.Wait()
	// Close the pseudo-console writer first. The runner keeps draining the
	// output pipe to EOF, so buffered tail output is not discarded.
	windows.ClosePseudoConsole(windows.Handle(p.conPTY.Fd()))
	_ = p.conPTY.InputPipe().Close()
	// go-pty v0.2.2 exposes the exit code in ProcessState but does not turn a
	// non-zero Windows exit into an error as exec.Cmd does.
	if waitErr == nil && p.conPTYCommand.ProcessState != nil && !p.conPTYCommand.ProcessState.Success() {
		return fmt.Errorf("exit status %d", p.conPTYCommand.ProcessState.ExitCode())
	}
	return waitErr
}

func (p *windowsProcess) PID() int {
	if p.conPTY != nil {
		return p.conPTYCommand.Process.Pid
	}
	return p.fallback.PID()
}

func isConPTYAvailable() bool {
	return windows.NewLazySystemDLL("kernel32.dll").NewProc("CreatePseudoConsole").Find() == nil
}

func stopProcessGracefully(process commandProcess, done <-chan struct{}) error {
	select {
	case <-done:
		return nil
	default:
	}

	pid := strconv.Itoa(process.PID())
	if err := runTaskkill("/PID", pid); err != nil {
		select {
		case <-done:
			return nil
		default:
			return runTaskkill("/F", "/T", "/PID", pid)
		}
	}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		return runTaskkill("/F", "/T", "/PID", pid)
	case <-done:
		return nil
	}
}

func runTaskkill(arguments ...string) error {
	cmd := exec.Command("taskkill", arguments...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
