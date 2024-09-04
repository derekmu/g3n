// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// Align specifies the alignment of an object inside another.
type Align int

const (
	AlignNone = Align(iota)
	AlignLeftTop
	AlignLeftBottom
	AlignLeftCenter
	AlignRightTop
	AlignRightBottom
	AlignRightCenter
	AlignCenterTop
	AlignCenterBottom
	AlignCenterCenter
)

// CalculatePosition calculates where to place an element within another element.
//
// outsideWidth and outsideHeight usually consist of a Panel's ContentSize.
// insideWidth and insideHeight usually consist of a Panel's Size.
//
// AlignNone returns (0, 0).
func (a Align) CalculatePosition(outsideWidth, outsideHeight, insideWidth, insideHeight float32) (x float32, y float32) {
	switch a {
	case AlignLeftTop:
		x, y = 0, 0
	case AlignLeftBottom:
		x, y = 0, outsideHeight-insideHeight
	case AlignLeftCenter:
		x, y = 0, (outsideHeight-insideHeight)/2
	case AlignRightTop:
		x, y = outsideWidth-insideWidth, 0
	case AlignRightBottom:
		x, y = outsideWidth-insideWidth, outsideHeight-insideHeight
	case AlignRightCenter:
		x, y = outsideWidth-insideWidth, (outsideHeight-insideHeight)/2
	case AlignCenterTop:
		x, y = (outsideWidth-insideWidth)/2, 0
	case AlignCenterBottom:
		x, y = (outsideWidth-insideWidth)/2, outsideHeight-insideHeight
	case AlignCenterCenter:
		x, y = (outsideWidth-insideWidth)/2, (outsideHeight-insideHeight)/2
	default:
		x, y = 0, 0
	}
	return x, y
}
