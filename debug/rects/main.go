package main

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		w := new(app.Window)
		var ops op.Ops

		colors := []color.NRGBA{
			{0xFF, 0, 0, 0xFF},
			{0xFF, 0xFF, 0, 0xFF},
			{0xFF, 0, 0xFF, 0xFF},
			{0, 0xFF, 0xFF, 0xFF},
			{0xAA, 0xFF, 0xFF, 0xFF},
		}

		for {
			e := w.Event()
			if ev, ok := e.(app.FrameEvent); ok {
				gtx := app.NewContext(&ops, ev)

				for x := 0; x < 70; x++ {
					for y := 0; y < 70; y++ {
						stack := clip.Rect{
							Min: image.Point{X: x * 10, Y: y * 10},
							Max: image.Point{X: x*10 + 10, Y: y*10 + 10},
						}.Push(gtx.Ops)
						paint.Fill(gtx.Ops, colors[x*y%len(colors)])
						stack.Pop()
					}
				}

				gtx.Execute(op.InvalidateCmd{})
				ev.Frame(gtx.Ops)
			}
		}

	}()

	app.Main()
}
