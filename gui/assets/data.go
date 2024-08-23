package assets

import (
	"bytes"
	_ "embed"
	"github.com/derekmu/g3n/text"
	"image"
)

//go:embed fonts/FreeMono.ttf
var freeMonoTTF []byte

func NewFreeMonoFont() *text.Font {
	font, err := text.NewFontFromData(freeMonoTTF)
	if err != nil {
		panic(err)
	}
	return font
}

//go:embed fonts/FreeSans.ttf
var freeSansTTF []byte

func NewFreeSansFont() *text.Font {
	font, err := text.NewFontFromData(freeSansTTF)
	if err != nil {
		panic(err)
	}
	return font
}

//go:embed fonts/FreeSansBold.ttf
var freeSansBoldTTF []byte

func NewFreeSansBoldFont() *text.Font {
	font, err := text.NewFontFromData(freeSansBoldTTF)
	if err != nil {
		panic(err)
	}
	return font
}

//go:embed fonts/MaterialIcons-Regular.ttf
var materialIconsRegularTTF []byte

func NewMaterialIconsRegularFont() *text.Font {
	font, err := text.NewFontFromData(materialIconsRegularTTF)
	if err != nil {
		panic(err)
	}
	return font
}

//go:embed cursors/trbl.png
var cursorsTrblPng []byte

func NewCursorTrblImage() (image.Image, string, error) {
	return image.Decode(bytes.NewReader(cursorsTrblPng))
}

//go:embed cursors/tlbr.png
var cursorsTlbrPng []byte

func NewCursorTlbrImage() (image.Image, string, error) {
	return image.Decode(bytes.NewReader(cursorsTlbrPng))
}
