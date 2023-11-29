package main

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/profile"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		w := app.NewWindow()
		var ops op.Ops

		colors := []color.NRGBA{
			{0xFF, 0, 0, 0xFF},
			{0xFF, 0xFF, 0, 0xFF},
			{0xFF, 0, 0xFF, 0xFF},
			{0, 0xFF, 0xFF, 0xFF},
			{0xAA, 0xFF, 0xFF, 0xFF},
		}

		profileTag := new(int)
		for {
			e := w.NextEvent()
			if ev, ok := e.(app.FrameEvent); ok {
				gtx := layout.NewContext(&ops, ev)

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

				for _, ev := range gtx.Events(profileTag) {
					fmt.Println(ev)
				}

				profile.Op{Tag: profileTag}.Add(gtx.Ops)
				op.InvalidateOp{}.Add(gtx.Ops)
				ev.Frame(gtx.Ops)
			}
		}

	}()

	app.Main()
}
