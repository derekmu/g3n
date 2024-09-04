// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/derekmu/g3n/core"
	"github.com/derekmu/g3n/gui/assets/icon"
	"time"
)

// Menu is the menu GUI element
type Menu struct {
	Panel                // embedded panel
	bar      bool        // true for menu bar
	items    []*MenuItem // menu items
	autoOpen bool        // open sub menus when mouse over if true
	mitem    *MenuItem   // parent menu item for sub menu
}

// MenuItem is an option of a Menu
type MenuItem struct {
	Panel                     // embedded panel
	menu     *Menu            // pointer to parent menu
	licon    *Label           // optional left icon label
	label    *Label           // optional text label (nil for separators)
	shortcut *Label           // optional shorcut text label
	ricon    *Label           // optional right internal icon label for submenu
	id       string           // optional text id
	icode    int              // icon code (if icon is set)
	submenu  *Menu            // pointer to optional associated sub menu
	keyMods  core.ModifierKey // shortcut key modifier
	keyCode  core.Key         // shortcut key code
	disabled bool             // item disabled state
	selected bool             // selection state
}

var mapKeyModifier = map[core.ModifierKey]string{
	core.ModShift:   "Shift",
	core.ModControl: "Ctrl",
	core.ModAlt:     "Alt",
}
var mapKeyText = map[core.Key]string{
	core.KeyApostrophe: "'",
	core.KeyComma:      ",",
	core.KeyMinus:      "-",
	core.KeyPeriod:     ".",
	core.KeySlash:      "/",
	core.Key0:          "0",
	core.Key1:          "1",
	core.Key2:          "2",
	core.Key3:          "3",
	core.Key4:          "4",
	core.Key5:          "5",
	core.Key6:          "6",
	core.Key7:          "7",
	core.Key8:          "8",
	core.Key9:          "9",
	core.KeySemicolon:  ";",
	core.KeyEqual:      "=",
	core.KeyA:          "A",
	core.KeyB:          "B",
	core.KeyC:          "C",
	core.KeyD:          "D",
	core.KeyE:          "E",
	core.KeyF:          "F",
	core.KeyG:          "G",
	core.KeyH:          "H",
	core.KeyI:          "I",
	core.KeyJ:          "J",
	core.KeyK:          "K",
	core.KeyL:          "L",
	core.KeyM:          "M",
	core.KeyN:          "N",
	core.KeyO:          "O",
	core.KeyP:          "P",
	core.KeyQ:          "Q",
	core.KeyR:          "R",
	core.KeyS:          "S",
	core.KeyT:          "T",
	core.KeyU:          "U",
	core.KeyV:          "V",
	core.KeyW:          "W",
	core.KeyX:          "X",
	core.KeyY:          "Y",
	core.KeyZ:          "Z",
	core.KeyF1:         "F1",
	core.KeyF2:         "F2",
	core.KeyF3:         "F3",
	core.KeyF4:         "F4",
	core.KeyF5:         "F5",
	core.KeyF6:         "F6",
	core.KeyF7:         "F7",
	core.KeyF8:         "F8",
	core.KeyF9:         "F9",
	core.KeyF10:        "F10",
	core.KeyF11:        "F11",
	core.KeyF12:        "F12",
}

// NewMenuBar creates and returns a pointer to a new empty menu bar
func NewMenuBar() *Menu {
	m := NewMenu()
	m.bar = true
	m.Panel.Subscribe(OnMouseDownOut, m.onMouse)
	return m
}

// NewMenu creates and returns a pointer to a new empty vertical menu
func NewMenu() *Menu {
	m := new(Menu)
	m.Panel.InitPanel(m, 0, 0)
	m.items = make([]*MenuItem, 0)
	m.Panel.Subscribe(OnKeyDown, m.onKey)
	m.Panel.Subscribe(OnResize, m.onResize)
	return m
}

// AddOption creates and adds a new menu item to this menu with the
// specified text and returns the pointer to the created menu item.
func (m *Menu) AddOption(text string) *MenuItem {
	mi := newMenuItem(text)
	m.Panel.Add(mi)
	m.items = append(m.items, mi)
	mi.menu = m
	m.recalc()
	return mi
}

