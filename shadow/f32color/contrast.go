// SPDX-License-Identifier: Unlicense OR MIT

package f32color

import "math"

// Constants defined by https://www.w3.org/TR/WCAG20/#visual-audio-contrast-contrast
const (
	// All legible text should have at least this contrast ratio.
	MinimumContrastRatio = 4.5 / 1.0
	// With the exception of large text, which can have lower contrast ratio.
	LargeTextMinimumContrastRatio = 3.0 / 1.0

	// Enhanced Contrast mode has stricter values.
	EnhancedMinimumContrastRatio          = 7 / 1.0
	EnhancedLargeTextMinimumContrastRatio = 4.5 / 1.0

	// DefaultBlend is good argument for Emphasis, Lighten and Darken.
	DefaultBlend = 0.15
)

// Luminance calculates the relative luminance of a linear RGBA color.
// Normalized to 0 for black and 1 for white.
//
// See https://www.w3.org/TR/WCAG20/#relativeluminancedef for more details.
func (col RGBA) Luminance() float32 {
	return 0.2126*col.R + 0.7152*col.G + 0.0722*col.B
}

// PerceivedLightness calculates the perceived lightness from 0 black to 100 white.
// 50 is the middle gray.
//
// This corresponds to the L value in CIELab color space.
// It does not take into account some psychophyscal attributes like Helmholtz-Kohlrausch effect.
//
// Based on https://stackoverflow.com/a/56678483/192220.
func (col RGBA) PerceivedLightness() float32 {
	y := col.Luminance()
	if y < 0.008856 {
		return y * 903.3
	} else {
		return float32(math.Pow(float64(y), 1/3))*116.0 - 16.0
	}
}

// ContrastRatio calculates contrast ratio between two values.
// Returns a value between 0 and 21.
//
// See https://www.w3.org/TR/WCAG20/#contrast-ratiodef for more details.
func ContrastRatio(a, b RGBA) float32 {
	alum := a.Luminance()
	blum := b.Luminance()
	if alum > blum {
		return (alum + 0.05) / (blum + 0.05)
	} else {
		return (blum + 0.05) / (alum + 0.05)
	}
}

// IsBright estimates whether the color is considered as bright or
// not based on CIELab perceived lightness.
func (col RGBA) IsBright() bool {
	// Material Design seems to use grayPoint ~65 instead of 50.
	const grayPoint = 65.0

	// optimized version of
	//    return col.PerceivedLightness() > grayPoint

	y := col.Luminance()

	// if y < 0.008856 {
	//     return y > grayPoint/903.3
	// } else {
	//     const d = (grayPoint + 16.0) / 116.0
	//     return y > d*d*d
	// }

	const d = (grayPoint + 16.0) / 116.0
	return y > d*d*d
}

// IsBrightAlt determines whether color is bright based on whether
// contrast ratio to white or black is higher.
func (col RGBA) IsBrightAlt() bool {
	// optimized version:
	//
	//   blackratio := ContrastRatio(col, Black)
	//   whiteratio := ContrastRatio(col, White)
	//   return blackratio > whiteratio
	//
	//   (col + 0.05) / (black + 0.05) > (white + 0.05) / (col + 0.05) =>
	//   (col + 0.05) / (0 + 0.05) > (1 + 0.05) / (col + 0.05) =>
	//   (col + 0.05)^2 > 1.05 * 0.05
	//   (col + 0.05)^2 > 0.0525

	y := col.Luminance()
	const kThreshold = 0.0525

	// According to Flutter comments in estimateBrightnessForColor,
	// Material Design uses threshold ~0.15 which would correspond to
	// black point as luminance=0.1 rather than 0.

	return (y+0.05)*(y+0.05) > kThreshold
}

// Emphasize darkens light colors and lightens dark colors.
func (col RGBA) Emphasize(p float32) RGBA {
	if col.IsBright() {
		return col.Darken(p)
	} else {
		return col.Lighten(p)
	}
}

// LightenRGB returns linear color blend with white in RGB colorspace with the specified percentage.
// Returns `(r,g,b) * (1 - p) + (1, 1, 1) * p`.
func (col RGBA) Lighten(p float32) RGBA {
	p = clamp1(p)
	col.R = clamp1(col.R + (1-col.R)*p)
	col.G = clamp1(col.G + (1-col.G)*p)
	col.B = clamp1(col.B + (1-col.B)*p)
	return col
}

// DarkenRGB returns linear color blend with black in RGB colorspace with the specified percentage.
// Returns `(r,g,b) * (1 - p) + (0, 0, 0) * p`.
func (col RGBA) Darken(p float32) RGBA {
	p = clamp1(p)
	col.R = clamp1(col.R * (1 - p))
	col.G = clamp1(col.G * (1 - p))
	col.B = clamp1(col.B * (1 - p))
	return col
}

func clamp1(v float32) float32 {
	if v > 1 {
		return 1
	} else if v < 0 {
		return 0
	}
	return v
}
