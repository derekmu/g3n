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

/*********************************************

 Panel areas:
 +------------------------------------------+
 |  Margin area                             |
 |  +------------------------------------+  |
 |  |  Border area                       |  |
 |  |  +------------------------------+  |  |
 |  |  | Padding area                 |  |  |
 |  |  |  +------------------------+  |  |  |
 |  |  |  | Content area           |  |  |  |
 |  |  |  |                        |  |  |  |
 |  |  |  |                        |  |  |  |
 |  |  |  +------------------------+  |  |  |
 |  |  |                              |  |  |
 |  |  +------------------------------+  |  |
 |  |                                    |  |
 |  +------------------------------------+  |
 |                                          |
 +------------------------------------------+

*********************************************/

// IPanel is the interface for all panel types.
type IPanel interface {
	graphic.IGraphic
	GetPanel() *Panel
	Width() float32
	Height() float32
	Enabled() bool
	SetEnabled(bool)
	InsideBorders(x, y float32) bool
	SetZLayerDelta(zLayerDelta int)
	ZLayerDelta() int
	SetPosition(x, y float32)
	SetPositionX(x float32)
	SetPositionY(y float32)
	SetPositionZ(y float32)
}

// Panel is 2D rectangular graphic which by default has a quad geometry.
// When using the default geometry, a panel has margins, borders, paddings, and a content area.
// The content area can be associated with a texture.
// It is the building block of most GUI elements.
type Panel struct {
	*graphic.Graphic
	mat         *material.Material
	tex         *texture.Texture2D
	zLayerDelta int // Z-layer relative to parent

	bounded bool // Whether panel is bounded by its parent
	enabled bool // Whether events should be processed for this panel

	marginSizes  RectBounds // external margin sizes in pixel coordinates
	borderSizes  RectBounds // border sizes in pixel coordinates
	paddingSizes RectBounds // padding sizes in pixel coordinates
	content      Rect       // content rectangle in pixel coordinates

	pospix math32.Vector2 // Absolute screen position
	width  float32        // External size
	height float32        // External size

	xmin float32 // minimum absolute x this panel can use
	xmax float32 // maximum absolute x this panel can use
	ymin float32 // minimum absolute y this panel can use
	ymax float32 // maximum absolute y this panel can use

	// Uniforms sent to shader
	uniMatrix gls.Uniform // model matrix uniform location cache
	uniPanel  gls.Uniform // panel parameters uniform location cache
	udata     struct {    // Combined uniform data 8 * vec4
		bounds       math32.Vector4 // panel bounds in texture coordinates
		borders      math32.Vector4 // panel borders in texture coordinates
		paddings     math32.Vector4 // panel paddings in texture coordinates
		content      math32.Vector4 // panel content area in texture coordinates
		borderColor  math32.Color4  // panel border color
		paddingColor math32.Color4  // panel padding color
		contentColor math32.Color4  // panel content color
		textureValid float32        // texture valid flag (bool)
		dummy        [3]float32     // complete 8 * vec4
	}
}

// PanelStyle contains all the styling attributes of a Panel.
type PanelStyle struct {
	Margin      RectBounds
	Border      RectBounds
	Padding     RectBounds
	BorderColor math32.Color4
	BgColor     math32.Color4
}

// BasicStyle extends PanelStyle by adding a foreground color.
// Many GUI components can be styled using BasicStyle or redeclared versions thereof (e.g. ButtonStyle).
type BasicStyle struct {
	PanelStyle
	FgColor math32.Color4
}

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

// NewPanel a new panel with the specified dimensions.
func NewPanel(width, height float32) *Panel {
	p := new(Panel)
	p.InitPanel(p, width, height)
	return p
}

// InitPanel initializes this panel and is normally used by other types which embed a panel.
func (p *Panel) InitPanel(ipan IPanel, width, height float32) {
	p.width = width
	p.height = height

	// Initialize material
	p.mat = material.NewMaterial()
	p.mat.SetUseLights(material.UseLightNone)
	p.mat.SetShader("panel")
	p.mat.SetTransparent(true)

	// Initialize graphic
	p.Graphic = graphic.NewGraphic(ipan, panelQuadGeometry.Incref(), gls.TRIANGLES)
	p.AddMaterial(p, p.mat, 0, 0)

	// Initialize uniforms location caches
	p.uniMatrix.Init("ModelMatrix")
	p.uniPanel.Init("Panel")

	// Set defaults
	p.udata.borderColor = math32.Color4{0, 0, 0, 1}
	p.bounded = true
	p.enabled = true
	p.resize(width, height, true)
}