// AddSeparator creates and adds a new separator to the menu
func (m *Menu) AddSeparator() *MenuItem {
	mi := newMenuItem("")
	m.Panel.Add(mi)
	m.items = append(m.items, mi)
	mi.menu = m
	m.recalc()
	return mi
}

// AddMenu creates and adds a new menu item to this menu with the
// specified text and sub menu.
// Returns the pointer to the created menu item.
func (m *Menu) AddMenu(text string, subm *Menu) *MenuItem {
	mi := newMenuItem(text)
	mi.zLayerDelta = 1
	m.Panel.Add(mi)
	m.items = append(m.items, mi)
	mi.submenu = subm
	mi.submenu.SetVisible(false)
	mi.submenu.SetBounded(false)
	mi.submenu.mitem = mi
	mi.submenu.autoOpen = true
	mi.menu = m
	if !m.bar {
		mi.ricon = NewIconLabel(icon.PlayArrow)
		mi.Panel.Add(mi.ricon)
	}
	mi.Panel.Add(mi.submenu)
	mi.update()
	m.recalc()
	return mi
}

// onKey process subscribed key events
func (m *Menu) onKey(_ string, ev interface{}) {
	sel := m.selectedPos()
	kev := ev.(*core.KeyEvent)
	switch kev.Key {
	// Select next enabled menu item
	case core.KeyDown:
		if sel < 0 {
			return
		}
		mi := m.items[sel]
		// Select next enabled menu item
		if m.bar {
			// If selected item is not a sub menu, ignore
			if mi.submenu == nil {
				return
			}
			// Sets autoOpen and selects sub menu
			m.autoOpen = true
			mi.update()
			GetManager().SetKeyFocus(mi.submenu)
			mi.submenu.setSelectedPos(0)
			return
		}
		// Select next enabled menu item for vertical menu
		next := m.nextItem(sel)
		m.setSelectedPos(next)
	// Up -> Previous item for vertical menus
	case core.KeyUp:
		if sel < 0 {
			return
		}
		if m.bar {
			return
		}
		prev := m.prevItem(sel)
		m.setSelectedPos(prev)
	// Left -> Previous menu item for menu bar
	case core.KeyLeft:
		if sel < 0 {
			return
		}
		// For menu bar, select previous menu item
		if m.bar {
			prev := m.prevItem(sel)
			m.setSelectedPos(prev)
			return
		}
		// If menu has parent menu item
		if m.mitem != nil {
			if m.mitem.menu.bar {
				sel := m.mitem.menu.selectedPos()
				prev := m.mitem.menu.prevItem(sel)
				m.mitem.menu.setSelectedPos(prev)
			} else {
				m.mitem.menu.setSelectedItem(m.mitem)
			}
			GetManager().SetKeyFocus(m.mitem.menu)
			return
		}

	// Right -> Next menu bar item || Next sub menu
	case core.KeyRight:
		if sel < 0 {
			return
		}
		mi := m.items[sel]
		// For menu bar, select next menu item
		if m.bar {
			next := m.nextItem(sel)
			m.setSelectedPos(next)
			return
		}
		// Enter into sub menu
		if mi.submenu != nil {
			GetManager().SetKeyFocus(mi.submenu)
			mi.submenu.setSelectedPos(0)
			return
		}
		// If parent menu of this menu item is bar menu
		if m.mitem != nil && m.mitem.menu.bar {
			sel := m.mitem.menu.selectedPos()
			next := m.mitem.menu.nextItem(sel)
			m.mitem.menu.setSelectedPos(next)
			GetManager().SetKeyFocus(m.mitem.menu)
		}
	// Enter -> Select menu option
	case core.KeyEnter:
		if sel < 0 {
			return
		}
		mi := m.items[sel]
		mi.activate()
	// Check for menu items shortcuts
	default:
		var root *Menu
		if sel < 0 {
			root = m
		} else {
			mi := m.items[sel]
			root = mi.rootMenu()
		}
		found := root.checkKey(kev)
		if found == nil {
			return
		}
		if found.submenu == nil {
			found.activate()
			return
		}
		if found.menu.bar {
			found.menu.autoOpen = true
		}
		found.menu.setSelectedItem(found)
	}
}

