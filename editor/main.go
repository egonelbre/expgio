package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/egonelbre/expgio/editor/rig"
)

func main() {
	ui := NewUI()

	go func() {
		var w app.Window
		w.Option(app.Title("Font Demo"), app.Size(unit.Dp(1024), unit.Dp(1024)))
		if err := ui.Run(&w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme  *material.Theme
	Screen *rig.Screen
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Screen = rig.NewScreen()
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return ui.Screen.Layout(gtx)
}
