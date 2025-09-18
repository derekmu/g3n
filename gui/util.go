// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// RectBounds specifies the size of the boundaries of a rectangle.
// It can represent the thickness of the borders, the margins, or the padding of a rectangle.
type RectBounds struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// Rect represents a rectangle.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Rect) Clip(clip Rect) Rect {
	clipped := Rect{
		X: max(r.X, clip.X),
		Y: max(r.Y, clip.Y),
	}
	clipped.Width = max(0, min(r.X+r.Width, clip.X+clip.Width)-clipped.X)
	clipped.Height = max(0, min(r.Y+r.Height, clip.Y+clip.Height)-clipped.Y)
	return clipped
}

// Contains returns whether this rect contains a point.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && y >= r.Y && x < (r.X+r.Width) && y < (r.Y+r.Height)
}

// Intersects returns whether this Rect intersects with another Rect.
func (r Rect) Intersects(r2 Rect) bool {
	return r.X+r.Width > r2.X && r2.X+r2.Width > r.X &&
		r.Y+r.Height > r2.Y && r2.Y+r2.Height > r.Y
}

// RectF represents a rectangle.
type RectF struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}
