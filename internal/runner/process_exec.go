package runner

import (
	"io"
	"os/exec"
)

type execProcess struct {
	cmd *exec.Cmd
}

func (p *execProcess) Start() ([]io.ReadCloser, error) {
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		return nil, err
	}

	if err := p.cmd.Start(); err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		return nil, err
	}

	return []io.ReadCloser{stdout, stderr}, nil
}

func (p *execProcess) Wait() error {
	return p.cmd.Wait()
}

func (p *execProcess) PID() int {
	return p.cmd.Process.Pid
}
