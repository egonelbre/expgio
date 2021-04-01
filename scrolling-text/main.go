package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/image/math/fixed"

	"gioui.org/x/scroll"
)

func main() {
	ui := NewUI()

	go func() {
		w := app.NewWindow(
			app.Title("Scrolling Text"),
			app.Size(unit.Dp(600), unit.Dp(600)),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

type UI struct {
	Theme    *material.Theme
	Terminal Terminal
}

type Terminal struct {
	List   layout.List
	Scroll scroll.Scrollable
	Lines  []string
}

func NewUI() *UI {
	ui := &UI{}

	fonts := gofont.Collection()
	ui.Theme = material.NewTheme(fonts)
	ui.Theme.Shaper = NewCache(fonts)

	ui.Terminal.List.Axis = layout.Vertical
	totalStringSize := 0
	for i := 0; i < 10000; i++ {
		line := strings.Repeat(strconv.Itoa(i)+" | ", 10)
		ui.Terminal.Lines = append(ui.Terminal.Lines, line)
		totalStringSize += len(line)
	}
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)

		case key.Event:
			switch e.Name {
			case key.NameEscape:
				return nil
			}

		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

var defaultInset = layout.UniformInset(unit.Dp(8))

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return defaultInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return ui.Terminal.Layout(ui.Theme, gtx)
	})
}

func (term *Terminal) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if didScroll, progress := term.Scroll.Scrolled(); didScroll {
		term.List.Position.First = int(progress * float32(len(term.Lines)))
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return term.List.Layout(gtx, len(term.Lines), func(gtx layout.Context, index int) layout.Dimensions {
				return material.Body1(th, term.Lines[index]).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			progress := float32(term.List.Position.First) / float32(len(term.Lines))
			size := float32(term.List.Position.Count) / float32(len(term.Lines))
			return scroll.DefaultBar(&term.Scroll, progress, size).Layout(gtx)
		}),
	)
}

// Cache implements font rendering without caching.
type Cache struct {
	def   text.Typeface
	faces map[text.Font]*faceCache
}

type faceCache struct {
	face text.Face
}

func (c *Cache) lookup(font text.Font) *faceCache {
	f := c.faceForStyle(font)
	if f == nil {
		font.Typeface = c.def
		f = c.faceForStyle(font)
	}
	return f
}

func (c *Cache) faceForStyle(font text.Font) *faceCache {
	tf := c.faces[font]
	if tf == nil {
		font := font
		font.Weight = text.Normal
		tf = c.faces[font]
	}
	if tf == nil {
		font := font
		font.Style = text.Regular
		tf = c.faces[font]
	}
	if tf == nil {
		font := font
		font.Style = text.Regular
		font.Weight = text.Normal
		tf = c.faces[font]
	}
	return tf
}

func NewCache(collection []text.FontFace) *Cache {
	c := &Cache{
		faces: make(map[text.Font]*faceCache),
	}
	for i, ff := range collection {
		if i == 0 {
			c.def = ff.Font.Typeface
		}
		c.faces[ff.Font] = &faceCache{face: ff.Face}
	}
	return c
}

// Layout implements the Shaper interface.
func (s *Cache) Layout(font text.Font, size fixed.Int26_6, maxWidth int, txt io.Reader) ([]text.Line, error) {
	cache := s.lookup(font)
	return cache.face.Layout(size, maxWidth, txt)
}

// LayoutString is a caching implementation of the Shaper interface.
func (s *Cache) LayoutString(font text.Font, size fixed.Int26_6, maxWidth int, str string) []text.Line {
	cache := s.lookup(font)
	return cache.layout(size, maxWidth, str)
}

// Shape is a caching implementation of the Shaper interface. Shape assumes that the layout
// argument is unchanged from a call to Layout or LayoutString.
func (s *Cache) Shape(font text.Font, size fixed.Int26_6, layout text.Layout) op.CallOp {
	cache := s.lookup(font)
	return cache.shape(size, layout)
}

func (f *faceCache) layout(ppem fixed.Int26_6, maxWidth int, str string) []text.Line {
	if f == nil {
		return nil
	}
	l, _ := f.face.Layout(ppem, maxWidth, strings.NewReader(str))
	return l
}

func (f *faceCache) shape(ppem fixed.Int26_6, layout text.Layout) op.CallOp {
	if f == nil {
		return op.CallOp{}
	}
	return f.face.Shape(ppem, layout)
}
