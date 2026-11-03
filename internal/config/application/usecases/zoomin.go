package usecases

import "gomander/internal/config/domain"

type ZoomIn struct {
	zoom *Zoom
}

func NewZoomIn(zoom *Zoom) *ZoomIn {
	return &ZoomIn{zoom: zoom}
}

func (uc *ZoomIn) Execute() error {
	return uc.zoom.change(domain.ZoomLevel.In)
}
