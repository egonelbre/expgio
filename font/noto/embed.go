package noto

import (
	_ "embed"
	"fmt"
	"sync"

	"gioui.org/font"
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
		register("Noto Sans", font.Font{}, NotoSansRegular)
		register("Noto Sans", font.Font{Weight: font.Bold}, NotoSansBold)
		register("Noto Music", font.Font{}, NotoMusicRegular)
	})
	return collection
}

func register(typeface string, fnt font.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	fnt.Typeface = font.Typeface(typeface)
	collection = append(collection, font.FontFace{Font: fnt, Face: face})
}