// onMouse process subscribed mouse events for the menu
func (m *Menu) onMouse(_ string, _ interface{}) {
	// Clear menu bar after some time, to give time for menu items
	// to receive onMouse events.
	GetManager().SetTimeout(1*time.Millisecond, nil, func(arg interface{}) {
		m.autoOpen = false
		m.setSelectedPos(-1)
	})
}

// onResize process menu onResize events
func (m *Menu) onResize(_ string, _ interface{}) {
	if m.bar {
		m.recalcBar(false)
	}
}

// checkKey checks if this menu and any of its children contains
// a menu item with the specified key shortcut
func (m *Menu) checkKey(kev *core.KeyEvent) *MenuItem {
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.keyCode == kev.Key && mi.keyMods == kev.Mods {
			return mi
		}
		if mi.submenu != nil {
			found := mi.submenu.checkKey(kev)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

// setSelectedPos sets the menu item at the specified position as selected
// and all others as not selected.
func (m *Menu) setSelectedPos(pos int) {
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if i == pos {
			mi.selected = true
		} else {
			mi.selected = false
		}
		// If menu item has a sub menu, unselects the sub menu options recursively
		if mi.submenu != nil {
			mi.submenu.setSelectedPos(-1)
		}
		mi.update()
	}
}

// setSelectedItem sets the specified menu item as selected
// and all others as not selected
func (m *Menu) setSelectedItem(mitem *MenuItem) {
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi == mitem {
			mi.selected = true
		} else {
			mi.selected = false
		}
		// If menu item has a sub menu, unselects the sub menu options recursively
		if mi.submenu != nil {
			mi.submenu.setSelectedItem(nil)
		}
		mi.update()
	}
}

// selectedPos returns the position of the current selected menu item
// Returns -1 if no item selected
func (m *Menu) selectedPos() int {
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.selected {
			return i
		}
	}
	return -1
}

// nextItem returns the position of the next enabled option from the
// specified position
func (m *Menu) nextItem(pos int) int {
	res := 0
	for i := pos + 1; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.disabled || mi.label == nil {
			continue
		}
		res = i
		break
	}
	return res
}

// prevItem returns the position of previous enabled menu item from
// the specified position
func (m *Menu) prevItem(pos int) int {
	res := len(m.items) - 1
	for i := pos - 1; i >= 0 && i < len(m.items); i-- {
		mi := m.items[i]
		if mi.disabled || mi.label == nil {
			continue
		}
		res = i
		break
	}
	return res
}

// recalc recalculates the positions of this menu internal items
// and the content width and height of the menu
func (m *Menu) recalc() {
	if m.bar {
		m.recalcBar(true)
		return
	}

	// Find the maximum icon and label widths
	minWidth := float32(0)
	iconWidth := float32(0)
	labelWidth := float32(0)
	shortcutWidth := float32(0)
	riconWidth := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		minWidth = mi.MinWidth()
		// Separator
		if mi.label == nil {
			continue
		}
		// Left icon width
		if mi.licon != nil && mi.licon.width > iconWidth {
			iconWidth = mi.licon.width
		}
		// Option label width
		if mi.label.width > labelWidth {
			labelWidth = mi.label.width
		}
		// Shortcut label width
		if mi.shortcut != nil && mi.shortcut.width > shortcutWidth {
			shortcutWidth = mi.shortcut.width
		}
		// Right icon (submenu indicator) width
		if mi.ricon != nil && mi.ricon.width > riconWidth {
			riconWidth = mi.ricon.width
		}
	}
	width := minWidth + iconWidth + labelWidth + shortcutWidth + riconWidth

	// Sets the position and width of the menu items
	// The height is defined by the menu item itself
	px := float32(0)
	py := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		mi.SetPosition(px, py)
		mh := mi.minHeight()
		py += mh
		mi.SetSize(width, mh)
		mi.recalc(iconWidth, labelWidth, shortcutWidth)
	}
	m.SetContentSize(width, py)
}

