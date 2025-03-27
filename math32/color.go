// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "image/color"

type Color interface {
	Color3() Color3
	Color4() Color4
	NRGBA() color.NRGBA
	RGBA() color.RGBA
}

var _ Color = Color3{}

// Color3 describes an RGB color.
type Color3 struct {
	R float32
	G float32
	B float32
}

func (c Color3) Color3() Color3 {
	return c
}

func (c Color3) Color4() Color4 {
	return Color4{
		R: c.R,
		G: c.G,
		B: c.B,
		A: 1.0,
	}
}

func (c Color3) NRGBA() color.NRGBA {
	return color.NRGBA{
		R: uint8(c.R * 0xff),
		G: uint8(c.G * 0xff),
		B: uint8(c.B * 0xff),
		A: 0xff,
	}
}

func (c Color3) RGBA() color.RGBA {
	return color.RGBA{
		R: uint8(c.R * 0xff),
		G: uint8(c.G * 0xff),
		B: uint8(c.B * 0xff),
		A: 0xff,
	}
}

// MultiplyScalar returns a Color3 with the RGB components multiplied by a value.
func (c Color3) MultiplyScalar(v float32) Color3 {
	return Color3{
		R: c.R * v,
		G: c.G * v,
		B: c.B * v,
	}
}

var _ Color = Color4{}

// Color4 describes an RGBA color.
type Color4 struct {
	R float32
	G float32
	B float32
	A float32
}

func (c Color4) Color3() Color3 {
	return Color3{
		R: c.R,
		G: c.G,
		B: c.B,
	}
}

func (c Color4) Color4() Color4 {
	return c
}

func (c Color4) NRGBA() color.NRGBA {
	return color.NRGBA{
		R: uint8(c.R * 0xff),
		G: uint8(c.G * 0xff),
		B: uint8(c.B * 0xff),
		A: uint8(c.A * 0xff),
	}
}

func (c Color4) RGBA() color.RGBA {
	return color.RGBA{
		R: uint8(c.R * 0xff),
		G: uint8(c.G * 0xff),
		B: uint8(c.B * 0xff),
		A: uint8(c.A * 0xff),
	}
}
