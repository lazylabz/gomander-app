package main

import (
	"gomander/internal/app"
	"gomander/internal/config/domain"
)

type WailsControllers struct {
	useCases app.UseCases
}

func NewWailsControllers() *WailsControllers {
	return &WailsControllers{}
}

func (wc *WailsControllers) loadUseCases(useCases app.UseCases) {
	wc.useCases = useCases
}

func (wc *WailsControllers) GetUserConfigController() (*domain.Config, error) {
	return wc.useCases.GetUserConfig.Execute()
}
