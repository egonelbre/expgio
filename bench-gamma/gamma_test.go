package benchgamma_test

import (
	"math"
	"testing"
)

func gamma(r, g, b, a uint32) [4]float32 {
	color := [4]float32{float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff}
	// Assume that image.Uniform colors are in sRGB space. Linearize.
	for i := 0; i <= 2; i++ {
		c := color[i]
		// Use the formula from EXT_sRGB.
		if c <= 0.04045 {
			c = c / 12.92
		} else {
			c = float32(math.Pow(float64((c+0.055)/1.055), 2.4))
		}
		color[i] = c
	}
	return color
}

func gamma2(r, g, b, a uint32) [4]float32 {
	rf := float32(r) / 0xffff
	gf := float32(g) / 0xffff
	bf := float32(b) / 0xffff
	af := float32(a) / 0xffff

	if rf <= 0.04045 {
		rf = rf / 12.92
	} else {
		rf = float32(math.Pow(float64((rf+0.055)/1.055), 2.4))
	}

	if gf <= 0.04045 {
		gf = gf / 12.92
	} else {
		gf = float32(math.Pow(float64((gf+0.055)/1.055), 2.4))
	}

	if bf <= 0.04045 {
		bf = bf / 12.92
	} else {
		bf = float32(math.Pow(float64((bf+0.055)/1.055), 2.4))
	}

	return [4]float32{rf, gf, bf, af}
}

func gamma2x(r, g, b, a uint32) [4]float32 {
	return [4]float32{
		linearize(float32(r) / 0xffff),
		linearize(float32(g) / 0xffff),
		linearize(float32(b) / 0xffff),
		float32(a) / 0xffff,
	}
}

func linearize(v float32) float32 {
	if v <= 0.04045 {
		return v / 12.92
	} else {
		return float32(math.Pow(float64((v+0.055)/1.055), 2.4))
	}
}

type RGBA struct {
	R, G, B, A float32
}

func gamma3(r, g, b, a uint32) RGBA {
	rf := float32(r) / 0xffff
	gf := float32(g) / 0xffff
	bf := float32(b) / 0xffff
	af := float32(a) / 0xffff

	if rf <= 0.04045 {
		rf = rf / 12.92
	} else {
		rf = float32(math.Pow(float64((rf+0.055)/1.055), 2.4))
	}

	if gf <= 0.04045 {
		gf = gf / 12.92
	} else {
		gf = float32(math.Pow(float64((gf+0.055)/1.055), 2.4))
	}

	if bf <= 0.04045 {
		bf = bf / 12.92
	} else {
		bf = float32(math.Pow(float64((bf+0.055)/1.055), 2.4))
	}

	return RGBA{rf, gf, bf, af}
}

func BenchmarkGamma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for r := uint32(0); r < 256; r += 64 {
			for g := uint32(0); g < 256; g += 64 {
				for b := uint32(0); b < 256; b += 64 {
					for a := uint32(0); a < 256; a += 64 {
						_ = gamma(r, g, b, a)
					}
				}
			}
		}
	}
}

func BenchmarkGamma2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for r := uint32(0); r < 256; r += 64 {
			for g := uint32(0); g < 256; g += 64 {
				for b := uint32(0); b < 256; b += 64 {
					for a := uint32(0); a < 256; a += 64 {
						_ = gamma2(r, g, b, a)
					}
				}
			}
		}
	}
}

func BenchmarkGamma2x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for r := uint32(0); r < 256; r += 64 {
			for g := uint32(0); g < 256; g += 64 {
				for b := uint32(0); b < 256; b += 64 {
					for a := uint32(0); a < 256; a += 64 {
						_ = gamma2x(r, g, b, a)
					}
				}
			}
		}
	}
}
func BenchmarkGamma3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for r := uint32(0); r < 256; r += 64 {
			for g := uint32(0); g < 256; g += 64 {
				for b := uint32(0); b < 256; b += 64 {
					for a := uint32(0); a < 256; a += 64 {
						_ = gamma3(r, g, b, a)
					}
				}
			}
		}
	}
}
