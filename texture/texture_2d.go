// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package texture contains several types of textures which can be added to materials.
package texture

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gls"
	"image"
)

// Texture2D represents a texture.
type Texture2D struct {
	core.RefCount
	gs           *gls.GLS
	texname      uint32  // Texture handle
	magFilter    uint32  // magnification filter
	minFilter    uint32  // minification filter
	wrapS        uint32  // wrap mode for s coordinate
	wrapT        uint32  // wrap mode for t coordinate
	iformat      int32   // internal format
	width        int32   // texture width in pixels
	height       int32   // texture height in pixels
	format       uint32  // format of the pixel data
	formatType   uint32  // type of the pixel data
	updateData   bool    // texture data needs to be sent
	updateParams bool    // texture parameters needs to be sent
	size         int32   // the size of the texture data in bytes
	data         []uint8 // array with texture data
	uniUnit      gls.Uniform
	uniInfo      gls.Uniform
	udata        struct {
		offsetX float32
		offsetY float32
		repeatX float32
		repeatY float32
		flipY   float32
		visible float32
	}
}

// NewTexture2DFromRGBA creates a new texture from a pointer to an RGBA image object.
func NewTexture2DFromRGBA(rgba *image.RGBA) *Texture2D {
	t := new(Texture2D)
	t.InitTexture2D()
	t.SetFromRGBA(rgba)
	return t
}

func (t *Texture2D) InitTexture2D() {
	t.gs = nil
	t.texname = 0
	t.magFilter = gls.LINEAR
	t.minFilter = gls.LINEAR_MIPMAP_LINEAR
	t.wrapS = gls.CLAMP_TO_EDGE
	t.wrapT = gls.CLAMP_TO_EDGE
	t.updateData = false
	t.updateParams = true
	t.uniUnit.Init("uMatTexture")
	t.uniInfo.Init("uMatTexInfo")
	t.SetOffset(0, 0)
	t.SetRepeat(1, 1)
	t.SetFlipY(true)
	t.SetVisible(true)
}

// Dispose releases OpenGL resources associated with this texture.
func (t *Texture2D) Dispose() {
	if t.gs != nil {
		t.gs.DeleteTextures(t.texname)
	}
	t.InitTexture2D()
}

// TexName returns the texture handle for the texture
func (t *Texture2D) TexName() uint32 {
	return t.texname
}

// SetUniformNames sets the names of the uniforms in the shader for sampler and texture info.
func (t *Texture2D) SetUniformNames(sampler, info string) {
	t.uniUnit.Init(sampler)
	t.uniInfo.Init(info)
}

// GetUniformNames returns the names of the uniforms in the shader for sampler and texture info.
func (t *Texture2D) GetUniformNames() (sampler, info string) {
	return t.uniUnit.Name(), t.uniInfo.Name()
}

// SetFromRGBA sets the texture data from the specified image.RGBA object
func (t *Texture2D) SetFromRGBA(rgba *image.RGBA) {
	t.SetData(
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gls.RGBA,
		gls.UNSIGNED_BYTE,
		gls.RGBA8,
		rgba.Pix,
	)
}

// SetData sets the texture data
func (t *Texture2D) SetData(width, height int, format int, formatType, iformat int, data []uint8) {
	t.width = int32(width)
	t.height = int32(height)
	t.format = uint32(format)
	t.formatType = uint32(formatType)
	t.iformat = int32(iformat)
	t.data = data
	t.updateData = true
}

// SetVisible sets the visibility state of the texture
func (t *Texture2D) SetVisible(state bool) {
	if state {
		t.udata.visible = 1
	} else {
		t.udata.visible = 0
	}
}

// Visible returns the current visibility state of the texture
func (t *Texture2D) Visible() bool {
	return t.udata.visible != 0
}

// SetMagFilter sets the filter to be applied when the texture element
// covers more than on pixel. The default value is gls.Linear.
func (t *Texture2D) SetMagFilter(magFilter uint32) {
	t.magFilter = magFilter
	t.updateParams = true
}

// SetMinFilter sets the filter to be applied when the texture element
// covers less than on pixel. The default value is gls.Linear.
func (t *Texture2D) SetMinFilter(minFilter uint32) {
	t.minFilter = minFilter
	t.updateParams = true
}

