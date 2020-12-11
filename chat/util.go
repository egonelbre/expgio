package main

import "image/color"

func ifthen(v bool, a, b color.NRGBA) color.NRGBA {
	if v {
		return a
	} else {
		return b
	}
}
