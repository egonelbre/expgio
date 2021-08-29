package main

import (
	"image/color"

	"gioui.org/widget/material"
)

type Theme struct {
	Node Style
	Port Style
	Conn Style

	Selected Style
	Active   Style

	Grid color.NRGBA

	*material.Theme
}

type Style struct {
	Fill   color.NRGBA
	Mid    color.NRGBA
	Border color.NRGBA
}

func NewTheme(th *material.Theme) *Theme {
	return &Theme{
		Node: Tango[0],
		Port: Tango[1],
		Conn: Tango[3],

		Selected: Tango[4],
		Active:   Tango[5],

		Grid: color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},

		Theme: th,
	}
}

var Tango = []Style{
	0: { // Butter
		Fill:   hexRGB(0xFCE94F),
		Mid:    hexRGB(0xEDD400),
		Border: hexRGB(0xC4A000),
	},
	1: { // Orange
		Fill:   hexRGB(0xFCAF3E),
		Mid:    hexRGB(0xF57900),
		Border: hexRGB(0xCE5C00),
	},
	2: { // Chocolate
		Fill:   hexRGB(0xE9B96E),
		Mid:    hexRGB(0xC17D11),
		Border: hexRGB(0x8F5902),
	},
	3: { // Chameleon
		Fill:   hexRGB(0x8AE234),
		Mid:    hexRGB(0x73D216),
		Border: hexRGB(0x4E9A06),
	},
	4: { // Sky Blue
		Fill:   hexRGB(0x729FCF),
		Mid:    hexRGB(0x3465A4),
		Border: hexRGB(0x204A87),
	},
	5: { // Plum
		Fill:   hexRGB(0xAD7FA8),
		Mid:    hexRGB(0x75507B),
		Border: hexRGB(0x5C3566),
	},
	6: { // Scarlet Red
		Fill:   hexRGB(0xEF2929),
		Mid:    hexRGB(0xCC0000),
		Border: hexRGB(0xA40000),
	},
	7: { // Aluminium
		Fill:   hexRGB(0xEEEEEC),
		Mid:    hexRGB(0xD3D7CF),
		Border: hexRGB(0xBABDB6),
	},
	8: { // Dark Aluminium
		Fill:   hexRGB(0x888A85),
		Mid:    hexRGB(0x555753),
		Border: hexRGB(0x2E3436),
	},
}

func hexRGB(v uint32) color.NRGBA {
	return color.NRGBA{
		R: byte(v >> 16),
		G: byte(v >> 8),
		B: byte(v >> 0),
		A: 0xFF,
	}
}
