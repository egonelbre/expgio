package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Panel struct {
	Axis layout.Axis
	Size unit.Value

	Background  color.NRGBA
	Border      color.NRGBA
	BorderWidth unit.Value
}

func (panel *Panel) Layout(gtx layout.Context, widget layout.Widget) layout.Dimensions {
	minorSize := gtx.Px(panel.Size)

	size := gtx.Constraints.Max
	if panel.Axis == layout.Horizontal {
		size.Y = minorSize
	} else {
		size.X = minorSize
	}
	gtx.Constraints = layout.Exact(size)

	panel.fill(gtx)
	gtx = panel.border(gtx)

	_ = widget(gtx)

	return layout.Dimensions{Size: size}
}

func (panel *Panel) fill(gtx layout.Context) {
	paint.FillShape(gtx.Ops, panel.Background, clip.Rect{Max: gtx.Constraints.Max}.Op())
}

func (panel *Panel) border(gtx layout.Context) layout.Context {
	var bounds clip.Rect
	bounds.Max = gtx.Constraints.Max
	borderPx := gtx.Px(panel.BorderWidth)
	if panel.Axis == layout.Horizontal {
		bounds.Min.Y = bounds.Max.Y - borderPx
		gtx.Constraints.Min.Y -= borderPx
		gtx.Constraints.Max.Y -= borderPx
	} else {
		bounds.Min.X = bounds.Max.X - borderPx
		gtx.Constraints.Min.X -= borderPx
		gtx.Constraints.Max.X -= borderPx
	}

	paint.FillShape(gtx.Ops, panel.Border, bounds.Op())

	return gtx
}
