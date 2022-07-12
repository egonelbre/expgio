package rig

import (
	"image"
	"sort"

	"gioui.org/f32"
)

type Corner struct {
	Pos f32.Point // 0..1 position

	Horizontal *Border
	Vertical   *Border
}

func (a *Corner) Px(size image.Point) image.Point {
	return image.Point{
		X: int(float32(size.X) * a.Pos.X),
		Y: int(float32(size.Y) * a.Pos.Y),
	}
}

func (a *Corner) IsLocked() bool {
	return a.Horizontal.IsLocked() || a.Vertical.IsLocked()
}

func (a *Corner) SideLeft() *Border {
	if a == nil || a.Horizontal == nil || a.Horizontal.First() == a {
		return nil
	}
	return a.Horizontal
}

func (a *Corner) SideTop() *Border {
	if a == nil || a.Vertical == nil || a.Vertical.First() == a {
		return nil
	}
	return a.Vertical
}

func (a *Corner) SideRight() *Border {
	if a == nil || a.Horizontal == nil || a.Horizontal.Last() == a {
		return nil
	}
	return a.Horizontal
}

func (a *Corner) SideBottom() *Border {
	if a == nil || a.Vertical == nil || a.Vertical.Last() == a {
		return nil
	}
	return a.Vertical
}

func (a *Corner) CornerLeft() *Corner   { return a.SideLeft().Neighbor(a, -1) }
func (a *Corner) CornerTop() *Corner    { return a.SideTop().Neighbor(a, -1) }
func (a *Corner) CornerRight() *Corner  { return a.SideRight().Neighbor(a, 1) }
func (a *Corner) CornerBottom() *Corner { return a.SideBottom().Neighbor(a, 1) }

func (a *Corner) BlockingHorizontal(checkTop, checkBottom bool) (left, right *Corner) {
	index := a.Horizontal.Index(a)
	neighbors := a.Horizontal.Corners

	for k := index - 1; k >= 0; k-- {
		n := neighbors[k]
		if checkTop && n.SideTop() != nil {
			left = n
			break
		}
		if checkBottom && n.SideBottom() != nil {
			left = n
			break
		}
	}

	for k := index + 1; k < len(neighbors); k++ {
		n := neighbors[k]
		if checkTop && n.SideTop() != nil {
			right = n
			break
		}
		if checkBottom && n.SideBottom() != nil {
			right = n
			break
		}
	}
	return
}

func (a *Corner) BlockingVertical(checkLeft, checkRight bool) (top, bottom *Corner) {
	index := a.Vertical.Index(a)
	neighbors := a.Vertical.Corners

	for k := index - 1; k >= 0; k-- {
		n := neighbors[k]
		if checkLeft && n.SideLeft() != nil {
			top = n
			break
		}
		if checkRight && n.SideRight() != nil {
			top = n
			break
		}
	}

	for k := index + 1; k < len(neighbors); k++ {
		n := neighbors[k]
		if checkLeft && n.SideLeft() != nil {
			bottom = n
			break
		}
		if checkRight && n.SideRight() != nil {
			bottom = n
			break
		}
	}
	return
}

func (a *Corner) Less(b *Corner) bool {
	if a.Pos.X != b.Pos.X {
		return a.Pos.X < b.Pos.X
	}
	return a.Pos.Y < b.Pos.Y
}

func SortCorners(corners []*Corner) {
	sort.Slice(corners, func(i, k int) bool {
		return corners[i].Less(corners[k])
	})
}