// InitializeGraphic initializes this panel with an alternative graphic.
func (p *Panel) InitializeGraphic(width, height float32, gr *graphic.Graphic) {
	p.Graphic = gr
	p.width = width
	p.height = height

	// Initializes uniforms location caches
	p.uniMatrix.Init("ModelMatrix")
	p.uniPanel.Init("Panel")

	// Set defaults
	p.udata.borderColor = math32.Color4{0, 0, 0, 1}
	p.bounded = true
	p.enabled = true
	p.resize(width, height, true)
}

// GetPanel satisfies the IPanel interface and returns pointer to this panel.
func (p *Panel) GetPanel() *Panel {
	return p
}

// Material returns a pointer for the panel's material.
func (p *Panel) Material() *material.Material {
	return p.mat
}

// SetTexture changes the panel's texture.
// It returns a pointer to the previous texture.
func (p *Panel) SetTexture(tex *texture.Texture2D) *texture.Texture2D {
	prevtex := p.tex
	p.Material().RemoveTexture(prevtex)
	p.tex = tex
	if tex != nil {
		p.Material().AddTexture(p.tex)
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
		log.Printf("Invalid panel width:%v", width)
		width = 0
	}
	if height < 0 {
		log.Printf("Invalid panel height:%v", height)
		height = 0
	}
	p.resize(width, height, true)
}

// SetWidth sets this panel external width.
func (p *Panel) SetWidth(width float32) {
	p.SetSize(width, p.height)
}

// SetHeight sets this panel external height.
func (p *Panel) SetHeight(height float32) {
	p.SetSize(p.width, height)
}

// SetContentAspectWidth sets content width of this panel while maintaining the same aspect ratio.
func (p *Panel) SetContentAspectWidth(width float32) {
	aspect := p.content.Width / p.content.Height
	height := width / aspect
	p.SetContentSize(width, height)
}

// SetContentAspectHeight sets content height of this panel while maintaining the same aspect ratio.
func (p *Panel) SetContentAspectHeight(height float32) {
	aspect := p.content.Width / p.content.Height
	width := height / aspect
	p.SetContentSize(width, height)
}

// Size returns this panel external width and height.
func (p *Panel) Size() (float32, float32) {
	return p.width, p.height
}

// Width returns the panel external width.
func (p *Panel) Width() float32 {
	return p.width
}

// Height returns the panel external height.
func (p *Panel) Height() float32 {
	return p.height
}

// ContentWidth returns the width of the content area.
func (p *Panel) ContentWidth() float32 {
	return p.content.Width
}

// ContentHeight returns the height of the content area.
func (p *Panel) ContentHeight() float32 {
	return p.content.Height
}

