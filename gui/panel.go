// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/graphic"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/texture"
	"log"
)

// Quad geometry shared by all Panels.
var panelQuadGeometry *geometry.Geometry

func init() {
	// Builds array with vertex positions and texture coordinates
	positions := math32.NewArrayF32(0, 20)
	positions.Append(
		0, 0, 0, 0, 1,
		0, -1, 0, 0, 0,
		1, -1, 0, 1, 0,
		1, 0, 0, 1, 1,
	)
	// Builds array of indices
	indices := math32.NewArrayU32(0, 6)
	indices.Append(0, 1, 2, 0, 2, 3)

	// Creates geometry
	geom := geometry.NewGeometry()
	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(positions).
		AddAttrib(gls.VertexPosition).
		AddAttrib(gls.VertexTexcoord),
	)
	panelQuadGeometry = geom
}

// IPanel is the interface for all panel types.
type IPanel interface {
	graphic.IGraphic
	core.IDispatcher[core.GuiEvent]
	InsideBorders(x float32, y float32) bool
	SetPositionZ(z float32)
	ZLayerDelta() int
	Enabled() bool
	ContentArea() Rect
	ClipArea() Rect
}

var _ IPanel = &Panel{}

// Panel is 2D rectangular graphic which by default has a quad geometry.
// When using the default geometry, a panel has margins, borders, paddings, and a content area.
// The content area can be associated with a texture.
// It is the building block of most GUI elements.
type Panel struct {
	*graphic.Graphic
	core.Dispatcher[core.GuiEvent]
	material    *material.Material
	texture     *texture.Texture2D
	zLayerDelta int  // Z-layer relative to parent
	enabled     bool // Whether events should be processed for this panel

	paddings    RectBounds // pixels around the content
	panelArea   Rect       // including paddings, un-clipped
	contentArea Rect       // excluding paddings, un-clipped
	clipArea    Rect       // including paddings, clipped

	uniMatrix gls.Uniform // model matrix uniform location cache
	uniPanel  gls.Uniform // panel parameters uniform location cache
	udata     struct {
		bounds       Rect          // bounds in texture coordinates
		color        math32.Color4 // panel color
		textureValid [4]float32    // texture valid flag (bool, only first float32 is used, align to vec4)
	}
}

// NewPanel a new panel with the specified dimensions.
func NewPanel(width, height float32) *Panel {
	p := new(Panel)
	p.InitPanel(p, width, height)
	return p
}

// InitPanel initializes this panel and is normally used by other types which embed a panel.
func (p *Panel) InitPanel(ipan IPanel, width, height float32) {
	// Initialize material
	p.material = material.NewMaterial()
	p.material.SetUseLights(material.UseLightNone)
	p.material.SetShader("panel")
	p.material.SetTransparent(true)

	// Initialize graphic
	p.Graphic = graphic.NewGraphic(ipan, panelQuadGeometry.Incref(), gls.TRIANGLES)
	p.AddMaterial(p, p.material, 0, 0)

	// Initialize uniforms location caches
	p.uniMatrix.Init("uModelMatrix")
	p.uniPanel.Init("uPanel")

	// Set defaults
	p.enabled = true
	p.resize(width, height, true)
}

// Material returns a pointer for the panel's material.
func (p *Panel) Material() *material.Material {
	return p.material
}

// SetTexture changes the panel's texture.
// It returns a pointer to the previous texture.
func (p *Panel) SetTexture(tex *texture.Texture2D) *texture.Texture2D {
	prevtex := p.texture
	p.Material().RemoveTexture(prevtex)
	p.texture = tex
	if tex != nil {
		p.Material().AddTexture(p.texture)
	}
	return prevtex
}

// SetTopChild moves the specified panel to be the last child of this panel.
func (p *Panel) SetTopChild(ipan IPanel) {
	// Remove panel and if found appends to the end
	if p.Remove(ipan) {
		p.Add(ipan)
		p.SetChanged(true)
	}
}

// SetZLayerDelta sets the Z-layer of this panel relative to its parent.
func (p *Panel) SetZLayerDelta(zLayerDelta int) {
	p.zLayerDelta = zLayerDelta
}

// ZLayerDelta returns the Z-layer of this panel relative to its parent.
func (p *Panel) ZLayerDelta() int {
	return p.zLayerDelta
}

// SetPosition sets the panel's position in pixel coordinates from left to right and from top to bottom of the screen.
func (p *Panel) SetPosition(x, y float32) {
	p.Node.SetPositionX(math32.Round(x))
	p.Node.SetPositionY(math32.Round(y))
}

// SetSize sets this panel external width and height.
func (p *Panel) SetSize(width, height float32) {
	if width < 0 {
		log.Printf("Invalid panel width %v", width)
		width = 0
	}
	if height < 0 {
		log.Printf("Invalid panel height %v", height)
		height = 0
	}
	p.resize(width, height, true)
}

// SetWidth sets this panel external width.
func (p *Panel) SetWidth(width float32) {
	p.SetSize(width, p.panelArea.Height)
}

