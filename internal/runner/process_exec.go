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
		closeReaders(stdoutReader, stdoutWriter)
		return nil, err
	}

	p.cmd.Stdout = stdoutWriter
	p.cmd.Stderr = stderrWriter
	if err := p.cmd.Start(); err != nil {
		closeReaders(stdoutReader, stdoutWriter, stderrReader, stderrWriter)
		return nil, err
	}

	// The child inherited its own write handles. Closing the parent's copies
	// lets readers reach EOF without allowing Cmd.Wait to close them early.
	closeReaders(stdoutWriter, stderrWriter)
	return []io.ReadCloser{stdoutReader, stderrReader}, nil
}

func (p *execProcess) Wait() error {
	return p.cmd.Wait()
}

func (p *execProcess) PID() int {
	return p.cmd.Process.Pid
}

func closeReaders[T io.Closer](readers ...T) {
	for _, reader := range readers {
		_ = reader.Close()
	}
}
