package main

import (
	"image"
	"math"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Zoom struct {
	Level  int
	Center image.Point
}

func (zoom *Zoom) Multiplier() float32 {
	return 1 / ZoomLevels[zoom.Level]
}

const defaultZoom = 2

var ZoomLevels = [...]float32{
	defaultZoom - 2: 0.50,
	defaultZoom - 1: 0.75,
	defaultZoom:     1,
	defaultZoom + 1: 1.50,
	defaultZoom + 2: 2.00,
}

type ZoomHud struct {
	Zoom *Zoom

	slider widget.Float
}

func (hud *ZoomHud) Layout(gtx *Context) {
	layout.NW.Layout(gtx.Context, func(lgtx layout.Context) layout.Dimensions {
		lgtx.Constraints.Min.X = min(lgtx.Dp(100), lgtx.Constraints.Max.X)

		hud.slider.Value = float32(hud.Zoom.Level) / float32(len(ZoomLevels)-1)
		size := material.Slider(gtx.Theme.Theme, &hud.slider).Layout(lgtx)

		hud.Zoom.Level = int(math.Round(float64(hud.slider.Value) * float64(len(ZoomLevels)-1)))
		if hud.Zoom.Level < 0 {
			hud.Zoom.Level = 0
		} else if hud.Zoom.Level >= len(ZoomLevels) {
			hud.Zoom.Level = len(ZoomLevels) - 1
		}

		return size
	})
}
