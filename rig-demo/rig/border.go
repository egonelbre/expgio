package rig

import "gioui.org/f32"

type Border struct {
	Locked     bool
	Horizontal bool
	Corners    []*Corner
}

func (border *Border) IsLocked() bool {
	return (border != nil) && border.Locked
}

func (border *Border) Index(corner *Corner) int {
	if border == nil {
		return -1
	}
	for i, c := range border.Corners {
		if c == corner {
			return i
		}
	}
	return -1
}

func (border *Border) Insert(corner *Corner) {
	border.Corners = append(border.Corners, corner)
	border.Sort()
}

func (border *Border) Neighbor(corner *Corner, di int) *Corner {
	if border == nil {
		return nil
	}

	i := border.Index(corner)
	if i < 0 {
		return nil
	}

	ti := i + di
	if 0 <= ti && ti < len(border.Corners) {
		return border.Corners[ti]
	}

	return nil
}

func (border *Border) Center() f32.Point { return border.Min().Add(border.Max()).Mul(0.5) }

func (border *Border) First() *Corner { return border.Corners[0] }
func (border *Border) Last() *Corner  { return border.Corners[len(border.Corners)-1] }

func (border *Border) Min() f32.Point { return border.First().Pos }
func (border *Border) Max() f32.Point { return border.Last().Pos }

func (border *Border) Sort() { SortCorners(border.Corners) }
