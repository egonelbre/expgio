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

	r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
	s.layoutShadow(gtx, r, rr)
	clip.UniformRRect(r, rr).Add(gtx.Ops)

	background := s.Background
	if gtx.Queue == nil {
		// calculate disabled color
		background = f32color.MulAlpha(s.Background, 150)
	}
	if s.DarkMode {
		p := darkBlend(s.Elevation.V)
		background = f32color.RGBAFromSRGB(background).Lighten(p).SRGB()
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

	d := int(offset + 1)
	if d > 6 {
		d = 6
	}

	background := (f32color.RGBA{A: 0.14 / float32(d*d)}).SRGB()
	for x := 0; x <= d; x++ {
		for y := 0; y <= d; y++ {
			px, py := float32(x)/float32(d)-0.5, float32(y)/float32(d)-0.15
			stack := op.Push(gtx.Ops)
			op.Offset(f32.Pt(px*offset, py*offset)).Add(gtx.Ops)
			clip.UniformRRect(r, rr).Add(gtx.Ops)
			paint.Fill(gtx.Ops, background)
			stack.Pop()
		}
	}
}

func outset(r f32.Rectangle, y, s float32) f32.Rectangle {
	r.Min.X += s
	r.Min.Y += s + y
	r.Max.X += -s
	r.Max.Y += -s + y
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
