package main

import (
	"context"
	"image/color"
	"log"
	"os"
	"os/signal"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		w := app.NewWindow(app.Fullscreen.Option())
		err := run(ctx, w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(ctx context.Context, w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				title := material.H1(th, "Hello, Gio")
				maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
				title.Color = maroon
				title.Alignment = text.Middle
				title.Layout(gtx)

				e.Frame(gtx.Ops)
			}
		}
	}
}
