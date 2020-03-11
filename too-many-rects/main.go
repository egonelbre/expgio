package main

import (
	"log"
	"math"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/paint"

	"gioui.org/font/gofont"
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
	gtx := layout.NewContext(w.Queue())
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)

			const size = 2
			for y := 0; y < e.Size.Y; y += size {
				for x := 0; x < e.Size.X; x += size {
					hue := float32(x*y) * math.Phi / 1024.0
					paint.ColorOp{
						Color: HSL(hue, 0.6, 0.6),
					}.Add(gtx.Ops)
					paint.PaintOp{Rect: f32.Rectangle{
						Min: f32.Point{
							X: float32(x),
							Y: float32(y),
						},
						Max: f32.Point{
							X: float32(x + size),
							Y: float32(y + size),
						},
					}}.Add(gtx.Ops)
				}
			}

			e.Frame(gtx.Ops)
		}
	}
}
