// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// RectBounds specifies the size of the boundaries of a rectangle.
// It can represent the thickness of the borders, the margins, or the padding of a rectangle.
type RectBounds struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// Rect represents a rectangle.
type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}
