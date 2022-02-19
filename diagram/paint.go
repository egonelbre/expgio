package main

import (
	"image"
	"image/color"

	"gioui.org/layout"
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
			Path:  clip.RRect{Rect: rr}.Path(gtx.Ops),
			Width: w,
		}.Op())
}

func FillLine(gtx *Context, from, to image.Point, width int, c color.NRGBA) {
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(layout.FPt(from))
	p.LineTo(layout.FPt(to))
	paint.FillShape(gtx.Ops, c, clip.Stroke{
		Path:  p.End(),
		Width: float32(width),
	}.Op())
}
