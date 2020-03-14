package main

import (
	"math"

	"gioui.org/f32"
)

const (
	Tau = 2 * math.Pi
)

func Bounce(t, start, end, duration float32) float32 {
	t = t / duration
	ts := t * t
	tc := ts * t
	return start + end*(33*tc*ts+-106*ts*ts+126*tc+-67*ts+15*t)
}

func SegmentNormal(a, b f32.Point) f32.Point {
	return Rotate(b.Sub(a))
}

func Rotate(a f32.Point) f32.Point {
	return f32.Point{X: -a.Y, Y: a.X}
}

func ScaleTo(a f32.Point, r float32) f32.Point {
	x := Len(a)
	return a.Mul(r / x)
}

func LerpPoint(p float32, a, b f32.Point) f32.Point {
	return f32.Point{
		X: Lerp(p, a.X, b.X),
		Y: Lerp(p, a.Y, b.Y),
	}
}

func Lerp(p float32, a, b float32) float32 {
	return p*(b-a) + a
}

func Floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}

func Sin(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func Sincos(v float32) (float32, float32) {
	sn, cs := math.Sincos(float64(v))
	return float32(sn), float32(cs)
}

func Sqrt(p float32) float32 {
	return float32(math.Sqrt(float64(p)))
}

func Pow(a, b float32) float32 {
	return float32(math.Pow(float64(a), float64(b)))
}

func Neg(p f32.Point) f32.Point {
	return f32.Point{X: -p.X, Y: -p.Y}
}

func Len(p f32.Point) float32 {
	return Sqrt(p.X*p.X + p.Y*p.Y)
}

func Dot(a, b f32.Point) float32 {
	return a.X*b.X + a.Y*b.Y
}

func Mod(a, b float32) float32 {
	return float32(math.Mod(float64(a), float64(b)))
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
