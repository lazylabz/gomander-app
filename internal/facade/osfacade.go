package facade

import "os"

type OSFacade interface {
	Stat(name string) (os.FileInfo, error)
	TempDir() string
	Create(name string) (*os.File, error)
}

type DefaultOSFacade struct{}

func (d DefaultOSFacade) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (d DefaultOSFacade) TempDir() string {
	return os.TempDir()
}

func (d DefaultOSFacade) Create(name string) (*os.File, error) {
	return os.Create(name)
}
