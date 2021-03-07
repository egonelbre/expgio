package main

import (
	"image/color"
)

var Default = Tango[0]
var DefaultPort = Tango[1]
var DefaultConnection = Tango[3]

var FocusColor = Tango[4]
var ActiveColor = Tango[5]

type Style struct {
	Name string

	Fill   color.NRGBA
	Mid    color.NRGBA
	Border color.NRGBA
}

func WithAlpha(c color.NRGBA, v byte) color.NRGBA {
	c.A = v
	return c
}

type Styles []*Style

var Tango = Styles{
	{
		Name:   "Butter",
		Fill:   hexRGB(0xFCE94F),
		Mid:    hexRGB(0xEDD400),
		Border: hexRGB(0xC4A000),
	},
	{
		Name:   "Orange",
		Fill:   hexRGB(0xFCAF3E),
		Mid:    hexRGB(0xF57900),
		Border: hexRGB(0xCE5C00),
	},
	{
		Name:   "Chocolate",
		Fill:   hexRGB(0xE9B96E),
		Mid:    hexRGB(0xC17D11),
		Border: hexRGB(0x8F5902),
	},
	{
		Name:   "Chameleon",
		Fill:   hexRGB(0x8AE234),
		Mid:    hexRGB(0x73D216),
		Border: hexRGB(0x4E9A06),
	},
	{
		Name:   "Sky Blue",
		Fill:   hexRGB(0x729FCF),
		Mid:    hexRGB(0x3465A4),
		Border: hexRGB(0x204A87),
	},
	{
		Name:   "Plum",
		Fill:   hexRGB(0xAD7FA8),
		Mid:    hexRGB(0x75507B),
		Border: hexRGB(0x5C3566),
	},
	{
		Name:   "Scarlet Red",
		Fill:   hexRGB(0xEF2929),
		Mid:    hexRGB(0xCC0000),
		Border: hexRGB(0xA40000),
	},
	{
		Name:   "Aluminium",
		Fill:   hexRGB(0xEEEEEC),
		Mid:    hexRGB(0xD3D7CF),
		Border: hexRGB(0xBABDB6),
	},
	{
		Name:   "Dark Aluminium",
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
