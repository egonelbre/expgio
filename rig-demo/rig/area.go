package rig

import (
	"image/color"
	"math/rand"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/egonelbre/expgio/f32color"
)

const (
	JoinSplitSize    = 20
	AreaBorderRadius = 3
)

var (
	AreaBackground       = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	BorderColor          = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	BorderHighlightColor = color.NRGBA{R: 0xFF, G: 0x80, B: 0x80, A: 0xFF}
)

type Area struct {
	Screen *Screen
	Editor *Editor
}

func NewArea(screen *Screen) *Area {
	return &Area{
		Screen: screen,
	}
}

func (area *Area) Clone() *Area {
	clone := &Area{}
	clone.Screen = area.Screen
	clone.Editor = area.Editor.Clone()
	return clone
}

func (area *Area) Update(gtx layout.Context) layout.Dimensions {
	return area.Editor.Layout(gtx)
}

type Editor struct {
	Area  *Area
	Color color.NRGBA
}

func (editor *Editor) Layout(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, editor.Color)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (editor *Editor) Clone() *Editor {
	clone := &Editor{}
	clone.Area = editor.Area
	// clone.Color = editor.Color
	clone.Color = f32color.HSL(rand.Float32(), 0.7, 0.7)
	return clone
}
