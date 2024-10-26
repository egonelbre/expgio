package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"sort"
	"strconv"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Data struct {
	Metrics     []Metric
	Percentiles []Percentile
	Range       Range
}

type Percentile struct {
	Label string
	Value float64
}

type Metric struct {
	Label  string
	Values []float64
}

func RandomData() *Data {
	data := &Data{}

	data.Percentiles = []Percentile{
		{Label: "p0.1", Value: 0.001},
		{Label: "p1", Value: 0.01},
		{Label: "p5", Value: 0.05},
		{Label: "p10", Value: 0.10},
		{Label: "p25", Value: 0.25},
		{Label: "p50", Value: 0.50},
		{Label: "p75", Value: 0.75},
		{Label: "p90", Value: 0.90},
		{Label: "p95", Value: 0.95},
		{Label: "p99", Value: 0.99},
		{Label: "p99.9", Value: 0.999},
	}
	data.Metrics = make([]Metric, 32)

	for i := range data.Metrics {
		metric := &data.Metrics[i]
		metric.Label = strconv.Itoa(i)
		metric.Values = make([]float64, len(data.Percentiles))
		for i := range metric.Values {
			v := rand.ExpFloat64() * 100
			if v < 0 {
				v = -v
			}
			if v > 1000 {
				v = 1000
			}

			data.Range.Max = max(data.Range.Max, v)
			metric.Values[i] = v
		}
		sort.Float64s(metric.Values)
	}

	return data
}

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

type PercentileGridPlot struct {
	Theme   *material.Theme
	Palette *Palette
	Data    *Data

	Hover    bool
	HoverAt  image.Point
	Selected image.Point
}

func (w *PercentileGridPlot) Color(value float64) color.NRGBA {
	return w.Palette.Color(value, w.Data.Range)
}

func (w *PercentileGridPlot) Target() image.Point {
	if w.Hover {
		return w.HoverAt
	} else {
		return w.Selected
	}
}

