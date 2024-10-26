package main

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func main() {
	data := RandomData()
	ui := NewUI(data)

	go func() {
		var w app.Window
		w.Option(app.Title("histviz"))
		if err := ui.Run(&w); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	theme   *material.Theme
	palette *Palette

	data *Data

	percentileGrid *PercentileGridPlot
	percentilePlot *PercentilePlot
	trendPlot      *TrendPlot
	colorLegend    *ColorLegend
}

func NewUI(data *Data) *UI {
	palette := &Palette{}
	theme := material.NewTheme()
	return &UI{
		theme:   theme,
		palette: palette,

		data: data,

		percentileGrid: &PercentileGridPlot{
			Theme:   theme,
			Palette: palette,
			Data:    data,
		},

		percentilePlot: &PercentilePlot{
			Data: data,
		},
		trendPlot: &TrendPlot{
			Data: data,
		},
		colorLegend: &ColorLegend{
			Data:    data,
			Palette: palette,
		},
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
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return ui.percentileGrid.Layout(gtx)
}
