package facade

import "github.com/skratchdot/open-golang/open"

type OpenFacade interface {
	Run(input string) error
}

type DefaultOpenFacade struct{}

func (DefaultOpenFacade) Run(input string) error {
	return open.Run(input)
}
