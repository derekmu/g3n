// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"log"

	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/geometry"
	"github.com/derekmu/g3n/gls"
	"github.com/derekmu/g3n/graphic"
	"github.com/derekmu/g3n/material"
	"github.com/derekmu/g3n/math32"
	"github.com/derekmu/g3n/texture"
)

// Quad geometry shared by all Panels.
var panelQuadGeometry *geometry.Geometry

func init() {
	// Builds array with vertex positions and texture coordinates
	positions := math32.NewArrayF32(0, 20)
	positions.Append(
		0, 0, 0, 0, 0,
		0, -1, 0, 0, 1,
		1, -1, 0, 1, 1,
		1, 0, 0, 1, 0,
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
	geom.SetPermanent(true)
	panelQuadGeometry = geom
}

// IPanel is the interface for all panel types.
type IPanel interface {
	graphic.IGraphic
	core.IDispatcher[core.GuiEvent]
	SetPositionZ(z float32)
	ZLayerDelta() int
	Enabled() bool
	ContainsMouse(x, y float64) bool

	Width() int
	Height() int
	ContentWidth() int
	ContentHeight() int

	getContentArea() Rect
	getClipArea() Rect
}

var _ IPanel = &Panel{}

// Panel is 2D rectangular graphic which by default has a quad geometry.
// It is the building block of GUI elements.
type Panel struct {
	graphic.Graphic
	core.Dispatcher[core.GuiEvent]
	material    *material.Material
	texture     *texture.Texture2D
	zLayerDelta int
	enabled     bool

	paddings    RectBounds
	panelArea   Rect
	contentArea Rect
	clipArea    Rect

	uniMatrix gls.Uniform
	uniPanel  gls.Uniform
	udata     struct {
		bounds       RectF
		color        math32.Color4
		textureValid [4]float32
	}
}

// NewPanel a new panel with the specified dimensions.
func NewPanel(width, height int) *Panel {
	p := new(Panel)
	p.InitPanel(p, width, height)
	return p
}

// InitPanel initializes this panel and is normally used by other types which embed a panel.
func (p *Panel) InitPanel(ipan IPanel, width, height int) {
	// Initialize material
	p.material = material.NewMaterial()
	p.material.SetUseLights(material.UseLightNone)
	p.material.SetShader("panel")
	p.material.SetTransparent(true)

	// Initialize graphic
	p.InitGraphic(ipan, panelQuadGeometry, gls.TRIANGLES)
	p.AddMaterial(ipan, p.material, 0, 0)

	// Initialize uniforms location caches
	p.uniMatrix.Init("uModelMatrix")
	p.uniPanel.Init("uPanel")

	// Set defaults
	p.enabled = true
	p.resize(width, height)
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

// SetTopChild moves the specified panel to be the last child of this panel.
func (p *Panel) SetTopChild(ipan IPanel) {
	// Remove panel and if found appends to the end
	if p.Remove(ipan) {
		p.Add(ipan)
		p.SetChanged(true)
	}
}

// Material returns a pointer for the panel's material.
func (p *Panel) Material() *material.Material {
	return p.material
}

// SetTexture changes the panel's texture.
func (p *Panel) SetTexture(tex *texture.Texture2D) {
	if tex != p.texture {
		if p.texture != nil {
			p.Material().RemoveTexture(p.texture)
		}
		p.texture = tex
		if tex != nil {
			p.Material().AddTexture(p.texture)
		}
	}
}

// ZLayerDelta returns the Z-layer of this panel relative to its parent.
func (p *Panel) ZLayerDelta() int {
	return p.zLayerDelta
}

// SetZLayerDelta sets the Z-layer of this panel relative to its parent.
func (p *Panel) SetZLayerDelta(zLayerDelta int) {
	p.zLayerDelta = zLayerDelta
}

// Enabled is whether the panel processes events.
func (p *Panel) Enabled() bool {
	return p.enabled
}

// SetEnabled sets the panel's enabled state.
func (p *Panel) SetEnabled(state bool) {
	p.enabled = state
	p.Dispatch(core.GuiEnableEvent{Enabled: state})
}

func (p *Panel) ContainsMouse(x, y float64) bool {
	return p.clipArea.Contains(int(x), int(y))
}

// Paddings is the panel's padding sizes.
func (p *Panel) Paddings() RectBounds {
	return p.paddings
}

// SetPaddings sets the panel's padding sizes.
func (p *Panel) SetPaddings(src RectBounds) {
	p.paddings = src
	p.resize(p.contentArea.Width+p.paddings.Left+p.paddings.Right, p.contentArea.Height+p.paddings.Top+p.paddings.Bottom)
}

// Width is the panel's width.
func (p *Panel) Width() int {
	return p.panelArea.Width
}

// Height is the panel's height.
func (p *Panel) Height() int {
	return p.panelArea.Height
}

// MinWidth returns the minimum width of this panel (assuming content width was 0).
func (p *Panel) MinWidth() int {
	return p.paddings.Left + p.paddings.Right
}

// MinHeight returns the minimum height of this panel (assuming content height was 0).
func (p *Panel) MinHeight() int {
	return p.paddings.Top + p.paddings.Bottom
}

// SetPosition sets the panel's relative position.
func (p *Panel) SetPosition(x, y float32) {
	p.Node.SetPositionX(math32.Round(x))
	p.Node.SetPositionY(math32.Round(y))
}

// SetSize sets this panel external width and height.
func (p *Panel) SetSize(width, height int) {
	if width < 0 {
		log.Printf("Invalid panel width %v", width)
		width = 0
	}
	if height < 0 {
		log.Printf("Invalid panel height %v", height)
		height = 0
	}
	p.resize(width, height)
}

// SetWidth sets this panel external width.
func (p *Panel) SetWidth(width int) {
	p.SetSize(width, p.panelArea.Height)
}

// SetHeight sets this panel external height.
func (p *Panel) SetHeight(height int) {
	p.SetSize(p.panelArea.Width, height)
}

// ContentWidth is the panel's content width.
func (p *Panel) ContentWidth() int {
	return p.contentArea.Width
}

// ContentHeight is the panel's content height.
func (p *Panel) ContentHeight() int {
	return p.contentArea.Height
}

// SetContentSize sets the panel's content size.
func (p *Panel) SetContentSize(width, height int) {
	eWidth := width + p.paddings.Left + p.paddings.Right
	eHeight := height + p.paddings.Top + p.paddings.Bottom
	p.resize(eWidth, eHeight)
}

// SetContentWidth sets the panel's content width.
func (p *Panel) SetContentWidth(width int) {
	p.SetContentSize(width, p.contentArea.Height)
}

// SetContentHeight sets the panel's content height.
func (p *Panel) SetContentHeight(height int) {
	p.SetContentSize(p.contentArea.Width, height)
}

// PanelColor is the color of the panel if there is no texture.
func (p *Panel) PanelColor() math32.Color {
	return p.udata.color
}

// SetPanelColor sets the panel's color.
func (p *Panel) SetPanelColor(color math32.Color) *Panel {
	p.udata.color = color.ToColor4()
	p.SetChanged(true)
	return p
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

// ContentCoords converts the specified absolute coordinates to the panel's relative content coordinates.
func (p *Panel) ContentCoords(wx, wy int) (int, int) {
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

	// Get scale of window (for HiDPI support)
	sX, sY := GetManager().window.GetScale()
	// Get the viewport width and height
	_, _, width, height := gl.GetViewport()
	// Compute common factors
	fX := 2 * float32(sX) / float32(width)
	fY := 2 * float32(sY) / float32(height)
	// Calculate the model matrix
	// Convert pixel coordinates to standard OpenGL clip coordinates and scale the quad for the viewport
	var mm math32.Matrix4
	mm.Set(
		fX*float32(p.panelArea.Width), 0, 0, fX*float32(p.panelArea.X)-1,
		0, fY*float32(p.panelArea.Height), 0, 1-fY*float32(p.panelArea.Y),
		0, 0, 1, p.Position().Z,
		0, 0, 0, 1,
	)

	// Transfer model matrix uniform
	location := p.uniMatrix.Location(gl)
	gl.UniformMatrix4fv(location, 1, false, &mm[0])

	// Transfer panel parameters combined uniform
	location = p.uniPanel.Location(gl)
	const vec4count = 8
	gl.Uniform4fv(location, vec4count, &p.udata.bounds.X)
}

func (p *Panel) getContentArea() Rect {
	return p.contentArea
}

func (p *Panel) getClipArea() Rect {
	return p.clipArea
}

// updateBounds calculates the panel's bounds considering its parent's bounds.
func (p *Panel) updateBounds(parent IPanel) {
	if parent == nil {
		// If this panel has no parent, its position is its position
		p.panelArea.X = int(p.Position().X)
		p.panelArea.Y = int(p.Position().Y)
		p.contentArea.X = p.panelArea.X + p.paddings.Left
		p.contentArea.Y = p.panelArea.Y + p.paddings.Top
		// No clipping necessary
		p.clipArea = p.panelArea
		// Set default bounds to be entire panel texture
		p.udata.bounds = RectF{Width: 1, Height: 1}
	} else {
		// Coordinates are relative to the parent's content
		parentContentArea := parent.getContentArea()
		p.panelArea.X = parentContentArea.X + int(p.Position().X)
		p.panelArea.Y = parentContentArea.Y + int(p.Position().Y)
		p.contentArea.X = p.panelArea.X + p.paddings.Left
		p.contentArea.Y = p.panelArea.Y + p.paddings.Top
		// Clip the panel area by the parent's content clipped by the parent's clip area
		p.clipArea = p.panelArea.Clip(parentContentArea.Clip(parent.getClipArea()))
		// Update bounds in texture coordinates
		p.udata.bounds = RectF{
			X:      float32(p.clipArea.X-p.panelArea.X) / float32(p.panelArea.Width),
			Y:      float32(p.clipArea.Y-p.panelArea.Y) / float32(p.panelArea.Height),
			Width:  float32(p.clipArea.X+p.clipArea.Width-p.panelArea.X) / float32(p.panelArea.Width),
			Height: float32(p.clipArea.Y+p.clipArea.Height-p.panelArea.Y) / float32(p.panelArea.Height),
		}
	}
}

// resize tries to set the external size of the panel to the specified dimensions.
// It recalculates the size and positions of the internal areas.
// The padding sizes are kept and the content area size is adjusted.
// If the panel is decreased, its minimum size is determined by the paddings.
func (p *Panel) resize(width, height int) {
	contentWidth := max(0, width-p.paddings.Left-p.paddings.Right)
	contentHeight := max(0, height-p.paddings.Top-p.paddings.Bottom)
	panelWidth := p.paddings.Left + contentWidth + p.paddings.Right
	panelHeight := p.paddings.Top + contentHeight + p.paddings.Bottom
	if contentWidth != p.contentArea.Width || contentHeight != p.contentArea.Height ||
		panelWidth != p.panelArea.Width || panelHeight != p.panelArea.Height {
		p.contentArea.Width = contentWidth
		p.contentArea.Height = contentHeight
		p.panelArea.Width = panelWidth
		p.panelArea.Height = panelHeight
		p.SetChanged(true)
		p.Dispatch(core.GuiResizeEvent{})
	}
}

// AddLabel adds a new Label to this panel.
// The panel will dynamically resize to fit the label if the expand parameter is true.
// The position of the label within the panel is determined by the align parameter.
func (p *Panel) AddLabel(text string, expand bool, align Align) *Label {
	label := NewLabel(text)
	p.Add(label)
	handleResize := func(event core.GuiEvent) bool {
		switch event.GuiEventType() {
		case core.GuiResize:
			width, height := p.ContentWidth(), p.ContentHeight()
			labelWidth, labelHeight := label.Width(), label.Height()
			// Sets new content width and height if necessary
			if expand {
				resize := false
				if width < labelWidth {
					width = labelWidth
					resize = true
				}
				if height < labelHeight {
					height = labelHeight
					resize = true
				}
				if resize {
					p.SetContentSize(width, height)
				}
			}
			// Position the label as desired
			if align != AlignNone {
				lx, ly := align.CalculatePosition(width, height, labelWidth, labelHeight)
				label.SetPosition(float32(lx), float32(ly))
			}
		default:
			return false
		}
		return true
	}
	label.Subscribe(handleResize)
	p.Subscribe(handleResize)
	handleResize(core.GuiResizeEvent{})
	return label
}
