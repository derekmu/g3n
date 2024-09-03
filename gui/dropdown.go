// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/gui/assets/icon"
)

// Dropdown represents a dropdown GUI element.
type Dropdown struct {
	Panel                        // Embedded panel
	icon         *Label          // internal label with icon
	list         *List           // internal list
	styles       *DropdownStyles // pointer to dropdown styles
	litem        *ImageLabel     // Item shown in drop box (copy of selected)
	selItem      *ImageLabel     // selected item from list
	overDropdown bool
	overList     bool
	focus        bool
}

// DropdownStyle contains the styling of a Dropdown.
type DropdownStyle BasicStyle

// DropdownStyles contains a DropdownStyle for each valid GUI state.
type DropdownStyles struct {
	Normal   DropdownStyle
	Over     DropdownStyle
	Focus    DropdownStyle
	Disabled DropdownStyle
}

// NewDropdown creates and returns a pointer to a new drop down widget with the specified width.
func NewDropdown(width float32, item *ImageLabel) *Dropdown {
	dd := new(Dropdown)
	dd.styles = &StyleDefault().Dropdown
	dd.litem = item

	dd.Panel.InitPanel(dd, width, 0)
	dd.Panel.Subscribe(OnMouseDown, dd.onMouse)
	dd.Panel.Subscribe(OnCursorEnter, dd.onCursor)
	dd.Panel.Subscribe(OnCursorLeave, dd.onCursor)
	dd.Panel.Subscribe(OnResize, func(name string, ev interface{}) { dd.recalc() })

	// ListItem
	dd.Panel.Add(dd.litem)

	// Create icon
	dd.icon = NewIconLabel(" ")
	dd.icon.SetFontSize(StyleDefault().Label.PointSize * 1.3)
	dd.icon.SetText(string(icon.ArrowDropDown))
	dd.Panel.Add(dd.icon)

	/// Create list
	dd.list = NewVList(0, 0)
	dd.list.bounded = false
	dd.list.zLayerDelta = 1
	dd.list.dropdown = true
	dd.list.SetVisible(false)

	dd.Panel.Subscribe(OnKeyDown, dd.list.onKeyEvent)
	dd.Subscribe(OnMouseDownOut, func(s string, i interface{}) {
		// Hide list when clicked out
		if dd.list.Visible() {
			dd.list.SetVisible(false)
		}
	})

	dd.list.Subscribe(OnCursorEnter, func(evname string, ev interface{}) {
		dd.Dispatch(OnCursorLeave, ev)
	})
	dd.list.Subscribe(OnCursorLeave, func(evname string, ev interface{}) {
		dd.Dispatch(OnCursorEnter, ev)
	})

	dd.list.Subscribe(OnChange, dd.onListChangeEvent)
	dd.Panel.Add(dd.list)

	dd.update()
	// This will trigger recalc()
	dd.Panel.SetContentHeight(item.Height())
	return dd
}

// Add adds a list item at the end of the list
func (dd *Dropdown) Add(item *ImageLabel) {
	dd.list.Add(item)
}

// InsertAt inserts a list item at the specified position
// Returs true if the item was successfully inserted
func (dd *Dropdown) InsertAt(pos int, item *ImageLabel) {
	dd.list.InsertAt(pos, item)
}

// RemoveAt removes the list item from the specified position
// Returs true if the item was successfully removed
func (dd *Dropdown) RemoveAt(pos int) {
	dd.list.RemoveAt(pos)
}

// ItemAt returns the list item at the specified position
func (dd *Dropdown) ItemAt(pos int) *ImageLabel {
	return dd.list.ItemAt(pos).(*ImageLabel)
}

// Len returns the number of items in the dropdown's list.
func (dd *Dropdown) Len() int {
	return dd.list.Len()
}

// Selected returns the currently selected item or nil if no item was selected
func (dd *Dropdown) Selected() *ImageLabel {
	return dd.selItem
}

// SelectedPos returns the currently selected position or -1 if no item was selected
func (dd *Dropdown) SelectedPos() int {
	return dd.list.selected()
}

// SetSelected sets the selected item
func (dd *Dropdown) SetSelected(item *ImageLabel) {
	dd.list.SetSelected(dd.selItem, false)
	dd.list.SetSelected(item, true)
	dd.copySelected()
	dd.update()
}

// SelectPos selects the item at the specified position
func (dd *Dropdown) SelectPos(pos int) {
	dd.list.SetSelected(dd.selItem, false)
	dd.list.SelectPos(pos, true)
	dd.Dispatch(OnChange, nil)
}

// SetStyles sets the drop down styles overriding the default style
func (dd *Dropdown) SetStyles(dds *DropdownStyles) {
	dd.styles = dds
	dd.update()
}

// onMouse receives subscribed mouse events over the dropdown
func (dd *Dropdown) onMouse(evname string, _ interface{}) {
	GetManager().SetKeyFocus(dd.list)
	if evname == OnMouseDown {
		dd.list.SetVisible(!dd.list.Visible())
		return
	}
}

// onCursor receives subscribed cursor events over the dropdown
func (dd *Dropdown) onCursor(evname string, _ interface{}) {
	if evname == OnCursorEnter {
		dd.overDropdown = true
	}
	if evname == OnCursorLeave {
		dd.overDropdown = false
	}
	dd.update()
}

// copySelected copy to the dropdown panel the selected item
// from the list.
func (dd *Dropdown) copySelected() {
	selected := dd.list.Selected()
	if len(selected) > 0 {
		dd.selItem = selected[0].(*ImageLabel)
		dd.litem.CopyFields(dd.selItem)
		dd.litem.SetWidth(dd.selItem.Width())
		dd.recalc()
		dd.Dispatch(OnChange, nil)
	} else {
		return
	}
}

// onListChangeEvent is called when an item in the list is selected
func (dd *Dropdown) onListChangeEvent(_ string, _ interface{}) {
	dd.copySelected()
}

// recalc recalculates the dimensions and positions of the dropdown
// panel, children and list
func (dd *Dropdown) recalc() {
	// Dropdown icon position
	posx := dd.Panel.ContentWidth() - dd.icon.Width()
	dd.icon.SetPosition(posx, 0)

	// List item position and width
	ipan := dd.litem.GetPanel()
	ipan.SetPosition(0, 0)
	height := ipan.Height()

	// List position
	dd.list.SetWidth(dd.Panel.Width())
	dd.list.SetHeight(6*height + 1)
	dd.list.SetPositionX(0)
	dd.list.SetPositionY(dd.Panel.Height())
}

// update updates the visual state
func (dd *Dropdown) update() {
	if dd.overDropdown || dd.overList {
		dd.applyStyle(&dd.styles.Over)
		dd.list.ApplyStyle(StyleOver)
		return
	}
	if dd.focus {
		dd.applyStyle(&dd.styles.Focus)
		dd.list.ApplyStyle(StyleFocus)
		return
	}
	dd.applyStyle(&dd.styles.Normal)
	dd.list.ApplyStyle(StyleNormal)
}

// applyStyle applies the specified style
func (dd *Dropdown) applyStyle(s *DropdownStyle) {
	dd.Panel.ApplyStyle(&s.PanelStyle)
}
