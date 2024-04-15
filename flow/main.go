// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func main() {
	theme := NewTheme(material.NewTheme())
	diagram := NewDemoDiagram()
	ui := &UI{
		Theme:  theme,
		Editor: NewEditor(diagram),
	}
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Diagram"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Theme  *Theme
	Editor *Editor
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
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
	return ui.Editor.Layout(ui.Theme, gtx)
}
