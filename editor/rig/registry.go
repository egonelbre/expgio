package rig

import "image/color"

type Registry struct {
	Editors []EditorDef
}

func NewRegistry() *Registry {
	return &Registry{
		Editors: []EditorDef{
			ColorViewer("Red", color.NRGBA{R: 0xFF, A: 0xFF}),
			ColorViewer("RedGreen", color.NRGBA{R: 0xFF, G: 0xFF, A: 0xFF}),
			ColorViewer("Green", color.NRGBA{G: 0xFF, A: 0xFF}),
			ColorViewer("GreenBlue", color.NRGBA{G: 0xFF, B: 0xFF, A: 0xFF}),
			ColorViewer("Blue", color.NRGBA{B: 0xFF, A: 0xFF}),
			ColorViewer("BlueRed", color.NRGBA{R: 0xFF, B: 0xFF, A: 0xFF}),
		},
	}
}