// SetHeight sets this panel external height.
func (p *Panel) SetHeight(height float32) {
	p.SetSize(p.panelArea.Width, height)
}

// Size returns this panel external width and height.
func (p *Panel) Size() (float32, float32) {
	return p.panelArea.Width, p.panelArea.Height
}

// Width returns the panel external width.
func (p *Panel) Width() float32 {
	return p.panelArea.Width
}

// Height returns the panel external height.
func (p *Panel) Height() float32 {
	return p.panelArea.Height
}

// ContentWidth returns the width of the content area.
func (p *Panel) ContentWidth() float32 {
	return p.contentArea.Width
}

// ContentHeight returns the height of the content area.
func (p *Panel) ContentHeight() float32 {
	return p.contentArea.Height
}

// ContentArea returns the whole content area.
func (p *Panel) ContentArea() Rect {
	return p.contentArea
}

// ClipArea returns the whole clip area.
func (p *Panel) ClipArea() Rect {
	return p.clipArea
}

// SetPaddings sets the panel's padding sizes.
func (p *Panel) SetPaddings(src RectBounds) {
	p.paddings = src
	p.resize(p.contentArea.Width+p.paddings.Left+p.paddings.Right, p.contentArea.Height+p.paddings.Top+p.paddings.Bottom, true)
}

// Paddings returns the panel's padding sizes.
func (p *Panel) Paddings() RectBounds {
	return p.paddings
}

// SetColor sets the panel's color.
func (p *Panel) SetColor(color math32.Color4) *Panel {
	p.udata.color = color
	p.SetChanged(true)
	return p
}

// Color returns the panel's color.
func (p *Panel) Color() math32.Color4 {
	return p.udata.color
}

// SetContentSize sets the panel's content size.
func (p *Panel) SetContentSize(width, height float32) {
	p.setContentSize(width, height, true)
}

// SetContentWidth sets the panel's content width.
func (p *Panel) SetContentWidth(width float32) {
	p.SetContentSize(width, p.contentArea.Height)
}

// SetContentHeight sets the panel's content height.
func (p *Panel) SetContentHeight(height float32) {
	p.SetContentSize(p.contentArea.Width, height)
}

// MinWidth returns the minimum width of this panel (assuming content width was 0).
func (p *Panel) MinWidth() float32 {
	return p.paddings.Left + p.paddings.Right
}

// MinHeight returns the minimum height of this panel (assuming content height was 0).
func (p *Panel) MinHeight() float32 {
	return p.paddings.Top + p.paddings.Bottom
}

// PanelPosition returns the panel's absolute position.
func (p *Panel) PanelPosition() math32.Vector2 {
	return math32.Vector2{X: p.panelArea.X, Y: p.panelArea.Y}
}

// Add adds a child panel to this one.
// This overrides the Node method to enforce that IPanels can only have IPanels as children.
func (p *Panel) Add(ichild IPanel) *Panel {
	p.Node.Add(ichild)
	return p
}

// Remove removes the specified child from this panel.
func (p *Panel) Remove(ichild IPanel) bool {
	return p.Node.Remove(ichild)
}

// UpdateMatrixWorld overrides the core.Node function to update panel bounds instead.
func (p *Panel) UpdateMatrixWorld() {
	par := p.Parent()
	if par == nil {
		p.updateBounds(nil)
	} else {
		if par, ok := par.(IPanel); ok {
			p.updateBounds(par)
		} else {
			p.updateBounds(nil)
		}
	}
	// Update this panel children
	for _, ichild := range p.Children() {
		ichild.UpdateMatrixWorld()
	}
}

// ContainsPosition returns whether this panel contains the specified screen position.
func (p *Panel) ContainsPosition(x, y float32) bool {
	return x >= p.panelArea.X && y >= p.panelArea.Y && x < (p.panelArea.X+p.panelArea.Width) && y < (p.panelArea.Y+p.panelArea.Height)
}

// InsideBorders returns whether a screen position is inside the panel borders, including the border width.
// Unlike ContainsPosition, it does not consider the panel margins.
func (p *Panel) InsideBorders(x, y float32) bool {
	return x >= p.panelArea.X && x < p.panelArea.X+p.panelArea.Width &&
		y >= p.panelArea.Y && y < p.panelArea.Y+p.panelArea.Height
}

// Intersects returns whether this panel intersects with another panel.
func (p *Panel) Intersects(p2 *Panel) bool {
	return p.panelArea.X+p.panelArea.Width > p2.panelArea.X && p2.panelArea.X+p2.panelArea.Width > p.panelArea.X &&
		p.panelArea.Y+p.panelArea.Height > p2.panelArea.Y && p2.panelArea.Y+p2.panelArea.Height > p.panelArea.Y
}

// SetEnabled sets the panel's enabled state.
// A disabled panel does not process events.
func (p *Panel) SetEnabled(state bool) {
	p.enabled = state
	p.Dispatch(core.GuiEnableEvent{Enabled: state})
}

// Enabled returns the enabled state of this panel.
func (p *Panel) Enabled() bool {
	return p.enabled
}

