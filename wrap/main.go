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
					return Wrap{Gap: theme.TextSize}.Layout(gtx, 'Z'-'A',
						func(gtx layout.Context, index int) layout.Dimensions {
							rec := op.Record(gtx.Ops)
							dims := layout.UniformInset(unit.Dp(2)).Layout(gtx,
								material.H1(theme, string('A'+index)).Layout)
							call := rec.Stop()

							clip.Rect{Max: dims.Size}.Add(gtx.Ops)
							paint.Fill(gtx.Ops, color.NRGBA{R: byte(index), G: byte(index * index), B: byte(index * index * index), A: 0xFF})

							call.Add(gtx.Ops)

							return dims
						})
				})

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

type Wrap struct {
	Gap unit.Value
}

func (wrap Wrap) Layout(gtx layout.Context, itemCount int, w layout.ListElement) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()

	// calculate the pixel size of the gap.
	gap := gtx.Px(wrap.Gap)
	// well use the constraints as the max width rather than automatically determining.
	width := gtx.Constraints.Max.X

	// measured keeps track of a measured element.
	type measured struct {
		dims layout.Dimensions
		call op.CallOp
	}

	// line keeps track of the current line.
	var line []measured
	// lineWidth keeps track of the total width of the line + gap
	var lineWidth int
	// maxHeight is the highest item in the line.
	var maxHeight int
	// y keeps track of the y position of the line + gap.
	var y int

	// flush flushes the current line.
	flush := func() {
		if len(line) == 0 {
			return
		}

		// x keeps track of the current x position of the widgets
		x := 0
		for _, w := range line {
			// save the stack to avoid propagating offsets to the next items
			stack := op.Save(gtx.Ops)

			// adjust the drawing to the correct location.
			op.Offset(f32.Pt(
				float32(x),
				// we center each item on the Y axis,
				float32(y+maxHeight/2-w.dims.Size.Y/2),
			)).Add(gtx.Ops)

			// draw the widget
			w.call.Add(gtx.Ops)

			// restore previous offset
			stack.Load()

			// update the x position.
			x += w.dims.Size.X + gap
		}

		// clear the line
		line = line[:0]
		// update the y position.
		y += maxHeight + gap
		maxHeight = 0
		lineWidth = 0
	}

	// create a child context that does not have a minimum constraint
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{}

	for i := 0; i < itemCount; i++ {
		// record the drawing of the item at position i
		macro := op.Record(gtx.Ops)
		dims := w(cgtx, i)
		call := macro.Stop()

		// update the max height
		if dims.Size.Y > maxHeight {
			maxHeight = dims.Size.Y
		}

		// when the line would overflow, then flush the current line.
		if lineWidth+dims.Size.X > width {
			flush()
		}

		// update the lineWidth with the new widget and add item to list
		lineWidth += dims.Size.X + gap
		line = append(line, measured{
			call: call,
			dims: dims,
		})
	}
	// flush the last line
	flush()

	if y > 0 {
		y -= gap
	}

	return layout.Dimensions{
		Size: image.Point{
			X: width,
			Y: y,
		},
		Baseline: 0,
	}
}
