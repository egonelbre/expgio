package main

import (
	"fmt"
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type HudManager struct {
	Theme *material.Theme

	Zoom    Zoom
	Diagram *Diagram

	Huds      []*HudControl
	Exclusive Hud

	Control struct {
		HudControl layout.List
	}
}

type HudControl struct {
	Visible widget.Bool
	Hud     Hud
}

func NewHudManager(theme *material.Theme) *HudManager {
	m := &HudManager{
		Theme:   theme,
		Diagram: NewDiagram(),
	}
	m.Zoom.Level = defaultZoom

	m.Control.HudControl.Axis = layout.Vertical

	m.Add(&NavHud{&m.Zoom})
	m.Add(&GridHud{})
	m.Add(&NodeHud{})
	connectionCreation := &ConnectionCreationHud{}
	m.Add(&PortHud{ShowAll: &connectionCreation.drawing})
	m.Add(&ConnectionHud{})
	m.Add(&NodeCreationHud{})
	m.Add(&ManipulationHud{})
	m.Add(connectionCreation)
	m.Add(&ZoomHud{Zoom: &m.Zoom})

	return m
}

func (m *HudManager) Add(hud Hud) {
	control := &HudControl{Hud: hud}
	control.Visible.Value = true
	m.Huds = append(m.Huds, control)
}

func (m *HudManager) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(m.LayoutControl),
		layout.Flexed(1, m.LayoutHuds),
	)
}

func (m *HudManager) LayoutControl(gtx layout.Context) layout.Dimensions {
	th := *m.Theme
	th.TextSize.V *= 0.8
	th.FingerSize.V *= 0.8

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			r := image.Rectangle{
				Max: image.Point{
					X: gtx.Constraints.Min.X + gtx.Px(unit.Dp(4)),
					Y: gtx.Constraints.Max.Y,
				},
			}
			paint.FillShape(gtx.Ops, PanelBackground, clip.Rect(r).Op())
			return layout.Dimensions{Size: r.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return m.Control.HudControl.Layout(gtx, len(m.Huds),
				func(gtx layout.Context, index int) layout.Dimensions {
					hud := m.Huds[index]
					return material.CheckBox(&th, &hud.Visible, fmt.Sprintf("%T", hud.Hud)).Layout(gtx)
				})
		}),
	)
}

func (m *HudManager) LayoutHuds(gtx layout.Context) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	clip.Rect{Max: gtx.Constraints.Max}.Add(gtx.Ops)

	for _, hud := range m.Huds {
		if !hud.Visible.Value {
			continue
		}

		var lgtx layout.Context
		if m.Exclusive == nil || m.Exclusive == hud.Hud {
			lgtx = gtx
		} else {
			lgtx = gtx.Disabled()
		}
		dgtx := NewContext(lgtx, m.Theme, &m.Zoom, m.Diagram)
		hud.Hud.Layout(dgtx)
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

type Hud interface {
	Layout(gtx *Context)
}

type Context struct {
	Transform
	layout.Context
	Theme   *material.Theme
	Diagram *Diagram
}

func NewContext(gtx layout.Context, th *material.Theme, zoom *Zoom, diagram *Diagram) *Context {
	return &Context{
		Transform: NewTransform(gtx, zoom),
		Context:   gtx,
		Theme:     th,
		Diagram:   diagram,
	}
}

type Transform struct {
	Dp        int
	PxPerUnit int
}

func NewTransform(gtx layout.Context, zoom *Zoom) Transform {
	px := gtx.Px(unit.Dp(30))
	px = (px / 24) * 24 // make it divisible by 2,3,4,6,12
	px = int(float32(px) * zoom.Multiplier())
	return Transform{
		Dp:        gtx.Px(unit.Dp(1)),
		PxPerUnit: px,
	}
}

func (tr *Transform) Px(v Unit) int {
	return tr.PxPerUnit * int(v)
}

func (tr *Transform) Pt(v Vector) image.Point {
	return image.Point{
		X: int(v.X) * tr.PxPerUnit,
		Y: int(v.Y) * tr.PxPerUnit,
	}
}

func (tr *Transform) FPt(v Vector) f32.Point {
	return f32.Point{
		X: float32(int(v.X) * tr.PxPerUnit),
		Y: float32(int(v.Y) * tr.PxPerUnit),
	}
}

func (tr *Transform) Inv(p image.Point) Vector {
	return Vector{
		X: Unit(p.X / tr.PxPerUnit),
		Y: Unit(p.Y / tr.PxPerUnit),
	}
}

func (tr *Transform) FInv(p f32.Point) Vector {
	return Vector{
		X: Unit(int(p.X) / tr.PxPerUnit),
		Y: Unit(int(p.Y) / tr.PxPerUnit),
	}
}

func (tr *Transform) Bounds(box Box) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(box.Pos.X) * tr.PxPerUnit,
			Y: int(box.Pos.Y) * tr.PxPerUnit,
		},
		Max: image.Point{
			X: int(box.Pos.X+box.Size.X) * tr.PxPerUnit,
			Y: int(box.Pos.Y+box.Size.Y) * tr.PxPerUnit,
		},
	}
}
