package assets

import (
	"bytes"
	_ "embed"
	"github.com/derekmu/g3n/text"
	"image"
)

//go:embed fonts/FreeMono.ttf
var freeMonoTTF []byte

func NewFreeMonoFont() (*text.Font, error) {
	return text.NewFontFromData(freeMonoTTF)
}

//go:embed fonts/FreeSans.ttf
var freeSansTTF []byte

func NewFreeSansFont() (*text.Font, error) {
	return text.NewFontFromData(freeSansTTF)
}

//go:embed fonts/FreeSansBold.ttf
var freeSansBoldTTF []byte

func NewFreeSansBoldFont() (*text.Font, error) {
	return text.NewFontFromData(freeSansBoldTTF)
}

//go:embed fonts/MaterialIcons-Regular.ttf
var materialIconsRegularTTF []byte

func NewMaterialIconsRegularFont() (*text.Font, error) {
	return text.NewFontFromData(materialIconsRegularTTF)
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
