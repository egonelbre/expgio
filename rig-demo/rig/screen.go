package rig

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

const (
	EditorMinSize = 50
)

type Screen struct {
	Registry *Registry
	Bounds   Rectangle
	Rig      *Rig
	Areas    []*Area
}

func New() *Screen {
	screen := &Screen{}
	screen.Registry = NewRegistry()

	screen.Rig = NewRig()
	return screen
}

func (screen *Screen) Layout(gtx layout.Context) layout.Dimensions {
	size := gtx.Constraints.Min
	gtx.Constraints = layout.Exact(size)
	cornerRadius := gtx.Metric.Dp(RigCornerRadius)
	joinerSize := gtx.Metric.Dp(RigJoinerSize)
	borderRadius := gtx.Metric.Dp(RigBorderRadius)

	_ = cornerRadius
	_ = joinerSize
	_ = borderRadius

	for _, corner := range screen.Rig.Corners {
		p := corner.Px(size)

		/*
			r := inflate(p, cornerRadius)
			canCapture := !corner.IsLocked() && gtx.Input.Mouse.Capture == nil && r.Contains(gtx.Input.Mouse.Pos)
			if canCapture {
				gtx.Input.Mouse.SetCaptureCursor(CrosshairCursor)
				if gtx.Input.Mouse.Pressed {
					gtx.Input.Mouse.Capture = (&Resizer{
						Screen:     screen,
						Start:      corner.Pos,
						Horizontal: corner.Horizontal,
						Vertical:   corner.Vertical,
					}).Capture
				}
			}
		*/

		// Draw the action corner, when there's an area to the bottom left.
		if corner.SideLeft() != nil && corner.SideBottom() != nil {
			r := image.Rectangle{
				Min: image.Point{
					X: p.X - joinerSize,
					Y: p.Y - borderRadius,
				},
				Max: image.Point{
					X: p.X - borderRadius,
					Y: p.Y + joinerSize,
				},
			}

			//canCapture := gtx.Input.Mouse.Capture == nil && r.Contains(gtx.Input.Mouse.Pos)
			canCapture := false
			if canCapture {
				// gtx.Input.Mouse.SetCaptureCursor(CrosshairCursor)
				// if !gtx.Input.Mouse.Pressed {
				// 	gtx.Draw.FillRect(&r, RigJoinerHighlightColor)
				// } else {
				// 	gtx.Input.Mouse.Capture = (&Joiner{
				// 		Screen: screen,
				// 		Corner: corner,
				// 	}).Capture
				// }
			} else {
				paint.FillShape(gtx.Ops, RigJoinerColor, clip.Rect(r).Op())
			}
		}
	}

	for _, border := range screen.Rig.Borders {
		min := border.First().Px(size)
		max := border.Last().Px(size)

		r := image.Rectangle{
			Min: min,
			Max: max,
		}.Inset(-borderRadius) // TODO: based on horizontal vs vertical

		paint.FillShape(gtx.Ops, RigBorderColor, clip.Rect(r).Op())

		/*
			min := gtx.Area.ToGlobal(border.Min())
			max := gtx.Area.ToGlobal(border.Max())

			r := Rectangle{min, max}.Inflate(RigBorderRadius)

			canCapture := !border.Locked && gtx.Input.Mouse.Capture == nil && r.Contains(gtx.Input.Mouse.Pos)
			if canCapture {
				var horz, vert *Border
				if border.Horizontal {
					gtx.Input.Mouse.SetCaptureCursor(VResizeCursor)
					horz = border
				} else {
					gtx.Input.Mouse.SetCaptureCursor(HResizeCursor)
					vert = border
				}

				gtx.Draw.FillRect(&r, RigBorderHighlightColor)
				if gtx.Input.Mouse.Pressed {
					gtx.Input.Mouse.Capture = (&Resizer{
						Screen:     screen,
						Start:      border.Center(),
						Horizontal: horz,
						Vertical:   vert,
					}).Capture
				}
			} else {
				gtx.Draw.FillRect(&r, RigBorderColor)
			}
		*/
	}

	return layout.Dimensions{
		Size: size,
	}
}
