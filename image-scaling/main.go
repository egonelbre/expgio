// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/draw"
)

func main() {
	img, err := createImage()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	th := material.NewTheme(gofont.Collection())
	ui := &UI{
		Theme: th,
		Image: paint.NewImageOp(img),
	}
	go func() {
		w := app.NewWindow(app.Title("Image Viewer"))
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Theme *material.Theme
	Image paint.ImageOp

	Direction widget.Enum
	Fit       widget.Enum
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				return nil
			}

		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

func (ui *UI) radio(enum *widget.Enum, label string) layout.FlexChild {
	return layout.Rigid(material.RadioButton(ui.Theme, enum, label, label).Layout)
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 1e6
			return layout.Flex{}.Layout(gtx,
				ui.radio(&ui.Direction, "NW"),
				ui.radio(&ui.Direction, "N"),
				ui.radio(&ui.Direction, "NE"),
				ui.radio(&ui.Direction, "E"),
				ui.radio(&ui.Direction, "SE"),
				ui.radio(&ui.Direction, "S"),
				ui.radio(&ui.Direction, "SW"),
				ui.radio(&ui.Direction, "W"),
				ui.radio(&ui.Direction, "Center"),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 1e6
			return layout.Flex{}.Layout(gtx,
				ui.radio(&ui.Fit, "Unscaled"),
				ui.radio(&ui.Fit, "Contain"),
				ui.radio(&ui.Fit, "Cover"),
				ui.radio(&ui.Fit, "Fill"),
				ui.radio(&ui.Fit, "ScaleDown"),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0xFF, A: 0xFF},
				clip.Rect{Max: gtx.Constraints.Max}.Op(),
			)

			var pos layout.Direction
			switch ui.Direction.Value {
			case "NW":
				pos = layout.NW
			case "N":
				pos = layout.N
			case "NE":
				pos = layout.NE
			case "E":
				pos = layout.E
			case "SE":
				pos = layout.SE
			case "S":
				pos = layout.S
			case "SW":
				pos = layout.SW
			case "W":
				pos = layout.W
			case "", "Center":
				pos = layout.Center
			}

			var fit widget.Fit
			switch ui.Fit.Value {
			case "", "Unscaled":
				fit = widget.Unscaled
			case "Contain":
				fit = widget.Contain
			case "Cover":
				fit = widget.Cover
			case "Fill":
				fit = widget.Fill
			case "ScaleDown":
				fit = widget.ScaleDown
			default:
				fit = widget.Unscaled
			}

			return layout.UniformInset(64).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					paint.FillShape(gtx.Ops,
						color.NRGBA{G: 0xFF, A: 0xFF},
						clip.Rect{Max: gtx.Constraints.Max}.Op(),
					)
					return widget.Image{Src: ui.Image, Fit: fit, Position: pos}.Layout(gtx)
				})
		}),
	)
}

func createImage() (image.Image, error) {
	gif, err := gif.Decode(bytes.NewReader(gifData[:]))
	if err != nil {
		return nil, err
	}
	scaled := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.NearestNeighbor.Scale(scaled, scaled.Bounds(), gif, gif.Bounds(), draw.Over, nil)
	return scaled, nil
}

var gifData = [...]byte{
	0x47, 0x49, 0x46, 0x38, 0x37, 0x61, 0x19, 0x00, 0x19, 0x00, 0xa2, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x57, 0x72, 0x82, 0x39, 0x46,
	0x54, 0x96, 0xd6, 0xff, 0x00, 0x00, 0x00, 0x6c, 0x55, 0x19, 0x00, 0x00,
	0x00, 0x21, 0xf9, 0x04, 0x09, 0x0a, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00,
	0x00, 0x00, 0x19, 0x00, 0x19, 0x00, 0x00, 0x03, 0x81, 0x18, 0xba, 0xdc,
	0xfe, 0x50, 0x08, 0xf8, 0xc6, 0x90, 0x4b, 0x10, 0x32, 0x29, 0x1b, 0x9b,
	0x24, 0x86, 0xde, 0x07, 0x6e, 0x28, 0xd7, 0x95, 0x67, 0x9a, 0xae, 0x54,
	0xeb, 0xa2, 0x98, 0xd7, 0x2e, 0x68, 0xc6, 0x6e, 0x0d, 0xd1, 0xc0, 0x8e,
	0x93, 0xa2, 0x50, 0x08, 0xf8, 0x02, 0xc4, 0x00, 0xb0, 0x21, 0x44, 0x16,
	0x8f, 0xc9, 0xe5, 0x27, 0xa7, 0x20, 0x24, 0x15, 0x83, 0xdd, 0xac, 0xba,
	0xc9, 0xda, 0x66, 0x33, 0x6f, 0x6c, 0x63, 0x00, 0x77, 0x0d, 0x5f, 0x82,
	0xa1, 0x0c, 0x1e, 0xac, 0xc7, 0x64, 0xb6, 0x4b, 0x0c, 0x07, 0x13, 0x8b,
	0x4e, 0x88, 0x2c, 0x75, 0xc7, 0x5f, 0x2b, 0x76, 0x45, 0x78, 0x0b, 0x7d,
	0x4c, 0x16, 0x87, 0x03, 0x7f, 0x79, 0x85, 0x0d, 0x6f, 0x0a, 0x06, 0x49,
	0x91, 0x82, 0x43, 0x0e, 0x8e, 0x01, 0x6b, 0x85, 0x92, 0x8a, 0x14, 0x99,
	0x82, 0x77, 0x25, 0xa0, 0x0b, 0x09, 0x00, 0x3b,
}
