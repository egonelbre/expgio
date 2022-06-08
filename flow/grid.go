package main

import (
	"image"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type BackgroundLayer struct{}

func (*BackgroundLayer) Layout(gtx *Context) {
	paint.FillShape(gtx.Ops, gtx.Theme.Background, clip.Rect{Max: gtx.Constraints.Max}.Op())
}

type GridLayer struct{}

func (*GridLayer) Layout(gtx *Context) {
	paint.ColorOp{Color: gtx.Theme.Grid}.Add(gtx.Ops)

	min := image.Point{X: gtx.Transform.Dp / 2, Y: gtx.Transform.Dp / 2}
	max := image.Point{X: min.X + gtx.Transform.Dp, Y: min.Y + gtx.Transform.Dp}

	scalePx := gtx.PxPerUnit
	var p image.Point
	for p.X = 0; p.X < gtx.Constraints.Max.X; p.X += scalePx {
		for p.Y = 0; p.Y < gtx.Constraints.Max.Y; p.Y += scalePx {
			stack := clip.Rect{Min: p.Sub(min), Max: p.Add(max)}.Push(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Pop()
		}
	}
}
