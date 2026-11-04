package usecases

import (
	"sync"

	"gomander/internal/config/domain"
)

type Window interface {
	SetZoom(level float64)
}

// Zoom is shared by the zoom use cases so their presses are serialised: each
// reaches the backend on its own goroutine, and a held key repeats.
type Zoom struct {
	mutex      sync.Mutex
	repository domain.Repository
	window     Window
}

func NewZoom(repository domain.Repository, window Window) *Zoom {
	return &Zoom{repository: repository, window: window}
}

func (z *Zoom) change(next func(domain.ZoomLevel) domain.ZoomLevel) error {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	config, err := z.repository.GetOrCreate()
	if err != nil {
		return err
	}

	level := next(config.ZoomLevel)
	if err := z.repository.UpdateZoomLevel(level); err != nil {
		return err
	}

	z.window.SetZoom(float64(level))

	return nil
}
