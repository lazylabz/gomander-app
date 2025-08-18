package facade

import "os"

type FsFacade interface {
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) ([]byte, error)
}

type DefaultFsFacade struct{}

func (d DefaultFsFacade) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (d DefaultFsFacade) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
