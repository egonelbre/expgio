package main

import (
	"image"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Size(150*6+50, 150*6-50))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	theme := material.NewTheme()

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			RelativeY{Pos: 0.3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H1(theme, "Hello").Layout(gtx)
			})
			e.Frame(gtx.Ops)
		}
	}
}

type RelativeY struct {
	// Pos is relative to the size.
	// 0 is top, 1 is bottom and 0.5 centered.
	Pos float32
}

func (ry RelativeY) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	min := gtx.Constraints.Min
	gtx.Constraints.Min.Y = 0
	dims := w(gtx)
	call := macro.Stop()

	sz := dims.Size
	if sz.X < min.X {
		sz.X = min.X
	}
	if sz.Y < min.Y {
		sz.Y = min.Y
	}

	fy := sz.Y - dims.Size.Y
	y := int(math.Round(float64(float32(fy) * ry.Pos)))
	defer op.Offset(image.Point{Y: y}).Push(gtx.Ops).Pop()

	call.Add(gtx.Ops)
	return layout.Dimensions{
		Size:     sz,
		Baseline: dims.Baseline + sz.Y - dims.Size.Y - y,
	}
}
