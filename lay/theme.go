package lay

import (
	"image/color"
	"math"

	"gioui.org/text"
	"gioui.org/unit"
)

type Theme struct {
	Shaper *text.Shaper

	Base      unit.Sp
	BaseRatio ScaleRatio

	TextSize   unit.Sp
	LineHeight unit.Sp
	TextRatio  ScaleRatio

	Palette

	FingerSize unit.Dp
}

func NewTheme(fontCollection []text.FontFace) *Theme {
	th := &Theme{
		Shaper: text.NewShaper(text.WithCollection(fontCollection)),
	}

	th.Palette = Palette{
		Fg: color.NRGBA{A: 0xFF},
		Bg: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}

	th.Base = unit.Sp(4)
	th.BaseRatio = 1.5

	th.TextSize = unit.Sp(16)
	th.LineHeight = unit.Sp(18)
	th.TextRatio = 1.4

	// 38dp is on the lower end of possible finger size.
	th.FingerSize = unit.Dp(38)

	return th
}

func (th *Theme) Scale(modifier float32) *Theme {
	scaled := *th
	scaled.Base *= unit.Sp(modifier)
	scaled.TextSize *= unit.Sp(modifier)
	scaled.LineHeight *= unit.Sp(modifier)
	return &scaled
}

func (th *Theme) Colorize(fg, bg color.NRGBA) *Theme {
	scaled := *th
	scaled.Fg = fg
	scaled.Bg = bg
	return &scaled
}

type Palette struct {
	Fg color.NRGBA
	Bg color.NRGBA
}

type Scale int8

type ScaleRatio = unit.Sp

const (
	None Scale = -0x7f

	Smaller Scale = -2
	Small   Scale = -1
	Default Scale = 0
	Big     Scale = 1
	Bigger  Scale = 1
)

func (s Scale) Value(base unit.Sp, ratio ScaleRatio) unit.Sp {
	if s == Default {
		return base
	}
	if s == None {
		return base
	}
	return base * unit.Sp(math.Pow(float64(ratio), float64(s)))
}
