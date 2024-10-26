package main

import (
	"gioui.org/layout"
)

type PercentilePlot struct {
	Data *Data
	Col  int
}

func (plot *PercentilePlot) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	size := gtx.Constraints.Max

	return layout.Dimensions{Size: size}
}
