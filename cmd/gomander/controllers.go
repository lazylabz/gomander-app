package main

import (
	"gomander/internal/app"
	configdomain "gomander/internal/config/domain"
	projectdomain "gomander/internal/project/domain"
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

// User config controllers

func (wc *WailsControllers) GetUserConfigController() (*configdomain.Config, error) {
	return wc.useCases.GetUserConfig.Execute()
}

func (wc *WailsControllers) SaveUserConfigController(newConfig configdomain.Config) error {
	return wc.useCases.SaveUserConfig.Execute(newConfig)
}

// Project controllers

func (wc *WailsControllers) GetCurrentProjectController() (*projectdomain.Project, error) {
	return wc.useCases.GetCurrentProject.Execute()
}
