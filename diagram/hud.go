package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type HudManager struct {
	Theme *material.Theme

	Diagram *Diagram

	Huds      []Hud
	Exclusive Hud
}

func NewHudManager(theme *material.Theme) *HudManager {
	m := &HudManager{
		Theme:   theme,
		Diagram: NewDiagram(),
	}

	m.Huds = append(m.Huds,
		&GridHud{},
		&NodeHud{},
		&PortHud{},
		&ConnectionHud{},
		&NodeCreationHud{},
		&ManipulationHud{},
		&ConnectionCreationHud{},
		// &NodeDeleter{}
		// &NodeOrderer{}
	)
	return m
}

func (m *HudManager) Layout(gtx layout.Context) layout.Dimensions {
	for _, hud := range m.Huds {
		var lgtx layout.Context
		if m.Exclusive == nil || m.Exclusive == hud {
			lgtx = gtx
		} else {
			lgtx = gtx.Disabled()
		}
		dgtx := NewContext(lgtx, m.Theme, m.Diagram)
		hud.Layout(dgtx)
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

type Hud interface {
	Layout(gtx *Context)
}

type Context struct {
	Zoom
	layout.Context
	Theme   *material.Theme
	Diagram *Diagram
}

func NewContext(gtx layout.Context, th *material.Theme, diagram *Diagram) *Context {
	px := gtx.Px(unit.Dp(30))
	px = (px / 24) * 24 // make it divisible by 2,3,4,6,12
	return &Context{
		Zoom:    Zoom{PxPerUnit: px},
		Context: gtx,
		Theme:   th,
		Diagram: diagram,
	}
}

type Zoom struct {
	PxPerUnit int
}

func (zoom *Zoom) Px(v Unit) int {
	return zoom.PxPerUnit * int(v)
}

func (zoom *Zoom) Pt(v Vector) image.Point {
	return image.Point{
		X: int(v.X) * zoom.PxPerUnit,
		Y: int(v.Y) * zoom.PxPerUnit,
	}
}

func (zoom *Zoom) FPt(v Vector) f32.Point {
	return f32.Point{
		X: float32(int(v.X) * zoom.PxPerUnit),
		Y: float32(int(v.Y) * zoom.PxPerUnit),
	}
}

func (zoom *Zoom) Inv(p image.Point) Vector {
	return Vector{
		X: Unit(p.X / zoom.PxPerUnit),
		Y: Unit(p.Y / zoom.PxPerUnit),
	}
}

func (zoom *Zoom) FInv(p f32.Point) Vector {
	return Vector{
		X: Unit(int(p.X) / zoom.PxPerUnit),
		Y: Unit(int(p.Y) / zoom.PxPerUnit),
	}
}

func (zoom *Zoom) Bounds(box Box) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(box.Pos.X) * zoom.PxPerUnit,
			Y: int(box.Pos.Y) * zoom.PxPerUnit,
		},
		Max: image.Point{
			X: int(box.Pos.X+box.Size.X) * zoom.PxPerUnit,
			Y: int(box.Pos.Y+box.Size.Y) * zoom.PxPerUnit,
		},
	}
}
