package main

import (
	"gioui.org/layout"
)

type TrendPlot struct {
	Data *Data
	Row  int
}

func (plot *TrendPlot) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	size := gtx.Constraints.Max

	return layout.Dimensions{Size: size}
}
