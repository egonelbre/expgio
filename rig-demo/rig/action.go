package rig

import (
	"gioui.org/f32"
	"gioui.org/layout"
)

type Joiner struct {
	Screen *Screen
	Corner *Corner

	inited    bool
	splitArea Rectangle
	limit     Rectangle

	canMergeVertically   bool
	canMergeHorizontally bool
	canFullscreen        bool
}

func (act *Joiner) init(gtx layout.Context) {
	left, _ := act.Corner.BlockingHorizontal(false, true)
	_, bottom := act.Corner.BlockingVertical(true, false)

	act.splitArea = Rectangle{
		Min: left.Pos,
		Max: bottom.Pos,
	}

	if act.Corner.SideRight() != nil {
		_, bottomBlockLeft := act.Corner.BlockingVertical(true, false)
		_, bottomBlockRight := act.Corner.BlockingVertical(false, true)
		act.canMergeHorizontally = bottomBlockLeft == bottomBlockRight
	}

	if act.Corner.SideTop() != nil {
		leftBlockTop, _ := act.Corner.BlockingHorizontal(true, false)
		leftBlockBottom, _ := act.Corner.BlockingHorizontal(false, true)
		act.canMergeVertically = leftBlockTop == leftBlockBottom
	}
}

func (act *Joiner) Capture(gtx layout.Context) bool {
	/*
		if !gtx.Input.Mouse.Down {
			return true
		}
		if !act.inited {
			act.inited = true
			act.init(gtx)
		}

		cornerPos := act.Screen.Bounds.ToGlobal(act.Corner.Pos)
		delta := gtx.Input.Mouse.Pos.Sub(cornerPos)

		if delta.Len() < RigJoinerSize {
			// for canceling
			return false
		}

		// splitting action
		if delta.X < 0 && delta.Y > 0 {
			return act.trySplit(gtx, delta)
		}

		// merge to right
		if 0 < delta.X && 0 < delta.Y && act.canMergeHorizontally {
			return act.tryMerge(gtx, delta)
		}

		// merge to top
		if delta.X < 0 && delta.Y < 0 && act.canMergeVertically {
			return act.tryMerge(gtx, delta)
		}

		// fullscreen on release
		if 0 < delta.X && delta.Y < 0 && act.canFullscreen {
			return false
		}
	*/
	return false
}

func (act *Joiner) trySplit(gtx layout.Context, delta f32.Point) bool {
	/*
		dx := -delta.X
		dy := delta.Y

		r := act.Screen.Bounds.Subset(act.splitArea)

		if dx < dy {
			// do we have enough room for two areas?
			if r.Dy() < 2*RigTriggerSize {
				return false
			}
			gtx.Input.Mouse.Cursor = CrosshairCursor

			r.Max.Y = r.Max.Y - RigTriggerSize
			if r.Min.Y+dy < r.Max.Y {
				r.Max.Y = r.Min.Y + dy
			}

			alpha := Sat8((dy - RigJoinerSize) / (RigTriggerSize - RigJoinerSize))
			gtx.Hover.FillRect(&r, RigBackground.WithAlpha(alpha))

			if dy > RigTriggerSize {
				// TODO: fix don't use mouse pos, it might be outside of limits
				rp := act.Screen.Bounds.ToRelative(gtx.Input.Mouse.Pos)
				split := act.Screen.Rig.SplitHorizontally(act.Corner, rp.Y)
				gtx.Input.Mouse.Capture = (&Resizer{
					Screen:     act.Screen,
					Start:      split.Center(),
					Horizontal: split,
					Vertical:   nil,
				}).Capture
			}
		} else {
			// do we have enough room for two areas?
			if r.Dx() < 2*RigTriggerSize {
				return false
			}
			gtx.Input.Mouse.Cursor = CrosshairCursor

			r.Min.X = r.Min.X + RigTriggerSize
			if r.Min.X < r.Max.X-dx {
				r.Min.X = r.Max.X - dx
			}

			alpha := Sat8((dx - RigJoinerSize) / (RigTriggerSize - RigJoinerSize))
			gtx.Hover.FillRect(&r, RigBackground.WithAlpha(alpha))

			if dx > RigTriggerSize {
				// TODO: fix don't use mouse pos, it might be outside of limits
				rp := act.Screen.Bounds.ToRelative(gtx.Input.Mouse.Pos)
				split := act.Screen.Rig.SplitVertically(act.Corner, rp.X)
				gtx.Input.Mouse.Capture = (&Resizer{
					Screen:     act.Screen,
					Start:      split.Center(),
					Horizontal: nil,
					Vertical:   split,
				}).Capture
			}
		}
	*/
	return false
}

