package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/egonelbre/expgio/surface/f32color"
)

type SurfaceLayoutStyle struct {
	DarkMode     bool
	Background   color.RGBA
	CornerRadius unit.Value
	Elevation    unit.Value
}

func (s *SurfaceLayoutStyle) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Constraints.Min
	rr := float32(gtx.Px(s.CornerRadius))

	r := f32.Rect(0, 0, sz.X, sz.Y)
	s.layoutShadow(gtx, r, rr)

	clip.UniformRRect(r, rr).Add(gtx.Ops)

	background := b.Background
	if gtx.Queue == nil {
		// calculate disabled color
		background = f32color.MulAlpha(b.Background, 150)
	}

	paint.Fill(gtx.Ops, background)

	return layout.Dimensions{Size: sz}
}

func (s *SurfaceLayoutStyle) layoutShadow(gtx layout.Context, r f32.Rect, rr float32) {
	if s.Elevation.V <= 0 {
		return
	}
	defer op.Push(ops).Pop()

	clip.UniformRRect(r, rr).Add(gtx.Ops)
	paint.Fill(gtx.Ops, color.RGBA{A: 0x15})
}

func outset(r f32.Rect, x, y float32) f32.Rect {
	r.Min.X -= x
	r.Min.Y -= y
	r.Max.X += x
	r.Max.Y += y
	return r
}
