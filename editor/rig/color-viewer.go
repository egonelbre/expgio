package rig

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func ColorViewer(name string, color color.NRGBA) EditorDef {
	return EditorDef{
		Name: name,
		New: func() layout.Widget {
			return Color{Color: color}.Layout
		},
	}
}

type Color struct {
	Color color.NRGBA
}

func (color Color) Layout(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
