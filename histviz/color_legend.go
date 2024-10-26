package main

import (
	"gioui.org/layout"
)

type ColorLegend struct {
	Data    *Data
	Palette *Palette
	Active  float64
}

func (plot *ColorLegend) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	size := gtx.Constraints.Max

	return layout.Dimensions{Size: size}
}
