package main

import (
	"fmt"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/profile"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Circles"),
		)

		var ops op.Ops

		for {
			switch e := w.NextEvent().(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				var p clip.Path
				p.Begin(gtx.Ops)
				for x := 0; x < 1000; x += 10 {
					for y := 0; y < 1000; y += 10 {
						p.MoveTo(f32.Pt(float32(x), float32(y)))
						p.LineTo(f32.Pt(float32(x+5), float32(y)))
						p.LineTo(f32.Pt(float32(x+5), float32(y+5)))
						p.LineTo(f32.Pt(float32(x), float32(y+5)))
						p.Close()
					}
				}
				paint.FillShape(gtx.Ops, color.NRGBA{A: 0xFF}, clip.Outline{Path: p.End()}.Op())

				for _, ev := range gtx.Events(0) {
					fmt.Println(ev)
				}
				profile.Op{Tag: 0}.Add(gtx.Ops)
				gtx.Execute(op.InvalidateCmd{})

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
