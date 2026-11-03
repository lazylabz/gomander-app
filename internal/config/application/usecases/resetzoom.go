package usecases

import "gomander/internal/config/domain"

type ResetZoom struct {
	zoom *Zoom
}

func NewResetZoom(zoom *Zoom) *ResetZoom {
	return &ResetZoom{zoom: zoom}
}

func (uc *ResetZoom) Execute() error {
	return uc.zoom.change(func(domain.ZoomLevel) domain.ZoomLevel { return domain.DefaultZoomLevel })
}
