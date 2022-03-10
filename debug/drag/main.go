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
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Px(150*6+50), unit.Px(150*6-50)))
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
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layoutBox(gtx)
			layoutDrag(gtx)
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

var drag gesture.Drag

func layoutDrag(gtx layout.Context) {
	defer op.Offset(f32.Point{X: 100, Y: 100}).Push(gtx.Ops).Pop()
	defer clip.Rect{Max: image.Point{X: 30, Y: 30}}.Push(gtx.Ops).Pop()

	pointer.CursorGrab.Add(gtx.Ops)
	drag.Add(gtx.Ops)

	paint.ColorOp{Color: color.NRGBA{G: 0xFF, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	for _, e := range drag.Events(gtx.Metric, gtx, gesture.Both) {
		if e.Type == pointer.Drag {
			p.X = e.Position.X
			p.Y = e.Position.Y
		}
	}
}

var p f32.Point

func layoutBox(gtx layout.Context) {
	defer op.Offset(p).Push(gtx.Ops).Pop()
	defer clip.Rect{Max: image.Point{X: 15, Y: 15}}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color.NRGBA{R: 0xFF, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
