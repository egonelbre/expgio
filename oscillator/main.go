package main

import (
	"context"
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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
	theme := material.NewTheme()
	return &UI{
		theme: theme,

		status:   NewStatus(),
		controls: NewControls(gen),
		scope:    scope.NewDisplay(),

		generator: gen,
	}
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for {
		select {
		case ui.status.Current = <-ui.generator.Status:
			w.Invalidate()
		case ui.scope.Data = <-ui.generator.Data:
			w.Invalidate()

		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				ui.Layout(gtx)
				e.Frame(gtx.Ops)
			}
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

	Previous generator.Config
	Next     generator.Config

	Function Spin[generator.Function]
	Scale    Spin[generator.Scale]
}

func NewControls(gen *generator.Client) *Controls {
	initial := gen.InitialConfig()

	panel := &Controls{
		Generator: gen,

		Previous: initial,
		Next:     initial,
	}

	panel.Function.Current = &panel.Next.Function
	panel.Scale.Current = &panel.Next.Scale

	return panel
}

func (controls *Controls) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	defer func() {
		if controls.Next != controls.Previous {
			controls.Generator.Update(controls.Next)
			controls.Previous = controls.Next
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}()

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return controls.Function.Layout(th, gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return controls.Scale.Layout(th, gtx)
		}),
	)
}
