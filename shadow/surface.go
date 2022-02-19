package main

import (
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/egonelbre/expgio/shadow/f32color"
)

type SurfaceLayoutStyle struct {
	DarkMode     bool
	Background   color.NRGBA
	CornerRadius unit.Value
	Elevation    unit.Value
}

func (s *SurfaceLayoutStyle) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Constraints.Min
	rr := float32(gtx.Px(s.CornerRadius))

	r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
	s.layoutShadow(gtx, r, rr)
	defer clip.UniformRRect(r, rr).Push(gtx.Ops).Pop()

	background := s.Background
	if s.DarkMode {
		p := darkBlend(s.Elevation.V)
		background = f32color.LinearFromSRGB(background).Lighten(p).SRGB()
	}
	paint.Fill(gtx.Ops, background)

	return layout.Dimensions{Size: sz}
}

func darkBlend(x float32) float32 {
	if x <= 0 {
		return 0
	}
	p := 15.77125 - 15.77125/float32(math.Pow(2, float64(x)/3.438155))
	if p <= 0 {
		return 0
	} else if p > 16 {
		return 16 * 0.01
	}
	return p * 0.01
}

func (s *SurfaceLayoutStyle) layoutShadow(gtx layout.Context, r f32.Rectangle, rr float32) {
	if s.Elevation.V <= 0 {
		return
	}

	offset := pxf(gtx.Metric, s.Elevation)
	if offset > 1 {
		offset = float32(math.Sqrt(float64(offset)))
	}

	ambient := r
	gradientBox(gtx.Ops, ambient, rr, offset/2, color.NRGBA{A: 0x05})

	penumbra := r.Add(f32.Pt(0, offset/2))
	gradientBox(gtx.Ops, penumbra, rr, offset, color.NRGBA{A: 0x15})

	umbra := outset(penumbra, -offset/2)
	gradientBox(gtx.Ops, umbra, rr/4, offset/2, color.NRGBA{A: 0x20})
}

func gradientBox(ops *op.Ops, r f32.Rectangle, rr, spread float32, col color.NRGBA) {
	paint.FillShape(ops, col, clip.RRect{
		Rect: outset(r, spread),
		SE:   rr + spread, SW: rr + spread, NW: rr + spread, NE: rr + spread,
	}.Op(ops))
}

func outset(r f32.Rectangle, rr float32) f32.Rectangle {
	r.Min.X -= rr
	r.Min.Y -= rr
	r.Max.X += rr
	r.Max.Y += rr
	return r
}

func pxf(c unit.Metric, v unit.Value) float32 {
	switch v.U {
	case unit.UnitPx:
		return v.V
	case unit.UnitDp:
		s := c.PxPerDp
		if s == 0 {
			s = 1
		}
		return s * v.V
	case unit.UnitSp:
		s := c.PxPerSp
		if s == 0 {
			s = 1
		}
		return s * v.V
	default:
		panic("unknown unit")
	}
}
