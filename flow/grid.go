package main

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type GridLayer struct{}

func (*GridLayer) Layout(gtx *Context) {
	defer op.Save(gtx.Ops).Load()
	paint.ColorOp{Color: gtx.Theme.Grid}.Add(gtx.Ops)

	min := image.Point{X: gtx.Dp / 2, Y: gtx.Dp / 2}
	max := image.Point{X: min.X + gtx.Dp, Y: min.Y + gtx.Dp}

	scalePx := gtx.PxPerUnit
	var p image.Point
	for p.X = 0; p.X < gtx.Constraints.Max.X; p.X += scalePx {
		for p.Y = 0; p.Y < gtx.Constraints.Max.Y; p.Y += scalePx {
			stack := op.Save(gtx.Ops)
			clip.Rect{Min: p.Sub(min), Max: p.Add(max)}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Load()
		}
	}
}
