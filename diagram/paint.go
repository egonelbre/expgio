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
	paint.FillShape(gtx.Ops, c, clip.Rect(r).Op())
}

func FillRectBorder(gtx *Context, r image.Rectangle, w float32, c color.NRGBA) {
	rr := layout.FRect(r)
	rr.Min.X += w / 2
	rr.Min.Y += w / 2
	rr.Max.X -= w / 2
	rr.Max.Y -= w / 2

	paint.FillShape(gtx.Ops, c,
		clip.Stroke{
			Path: clip.RRect{Rect: rr}.Path(gtx.Ops),
			Style: clip.StrokeStyle{
				Width: w,
			},
		}.Op())
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
