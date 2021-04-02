// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/egonelbre/expgio/async-loading/resource"
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
	theme  *material.Theme
	loader *resource.Loader

	reels *Reels
}

func NewUI() *UI {
	return &UI{
		theme:  material.NewTheme(gofont.Collection()),
		loader: resource.NewLoader(5),
		reels:  NewReels(),
	}
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {
	go ui.loader.Run(context.Background())

	var ops op.Ops
	for {
		select {
		case <-ui.loader.Updated():
			fmt.Printf("%#v\n", ui.loader.Stats())
			w.Invalidate()
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:

				gtx := layout.NewContext(&ops, e)

				ui.loader.Frame(gtx, func(gtx layout.Context) layout.Dimensions {
					return ui.reels.Layout(gtx, ui.theme, ui.loader)
				})

				e.Frame(gtx.Ops)
			}
		}
	}
}

type Reels struct {
	items []*Reel
	list  layout.List
}

func NewReels() *Reels {
	reels := &Reels{}
	reels.list.Axis = layout.Vertical
	for index := 0; index < 100; index++ {
		reels.items = append(reels.items, &Reel{
			index: index,
			count: 1000,
		})
	}
	return reels
}

func (reels *Reels) Layout(gtx layout.Context, th *material.Theme, loader *resource.Loader) layout.Dimensions {
	return reels.list.Layout(gtx, len(reels.items),
		func(gtx layout.Context, index int) layout.Dimensions {
			reel := reels.items[index]
			return reel.Layout(gtx, th, loader)
		})
}

type Reel struct {
	index int
	count int
	list  layout.List
}

func (reel *Reel) Layout(gtx layout.Context, th *material.Theme, loader *resource.Loader) layout.Dimensions {
	return reel.list.Layout(gtx, reel.count,
		func(gtx layout.Context, index int) layout.Dimensions {
			return defaultInset.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					px := gtx.Px(unit.Dp(64))
					size := image.Point{X: px, Y: px}
					gtx.Constraints = layout.Exact(size)

					data := Data{Reel: reel.index, Item: index}
					r := loader.Schedule(data, data.Load)

					switch r.State() {
					case resource.Queued:
						col := color.NRGBA{R: 0xFF, G: 0xC0, B: 0xC0, A: 0xFF}
						paint.FillShape(gtx.Ops, col, clip.Rect{Max: size}.Op())
					case resource.Loading:
						col := color.NRGBA{R: 0xC0, G: 0xFF, B: 0xC0, A: 0xFF}
						paint.FillShape(gtx.Ops, col, clip.Rect{Max: size}.Op())
					case resource.Loaded:
						col := color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF}
						paint.FillShape(gtx.Ops, col, clip.Rect{Max: size}.Op())

						data := r.Value().(Data)
						layout.Center.Layout(gtx, material.Body1(th, data.String()).Layout)
					}

					return layout.Dimensions{Size: size}
				})
		})
}

var defaultInset = layout.UniformInset(unit.Dp(8))

type Data struct {
	Reel int
	Item int
}

func (data *Data) Load(ctx context.Context) interface{} {
	time.Sleep(5 * time.Millisecond)
	return *data
}

func (data *Data) String() string {
	return strconv.Itoa(data.Reel) + ":" + strconv.Itoa(data.Item)
}
