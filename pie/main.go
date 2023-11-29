package main

import (
	"image/color"
	"math"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type C = layout.Context
type D = layout.Dimensions

func main() {
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(400), unit.Dp(700)),
		)

		var ops op.Ops
		start := time.Now()

		for {
			switch e := w.NextEvent().(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				seconds := gtx.Now.Sub(start).Seconds() * 0.5
				layout.Center.Layout(gtx,
					Pie{
						Slices: []Slice{{
							Weight: float32((math.Sin(11*seconds)+1)*0.5 + 0.5),
							Color:  color.NRGBA{R: 255, G: 200, B: 0, A: 255},
						}, {
							Weight: float32((math.Sin(7*seconds)+1)*0.5 + 0.5),
							Color:  color.NRGBA{R: 0, G: 200, B: 0, A: 255},
						}, {
							Weight: float32((math.Sin(9*seconds)+1)*0.5 + 0.5),
							Color:  color.NRGBA{R: 0, G: 200, B: 255, A: 255},
						}, {
							Weight: float32((math.Sin(13*seconds)+1)*0.5 + 0.5),
							Color:  color.NRGBA{R: 0, G: 0, B: 200, A: 255},
						}, {
							Weight: float32((math.Sin(13*seconds)+1)*0.5 + 0.5),
							Color:  color.NRGBA{R: 200, G: 0, B: 200, A: 255},
						}},
						Hole: float32(math.Sin(seconds)*0.4) + 0.5,
					}.Layout)

				op.InvalidateOp{}.Add(gtx.Ops)
				e.Frame(gtx.Ops)
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	app.Main()
}

type Pie struct {
	Slices []Slice
	Hole   float32
}

type Slice struct {
	Weight float32
	Color  color.NRGBA
}

func (pie Pie) Layout(gtx layout.Context) layout.Dimensions {
	size := gtx.Constraints.Max
	size.X = min(size.X, size.Y)
	size.Y = size.X

	center := layout.FPt(size.Div(2))
	radius := float32(center.X)
	hole := radius * pie.Hole

	polarToXY := func(angle, radius float32) f32.Point {
		sn, cs := math.Sincos(float64(angle))
		return f32.Point{
			X: center.X + float32(cs*float64(radius)),
			Y: center.Y + float32(sn*float64(radius)),
		}
	}

	total := float32(0)
	for i := range pie.Slices {
		total += pie.Slices[i].Weight
	}

	const tau = 2 * math.Pi
	const segmentAngle = tau / 12
	overlap := float32(1 / radius)

	startAngle := float32(0.0)
	for _, slice := range pie.Slices {
		endAngle := startAngle + slice.Weight*tau/total
		if endAngle >= tau {
			endAngle = tau
		}

		if slice.Color.A == 0 {
			startAngle = endAngle
			continue
		}

		var p clip.Path
		p.Begin(gtx.Ops)
		if hole <= 0 {
			p.MoveTo(center)
		} else {
			p.MoveTo(polarToXY(startAngle, hole))
		}
		p.LineTo(polarToXY(startAngle, float32(radius)))

		drawArcTo := func(from, to, radius float32) {
			sagitta := radius * float32(1-math.Cos(float64(to-from)/2))
			p.QuadTo(
				polarToXY((from+to)/2, radius+sagitta),
				polarToXY(to, radius),
			)
		}

		// next segment position rounded to segment angle
		lastAngle := startAngle
		for segmentAlpha := startAngle + segmentAngle; segmentAlpha < endAngle; segmentAlpha += segmentAngle {
			drawArcTo(lastAngle, segmentAlpha, radius)
			lastAngle = segmentAlpha
		}
		drawArcTo(lastAngle, endAngle+overlap, radius)

		if hole <= 0 {
			p.LineTo(center)
		} else {
			p.LineTo(polarToXY(endAngle+overlap, hole))
			segmentAlpha := lastAngle - segmentAngle
			lastAngle := endAngle + overlap
			for ; segmentAlpha > startAngle; segmentAlpha -= segmentAngle {
				drawArcTo(lastAngle, segmentAlpha, hole)
				lastAngle = segmentAlpha
			}
			drawArcTo(lastAngle, startAngle, hole)
		}

		paint.FillShape(gtx.Ops, slice.Color, clip.Outline{Path: p.End()}.Op())

		startAngle = endAngle
	}

	return layout.Dimensions{Size: size}
}