// recalcBar recalculates the positions of this MenuBar internal items
// If setSize is true it also sets the size of the menu bar
func (m *Menu) recalcBar(setSize bool) {
	// Calculate the maximum item height
	height := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		if mi.minHeight() > height {
			height = mi.minHeight()
		}
	}

	// Calculates the y position of the items to center inside the menu panel
	py := (m.ContentHeight() - height) / 2
	if py < 0 {
		py = 0
	}

	// Sets the position of each item
	px := float32(0)
	for i := 0; i < len(m.items); i++ {
		mi := m.items[i]
		mi.SetPosition(px, py)
		width := float32(0)
		width = mi.minWidth()
		mi.SetSize(width, height)
		px += mi.Width()
	}

	// Sets the size of this menu if requested
	if setSize {
		m.SetContentSize(px, height)
	}
}

// newMenuItem creates and returns a pointer to a new menu item
// with the specified text.
func newMenuItem(text string) *MenuItem {
	mi := new(MenuItem)
	mi.Panel.InitPanel(mi, 0, 0)
	if text != "" {
		mi.label = NewLabel(text)
		mi.Panel.Add(mi.label)
		mi.Panel.Subscribe(OnCursorEnter, mi.onCursor)
		mi.Panel.Subscribe(OnCursor, mi.onCursor)
		mi.Panel.Subscribe(OnMouseDown, mi.onMouse)
	}
	mi.update()
	return mi
}

// SetIcon sets the left icon of this menu item
// If an image was previously set it is replaced by this icon
func (mi *MenuItem) SetIcon(icon string) *MenuItem {
	// Remove and dispose previous icon
	if mi.licon != nil {
		mi.Panel.Remove(mi.licon)
		mi.Dispose()
		mi.licon = nil
	}
	// Sets the new icon
	mi.licon = NewIconLabel(icon)
	mi.Panel.Add(mi.licon)
	mi.update()
	return mi
}

// SetText sets the text of this menu item
func (mi *MenuItem) SetText(text string) *MenuItem {
	if mi.label == nil {
		return mi
	}
	mi.label.SetText(text)
	mi.update()
	mi.menu.recalc()
	return mi
}

// SetShortcut sets the keyboard shortcut of this menu item
func (mi *MenuItem) SetShortcut(mods core.ModifierKey, key core.Key) *MenuItem {
	if mapKeyText[key] == "" {
		panic("Invalid menu shortcut key")
	}
	mi.keyMods = mods
	mi.keyCode = key

	// If parent menu is a menu bar, nothing more to do
	if mi.menu.bar {
		return mi
	}

	// Builds shortcut text
	text := ""
	if mi.keyMods&core.ModShift != 0 {
		text = mapKeyModifier[core.ModShift]
	}
	if mi.keyMods&core.ModControl != 0 {
		if text != "" {
			text += "+"
		}
		text += mapKeyModifier[core.ModControl]
	}
	if mi.keyMods&core.ModAlt != 0 {
		if text != "" {
			text += "+"
		}
		text += mapKeyModifier[core.ModAlt]
	}
	if text != "" {
		text += "+"
	}
	text += mapKeyText[key]

	// Creates and adds shortcut label
	mi.shortcut = NewLabel(text)
	mi.Panel.Add(mi.shortcut)
	mi.update()
	mi.menu.recalc()
	return mi
}

// SetEnabled sets the enabled state of this menu item
func (mi *MenuItem) SetEnabled(enabled bool) {
	mi.disabled = !enabled
	mi.update()
}

// SetId sets this menu item string id which can be used to identify
// the selected menu option.
func (mi *MenuItem) SetId(id string) *MenuItem {
	mi.id = id
	return mi
}

// Id returns this menu item current id
func (mi *MenuItem) Id() string {
	return mi.id
}

