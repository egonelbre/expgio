package main

import (
	"image/color"
)

type Palette struct{}

func (palette Palette) Color(value float64, span Range) color.NRGBA {
	p := (value - span.Min) / (span.Max - span.Min)

	sat, lit := float32(0.6), float32(0.6)
	// if highlight {
	// 	sat = 0.8
	// 	lit = 0.7
	// }
	return HSL(float32(p*0.5), sat, lit)
}
