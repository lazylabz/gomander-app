package domain

import "math"

// ZoomLevel is the magnification of the whole window, 1 being 100%.
type ZoomLevel float64

const (
	DefaultZoomLevel ZoomLevel = 1
	MinZoomLevel     ZoomLevel = 0.5
	MaxZoomLevel     ZoomLevel = 2

	zoomStepsPerUnit = 10
)

func (z ZoomLevel) In() ZoomLevel {
	return z.stepped(1)
}

func (z ZoomLevel) Out() ZoomLevel {
	return z.stepped(-1)
}

func (z ZoomLevel) Clamped() ZoomLevel {
	return min(max(z, MinZoomLevel), MaxZoomLevel)
}

// Counted in whole steps, so repeated zooms land on 1.1, 1.2... rather than on
// the float drift adding 0.1 accumulates.
func (z ZoomLevel) stepped(steps int) ZoomLevel {
	return ZoomLevel((math.Round(float64(z)*zoomStepsPerUnit) + float64(steps)) / zoomStepsPerUnit).Clamped()
}
