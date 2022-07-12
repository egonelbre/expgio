package rig

import (
	"image"

	"gioui.org/f32"
)

type Rectangle struct{ Min, Max f32.Point }

func inflate(p image.Point, size int) image.Rectangle {
	return image.Rectangle{Min: p, Max: p}.Inset(-size)
}
