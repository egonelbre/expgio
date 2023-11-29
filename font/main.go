package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/egonelbre/expgio/font/noto"
)

func main() {
	ui := NewUI()

	go func() {
		w := app.NewWindow(app.Title("Font Demo"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme *material.Theme

	List widget.List
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(noto.Collection()))
	ui.List.List.Axis = layout.Vertical
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.NextEvent().(type) {
		case app.FrameEvent:

			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				return nil
			}

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return material.List(ui.Theme, &ui.List).Layout(gtx, 3,
		func(gtx layout.Context, index int) layout.Dimensions {
			switch index {
			case 0:
				lbl := material.H2(ui.Theme, "Ebm")
				lbl.Font.Typeface = "Noto Sans"
				return lbl.Layout(gtx)
			case 1:
				lbl := material.H2(ui.Theme, "Ebm")
				lbl.Font.Typeface = "Noto Sans"
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			case 2:
				lbl := material.H2(ui.Theme, "E♭m")
				lbl.Font.Typeface = "Noto Music"
				return lbl.Layout(gtx)
			}
			return layout.Dimensions{}
		})
}
