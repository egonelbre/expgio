package rig

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/egonelbre/expgio/f32color"
)

const (
	BorderWidth  = unit.Dp(4)
	AreaRadius   = unit.Dp(3)
	CornerRadius = unit.Dp(10)
)

var (
	BackgroundColor    = color.NRGBA{R: 0x10, G: 0x10, B: 0x10, A: 0xFF}
	CornerDragColor    = color.NRGBA{R: 0xFF, G: 0x88, B: 0x88, A: 0xFF}
	CornerHoverColor   = color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	CornerPassiveColor = color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0x30}
)

type EditorID string

type EditorDef struct {
	Name string
	New  func() layout.Widget
}

type Screen struct {
	Registry *Registry
	Bounds   image.Rectangle
	Areas    []*Area
}

func NewScreen() *Screen {
	return &Screen{
		Registry: NewRegistry(),
		Bounds:   image.Rect(0, 0, 1024, 1024),
		Areas: []*Area{
			{
				Bounds: image.Rect(0, 0, 512, 512),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0, 0.6, 0.6)}.Layout,
				},
			},
			{
				Bounds: image.Rect(512, 0, 1024, 512),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0.2, 0.6, 0.6)}.Layout,
				},
			},
			{
				Bounds: image.Rect(0, 512, 512, 1024),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0.4, 0.6, 0.6)}.Layout,
				},
			},
			{
				Bounds: image.Rect(512, 512, 1024, 1024),
				Editor: &Editor{
					Widget: Color{Color: f32color.HSL(0.6, 0.6, 0.6)}.Layout,
				},
			},
		},
	}
}

func dprect(gtx layout.Context, r image.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: gtx.Metric.Dp(unit.Dp(r.Min.X)),
			Y: gtx.Metric.Dp(unit.Dp(r.Min.Y)),
		},
		Max: image.Point{
			X: gtx.Metric.Dp(unit.Dp(r.Max.X)),
			Y: gtx.Metric.Dp(unit.Dp(r.Max.Y)),
		},
	}
}

func (screen *Screen) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(gtx.Constraints.Max)

	defer clip.Rect(dprect(gtx, screen.Bounds)).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: BackgroundColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	for _, area := range screen.Areas {
		area.Layout(screen, gtx)
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

type Area struct {
	Bounds image.Rectangle
	Editor *Editor

	NE, NW, SE, SW Sizer
}

type Editor struct {
	Widget layout.Widget
}

func (area *Area) Layout(screen *Screen, gtx layout.Context) {
	inside := dprect(gtx, area.Bounds).Inset(gtx.Dp(BorderWidth) / 2)
	defer clip.UniformRRect(inside, gtx.Dp(AreaRadius)).Push(gtx.Ops).Pop()
	defer op.Offset(inside.Min).Push(gtx.Ops).Pop()

	egtx := gtx
	egtx.Constraints = layout.Exact(inside.Size())
	area.Editor.Widget(egtx)

	size := inside.Size()

	cornerRadius := gtx.Dp(CornerRadius)

	Corner{
		Rect:   offsetRect(0, 0, cornerRadius, cornerRadius),
		Cursor: pointer.CursorSouthEastResize,
		Sizer:  &area.NW,
	}.Layout(gtx)

	Corner{
		Rect:   offsetRect(size.X, 0, -cornerRadius, cornerRadius),
		Cursor: pointer.CursorSouthWestResize,
		Sizer:  &area.NE,
	}.Layout(gtx)

	Corner{
		Rect:   offsetRect(size.X, size.Y, -cornerRadius, -cornerRadius),
		Cursor: pointer.CursorNorthWestResize,
		Sizer:  &area.SE,
	}.Layout(gtx)

	Corner{
		Rect:   offsetRect(0, size.Y, cornerRadius, -cornerRadius),
		Cursor: pointer.CursorNorthEastResize,
		Sizer:  &area.SW,
	}.Layout(gtx)
}

func offsetRect(x, y int, dx, dy int) image.Rectangle {
	return image.Rect(x, y, x+dx, y+dy).Canon()
}

type Corner struct {
	Rect   image.Rectangle
	Cursor pointer.Cursor
	Sizer  *Sizer
}

func (corner Corner) Layout(gtx layout.Context) {
	defer op.Offset(corner.Rect.Min).Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{Max: corner.Rect.Size()}).Push(gtx.Ops).Pop()
	gtx.Constraints = layout.Exact(corner.Rect.Size())

	for _, ev := range corner.Sizer.Events(gtx.Metric, gtx.Queue, gesture.Both) {
		fmt.Println(ev)
	}
	corner.Cursor.Add(gtx.Ops)
	corner.Sizer.Add(gtx.Ops)

	var color color.NRGBA
	switch {
	case corner.Sizer.Dragging():
		color = CornerDragColor
	case corner.Sizer.Hovered():
		color = CornerHoverColor
	default:
		color = CornerPassiveColor
	}
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
