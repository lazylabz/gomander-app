package facade

import (
	"io"
)

type IOFacade interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	ReadAll(r io.Reader) ([]byte, error)
}

type DefaultIOFacade struct {
}

func (d DefaultIOFacade) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (d DefaultIOFacade) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
