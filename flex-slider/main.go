//+build ignore

package main

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func main() {
	w := app.NewWindow()
	go loop(w)
	app.Main()
}

func loop(w *app.Window) {
	ops := new(op.Ops)
	box := new(box)
	bar := new(bar)
	var weight float32 = 0.5

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e.Queue, e.Config, e.Size)

			q := gtx.Queue
			// Check for drag event
			for _, ev := range q.Events(tag) {
				if x, ok := ev.(pointer.Event); ok {
					switch x.Type {
					case pointer.Drag:
						weight += x.Position.X / float32(gtx.Constraints.Min.X)
					}
				}
			}

			flex := layout.Flex{}
			flex.Layout(gtx,
				layout.Flexed(weight, box.Layout),
				layout.Rigid(bar.Layout),
			)
			e.Frame(gtx.Ops)
		}

	}
}

type box struct{}

func (b *box) Layout(gtx layout.Context) layout.Dimensions {
	yellow := color.NRGBA{R: 0xEE, G: 0xEE, B: 0x9E, A: 0xFF}
	paint.ColorOp{Color: yellow}.Add(gtx.Ops)
	size := gtx.Constraints.Max
	bounds := image.Rect(0, 0, size.X, size.Y)
	paint.PaintOp{Rect: layout.FRect(bounds)}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

var tag = new(bool)

type bar struct{}

func (b *bar) Layout(gtx layout.Context) layout.Dimensions {
	ops := gtx.Ops

	size := image.Pt(
		gtx.Px(unit.Dp(10)),
		gtx.Constraints.Max.Y,
	)
	bounds := image.Rect(0, 0, size.X, size.Y)

	// Make draggable
	pointer.Rect(bounds).Add(ops)
	pointer.InputOp{Tag: tag, Types: pointer.Drag}.Add(ops)

	black := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	paint.ColorOp{Color: black}.Add(ops)
	paint.PaintOp{Rect: layout.FRect(bounds)}.Add(ops)
	return layout.Dimensions{Size: size}
}

// ColorBox creates a widget with the specified dimensions and color.
func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	bounds := f32.Rect(0, 0, float32(size.X), float32(size.Y))
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: bounds}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