// ContentCoords converts the specified absolute coordinates to the panel's relative content coordinates.
func (p *Panel) ContentCoords(wx, wy float32) (float32, float32) {
	cx := wx - p.panelArea.X - p.paddings.Left
	cy := wy - p.panelArea.Y - p.paddings.Top
	return cx, cy
}

// RenderSetup is called by the engine before drawing the object.
func (p *Panel) RenderSetup(gl *gls.GLS, _ *core.RenderInfo) {
	// Sets texture valid flag in uniforms if the material has texture
	if p.material.TextureCount() > 0 {
		p.udata.textureValid[0] = 1
	} else {
		p.udata.textureValid[0] = 0
	}

	// Sets model matrix
	var mm math32.Matrix4
	p.SetModelMatrix(gl, &mm)

	// Transfer model matrix uniform
	location := p.uniMatrix.Location(gl)
	gl.UniformMatrix4fv(location, 1, false, &mm[0])

	// Transfer panel parameters combined uniform
	location = p.uniPanel.Location(gl)
	const vec4count = 8
	gl.Uniform4fv(location, vec4count, &p.udata.bounds.X)
}

// SetModelMatrix calculates and sets the specified matrix with the model matrix for this panel.
func (p *Panel) SetModelMatrix(gl *gls.GLS, mm *math32.Matrix4) {
	// Get scale of window (for HiDPI support)
	sX, sY := GetManager().window.GetScale()

	// Get the viewport width and height
	_, _, width, height := gl.GetViewport()

	// Compute common factors
	fX := 2 * float32(sX) / float32(width)
	fY := 2 * float32(sY) / float32(height)

	// Calculate the model matrix
	// Convert pixel coordinates to standard OpenGL clip coordinates and scale the quad for the viewport
	mm.Set(
		fX*p.panelArea.Width, 0, 0, fX*p.panelArea.X-1,
		0, fY*p.panelArea.Height, 0, 1-fY*p.panelArea.Y,
		0, 0, 1, p.Position().Z,
		0, 0, 0, 1,
	)
}

// setContentSize is an internal version of SetContentSize() which allows choosing if the panel will dispatch an event.
func (p *Panel) setContentSize(width, height float32, dispatch bool) {
	eWidth := width + p.paddings.Left + p.paddings.Right
	eHeight := height + p.paddings.Top + p.paddings.Bottom
	p.resize(eWidth, eHeight, dispatch)
}

// updateBounds calculates the panel's bounds considering its parent's bounds.
func (p *Panel) updateBounds(parent IPanel) {
	// Set default bounds to be entire panel texture
	p.udata.bounds = Rect{Width: 1, Height: 1}
	if parent == nil {
		// If this panel has no parent, its position is its position
		p.panelArea.X = p.Position().X
		p.panelArea.Y = p.Position().Y
		// No clipping necessary
		p.clipArea = p.panelArea
	} else {
		// Coordinates are relative to the parent internal content rectangle
		parentContentArea := parent.ContentArea()
		p.panelArea.X = p.Position().X + parentContentArea.X
		p.panelArea.Y = p.Position().Y + parentContentArea.Y
		// Clip the panel area by the parent's content clipped by the parent's clip area
		p.clipArea = p.panelArea.Clip(parentContentArea.Clip(parent.ClipArea()))
		// Update bounds in texture coordinates
		p.udata.bounds = Rect{
			X:      (p.clipArea.X - p.panelArea.X) / p.panelArea.Width,
			Y:      (p.clipArea.Y - p.panelArea.Y) / p.panelArea.Height,
			Width:  (p.clipArea.X + p.clipArea.Width - p.panelArea.X) / p.panelArea.Width,
			Height: (p.clipArea.Y + p.clipArea.Height - p.panelArea.Y) / p.panelArea.Height,
		}
	}
}

// resize tries to set the external size of the panel to the specified dimensions.
// It recalculates the size and positions of the internal areas.
// The padding sizes are kept and the content area size is adjusted.
// If the panel is decreased, its minimum size is determined by the paddings.
// If dispatch is true, an OnResize event will be emitted.
func (p *Panel) resize(width, height float32, dispatch bool) {
	width = math32.Round(width)
	height = math32.Round(height)
	// update content area
	p.contentArea.X = p.paddings.Left
	p.contentArea.Y = p.paddings.Top
	p.contentArea.Width = width - p.paddings.Left - p.paddings.Right
	if p.contentArea.Width < 0 {
		p.contentArea.Width = 0
	}
	p.contentArea.Height = height - p.paddings.Top - p.paddings.Bottom
	if p.contentArea.Height < 0 {
		p.contentArea.Height = 0
	}
	// Set final panel dimensions
	p.panelArea.Width = p.paddings.Left + p.contentArea.Width + p.paddings.Right
	p.panelArea.Height = p.paddings.Top + p.contentArea.Height + p.paddings.Bottom
	p.SetChanged(true)
	if dispatch {
		p.Dispatch(core.GuiResizeEvent{})
	}
}
