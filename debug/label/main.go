//go:build ignore

package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(
			app.Title("Drawing Test"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	var list widget.List
	list.Axis = layout.Vertical
	ops := new(op.Ops)

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(ops, e)

			material.List(th, &list).Layout(gtx, 2,
				func(gtx layout.Context, index int) layout.Dimensions {
					switch index {
					case 0:
						style := material.Label(th, 14, "gioui.org/reallylongnamethatshouldntbewrappedunderanycondition")
						style.MaxLines = 1
						return style.Layout(gtx)
					case 1:
						style := material.Label(th, 14, "gioui.org/reallylongnamethatshouldntbewrappedunderanycondition")
						style.MaxLines = 1
						return component.TruncatingLabelStyle(style).Layout(gtx)
					}
					return layout.Dimensions{}
				})

			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}
