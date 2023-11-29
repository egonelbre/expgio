// SPDX-License-Identifier: Unlicense OR MIT

package main

// This diagram editor is based on the design described in:
// Game Programming Gems 5 - "Context-Sensitive HUDs for Editors" by Adam Martin.

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func main() {
	th := material.NewTheme()
	ui := &UI{
		Theme: th,
		Hud:   NewHudManager(th),
	}
	go func() {
		w := app.NewWindow(app.Title("Diagram"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Theme *material.Theme
	Hud   *HudManager
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.NextEvent().(type) {
		case app.FrameEvent:
			gtx := layout.NewContext(&ops, e)
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
	return ui.Hud.Layout(gtx)
}
