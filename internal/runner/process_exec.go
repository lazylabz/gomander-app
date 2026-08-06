package runner

import (
	"io"
	"os"
	"os/exec"
)

type execProcess struct {
	cmd *exec.Cmd
}

func (p *execProcess) Start() ([]io.ReadCloser, error) {
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
		return nil, err
	}

	p.cmd.Stdout = stdoutWriter
	p.cmd.Stderr = stderrWriter
	if err := p.cmd.Start(); err != nil {
		_ = stdoutReader.Close()
		_ = stdoutWriter.Close()
		_ = stderrReader.Close()
		_ = stderrWriter.Close()
		return nil, err
	}

	// The child inherited its own write handles. Closing the parent's copies
	// lets readers reach EOF after the child exits without letting Cmd.Wait
	// close the reader sides before trailing output is drained.
	_ = stdoutWriter.Close()
	_ = stderrWriter.Close()

	return []io.ReadCloser{stdoutReader, stderrReader}, nil
}

func (p *execProcess) Wait() error {
	return p.cmd.Wait()
}

func (p *execProcess) PID() int {
	return p.cmd.Process.Pid
}
