package main

import (
	"gioui.org/layout"
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
		Diagram: NewDemoDiagram(),
	}

	m.Huds = append(m.Huds,
		&GridDisplay{},
		&NodeDisplay{},
		// &Selecter{}
		&PortDisplay{},
		&ConnectionDisplay{},
		// &ConnectionCreator{}
		// &NodeMover{}
		// &NodeCreator{}
		// &NodeDeleter{}
		// &NodeOrderer{}
	)
	return m
}

func (m *HudManager) Layout(gtx layout.Context) layout.Dimensions {
	for _, hud := range m.Huds {
		if !hud.Enabled() {
			continue
		}
		if m.Exclusive == nil || m.Exclusive == hud {
			hud.Layout(m.Diagram, gtx)
		} else {
			hud.Layout(m.Diagram, gtx.Disabled())
		}
	}

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

type Hud interface {
	Enabled() bool
	Layout(diagram *Diagram, gtx layout.Context)
}
