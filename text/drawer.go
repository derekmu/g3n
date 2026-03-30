package text

import (
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Drawer is an extension of font.Drawer.
type Drawer font.Drawer

// MeasureRunes returns how far dot would advance by drawing s.
// widths will be updated with the width of each rune, if there is space.
func (d *Drawer) MeasureRunes(s []rune, widths []fixed.Int26_6) (advance fixed.Int26_6) {
	return MeasureRunes(d.Face, s, widths)
}

// MeasureRunes returns how far dot would advance by drawing s with f.
func MeasureRunes(f font.Face, s []rune, widths []fixed.Int26_6) (advance fixed.Int26_6) {
	prevC := rune(-1)
	for i, c := range s {
		if prevC >= 0 {
			advance += f.Kern(prevC, c)
		}
		a, _ := f.GlyphAdvance(c)
		if i < len(widths) {
			widths[i] = a
		}
		advance += a
		prevC = c
	}
	return advance
}

// DrawRunes draws s at the dot and advances the dot's location.
func (d *Drawer) DrawRunes(s []rune) {
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, _ := d.Face.Glyph(d.Dot, c)
		if !dr.Empty() {
			draw.DrawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp, draw.Over)
		}
		d.Dot.X += advance
		prevC = c
	}
}
