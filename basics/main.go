package main

import (
	"image"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
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
	tag := new(int)

	rect := f32.Rectangle{Max: f32.Point{X: 10, Y: 10}}
	bounds := image.Rect(0, 0, 10, 10)

	red := color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	blue := color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}

	color := red

	ops := new(op.Ops)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()

			for _, ev := range e.Queue.Events(tag) {
				if x, ok := ev.(pointer.Event); ok {
					switch x.Type {
					case pointer.Press:
						color = blue
					case pointer.Release:
						color = red
					}
				}
			}

			// register for listening events.
			pointer.Rect(bounds).Add(ops)
			pointer.InputOp{Tag: tag}.Add(ops)

			// draw the colored rect
			paint.ColorOp{Color: color}.Add(ops)
			paint.PaintOp{Rect: rect}.Add(ops)

			// render the frame
			e.Frame(ops)
		}
	}
}
