package main

import (
	"fmt"
	"image"

	"gioui.org/io/pointer"
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

	View struct {
		Huds   layout.List
		Styles layout.List
	}
}

type HudControl struct {
	Visible widget.Bool
	Hud     Hud
}

type Hud interface {
	Layout(gtx *Context)
}

func NewHudManager(theme *material.Theme) *HudManager {
	m := &HudManager{
		Theme:   theme,
		Diagram: NewDiagram(),
	}
	m.Zoom.Level = defaultZoom

	m.View.Huds.Axis = layout.Vertical
	m.View.Styles.Axis = layout.Vertical

	// m.Add(&NavHud{&m.Zoom})
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
		layout.Rigid(m.LayoutPaint),
	)
}

func (m *HudManager) LayoutPaint(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			r := image.Rectangle{
				Max: image.Point{
					X: gtx.Constraints.Min.X,
					Y: gtx.Constraints.Max.Y,
				},
			}
			paint.FillShape(gtx.Ops, PanelBackground, clip.Rect(r).Op())
			return layout.Dimensions{Size: r.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return m.View.Styles.Layout(gtx, len(Tango),
				func(gtx layout.Context, index int) layout.Dimensions {
					defer op.Save(gtx.Ops).Load()

					type paintTag *Style
					style := Tango[index]
					tag := paintTag(style)

					size := image.Point{
						X: gtx.Px(m.Theme.FingerSize),
						Y: gtx.Px(m.Theme.FingerSize),
					}
					paint.FillShape(gtx.Ops, style.Fill, clip.Rect{Max: size}.Op())

					pointer.Rect(image.Rectangle{Max: size}).Add(gtx.Ops)
					pointer.InputOp{
						Tag:   tag,
						Types: pointer.Press,
					}.Add(gtx.Ops)

					for _, ev := range gtx.Events(tag) {
						if ev, ok := ev.(pointer.Event); ok {
							switch ev.Type {
							case pointer.Press:
								for selected := range m.Diagram.Selection.Selected {
									switch sel := selected.(type) {
									case *Node:
										sel.Style = style
									}
								}
							}
						}
					}

					return layout.Dimensions{
						Size: size,
					}
				})
		}),
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
			return m.View.Huds.Layout(gtx, len(m.Huds),
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