// SetWrapS set the wrapping mode for texture S coordinate
// The default value is GL_CLAMP_TO_EDGE;
func (t *Texture2D) SetWrapS(wrapS uint32) {
	t.wrapS = wrapS
	t.updateParams = true
}

// SetWrapT set the wrapping mode for texture T coordinate
// The default value is GL_CLAMP_TO_EDGE;
func (t *Texture2D) SetWrapT(wrapT uint32) {
	t.wrapT = wrapT
	t.updateParams = true
}

// SetRepeat set the repeat factor
func (t *Texture2D) SetRepeat(x, y float32) {
	t.udata.repeatX = x
	t.udata.repeatY = y
}

// Repeat returns the current X and Y repeat factors
func (t *Texture2D) Repeat() (float32, float32) {
	return t.udata.repeatX, t.udata.repeatY
}

// SetOffset sets the offset factor
func (t *Texture2D) SetOffset(x, y float32) {
	t.udata.offsetX = x
	t.udata.offsetY = y
}

// Offset returns the current X and Y offset factors
func (t *Texture2D) Offset() (float32, float32) {
	return t.udata.offsetX, t.udata.offsetY
}

// SetFlipY set the state for flipping the Y coordinate
func (t *Texture2D) SetFlipY(state bool) {
	if state {
		t.udata.flipY = 1
	} else {
		t.udata.flipY = 0
	}
}

// Width returns the texture width in pixels
func (t *Texture2D) Width() int {
	return int(t.width)
}

// Height returns the texture height in pixels
func (t *Texture2D) Height() int {
	return int(t.height)
}

func (t *Texture2D) genMipmap() bool {
	return t.minFilter >= gls.NEAREST_MIPMAP_NEAREST && t.minFilter <= gls.LINEAR_MIPMAP_LINEAR
}

// RenderSetup binds the texture to the active texture unit and assigns it to a sampler slot.
func (t *Texture2D) RenderSetup(gs *gls.GLS, slotIdx, uniIdx int) {
	gs.ActiveTexture(uint32(gls.TEXTURE0 + slotIdx))
	t.BindAndTransfer(gs)

	// Transfer texture unit uniform
	var location int32
	if uniIdx == 0 {
		location = t.uniUnit.Location(gs)
	} else {
		location = t.uniUnit.LocationIdx(gs, int32(uniIdx))
	}
	gs.Uniform1i(location, int32(slotIdx))

	// Transfer texture info combined uniform
	if t.uniInfo.Name() != "" {
		const vec2count = 3
		location = t.uniInfo.LocationIdx(gs, vec2count*int32(uniIdx))
		gs.Uniform2fv(location, vec2count, &t.udata.offsetX)
	}
}

// BindAndTransfer binds the texture to the active texture unit and transfers data and parameters.
func (t *Texture2D) BindAndTransfer(gs *gls.GLS) {
	t.bindTexture(gs)
	t.transferData()
	t.transferParameters()
}

func (t *Texture2D) bindTexture(gs *gls.GLS) {
	if t.texname == 0 {
		t.texname = gs.GenTexture()
		t.gs = gs
	}
	t.gs.BindTexture(gls.TEXTURE_2D, t.texname)
}

func (t *Texture2D) transferData() {
	if t.updateData {
		t.gs.TexImage2D(
			gls.TEXTURE_2D, // texture type
			0,              // level of detail
			t.iformat,      // internal format
			t.width,        // width in texels
			t.height,       // height in texels
			t.format,       // format of supplied texture data
			t.formatType,   // type of external format color component
			t.data,         // image data
		)
		if t.genMipmap() {
			t.gs.GenerateMipmap(gls.TEXTURE_2D)
		}
		t.updateData = false
	}
}

func (t *Texture2D) transferParameters() {
	if t.updateParams {
		t.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MAG_FILTER, int32(t.magFilter))
		t.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MIN_FILTER, int32(t.minFilter))
		t.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_S, int32(t.wrapS))
		t.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_T, int32(t.wrapT))
		t.updateParams = false
	}
}
