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
	conPTYWidth               = 32766
	conPTYHeight              = 1000
	windows10MinBuild         = 17763
	windows11MinBuild         = 22000
	windowsProductWorkstation = 1
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
		// ConPTY initializes its screen before command output. Gomander has
		// already rendered the command header, so preserve that existing UI.
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

type windowsConPTYProcess struct {
	command          string
	workingDirectory string
	environment      []string
	conPTY           pty.ConPty
	conPTYCommand    *pty.Cmd
}

func newCommandProcess(command, workingDirectory string, environment []string, config Config) commandProcess {
	hostEnvironment := currentWindowsHostEnvironment()
	if shouldUseConPTY(config, hostEnvironment, isConPTYAvailable()) {
		return &windowsConPTYProcess{
			command:          command,
			workingDirectory: workingDirectory,
			environment:      environment,
		}
	}

	// Keep the established non-ConPTY execution path unchanged.
	cmd := exec.Command("cmd", "/C", command)
	cmd.Dir = workingDirectory
	cmd.Env = environment
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	return &execProcess{cmd: cmd}
}

func (p *windowsConPTYProcess) Start() ([]io.ReadCloser, error) {
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

	command := pseudoConsole.Command(commandInterpreter(), "/D", "/S", "/C", p.command)
	command.Dir = p.workingDirectory
	command.Env = p.environment
	if err := command.Start(); err != nil {
		_ = pseudoConsole.Close()
		return nil, err
	}

	p.conPTY = pseudoConsole
	p.conPTYCommand = command
	outputPipe := pseudoConsole.OutputPipe()
	return []io.ReadCloser{&conPTYOutput{
		reader: bufio.NewReader(outputPipe),
		pipe:   outputPipe,
	}}, nil
}

func (p *windowsConPTYProcess) Wait() error {
	waitErr := p.conPTYCommand.Wait()
	// Close the pseudo-console writer first while the runner drains its output.
	windows.ClosePseudoConsole(windows.Handle(p.conPTY.Fd()))
	_ = p.conPTY.InputPipe().Close()
	// go-pty exposes the exit code but does not convert it to exec.ExitError.
	if waitErr == nil && p.conPTYCommand.ProcessState != nil && !p.conPTYCommand.ProcessState.Success() {
		return fmt.Errorf("exit status %d", p.conPTYCommand.ProcessState.ExitCode())
	}
	return waitErr
}

func (p *windowsConPTYProcess) PID() int {
	return p.conPTYCommand.Process.Pid
}

func (p *windowsConPTYProcess) shouldSkipOutputLine(line string) bool {
	if !strings.HasPrefix(line, "\x1b]0;") {
		return false
	}
	titleEnd := strings.IndexByte(line, '\a')
	return titleEnd >= 0 && line[titleEnd+1:] == "\x1b[?25h"
}

func currentWindowsHostEnvironment() HostEnvironment {
	version := windows.RtlGetVersion()
	return classifyWindowsVersion(version.MajorVersion, version.BuildNumber, version.ProductType)
}

func classifyWindowsVersion(majorVersion, buildNumber uint32, productType byte) HostEnvironment {
	// Product type 1 is VER_NT_WORKSTATION. Server releases can share Windows
	// 10's major version and build range but are outside the validated scope.
	if majorVersion == 10 &&
		buildNumber >= windows10MinBuild &&
		buildNumber < windows11MinBuild &&
		productType == windowsProductWorkstation {
		return HostEnvironmentWindows10
	}
	return ""
}

func shouldUseConPTY(config Config, environment HostEnvironment, available bool) bool {
	return available && environment != "" && config.enablesConPTY(environment)
}

func isConPTYAvailable() bool {
	return windows.NewLazySystemDLL("kernel32.dll").NewProc("CreatePseudoConsole").Find() == nil
}

func commandInterpreter() string {
	interpreter := os.Getenv("COMSPEC")
	if interpreter == "" {
		return "cmd.exe"
	}
	return interpreter
}

func stopProcessGracefully(process commandProcess, exited <-chan struct{}) error {
	if alreadyExited(exited) {
		return nil
	}

	pid := strconv.Itoa(process.PID())
	if err := runTaskkill("/PID", pid); err != nil {
		return forceKill(pid, exited)
	}

	select {
	case <-time.After(5 * time.Second):
		return forceKill(pid, exited)
	case <-exited:
		return nil
	}
}

func forceKill(pid string, exited <-chan struct{}) error {
	err := runTaskkill("/F", "/T", "/PID", pid)
	if err != nil && alreadyExited(exited) {
		return nil
	}
	return err
}

func runTaskkill(arguments ...string) error {
	cmd := exec.Command("taskkill", arguments...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
