// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package texture contains several types of textures which can be added to materials.
package texture

import (
	"fmt"
	"github.com/derekmu/g3n/gls"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// Texture2D represents a texture
type Texture2D struct {
	gs           *gls.GLS
	refcount     int         // Current number of references
	texname      uint32      // Texture handle
	magFilter    uint32      // magnification filter
	minFilter    uint32      // minification filter
	wrapS        uint32      // wrap mode for s coordinate
	wrapT        uint32      // wrap mode for t coordinate
	iformat      int32       // internal format
	width        int32       // texture width in pixels
	height       int32       // texture height in pixels
	format       uint32      // format of the pixel data
	formatType   uint32      // type of the pixel data
	updateData   bool        // texture data needs to be sent
	updateParams bool        // texture parameters needs to be sent
	compressed   bool        // whether the texture is compressed
	size         int32       // the size of the texture data in bytes
	data         any         // array with texture data
	uniUnit      gls.Uniform // Texture unit uniform location cache
	uniInfo      gls.Uniform // Texture info uniform location cache
	udata        struct {    // Combined uniform data in 3 vec2:
		offsetX float32
		offsetY float32
		repeatX float32
		repeatY float32
		flipY   float32
		visible float32
	}
}

func newTexture2D() *Texture2D {
	t := new(Texture2D)
	t.gs = nil
	t.refcount = 1
	t.texname = 0
	t.magFilter = gls.LINEAR
	t.minFilter = gls.LINEAR_MIPMAP_LINEAR
	t.wrapS = gls.CLAMP_TO_EDGE
	t.wrapT = gls.CLAMP_TO_EDGE
	t.updateData = false
	t.updateParams = true

	// Initialize Uniform elements
	t.uniUnit.Init("MatTexture")
	t.uniInfo.Init("MatTexinfo")
	t.SetOffset(0, 0)
	t.SetRepeat(1, 1)
	t.SetFlipY(true)
	t.SetVisible(true)
	return t
}

// NewTexture2DFromImage creates and returns a pointer to a new Texture2D
// using the specified image file as data.
// Supported image formats are: PNG, JPEG and GIF.
func NewTexture2DFromImage(imgfile string) (*Texture2D, error) {
	// Decodes image file into RGBA8
	rgba, err := DecodeImage(imgfile)
	if err != nil {
		return nil, err
	}

	t := newTexture2D()
	t.SetFromRGBA(rgba)
	return t, nil
}

// NewTexture2DFromRGBA creates a new texture from a pointer to an RGBA image object.
func NewTexture2DFromRGBA(rgba *image.RGBA) *Texture2D {
	t := newTexture2D()
	t.SetFromRGBA(rgba)
	return t
}

// NewTexture2DFromData creates a new texture from data
func NewTexture2DFromData(width, height int, format int, formatType, iformat int, data any) *Texture2D {
	t := newTexture2D()
	t.SetData(width, height, format, formatType, iformat, data)
	return t
}

// NewTexture2DFromCompressedData creates a new compressed texture from data
func NewTexture2DFromCompressedData(width, height int, iformat int32, size int32, data any) *Texture2D {
	t := newTexture2D()
	t.SetCompressedData(width, height, iformat, size, data)
	return t
}

// Incref increments the reference count for this texture
// and returns a pointer to the geometry.
// It should be used when this texture is shared by another
// material.
func (t *Texture2D) Incref() *Texture2D {
	t.refcount++
	return t
}

// Dispose decrements this texture reference count and
// if necessary releases OpenGL resources and C memory
// associated with this texture.
func (t *Texture2D) Dispose() {
	if t.refcount > 1 {
		t.refcount--
		return
	}
	if t.gs != nil {
		t.gs.DeleteTextures(t.texname)
		t.gs = nil
	}
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

// SetImage sets a new image for this texture
func (t *Texture2D) SetImage(imgfile string) error {
	// Decodes image file into RGBA8
	rgba, err := DecodeImage(imgfile)
	if err != nil {
		return err
	}
	t.SetFromRGBA(rgba)
	return nil
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
func (t *Texture2D) SetData(width, height int, format int, formatType, iformat int, data any) {
	t.width = int32(width)
	t.height = int32(height)
	t.format = uint32(format)
	t.formatType = uint32(formatType)
	t.iformat = int32(iformat)
	t.compressed = false
	t.data = data
	t.updateData = true
}

// SetCompressedData sets the compressed texture data
func (t *Texture2D) SetCompressedData(width, height int, iformat int32, size int32, data any) {
	t.width = int32(width)
	t.height = int32(height)
	t.iformat = iformat
	t.compressed = true
	t.size = size
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

// Compressed returns whether this texture is compressed
func (t *Texture2D) Compressed() bool {
	return t.compressed
}

// DecodeImage reads and decodes the specified image file into RGBA8.
// The supported image files are PNG, JPEG and GIF.
func DecodeImage(imgfile string) (*image.RGBA, error) {
	// Open image file
	file, err := os.Open(imgfile)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Decodes image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Converts image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba, nil
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
		if t.compressed {
			t.gs.CompressedTexImage2D(
				gls.TEXTURE_2D,
				0,
				uint32(t.iformat),
				t.width,
				t.height,
				t.size,
				t.data,
			)
		} else {
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
		}
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
