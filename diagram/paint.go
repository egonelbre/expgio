package main

import (
	"image"
	"image/color"

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
