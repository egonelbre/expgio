package rig

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/unit"
)

var (
	RigJoinerColor          = color.NRGBA{R: 0xA0, G: 0xA0, B: 0xA0, A: 0xFF}
	RigJoinerHighlightColor = color.NRGBA{R: 0xE0, G: 0xA0, B: 0xA0, A: 0xFF}
	RigBorderColor          = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	RigBorderHighlightColor = color.NRGBA{R: 0xA0, G: 0x80, B: 0x80, A: 0xFF}
	RigCornerColor          = color.NRGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xFF}
	RigCornerHighlightColor = color.NRGBA{R: 0xA0, G: 0x40, B: 0x40, A: 0xFF}

	RigBackground = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}

	RigJoinerSize    = unit.Dp(20)
	RigTriggerSize   = unit.Dp(40)
	RigTriggerRadius = RigTriggerSize
	RigBorderRadius  = unit.Dp(2)
	RigCornerRadius  = unit.Dp(3)
)

type Rig struct {
	Corners []*Corner
	Borders []*Border
}

func NewRig() *Rig {
	topLeft := &Corner{Pos: f32.Point{X: 0, Y: 0}}
	topRight := &Corner{Pos: f32.Point{X: 1, Y: 0}}
	bottomLeft := &Corner{Pos: f32.Point{X: 0, Y: 1}}
	bottomRight := &Corner{Pos: f32.Point{X: 1, Y: 1}}

	leftSide := &Border{Locked: true, Horizontal: false, Corners: []*Corner{bottomLeft, topLeft}}
	topSide := &Border{Locked: true, Horizontal: true, Corners: []*Corner{topLeft, topRight}}
	rightSide := &Border{Locked: true, Horizontal: false, Corners: []*Corner{topRight, bottomRight}}
	bottomSide := &Border{Locked: true, Horizontal: true, Corners: []*Corner{bottomRight, bottomLeft}}

	leftSide.Sort()
	topSide.Sort()
	rightSide.Sort()
	bottomSide.Sort()

	topLeft.Vertical = leftSide
	topLeft.Horizontal = topSide
	topRight.Horizontal = topSide
	topRight.Vertical = rightSide
	bottomRight.Vertical = rightSide
	bottomRight.Horizontal = bottomSide
	bottomLeft.Horizontal = bottomSide
	bottomLeft.Vertical = leftSide

	return &Rig{
		Corners: []*Corner{topLeft, topRight, bottomLeft, bottomRight},
		Borders: []*Border{leftSide, topSide, rightSide, bottomSide},
	}
}

func (rig *Rig) SplitVertically(topRight *Corner, posX float32) *Border {
	_, bottomRight := topRight.BlockingVertical(true, false)
	topSide := topRight.Horizontal
	bottomSide := bottomRight.Horizontal

	centerTop := &Corner{Pos: f32.Point{X: posX, Y: topRight.Pos.Y}}
	centerBottom := &Corner{Pos: f32.Point{X: posX, Y: bottomRight.Pos.Y}}
	split := &Border{Locked: false, Horizontal: false, Corners: []*Corner{centerTop, centerBottom}}
	split.Sort()

	centerTop.Horizontal, centerTop.Vertical = topSide, split
	topSide.Insert(centerTop)

	centerBottom.Horizontal, centerBottom.Vertical = bottomSide, split
	bottomSide.Insert(centerBottom)

	rig.Corners = append(rig.Corners, split.Corners...)
	SortCorners(rig.Corners)
	rig.Borders = append(rig.Borders, split)

	return split
}

func (rig *Rig) SplitHorizontally(topRight *Corner, posY float32) *Border {
	topLeft, _ := topRight.BlockingHorizontal(false, true)
	rightSide := topRight.Vertical
	leftSide := topLeft.Vertical

	centerLeft := &Corner{Pos: f32.Point{X: topLeft.Pos.X, Y: posY}}
	centerRight := &Corner{Pos: f32.Point{X: topRight.Pos.X, Y: posY}}

	split := &Border{Locked: false, Horizontal: true, Corners: []*Corner{centerLeft, centerRight}}
	split.Sort()

	centerLeft.Vertical, centerLeft.Horizontal = leftSide, split
	leftSide.Insert(centerLeft)

	centerRight.Vertical, centerRight.Horizontal = rightSide, split
	rightSide.Insert(centerRight)

	rig.Corners = append(rig.Corners, split.Corners...)
	SortCorners(rig.Corners)
	rig.Borders = append(rig.Borders, split)

	return split
}
