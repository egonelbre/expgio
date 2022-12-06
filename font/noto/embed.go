package noto

import (
	_ "embed"
	"fmt"
	"sync"

	"gioui.org/font/opentype"
	"gioui.org/text"
)

//go:embed NotoMusic-Regular.ttf
var NotoMusicRegular []byte

//go:embed NotoSans-Regular.ttf
var NotoSansRegular []byte

//go:embed NotoSans-Bold.ttf
var NotoSansBold []byte

var (
	once       sync.Once
	collection []text.FontFace
)

func Collection() []text.FontFace {
	once.Do(func() {
		register("Noto Sans", text.Font{}, NotoSansRegular)
		register("Noto Sans", text.Font{Weight: text.Bold}, NotoSansBold)
		register("Noto Music", text.Font{}, NotoMusicRegular)
	})
	return collection
}

func register(typeface string, fnt text.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	fnt.Typeface = text.Typeface(typeface)
	collection = append(collection, text.FontFace{Font: fnt, Face: face})
}
