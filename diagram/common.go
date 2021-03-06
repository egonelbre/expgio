package main

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
)

type Vector = image.Point

var ScaleDp = unit.Dp(30)

type Set map[interface{}]struct{}

func NewSet() Set { return make(Set) }

func (s Set) Contains(v interface{}) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Include(v interface{}) { s[v] = struct{}{} }
func (s Set) Exclude(v interface{}) { delete(s, v) }

func VectorPx(v Vector, gtx layout.Context) image.Point {
	px := gtx.Px(ScaleDp)
	return image.Point{
		X: px * v.X,
		Y: px * v.Y,
	}
}
