package usecases

import "gomander/internal/config/domain"

// RestoreZoom puts the window back at the stored Zoom Level. It has to run once
// the window is on screen: macOS ignores a zoom given before that.
type RestoreZoom struct {
	repository domain.Repository
	window     Window
}

func NewRestoreZoom(repository domain.Repository, window Window) *RestoreZoom {
	return &RestoreZoom{repository: repository, window: window}
}

func (uc *RestoreZoom) Execute() error {
	config, err := uc.repository.GetOrCreate()
	if err != nil {
		return err
	}

	uc.window.SetZoom(float64(config.ZoomLevel.Clamped()))

	return nil
}
