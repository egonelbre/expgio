// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/widget/material"
)

const cursorCount = pointer.CursorNorthWestSouthEastResize + 1

func main() {
	ui := &UI{Theme: material.NewTheme()}
	go func() {
		w := &app.Window{}
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Theme *material.Theme
}

type ValueEvent struct {
	Int int
}

type EventQueue struct {
	window *app.Window
	mu     sync.Mutex
	items  []any
}

func (q *EventQueue) Enqueue(event any) {
	q.mu.Lock()
	q.items = append(q.items, event)
	q.mu.Unlock()

	q.window.Invalidate()
}

func (q *EventQueue) Events() []any {
	q.mu.Lock()
	defer q.mu.Unlock()
	r := q.items
	q.items = nil
	return r
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	customEvents := EventQueue{
		window: w,
	}

	go func() {
		for range time.Tick(time.Second) {
			customEvents.Enqueue(ValueEvent{Int: rand.Int()})
		}
	}()

	state := 0

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:

			for _, ce := range customEvents.Events() {
				switch e := ce.(type) {
				case ValueEvent:
					state = e.Int
				}
			}

			gtx := app.NewContext(&ops, e)

			material.H2(ui.Theme, strconv.Itoa(state)).Layout(gtx)

			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}
