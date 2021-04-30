package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
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
	theme := material.NewTheme(gofont.Collection())

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.UniformInset(theme.TextSize).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.N.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return Wrap{
								Gap:        theme.TextSize,
								ItemWidth:  unit.Dp(256),
								ItemHeight: unit.Dp(64),
							}.Layout(gtx, 'Z'-'A',
								func(gtx layout.Context, index int) layout.Dimensions {
									paint.FillShape(gtx.Ops,
										color.NRGBA{R: byte(index), G: byte(index * index), B: byte(index * index * index), A: 0xFF},
										clip.Rect{Max: gtx.Constraints.Max}.Op(),
									)
									return layout.Center.Layout(gtx, material.H3(theme, string('A'+index)).Layout)
								})
						})
				})
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

type Wrap struct {
	Gap unit.Value

	ItemWidth  unit.Value
	ItemHeight unit.Value
}

type Size struct {
	X unit.Value
	Y unit.Value
}

func (p Size) Px(gtx layout.Context) image.Point {
	return image.Point{
		X: gtx.Px(p.X),
		Y: gtx.Px(p.Y),
	}
}

func (wrap Wrap) Layout(gtx layout.Context, itemCount int, w layout.ListElement) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()

	gap := gtx.Px(wrap.Gap)
	itemWidth := gtx.Px(wrap.ItemWidth)
	itemHeight := gtx.Px(wrap.ItemHeight)

	itemsPerLine := (gtx.Constraints.Max.X + gap) / (itemWidth + gap)
	if itemsPerLine <= 0 {
		itemsPerLine = 1
	}

	index := 0

	cgtx := gtx
	cgtx.Constraints.Min = image.Point{X: itemWidth, Y: itemHeight}
	cgtx.Constraints.Max = image.Point{X: itemWidth, Y: itemHeight}

	maxY := 0
renderItems:
	for y := 0; ; y++ {
		for x := 0; x < itemsPerLine; x++ {
			stack := op.Save(gtx.Ops)

			op.Offset(f32.Pt(
				float32(x*(itemWidth+gap)),
				float32(y*(itemHeight+gap)),
			)).Add(gtx.Ops)

			w(cgtx, index)

			stack.Load()

			index++
			if index >= itemCount {
				maxY = y
				break renderItems
			}
		}
	}

	return layout.Dimensions{
		Size: image.Point{
			X: itemsPerLine*(itemWidth+gap) - gap,
			Y: (maxY+1)*(itemHeight+gap) - gap,
		},
		Baseline: 0,
	}
}