func (act *Joiner) tryMerge(gtx layout.Context, delta f32.Point) bool {
	/*
		dx := delta.X
		dy := -delta.Y

		r := act.Screen.Bounds.Subset(act.splitArea)
		if dx >= 0 {
			// do we have enough room for two areas?
			gtx.Input.Mouse.Cursor = CrosshairCursor

			r.Min.X = r.Max.X
			r.Max.X += dx

			alpha := Sat8((dx - RigJoinerSize) / (RigTriggerSize - RigJoinerSize))
			gtx.Hover.FillRect(&r, RigBackground.WithAlpha(alpha))
			// TODO: draw arrow

			if false && dy > RigTriggerSize {
				// TODO: fix don't use mouse pos, it might be outside of limits
				rp := act.Screen.Bounds.ToRelative(gtx.Input.Mouse.Pos)
				split := act.Screen.Rig.SplitHorizontally(act.Corner, rp.Y)
				gtx.Input.Mouse.Capture = (&Resizer{
					Screen:     act.Screen,
					Start:      split.Center(),
					Horizontal: split,
					Vertical:   nil,
				}).Capture
			}
		}

		if dy >= 0 {
			gtx.Input.Mouse.Cursor = CrosshairCursor

			r.Max.Y = r.Min.Y
			r.Min.Y -= dy

			alpha := Sat8((dy - RigJoinerSize) / (RigTriggerSize - RigJoinerSize))
			gtx.Hover.FillRect(&r, RigBackground.WithAlpha(alpha))
			// TODO: draw arrow

			if false && dx > RigTriggerSize {
				// TODO: fix don't use mouse pos, it might be outside of limits
				rp := act.Screen.Bounds.ToRelative(gtx.Input.Mouse.Pos)
				split := act.Screen.Rig.SplitVertically(act.Corner, rp.X)
				gtx.Input.Mouse.Capture = (&Resizer{
					Screen:     act.Screen,
					Start:      split.Center(),
					Horizontal: nil,
					Vertical:   split,
				}).Capture
			}
		}
	*/
	return false
}

type Resizer struct {
	Screen     *Screen
	Start      f32.Point
	Horizontal *Border
	Vertical   *Border

	inited bool
	area   Rectangle
}

func (act *Resizer) init(gtx layout.Context) {
	act.area = Rectangle{Max: f32.Pt(1, 1)}
	if act.Vertical != nil {
		for _, corner := range act.Vertical.Corners {
			checkTop := act.Vertical.First() != corner
			checkBottom := act.Vertical.Last() != corner

			left, right := corner.BlockingHorizontal(checkTop, checkBottom)
			if left != nil && act.area.Min.X < left.Pos.X {
				act.area.Min.X = left.Pos.X
			}
			if right != nil && right.Pos.X < act.area.Max.X {
				act.area.Max.X = right.Pos.X
			}
		}
	}
	if act.Horizontal != nil {
		for _, corner := range act.Horizontal.Corners {
			checkLeft := act.Horizontal.First() != corner
			checkRight := act.Horizontal.Last() != corner

			top, bottom := corner.BlockingVertical(checkLeft, checkRight)
			if top != nil && act.area.Min.Y < top.Pos.Y {
				act.area.Min.Y = top.Pos.Y
			}
			if bottom != nil && bottom.Pos.Y < act.area.Max.Y {
				act.area.Max.Y = bottom.Pos.Y
			}
		}
	}
}

func (act *Resizer) Capture(gtx layout.Context) bool {
	/*
		if !gtx.Input.Mouse.Down {
			for _, border := range act.Screen.Rig.Borders {
				border.Sort()
			}
			return true
		}

		if !act.inited {
			act.inited = true
			act.init(gtx)
		}

		limiter := act.Screen.Bounds.Subset(act.area).Deflate(RigTriggerRadius)
		if limiter.Max.X < limiter.Min.X || limiter.Max.Y < limiter.Min.Y {
			return true
		}

		gtx.Hover.StrokeRect(&limiter, 2, color.NRGBA{R: 0, G: 0, B: 0, A: 0x20})

		clampedMouse := limiter.ClosestPoint(gtx.Input.Mouse.Pos)
		p := act.Screen.Bounds.ToRelative(clampedMouse)

		if act.Horizontal != nil {
			for _, corner := range act.Horizontal.Corners {
				corner.Pos.Y = p.Y
			}
		}

		if act.Vertical != nil {
			for _, corner := range act.Vertical.Corners {
				corner.Pos.X = p.X
			}
		}
	*/
	return false
}

func Sat8(v float32) uint8 {
	v *= 255.0
	if v >= 255 {
		return 255
	} else if v <= 0 {
		return 0
	}
	return uint8(v)
}
