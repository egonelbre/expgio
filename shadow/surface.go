package main

import (
	"image"
	"image/color"
	"math"

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
	CornerRadius unit.Dp
	Elevation    unit.Dp
}

func (s *SurfaceLayoutStyle) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Constraints.Min
	rr := gtx.Dp(s.CornerRadius)

	r := image.Rect(0, 0, sz.X, sz.Y)
	s.layoutShadow(gtx, r, rr)
	defer clip.UniformRRect(r, rr).Push(gtx.Ops).Pop()

	background := s.Background
	if s.DarkMode {
		p := darkBlend(float32(s.Elevation))
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

func (s *SurfaceLayoutStyle) layoutShadow(gtx layout.Context, r image.Rectangle, rr int) {
	if s.Elevation <= 0 {
		return
	}

	offsetf := dpf(gtx.Metric, s.Elevation)
	if offsetf > 1 {
		offsetf = float32(math.Sqrt(float64(offsetf)))
	}
	offset := int(offsetf)

	ambient := r
	gradientBox(gtx.Ops, ambient, rr, offset/2, color.NRGBA{A: 0x05})

	penumbra := r.Add(image.Pt(0, offset/2))
	gradientBox(gtx.Ops, penumbra, rr, offset, color.NRGBA{A: 0x15})

	umbra := outset(penumbra, -offset/2)
	gradientBox(gtx.Ops, umbra, rr/4, offset/2, color.NRGBA{A: 0x20})
}

func gradientBox(ops *op.Ops, r image.Rectangle, rr, spread int, col color.NRGBA) {
	paint.FillShape(ops, col, clip.RRect{
		Rect: outset(r, spread),
		SE:   rr + spread, SW: rr + spread, NW: rr + spread, NE: rr + spread,
	}.Op(ops))
}

func outset(r image.Rectangle, rr int) image.Rectangle {
	r.Min.X -= rr
	r.Min.Y -= rr
	r.Max.X += rr
	r.Max.Y += rr
	return r
}

func dpf(c unit.Metric, v unit.Dp) float32 {
	s := c.PxPerDp
	if s == 0 {
		s = 1
	}
	return s * float32(v)
}
