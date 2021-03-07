package main

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func FillRect(gtx *Context, r image.Rectangle, c color.NRGBA) {
	defer op.Save(gtx.Ops).Load()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	clip.Rect(r).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func FillRectBorder(gtx *Context, r image.Rectangle, c color.NRGBA) {
	defer op.Save(gtx.Ops).Load()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	clip.Rect(r).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func FillLine(gtx *Context, from, to image.Point, width int, c color.NRGBA) {
	defer op.Save(gtx.Ops).Load()

	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(layout.FPt(from))
	p.LineTo(layout.FPt(to))
	clip.Stroke{
		Path: p.End(),
		Style: clip.StrokeStyle{
			Width: float32(width),
		},
	}.Op().Add(gtx.Ops)

	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
