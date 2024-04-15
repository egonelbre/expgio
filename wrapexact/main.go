package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
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

			layout.UniformInset(unit.Dp(theme.TextSize)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.N.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return Wrap{
								Gap:        unit.Dp(theme.TextSize),
								ItemWidth:  unit.Dp(256),
								ItemHeight: unit.Dp(64),
							}.Layout(gtx, 'Z'-'A',
								func(gtx layout.Context, index int) layout.Dimensions {
									paint.FillShape(gtx.Ops,
										color.NRGBA{R: byte(index), G: byte(index * index), B: byte(index * index * index), A: 0xFF},
										clip.Rect{Max: gtx.Constraints.Max}.Op(),
									)
									return layout.Center.Layout(gtx, material.H3(theme, string(rune('A'+index))).Layout)
								})
						})
				})
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

type Wrap struct {
	Gap unit.Dp

	ItemWidth  unit.Dp
	ItemHeight unit.Dp
}

type Size struct {
	X unit.Dp
	Y unit.Dp
}

func (p Size) Px(gtx layout.Context) image.Point {
	return image.Point{
		X: gtx.Dp(p.X),
		Y: gtx.Dp(p.Y),
	}
}

func (wrap Wrap) Layout(gtx layout.Context, itemCount int, w layout.ListElement) layout.Dimensions {
	gap := gtx.Dp(wrap.Gap)
	itemWidth := gtx.Dp(wrap.ItemWidth)
	itemHeight := gtx.Dp(wrap.ItemHeight)

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
			stack := op.Offset(image.Pt(
				x*(itemWidth+gap),
				y*(itemHeight+gap),
			)).Push(gtx.Ops)

			w(cgtx, index)

			stack.Pop()

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
