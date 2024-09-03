// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "image/color"

var _ color.Color = Color{}

// Color describes an RGB color.
type Color struct {
	R float32
	G float32
	B float32
}

// ToColor4 returns a Color4 with this Color's RGB components and a specified alpha.
func (c Color) ToColor4(a float32) Color4 {
	return Color4{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

// RGBA implements color.Color.
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R * 0xFF)
	r |= r << 8
	g = uint32(c.G * 0xFF)
	g |= g << 8
	b = uint32(c.B * 0xFF)
	b |= b << 8
	a = uint32(0xFF)
	a |= a << 8
	return
}

// MultiplyScalar returns a Color with the RGB components multiplied by a value.
func (c Color) MultiplyScalar(v float32) Color {
	return Color{
		R: c.R * v,
		G: c.G * v,
		B: c.B * v,
	}
}

var _ color.Color = Color4{}

// Color4 describes an RGBA color
type Color4 struct {
	R float32
	G float32
	B float32
	A float32
}

// ToColor returns a Color with this Color4's RGB components.
func (c Color4) ToColor() Color {
	return Color{c.R, c.G, c.B}
}

// RGBA implements color.Color.
func (c Color4) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R * 0xFF)
	r |= r << 8
	g = uint32(c.G * 0xFF)
	g |= g << 8
	b = uint32(c.B * 0xFF)
	b |= b << 8
	a = uint32(c.A * 0xFF)
	a |= a << 8
	return
}