// IdPath returns a slice with the path of menu items ids to this menu item
func (mi *MenuItem) IdPath() []string {
	// Builds lists of menu items ids
	path := []string{mi.id}
	menu := mi.menu
	for menu.mitem != nil {
		path = append(path, menu.mitem.id)
		menu = menu.mitem.menu
	}
	// Reverse and returns id list
	res := make([]string, 0, len(path))
	for i := len(path) - 1; i >= 0; i-- {
		res = append(res, path[i])
	}
	return res
}

// onCursor processes subscribed cursor events over the menu item
func (mi *MenuItem) onCursor(evname string, _ interface{}) {
	switch evname {
	case OnCursorEnter:
		mi.menu.setSelectedItem(mi)
	}
}

// onMouse processes subscribed mouse events over the menu item
func (mi *MenuItem) onMouse(evname string, _ interface{}) {
	switch evname {
	case OnMouseDown:
		// MenuBar option
		if mi.menu.bar {
			mi.menu.autoOpen = !mi.menu.autoOpen
			if mi.submenu != nil && mi.submenu.Visible() {
				mi.submenu.SetVisible(false)
				GetManager().SetKeyFocus(mi.menu)
			} else {
				mi.update()
			}
		}
		if mi.submenu != nil {
			return
		}
		mi.activate()
	}
}

// activate activates this menu item dispatching OnClick events
func (mi *MenuItem) activate() {
	rm := mi.rootMenu()
	if rm.bar {
		rm.autoOpen = false
	}
	rm.setSelectedPos(-1)
	GetManager().SetKeyFocus(rm)
	mi.dispatchAll(OnClick, mi)
}

// rootMenu returns the root menu for this menu item
func (mi *MenuItem) rootMenu() *Menu {
	root := mi.menu
	for root.mitem != nil {
		root = root.mitem.menu
	}
	return root
}

// dispatchAll dispatch the specified event for this menu item
// and all its parents
func (mi *MenuItem) dispatchAll(evname string, ev interface{}) {
	mi.Dispatch(evname, ev)
	pmenu := mi.menu
	for {
		pmenu.Dispatch(evname, ev)
		if pmenu.mitem == nil {
			break
		}
		pmenu = pmenu.mitem.menu
	}
}

// update updates the menu item visual state
func (mi *MenuItem) update() {
	// Separator
	if mi.label == nil {
		return
	}
	// Disabled item
	if mi.disabled {
		return
	}
	// Selected item
	if mi.selected {
		if mi.submenu != nil && mi.menu.autoOpen {
			mi.menu.SetTopChild(mi)
			mi.submenu.SetVisible(true)
			if mi.menu != nil && mi.menu.bar {
				mi.submenu.SetPosition(0, mi.Height()-2)
			} else {
				mi.submenu.SetPosition(mi.Width()-2, 0)
			}
		}
		return
	}
	// If this menu item has a sub menu and the sub menu is not active,
	// hides the sub menu
	if mi.submenu != nil {
		mi.submenu.SetVisible(false)
	}
}

// recalc recalculates the positions of this menu item internal panels
func (mi *MenuItem) recalc(iconWidth, labelWidth, shortcutWidth float32) {
	// Separator
	if mi.label == nil {
		return
	}
	if mi.licon != nil {
		py := (mi.label.height - mi.licon.height) / 2
		mi.licon.SetPosition(0, py)
	}
	mi.label.SetPosition(iconWidth, 0)
	if mi.shortcut != nil {
		mi.shortcut.SetPosition(iconWidth+labelWidth, 0)
	}
	if mi.ricon != nil {
		mi.ricon.SetPosition(iconWidth+labelWidth+shortcutWidth, 0)
	}
}

// minHeight returns the minimum height of this menu item
func (mi *MenuItem) minHeight() float32 {
	mh := mi.MinHeight()
	if mi.label == nil {
		return mh + 1
	}
	mh += mi.label.height
	return mh
}

// minWidth returns the minimum width of this menu item
func (mi *MenuItem) minWidth() float32 {
	mw := mi.MinWidth()
	if mi.label == nil {
		return mw + 1
	}
	mw += mi.label.width
	return mw
}