// SetMargins sets the panel's margin sizes.
func (p *Panel) SetMargins(src RectBounds) {
	p.marginSizes = src
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// Margins returns the panel's margin sizes.
func (p *Panel) Margins() RectBounds {
	return p.marginSizes
}

// SetBorders sets the panel's border sizes.
func (p *Panel) SetBorders(src RectBounds) {
	p.borderSizes = src
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// Borders returns the panel's border sizes.
func (p *Panel) Borders() RectBounds {
	return p.borderSizes
}

// SetPaddings sets the panel's padding sizes.
func (p *Panel) SetPaddings(src RectBounds) {
	p.paddingSizes = src
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// Paddings returns the panel's padding sizes.
func (p *Panel) Paddings() RectBounds {
	return p.paddingSizes
}

// SetBorderColor sets the panel's border color.
func (p *Panel) SetBorderColor(color math32.Color4) {
	p.udata.borderColor = color
	p.SetChanged(true)
}

// BorderColor returns the panel's border color.
func (p *Panel) BorderColor() math32.Color4 {
	return p.udata.borderColor
}

// SetPaddingColor sets the panel's padding color.
func (p *Panel) SetPaddingColor(color math32.Color4) {
	p.udata.paddingColor = color
	p.SetChanged(true)
}

// SetColor sets the panel's padding and content color.
func (p *Panel) SetColor(color math32.Color4) *Panel {
	p.udata.paddingColor = color
	p.udata.contentColor = color
	p.SetChanged(true)
	return p
}

// SetContentColor sets the panel's content color.
func (p *Panel) SetContentColor(color math32.Color4) *Panel {
	p.udata.contentColor = color
	p.SetChanged(true)
	return p
}

// ContentColor returns the panel's content color.
func (p *Panel) ContentColor() math32.Color4 {
	return p.udata.contentColor
}

// ApplyStyle applies the specified style to the panel.
func (p *Panel) ApplyStyle(ps *PanelStyle) {
	p.udata.borderColor = ps.BorderColor
	p.udata.paddingColor = ps.BgColor
	p.udata.contentColor = ps.BgColor
	p.marginSizes = ps.Margin
	p.borderSizes = ps.Border
	p.paddingSizes = ps.Padding
	p.resize(p.calcWidth(), p.calcHeight(), true)
}

// SetContentSize sets the panel's content size.
func (p *Panel) SetContentSize(width, height float32) {
	p.setContentSize(width, height, true)
}

// SetContentWidth sets the panel's content width.
func (p *Panel) SetContentWidth(width float32) {
	p.SetContentSize(width, p.content.Height)
}

// SetContentHeight sets the panel's content height.
func (p *Panel) SetContentHeight(height float32) {
	p.SetContentSize(p.content.Width, height)
}

// MinWidth returns the minimum width of this panel (assuming content width was 0).
func (p *Panel) MinWidth() float32 {
	return p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
}

// MinHeight returns the minimum height of this panel (assuming content height was 0).
func (p *Panel) MinHeight() float32 {
	return p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
}

// Pospix returns the panel's absolute coordinate.
func (p *Panel) Pospix() math32.Vector2 {
	return p.pospix
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

// Bounded returns the panel's bounded state.
func (p *Panel) Bounded() bool {
	return p.bounded
}

// SetBounded sets the panel's bounded state.
func (p *Panel) SetBounded(bounded bool) {
	p.bounded = bounded
	p.SetChanged(true)
}

// UpdateMatrixWorld overrides the standard core.Node version which is called by the Engine before rendering the frame.
func (p *Panel) UpdateMatrixWorld() {
	par := p.Parent()
	if par == nil {
		p.updateBounds(nil)
	} else {
		// Panel has parent
		par, ok := par.(IPanel)
		if ok {
			p.updateBounds(par.GetPanel())
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
	return x >= p.pospix.X && y >= p.pospix.Y && x < (p.pospix.X+p.width) && y < (p.pospix.Y+p.height)
}

// InsideBorders returns whether a screen position is inside the panel borders, including the border width.
// Unlike ContainsPosition, it does not consider the panel margins.
func (p *Panel) InsideBorders(x, y float32) bool {
	return x >= (p.pospix.X+p.marginSizes.Left) && x < (p.pospix.X+p.width-p.marginSizes.Right) &&
		y >= (p.pospix.Y+p.marginSizes.Top) && y < (p.pospix.Y+p.height-p.marginSizes.Bottom)
}

// Intersects returns whether this panel intersects with another panel.
func (p *Panel) Intersects(p2 *Panel) bool {
	return p.pospix.X+p.width > p2.pospix.X && p2.pospix.X+p2.width > p.pospix.X &&
		p.pospix.Y+p.height > p2.pospix.Y && p2.pospix.Y+p2.height > p.pospix.Y
}

// SetEnabled sets the panel's enabled state.
// A disabled panel does not process events.
func (p *Panel) SetEnabled(state bool) {
	p.enabled = state
	p.Dispatch(OnEnable, nil)
}

// Enabled returns the enabled state of this panel.
func (p *Panel) Enabled() bool {
	return p.enabled
}

// ContentCoords converts the specified absolute coordinates to the panel's relative content coordinates.
func (p *Panel) ContentCoords(wx, wy float32) (float32, float32) {
	cx := wx - p.pospix.X -
		p.paddingSizes.Left -
		p.borderSizes.Left -
		p.marginSizes.Left
	cy := wy - p.pospix.Y -
		p.paddingSizes.Top -
		p.borderSizes.Top -
		p.marginSizes.Top
	return cx, cy
}

// setContentSize is an internal version of SetContentSize() which allows choosing if the panel will recalculate its layout and dispatch event.
// It is normally used by layout managers when setting the panel content size to avoid another invocation of the layout eventManager.
func (p *Panel) setContentSize(width, height float32, dispatch bool) {
	// Calculates the new desired external width and height
	eWidth := width +
		p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
	eHeight := height +
		p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
	p.resize(eWidth, eHeight, dispatch)
}

// updateBounds is called by UpdateMatrixWorld() to calculate the panel's bounds considering the bounds of its parent.
func (p *Panel) updateBounds(par *Panel) {
	if par == nil {
		// If this panel has no parent, its position is its position
		p.pospix.X = p.Position().X
		p.pospix.Y = p.Position().Y
	} else if p.bounded {
		// If this panel is bounded to its parent, its coordinates are relative to the parent internal content rectangle.
		p.pospix.X = p.Position().X + par.pospix.X + par.marginSizes.Left + par.borderSizes.Left + par.paddingSizes.Left
		p.pospix.Y = p.Position().Y + par.pospix.Y + par.marginSizes.Top + par.borderSizes.Top + par.paddingSizes.Top
	} else {
		// Otherwise its coordinates are relative to the parent outer coordinates.
		p.pospix.X = p.Position().X + par.pospix.X
		p.pospix.Y = p.Position().Y + par.pospix.Y
	}
	// Maximum x,y coordinates for this panel
	p.xmin = p.pospix.X
	p.ymin = p.pospix.Y
	p.xmax = p.pospix.X + p.width
	p.ymax = p.pospix.Y + p.height
	// Set default bounds to be entire panel texture
	p.udata.bounds = math32.Vector4{0, 0, 1, 1}
	// If this panel has no parent or is unbounded then the default bounds are correct
	if par == nil || !p.bounded {
		return
	}
	// From here on panel has parent and is bounded by parent
	// Get the parent content area minimum and maximum absolute coordinates
	pxmin := par.pospix.X + par.marginSizes.Left + par.borderSizes.Left + par.paddingSizes.Left
	if pxmin < par.xmin {
		pxmin = par.xmin
	}
	pymin := par.pospix.Y + par.marginSizes.Top + par.borderSizes.Top + par.paddingSizes.Top
	if pymin < par.ymin {
		pymin = par.ymin
	}
	pxmax := par.pospix.X + par.width - (par.marginSizes.Right + par.borderSizes.Right + par.paddingSizes.Right)
	if pxmax > par.xmax {
		pxmax = par.xmax
	}
	pymax := par.pospix.Y + par.height - (par.marginSizes.Bottom + par.borderSizes.Bottom + par.paddingSizes.Bottom)
	if pymax > par.ymax {
		pymax = par.ymax
	}
	// Update this panel minimum x and y coordinates.
	if p.xmin < pxmin {
		p.xmin = pxmin
	}
	if p.ymin < pymin {
		p.ymin = pymin
	}
	// Update this panel maximum x and y coordinates.
	if p.xmax > pxmax {
		p.xmax = pxmax
	}
	if p.ymax > pymax {
		p.ymax = pymax
	}
	// If this panel is bounded to its parent, calculates the bounds
	// for clipping in texture coordinates
	if p.pospix.X < p.xmin {
		p.udata.bounds.X = (p.xmin - p.pospix.X) / p.width
	}
	if p.pospix.Y < p.ymin {
		p.udata.bounds.Y = (p.ymin - p.pospix.Y) / p.height
	}
	if p.pospix.X+p.width > p.xmax {
		p.udata.bounds.Z = (p.xmax - p.pospix.X) / p.width
	}
	if p.pospix.Y+p.height > p.ymax {
		p.udata.bounds.W = (p.ymax - p.pospix.Y) / p.height
	}
}

// calcWidth calculates the panel's external width.
func (p *Panel) calcWidth() float32 {
	return p.content.Width +
		p.paddingSizes.Left + p.paddingSizes.Right +
		p.borderSizes.Left + p.borderSizes.Right +
		p.marginSizes.Left + p.marginSizes.Right
}

// calcHeight calculates the panel's external height.
func (p *Panel) calcHeight() float32 {
	return p.content.Height +
		p.paddingSizes.Top + p.paddingSizes.Bottom +
		p.borderSizes.Top + p.borderSizes.Bottom +
		p.marginSizes.Top + p.marginSizes.Bottom
}

// resize tries to set the external size of the panel to the specified dimensions.
// It recalculates the size and positions of the internal areas.
// The margins, borders and padding sizes are kept and the content area size is adjusted.
// If the panel is decreased, its minimum size is determined by the margins, borders and paddings.
// If dispatch is true, the layout will be recalculated and an OnResize event will be emitted.
func (p *Panel) resize(width, height float32, dispatch bool) {
	var padding Rect
	var border Rect

	width = math32.Round(width)
	height = math32.Round(height)

	// Adjust content width
	p.content.Width = width - p.marginSizes.Left - p.marginSizes.Right - p.borderSizes.Left - p.borderSizes.Right - p.paddingSizes.Left - p.paddingSizes.Right
	if p.content.Width < 0 {
		p.content.Width = 0
	}
	// Adjusts content height
	p.content.Height = height - p.marginSizes.Top - p.marginSizes.Bottom - p.borderSizes.Top - p.borderSizes.Bottom - p.paddingSizes.Top - p.paddingSizes.Bottom
	if p.content.Height < 0 {
		p.content.Height = 0
	}

	// Adjust other area widths
	padding.Width = p.paddingSizes.Left + p.content.Width + p.paddingSizes.Right
	border.Width = p.borderSizes.Left + padding.Width + p.borderSizes.Right
	// Adjust other area heights
	padding.Height = p.paddingSizes.Top + p.content.Height + p.paddingSizes.Bottom
	border.Height = p.borderSizes.Top + padding.Height + p.borderSizes.Bottom

	// Set area positions
	border.X = p.marginSizes.Left
	border.Y = p.marginSizes.Top
	padding.X = border.X + p.borderSizes.Left
	padding.Y = border.Y + p.borderSizes.Top
	p.content.X = padding.X + p.paddingSizes.Left
	p.content.Y = padding.Y + p.paddingSizes.Top

	// Set final panel dimensions
	p.width = p.marginSizes.Left + border.Width + p.marginSizes.Right
	p.height = p.marginSizes.Top + border.Height + p.marginSizes.Bottom

	// Update border uniform in texture coordinates (0,0 -> 1,1)
	p.udata.borders = math32.Vector4{
		border.X / p.width,
		border.Y / p.height,
		border.Width / p.width,
		border.Height / p.height,
	}
	// Update padding uniform in texture coordinates (0,0 -> 1,1)
	p.udata.paddings = math32.Vector4{
		padding.X / p.width,
		padding.Y / p.height,
		padding.Width / p.width,
		padding.Height / p.height,
	}
	// Update content uniform in texture coordinates (0,0 -> 1,1)
	p.udata.content = math32.Vector4{
		p.content.X / p.width,
		p.content.Y / p.height,
		p.content.Width / p.width,
		p.content.Height / p.height,
	}
	p.SetChanged(true)

	// Update layout and dispatch event
	if !dispatch {
		return
	}
	p.Dispatch(OnResize, nil)
}

// RenderSetup is called by the engine before drawing the object.
func (p *Panel) RenderSetup(gl *gls.GLS, _ *core.RenderInfo) {
	// Sets texture valid flag in uniforms if the material has texture
	if p.mat.TextureCount() > 0 {
		p.udata.textureValid = 1
	} else {
		p.udata.textureValid = 0
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
		fX*p.width, 0, 0, fX*p.pospix.X-1,
		0, fY*p.height, 0, 1-fY*p.pospix.Y,
		0, 0, 1, p.Position().Z,
		0, 0, 0, 1,
	)
}
