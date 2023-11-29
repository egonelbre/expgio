package main

import (
	"image/color"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type (
	D = layout.Dimensions
	C = layout.Context
)

var defaultColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

func main() {
	go func() {
		w := app.NewWindow(app.Size(600, 600))
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops
	var angle float32
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.ColorOp{Color: defaultColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			angle += 0.01
			drawLine(gtx, angle)
			if angle > 2*math.Pi {
				angle -= 2 * math.Pi
			}

			op.InvalidateOp{}.Add(gtx.Ops)
			e.Frame(gtx.Ops)
		}
	}
}

func drawLine(gtx layout.Context, t float32) {
	center := gtx.Constraints.Max.Div(2)
	radius := center.X
	if radius > center.Y {
		radius = center.Y
	}
	r := float32(radius * 3 / 4)
	w := r * 0.2
	a := layout.FPt(center)

	y, x := math.Sincos(float64(t))
	d := f32.Pt(float32(x), float32(y))
	b := a.Add(d.Mul(r))

	var p clip.Path
	/*
		p.MoveTo(a)
		p.LineTo(b)
		p.LineTo(b.Add(f32.Pt(w, 0)))
		p.LineTo(a.Add(f32.Pt(w, 0)))
		p.Close()
	*/

	p.Begin(gtx.Ops)
	drawArrow(&p, a, b, 5)
	drawArrow(&p, b, b.Add(f32.Pt(w, 0)), 5)
	drawArrow(&p, b.Add(f32.Pt(w, 0)), a.Add(f32.Pt(w, 0)), 5)
	drawArrow(&p, a.Add(f32.Pt(w, 0)), a, 5)
	paint.FillShape(gtx.Ops, color.NRGBA{A: 0xFF}, clip.Outline{
		Path: p.End(),
	}.Op())

	p.Begin(gtx.Ops)
	p.MoveTo(a)
	p.LineTo(b)
	p.LineTo(b.Add(f32.Pt(w, 0)))
	p.LineTo(a.Add(f32.Pt(w, 0)))
	p.Close()
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0xFF, A: 0x80}, clip.Outline{
		Path: p.End(),
	}.Op())

	// enable-disable this part
	/*
		m := layout.FPt(gtx.Constraints.Max)
		p.MoveTo(f32.Pt(0, 0))
		p.LineTo(f32.Pt(m.X, 0))
		p.LineTo(f32.Pt(m.X, m.Y))
		p.LineTo(f32.Pt(0, m.Y))
		p.Close()
	*/
}

func drawArrow(p *clip.Path, a, b f32.Point, r float32) {
	n, dir := normal(a, b, r)
	dir = dir.Mul(10)
	p.MoveTo(a.Add(n))
	p.LineTo(b.Sub(dir).Add(n))
	p.LineTo(b.Sub(dir).Add(n.Mul(2)))
	p.LineTo(b)
	p.LineTo(b.Sub(dir).Sub(n.Mul(2)))
	p.LineTo(b.Sub(dir).Sub(n))
	p.LineTo(a.Sub(n))
	p.Close()
}

func normal(a, b f32.Point, w float32) (normal, dir f32.Point) {
	dir = b.Sub(a)
	normal.X, normal.Y = +dir.Y, -dir.X
	d := math.Hypot(float64(normal.X), float64(normal.Y))
	if math.Abs(d) < 1e-5 {
		return f32.Point{}, f32.Point{}
	}
	return normal.Mul(w / float32(d)), dir.Mul(w / float32(d))
}
