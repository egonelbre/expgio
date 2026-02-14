package main

import (
	"testing"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func BenchmarkLayout(b *testing.B) {
	// uncomment to use with commit 43c47f0
	t := material.NewTheme()
	t.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	// uncomment to use with commit babe7a2
	// t := material.NewTheme(gofont.Collection())
	ops := &op.Ops{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		gtx := app.NewContext(ops, app.FrameEvent{})
		for range 10 {
			material.Label(t, 10, "abcdefghijklmnopqrstuvwxyz").Layout(gtx)
			material.Label(t, 10, "oifajmorfj983 4mroaermfnkli").Layout(gtx)
			material.Label(t, 10, "1234 1234 5434 1234 41234").Layout(gtx)
		}
	}
}
