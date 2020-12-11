package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Fill struct {
	Color color.NRGBA
}

func (fill Fill) Layout(gtx layout.Context) layout.Dimensions {
	paint.FillShape(gtx.Ops, fill.Color, clip.Rect{Max: gtx.Constraints.Max}.Op())
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
