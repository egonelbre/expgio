package main

import (
	"image/color"

	"gioui.org/unit"
)

type Style struct {
	Name string

	Fill   color.NRGBA
	Mid    color.NRGBA
	Border color.NRGBA

	BorderWidth unit.Value
}

type Styles []*Style

var Tango = Styles{
	{Name: "Butter",
		Fill:        hexRGB(0xFCE94F),
		Mid:         hexRGB(0xEDD400),
		Border:      hexRGB(0xC4A000),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Orange",
		Fill:        hexRGB(0xFCAF3E),
		Mid:         hexRGB(0xF57900),
		Border:      hexRGB(0xCE5C00),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Chocolate",
		Fill:        hexRGB(0xE9B96E),
		Mid:         hexRGB(0xC17D11),
		Border:      hexRGB(0x8F5902),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Chameleon",
		Fill:        hexRGB(0x8AE234),
		Mid:         hexRGB(0x73D216),
		Border:      hexRGB(0x4E9A06),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Sky Blue",
		Fill:        hexRGB(0x729FCF),
		Mid:         hexRGB(0x3465A4),
		Border:      hexRGB(0x204A87),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Plum",
		Fill:        hexRGB(0xAD7FA8),
		Mid:         hexRGB(0x75507B),
		Border:      hexRGB(0x5C3566),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Scarlet Red",
		Fill:        hexRGB(0xEF2929),
		Mid:         hexRGB(0xCC0000),
		Border:      hexRGB(0xA40000),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Aluminium",
		Fill:        hexRGB(0xEEEEEC),
		Mid:         hexRGB(0xD3D7CF),
		Border:      hexRGB(0xBABDB6),
		BorderWidth: unit.Dp(2),
	},
	{Name: "Dark Aluminium",
		Fill:        hexRGB(0x888A85),
		Mid:         hexRGB(0x555753),
		Border:      hexRGB(0x2E3436),
		BorderWidth: unit.Dp(2),
	},
}

func hexRGB(v uint32) color.NRGBA {
	return color.NRGBA{
		R: byte(v >> 24),
		G: byte(v >> 16),
		B: byte(v >> 8),
		A: 0xFF,
	}
}
