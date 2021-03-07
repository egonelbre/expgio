package main

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type GridHud struct{}

func (*GridHud) Layout(gtx *Context) {
	defer op.Save(gtx.Ops).Load()
	paint.ColorOp{Color: color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}}.Add(gtx.Ops)

	scalePx := gtx.PxPerUnit
	var p image.Point
	for p.X = 0; p.X < gtx.Constraints.Max.X; p.X += scalePx {
		for p.Y = 0; p.Y < gtx.Constraints.Max.Y; p.Y += scalePx {
			stack := op.Save(gtx.Ops)
			clip.Rect{Min: p, Max: p.Add(image.Point{X: 1, Y: 1})}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Load()
		}
	}
}
