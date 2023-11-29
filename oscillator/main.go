package main

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/egonelbre/expgio/oscillator/generator"
	"github.com/egonelbre/expgio/oscillator/scope"
)

func main() {
	ctx := context.TODO()
	gen := generator.NewClient(generator.DefaultConfig)

	ui := NewUI(gen)

	go func() { gen.Run(ctx) }()

	go func() {
		w := app.NewWindow(
			app.Title("Oscillator"),
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
	theme    *material.Theme
	status   *StatusBar
	controls *Controls
	scope    *scope.Display

	generator *generator.Client
}

func NewUI(gen *generator.Client) *UI {
	return &UI{
		theme: material.NewTheme(),

		status:   NewStatus(),
		controls: NewControls(gen),
		scope:    scope.NewDisplay(),

		generator: gen,
	}
}

type AsyncValue[T any] struct {
	mu      sync.Mutex
	updated bool
	value   T
}

func AsyncRead[T any](data <-chan T, updated func()) *AsyncValue[T] {
	async := &AsyncValue[T]{}
	go async.Read(data, updated)
	return async
}

func (async *AsyncValue[T]) Read(data <-chan T, updated func()) {
	for v := range data {
		async.mu.Lock()
		async.updated = true
		async.value = v
		async.mu.Unlock()

		updated()
	}
}

func (async *AsyncValue[T]) Check(onChange func(v T)) {
	async.mu.Lock()
	defer async.mu.Unlock()

	if !async.updated {
		return
	}
	async.updated = false
	onChange(async.value)
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	asyncStatus := AsyncRead(ui.generator.Status, w.Invalidate)
	asyncData := AsyncRead(ui.generator.Data, w.Invalidate)

	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:

			asyncStatus.Check(func(v generator.Status) {
				ui.status.Current = v
			})
			asyncData.Check(func(v generator.Data) {
				ui.scope.Data = v
			})

			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.status.Layout(ui.theme, gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return ui.scope.Layout(ui.theme, gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Let's hardcode the height
			gtx.Constraints.Max.Y = gtx.Metric.Sp(150)
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

			return ui.controls.Layout(ui.theme, gtx)
		}),
	)
}

type StatusBar struct {
	Current generator.Status
}

func NewStatus() *StatusBar {
	return &StatusBar{Current: "Initializing"}
}

func (status *StatusBar) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return material.H4(th, string(status.Current)).Layout(gtx)
}

type Controls struct {
	Generator *generator.Client

	Active  generator.Config
	Pending generator.Config

	Function Spin[generator.Function]
	Scale    Spin[generator.Scale]

	Ping  widget.Clickable
	Trace widget.Clickable
	Tune  widget.Clickable
}

func NewControls(gen *generator.Client) *Controls {
	initial := gen.InitialConfig()

	panel := &Controls{
		Generator: gen,

		Active:  initial,
		Pending: initial,
	}

	panel.Function.Current = &panel.Pending.Function
	panel.Scale.Current = &panel.Pending.Scale

	return panel
}

func (controls *Controls) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	defer func() {
		if controls.Pending != controls.Active {
			controls.Generator.Reconf(controls.Pending)
			controls.Active = controls.Pending
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}()

	// wrapper to avoid needing to repeat passing theme into Layout
	themed := func(w func(th *material.Theme, gtx layout.Context) layout.Dimensions) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			return w(th, gtx)
		}
	}

	if controls.Ping.Clicked(gtx) {
		controls.Generator.Control(generator.Ping)
	}
	if controls.Tune.Clicked(gtx) {
		controls.Generator.Control(generator.Tune)
	}
	if controls.Trace.Clicked(gtx) {
		controls.Generator.Control(generator.Trace)
	}

	return Grid{
		Row: []float64{1, 1, 1},
		Col: []float64{1, 1, 1},
		Gap: 8, Margin: 8,
	}.Layout(gtx,
		CellAt(0, 0, themed(controls.Function.Layout)),
		CellAt(1, 0, themed(controls.Scale.Layout)),

		// non functional buttons for demo
		CellRows(0, 2, 1, ColorBox{G: 0x88, B: 0x88, A: 0x88}.Layout),

		CellAt(0, 2, material.Button(th, &controls.Ping, "Ping").Layout),
		CellAt(1, 2, material.Button(th, &controls.Trace, "Trace").Layout),
		CellAt(2, 2, material.Button(th, &controls.Tune, "Tune").Layout),
	)
}

// useful for stubbig out things

type ColorBox color.NRGBA

func (c ColorBox) Layout(gtx layout.Context) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA(c), clip.Rect{Max: gtx.Constraints.Min}.Op())
	return layout.Dimensions{Size: gtx.Constraints.Min}
}
