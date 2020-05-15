package main

import (
	"image"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	gofont.Register()

	th := material.NewTheme()
	gtx := new(layout.Context)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Queue, e.Config, e.Size)

			split.Layout(gtx, func() {
				fill(gtx, color.RGBA{R: 0xC0, G: 0x30, B: 0x30, A: 0xFF})

				layout.Center.Layout(gtx, func() {
					material.H1(th, "Left").Layout(gtx)
				})
			}, func() {
				fill(gtx, color.RGBA{R: 0x30, G: 0x30, B: 0xC0, A: 0xFF})

				layout.Center.Layout(gtx, func() {
					material.H1(th, "Right").Layout(gtx)
				})
			})

			e.Frame(gtx.Ops)
		}
	}
}

var split Split

func bounds(gtx *layout.Context) f32.Rectangle {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	return f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
}

func fill(gtx *layout.Context, col color.RGBA) {
	dr := bounds(gtx)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
}
