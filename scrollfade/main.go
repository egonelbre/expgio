// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	ui := NewUI()

	go func() {
		w := app.NewWindow(
			app.Title("Loader"),
		)
		if err := ui.Run(w); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

type UI struct {
	theme *material.Theme
	list  widget.List
}

func NewUI() *UI {
	ui := &UI{
		theme: material.NewTheme(),
	}
	ui.list.Axis = layout.Vertical
	return ui
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:

			gtx := layout.NewContext(&ops, e)

			pos := &ui.list.List.Position

			// TODO: you need to take into account pos.Offset and pos.OffsetLast as well
			// also first draw is invalid atm
			center := pos.First + pos.Count/2

			material.List(ui.theme, &ui.list).Layout(gtx, len(linesOfText),
				func(gtx layout.Context, index int) layout.Dimensions {
					line := material.Body1(ui.theme, linesOfText[index])
					line.Alignment = text.Middle

					distance := index - center
					if distance < 0 {
						distance = -distance
					}
					distance -= 3
					if distance > 0 {
						line.Color.A = uint8(max(255-distance*20, 0))
					}

					return line.Layout(gtx)
				})

			e.Frame(gtx.Ops)
		}
	}
}

var linesOfText = (func() []string {
	xs := make([]string, 100)
	for i := range xs {
		xs[i] = strconv.Itoa(i)
	}
	return xs
})()
