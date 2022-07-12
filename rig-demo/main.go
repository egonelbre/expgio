// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/egonelbre/expgio/rig-demo/rig"
)

const cursorCount = pointer.CursorNorthWestSouthEastResize + 1

func main() {
	go func() {
		th := material.NewTheme(gofont.Collection())
		_ = th
		w := app.NewWindow(app.Title("Rig Demo"))
		screen := rig.New()

		screen.Rig.SplitHorizontally(screen.Rig.Corners[1], 0.5)
		screen.Rig.SplitVertically(screen.Rig.Corners[3], 0.7)
		screen.Rig.SplitVertically(screen.Rig.Corners[3], 0.5)

		var ops op.Ops

		for e := range w.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				screen.Layout(gtx)
				e.Frame(gtx.Ops)

			case system.DestroyEvent:
				os.Exit(0)
			}
		}
		os.Exit(1)
	}()

	app.Main()
}
