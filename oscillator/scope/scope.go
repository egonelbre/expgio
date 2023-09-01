package scope

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/egonelbre/expgio/oscillator/generator"
)

type Display struct {
	Data generator.Data
}

func NewDisplay() *Display {
	return &Display{}
}

func (display *Display) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(gtx.Constraints.Max)

	// NOTE: you can still use giocanvas here, if you want

	size := gtx.Constraints.Max
	paint.FillShape(gtx.Ops, black, clip.Rect{Max: size}.Op())

	data := &display.Data

	viewmin := data.Min
	viewmax := data.Max

	toDisplay := func(v generator.Point) f32.Point {
		px := (viewmax.X - v.X) / (viewmax.X - viewmin.X)
		py := (viewmax.Y - v.Y) / (viewmax.Y - viewmin.Y)

		return f32.Point{
			X: px * float32(size.X),
			Y: py * float32(size.Y),
		}
	}

	if len(data.Values) > 0 {
		path := clip.Path{}
		path.Begin(gtx.Ops)

		path.MoveTo(toDisplay(data.Values[0]))
		for _, v := range data.Values {
			path.LineTo(toDisplay(v))
		}

		paint.FillShape(gtx.Ops, white, clip.Stroke{
			Path:  path.End(),
			Width: float32(gtx.Metric.Dp(1)),
		}.Op())
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

var (
	white = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	black = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)
