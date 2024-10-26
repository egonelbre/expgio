package main

import (
	"math"
)

// niceNumber calculates nicely rounded number.
func niceNumber(span float64, round bool) float64 {
	exp := math.Floor(math.Log10(span))
	frac := span / math.Pow(10, exp)
	var nice float64
	if round {
		switch {
		case frac < 1.5:
			nice = 1
		case frac < 3:
			nice = 2
		case frac < 7:
			nice = 5
		default:
			nice = 10
		}
	} else {
		switch {
		case frac <= 1:
			nice = 1
		case frac <= 2:
			nice = 2
		case frac <= 5:
			nice = 5
		default:
			nice = 10
		}
	}
	return nice * math.Pow(10, exp)
}

// lerp interpolates between min and max using p=[0,1].
func lerp(p, min, max float64) float64 {
	return p*(max-min) + min
}

// invlerp return inverse lerp from min, max and the value.
func invlerp(v, min, max float64) float64 {
	return (v - min) / (max - min)
}

// lerpUnit interpolates between min and max using p=[-1,1].
func lerpUnit(p, min, max float64) float64 {
	pu := (p + 1) * 0.5
	return pu*(max-min) + min
}

// cubicPulse calculates cubic-pulse function at a given location.
func cubicPulse(center, radius, invradius, at float64) float64 {
	at = at - center
	if at < 0 {
		at = -at
	}
	if at > radius {
		return 0
	}
	at *= invradius
	return 1 - at*at*(3-2*at)
}

// IntssToFloat64s convers a slice of ints to float64-s.
func IntsToFloat64s(xs []int) []float64 {
	r := make([]float64, len(xs))
	for i := range xs {
		r[i] = float64(xs[i])
	}
	return r
}

// Int32sToFloat64s convers a slice of int32-s to float64-s.
func Int32sToFloat64s(xs []int32) []float64 {
	r := make([]float64, len(xs))
	for i := range xs {
		r[i] = float64(xs[i])
	}
	return r
}

// Int64sToFloat64s convers a slice of int64-s to float64-s.
func Int64sToFloat64s(xs []int64) []float64 {
	r := make([]float64, len(xs))
	for i := range xs {
		r[i] = float64(xs[i])
	}
	return r
}
