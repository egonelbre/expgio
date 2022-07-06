package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow(
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
	th := material.NewTheme(gofont.Collection())
	var list widget.List
	list.Axis = layout.Vertical
	ops := new(op.Ops)

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(ops, e)

			material.List(th, &list).Layout(ElementsGap(gtx, 8,
				material.H1(th, "Hello, world!").Layout,
				material.H2(th, "This is an example").Layout,
				material.Body1(th, lipsum).Layout,
			))

			e.Frame(gtx.Ops)

		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

func Elements(gtx layout.Context, widgets ...layout.Widget) (layout.Context, int, layout.ListElement) {
	return gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	}
}

func ElementsGap(gtx layout.Context, gap unit.Dp, widgets ...layout.Widget) (layout.Context, int, layout.ListElement) {
	return gtx, len(widgets)*2 + 1, func(gtx layout.Context, index int) layout.Dimensions {
		if index%2 == 0 {
			return layout.Spacer{Width: gap, Height: gap}.Layout(gtx)
		}
		return widgets[index/2](gtx)
	}
}

const lipsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque pulvinar varius congue. In pharetra velit non nisi vehicula ullamcorper. Vivamus sed nulla vitae massa iaculis faucibus. Fusce erat dolor, hendrerit in efficitur sed, faucibus eget nulla. Phasellus at dolor ac justo gravida tincidunt a eu est. Ut in ante sapien. Cras magna urna, dictum id neque eu, ornare varius turpis. Phasellus feugiat ex sit amet consectetur dapibus. Nullam pulvinar mollis elementum. Donec lacinia velit eu orci semper tempor. Aenean fermentum a lacus sit amet facilisis. Duis dignissim sed quam non feugiat. Cras aliquet, eros vel venenatis dapibus, mi magna interdum ante, ut vulputate libero elit et magna.`
