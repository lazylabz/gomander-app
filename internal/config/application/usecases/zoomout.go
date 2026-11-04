package usecases

import "gomander/internal/config/domain"

type ZoomOut struct {
	zoom *Zoom
}

func NewZoomOut(zoom *Zoom) *ZoomOut {
	return &ZoomOut{zoom: zoom}
}

func (uc *ZoomOut) Execute() error {
	return uc.zoom.change(domain.ZoomLevel.Out)
}
