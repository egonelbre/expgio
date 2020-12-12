package main

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/egonelbre/expgio/f32color"
)

var (
	groupIconPadding     = unit.Dp(3)
	iconCornerRadiusPx   = float32(8)
	maxIconBorderWidthPx = float32(8)

	iconActiveBorder   = f32color.HSL(0.5, 0.30, 0.40)
	iconInactiveBorder = f32color.HSL(0.5, 0.16, 0.20)

	groupsPanel = Panel{
		Axis: layout.Vertical,
		Size: unit.Dp(60),

		Background:  panelBackground,
		Border:      panelBorder,
		BorderWidth: borderWidth,
	}
)

type Groups struct {
	Active *Group
	Groups []*Group
	List   layout.List
}

func NewGroups(groups ...*Group) *Groups {
	gs := &Groups{}
	gs.Groups = groups
	gs.Active = groups[0]
	gs.List.Axis = layout.Vertical
	return gs
}

type Group struct {
	Icon     string
	Name     string
	Click    widget.Clickable
	Hover    Hoverable
	Hovering AnimationTimer
}

func NewGroup(icon, name string) *Group {
	g := &Group{Icon: icon, Name: name}
	g.Hovering.Duration = 150 * time.Millisecond
	return g
}

func (groups *Groups) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return groupsPanel.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return groups.List.Layout(gtx, len(groups.Groups), func(gtx layout.Context, i int) layout.Dimensions {
			return groups.layoutGroup(th, gtx, groups.Groups[i])
		})
	})
}

func (groups *Groups) layoutGroup(th *material.Theme, gtx layout.Context, group *Group) layout.Dimensions {
	sz := gtx.Constraints.Max.X
	gtx.Constraints = layout.Exact(image.Pt(sz, sz))

	if group.Click.Clicked() {
		groups.Active = group
	}
	isActive := group == groups.Active
	progress := group.Hovering.Update(gtx, isActive || group.Hover.Active())

	btn := material.Button(th, &group.Click, group.Icon)
	btn.Background = color.NRGBA{}

	inset := layout.UniformInset(unit.Px(maxIconBorderWidthPx))
	dimensions := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return BorderSmooth{
			Color:        ifthen(isActive, iconActiveBorder, iconInactiveBorder),
			CornerRadius: unit.Px(iconCornerRadiusPx),
			Width:        easeInOutQuad(progress) * maxIconBorderWidthPx,
		}.Layout(gtx,
			btn.Layout,
		)
	})

	_ = group.Hover.Layout(gtx)

	return dimensions
}

func easeInOutQuad(t float32) float32 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}
