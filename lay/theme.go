package lay

import (
	"image/color"
	"math"

	"gioui.org/text"
	"gioui.org/unit"
)

type Theme struct {
	Shaper text.Shaper

	Base      unit.Value
	BaseRatio ScaleRatio

	TextSize   unit.Value
	LineHeight unit.Value
	TextRatio  ScaleRatio

	Palette

	FingerSize unit.Value
}

func NewTheme(fontCollection []text.FontFace) *Theme {
	th := &Theme{
		Shaper: text.NewCache(fontCollection),
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
	scaled.Base.V *= modifier
	scaled.TextSize.V *= modifier
	scaled.LineHeight.V *= modifier
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

type ScaleRatio = float32

const (
	None Scale = -0x7f

	Smaller Scale = -2
	Small   Scale = -1
	Default Scale = 0
	Big     Scale = 1
	Bigger  Scale = 1
)

func (s Scale) Value(base unit.Value, ratio ScaleRatio) unit.Value {
	if s == Default {
		return base
	}
	if s == None {
		return unit.Value{}
	}
	return unit.Value{
		V: base.V * float32(math.Pow(float64(ratio), float64(s))),
		U: base.U,
	}
}
