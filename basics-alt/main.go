package main

import (
	"image"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {

	ops := new(op.Ops)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()

			doButton(ops, e.Queue)

			// render the frame
			e.Frame(ops)
		}
	}
}

var tag = new(bool) // We could use &pressed for this instead.
var pressed = false

func doButton(ops *op.Ops, q event.Queue) {
	// Make sure we don’t pollute the graphics context.
	var stack op.StackOp
	stack.Push(ops)
	defer stack.Pop()

	for _, ev := range q.Events(tag) {
		if x, ok := ev.(pointer.Event); ok {
			switch x.Type {
			case pointer.Press:
				pressed = true
			case pointer.Release:
				pressed = false
			}
		}
	}

	pointer.Rect(image.Rect(0, 0, 100, 100)).Add(ops)
	pointer.InputOp{Tag: tag}.Add(ops)

	var c color.RGBA
	if pressed {
		c = color.RGBA{R: 0xFF, A: 0xFF}
	} else {
		c = color.RGBA{G: 0xFF, A: 0xFF}
	}
	paint.ColorOp{Color: c}.Add(ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: 100, Y: 100}}}.Add(ops)
}