func (w *PercentileGridPlot) Layout(gtx layout.Context) layout.Dimensions {
	totalSize := gtx.Constraints.Constrain(gtx.Constraints.Max)

	axisLineHeight := unit.Sp(20)

	axisSize := image.Point{
		X: gtx.Sp(80),
		Y: gtx.Sp(axisLineHeight),
	}

	gridSize := image.Point{
		X: totalSize.X - axisSize.X,
		Y: totalSize.Y - axisSize.Y,
	}

	cellcount := image.Point{
		X: len(w.Data.Metrics),
		Y: len(w.Data.Percentiles),
	}

	// calculate integer cell sizes
	cellSize := image.Point{
		X: gridSize.X / cellcount.X,
		Y: gridSize.Y / cellcount.Y,
	}

	// add any left over pixels for the header
	axisSize.X = totalSize.X - cellSize.X*cellcount.X
	axisSize.Y = totalSize.Y - cellSize.Y*cellcount.Y

	// final size of the grid
	gridSize = image.Point{
		X: totalSize.X - int(axisSize.X),
		Y: totalSize.Y - int(axisSize.Y),
	}

	func() {
		defer op.Offset(axisSize).Push(gtx.Ops).Pop()

		area := clip.Rect{
			Min: axisSize.Mul(-1),
			Max: gridSize,
		}
		defer area.Push(gtx.Ops).Pop()
		event.Op(gtx.Ops, w)

		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: w,
				Kinds:  pointer.Move | pointer.Enter | pointer.Leave | pointer.Cancel | pointer.Press,
			})
			if !ok {
				break
			}

			switch ev := ev.(type) {
			case pointer.Event:
				target := image.Point{
					X: int(ev.Position.X / float32(cellSize.X)),
					Y: int(ev.Position.Y / float32(cellSize.Y)),
				}

				target.X = min(max(0, target.X), cellcount.X)
				target.Y = min(max(0, target.Y), cellcount.Y)

				w.HoverAt = target

				if ev.Buttons == pointer.ButtonPrimary {
					w.Selected = w.HoverAt
				}

				switch ev.Kind {
				case pointer.Enter, pointer.Press, pointer.Move, pointer.Drag:
					w.Hover = true
				case pointer.Leave, pointer.Cancel:
					w.Hover = false
				}
			}
		}

		// draw cells
		for x := range w.Data.Metrics {
			metric := &w.Data.Metrics[x]
			for y, value := range metric.Values {
				zero := mulpoint(cellSize, image.Point{X: x, Y: y})
				cell := clip.Rect{
					Min: zero,
					Max: zero.Add(cellSize),
				}

				stack := cell.Push(gtx.Ops)
				paint.ColorOp{Color: w.Color(value)}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)

				stack.Pop()
			}
		}

		target := w.Target()

		func() { // draw highlighted column
			zero := image.Point{
				X: cellSize.X * target.X,
				Y: -axisSize.Y,
			}

			mulpoint(cellSize, target)

			StrokeRect{
				Rect: image.Rectangle{
					Min: zero,
					Max: image.Point{
						X: zero.X + cellSize.X,
						Y: gridSize.Y,
					},
				},
				Inset: gtx.Dp(-4),
				Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x40},
			}.Add(gtx.Ops)
		}()

		func() { // draw highlighted row
			zero := image.Point{
				X: -axisSize.X,
				Y: cellSize.Y * target.Y,
			}

			mulpoint(cellSize, target)

			StrokeRect{
				Rect: image.Rectangle{
					Min: zero,
					Max: image.Point{
						X: gridSize.X,
						Y: zero.Y + cellSize.Y,
					},
				},
				Inset: gtx.Dp(-4),
				Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x40},
			}.Add(gtx.Ops)
		}()

		func() { // draw highlighted element
			zero := mulpoint(cellSize, target)

			StrokeRect{
				Rect: image.Rectangle{
					Min: zero,
					Max: zero.Add(cellSize),
				},
				Inset: -gtx.Dp(4),
				Color: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x80},
			}.Add(gtx.Ops)
		}()

		func() { // draw column labels
			lineHeight := axisLineHeight
			textHeight := lineHeight * 2 / 3

			lgtx := gtx
			lgtx.Constraints.Max = image.Point{
				X: cellSize.X,
				Y: axisSize.Y,
			}
			lgtx.Constraints.Min = lgtx.Constraints.Max

			for x := range w.Data.Metrics {
				metric := &w.Data.Metrics[x]

				stack := op.Offset(image.Point{
					X: x * cellSize.X,
					Y: -axisSize.Y,
				}).Push(gtx.Ops)

				colMacro := op.Record(gtx.Ops)
				paint.ColorOp{
					Color: color.NRGBA{A: 0xFF},
				}.Add(gtx.Ops)
				colorOp := colMacro.Stop()

				widget.Label{
					Alignment:  text.Middle,
					MaxLines:   1,
					Truncator:  "...",
					LineHeight: lineHeight,
				}.Layout(lgtx, w.Theme.Shaper, font.Font{}, textHeight, metric.Label, colorOp)

				stack.Pop()
			}
		}()

		func() { // draw row labels
			lineHeight := axisLineHeight
			textHeight := lineHeight * 2 / 3

			lgtx := gtx
			lgtx.Constraints.Max = image.Point{
				X: axisSize.X,
				Y: cellSize.Y,
			}
			lgtx.Constraints.Min = lgtx.Constraints.Max

			for y, value := range w.Data.Percentiles {
				label := value.Label

				stack := op.Offset(image.Point{
					X: -axisSize.X,
					Y: y * cellSize.Y,
				}).Push(gtx.Ops)

				colMacro := op.Record(gtx.Ops)
				paint.ColorOp{
					Color: color.NRGBA{A: 0xFF},
				}.Add(gtx.Ops)
				colorOp := colMacro.Stop()

				widget.Label{
					Alignment:  text.Middle,
					MaxLines:   1,
					Truncator:  "...",
					LineHeight: lineHeight,
				}.Layout(lgtx, w.Theme.Shaper, font.Font{}, textHeight, label, colorOp)

				stack.Pop()
			}
		}()
	}()

	return layout.Dimensions{Size: totalSize}
}

type PercentilePlot struct {
	Data *Data
	Col  int
}

func (plot *PercentilePlot) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	size := gtx.Constraints.Max

	return layout.Dimensions{Size: size}
}

type TrendPlot struct {
	Data *Data
	Row  int
}

func (plot *TrendPlot) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	size := gtx.Constraints.Max

	return layout.Dimensions{Size: size}
}

type ColorLegend struct {
	Data    *Data
	Palette *Palette
	Active  float64
}

func (plot *ColorLegend) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	size := gtx.Constraints.Max

	return layout.Dimensions{Size: size}
}

func mulpoint(a, b image.Point) image.Point {
	return image.Point{
		X: a.X * b.X,
		Y: a.Y * b.Y,
	}
}
