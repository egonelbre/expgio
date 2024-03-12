package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(app.Size(150*6+50, 150*6-50))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			layoutBox(gtx)
			layoutDrag(gtx)
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

var drag gesture.Drag

func layoutDrag(gtx layout.Context) {
	defer op.Offset(image.Point{X: 100, Y: 100}).Push(gtx.Ops).Pop()
	defer clip.Rect{Max: image.Point{X: 30, Y: 30}}.Push(gtx.Ops).Pop()

	pointer.CursorGrab.Add(gtx.Ops)
	drag.Add(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{G: 0xFF, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	for {
		e, ok := drag.Update(gtx.Metric, gtx.Source, gesture.Both)
		if !ok {
			break
		}
		if e.Kind == pointer.Drag {
			p.X = e.Position.X
			p.Y = e.Position.Y
		}
	}
}

var p f32.Point

func layoutBox(gtx layout.Context) {
	defer op.Affine(f32.Affine2D{}.Offset(p)).Push(gtx.Ops).Pop()
	defer clip.Rect{Max: image.Point{X: 15, Y: 15}}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color.NRGBA{R: 0xFF, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
