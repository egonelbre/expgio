//go:build ignore
// +build ignore

package main

import (
	"image"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/egonelbre/expgio/f32color"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Size(unit.Px(150*6+50), unit.Px(150*6-50)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	t := 0.0

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, f32color.NRGBAHex(0xe5e5e5FF))
			//paint.Fill(gtx.Ops, f32color.NRGBAHex(0x121212ff))

			gtx.Execute(op.InvalidateCmd{})

			for dp := 0; dp < 24; dp++ {
				ix := dp % 6
				iy := dp / 6
				offset := f32.Pt(float32(50+150*ix), float32(50+150*iy))

				if dp == 0 {
					t += 0.1
					s := (float32(math.Sin(float64(t))) + 1) * 0.5
					drawSurface(gtx, offset, unit.Dp(s*24))
				} else {
					drawSurface(gtx, offset, unit.Dp(float32(dp)))
				}
			}

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

func drawSurface(gtx layout.Context, offset f32.Point, elevation unit.Value) {
	defer op.Save(gtx.Ops).Load()
	op.Offset(offset).Add(gtx.Ops)

	gtx.Constraints.Min = image.Pt(100, 100)
	gtx.Constraints.Max = image.Pt(100, 100)

	style := SurfaceLayoutStyle{
		//DarkMode:     true,
		//Background:   f32color.NRGBAHex(0x1e1e1eff),
		Background:   f32color.NRGBAHex(0xffffffff),
		CornerRadius: unit.Dp(5),
		Elevation:    elevation,
	}
	style.Layout(gtx)
}
