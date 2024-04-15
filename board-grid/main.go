package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"

	"gioui.org/op"
	"gioui.org/unit"
)

var (
	cellSize  = unit.Dp(5)
	boardSize = image.Pt(9, 9)
)

func main() {
	ui := NewUI()

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Board"),
			app.Size(800, 800),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Board *Board
}

func NewUI() *UI {
	return &UI{
		Board: &Board{Size: boardSize},
	}
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Center.Layout(gtx,
				BoardStyle{
					CellSize: 50,
					Gap:      10,
					Color:    color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF},
					Board:    ui.Board,
				}.Layout)

			e.Frame(gtx.Ops)
		}
	}
}
