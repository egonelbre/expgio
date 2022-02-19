package main

import (
	"image"
	"image/color"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type GridHud struct{}

func (*GridHud) Layout(gtx *Context) {
	paint.ColorOp{Color: color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}}.Add(gtx.Ops)

	min := image.Point{X: gtx.Dp / 2, Y: gtx.Dp / 2}
	max := image.Point{X: min.X + gtx.Dp, Y: min.Y + gtx.Dp}

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
