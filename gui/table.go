// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"github.com/derekmu/g3n/core"
	"sort"
	"strconv"
)

const (
	tableColMinWidth = 16
	tableErrInvRow   = "Invalid row index"
	tableErrInvCol   = "Invalid column id"
)

// Table implements a panel which can contains child panels
// organized in rows and columns.
type Table struct {
	Panel                   // Embedded panel
	header      tableHeader // table headers
	rows        []*tableRow // array of table rows
	firstRow    int         // index of the first visible row
	lastRow     int         // index of the last visible row
	statusPanel Panel       // optional bottom status panel
	statusLabel *Label      // status label
}

// TableColumn describes a table column
type TableColumn struct {
	Id         string          // Column id used to reference the column. Must be unique
	Header     string          // Column name shown in the table header
	Width      float32         // Initial column width in pixels
	Minwidth   float32         // Minimum width in pixels for this column
	Hidden     bool            // Hidden flag
	Align      Align           // Cell content alignment
	Format     string          // Format string for formatting the columns' cells
	FormatFunc TableFormatFunc // Format function (overrides Format string)
	Expand     float32         // Column width expansion factor (0 for no expansion)
}

// TableCell describes a table cell.
// It is used as a parameter for formatting function
type TableCell struct {
	Tab   *Table // Pointer to table
	Row   int    // Row index
	Col   string // Column id
	Value any    // Cell value
}

// TableFormatFunc is the type for formatting functions
type TableFormatFunc func(cell TableCell) string

// tableHeader is panel which contains the individual header panels for each column
type tableHeader struct {
	Panel                              // embedded panel
	cmap    map[string]*tableColHeader // maps column id with its panel/descriptor
	cols    []*tableColHeader          // array of individual column headers/descriptors
	lastPan Panel                      // last header panel not associated with a user column
}

// tableColHeader is panel for a column header
type tableColHeader struct {
	Panel                      // header panel
	label      *Label          // header label
	id         string          // column id
	width      float32         // initial column width
	minWidth   float32         // minimum width
	format     string          // column format string
	formatFunc TableFormatFunc // column format function
	align      Align           // column alignment
	expand     float32         // column expand factor
	order      int             // row columns order
	xl         float32         // left border coordinate in pixels
	xr         float32         // right border coordinate in pixels
}

// tableRow is panel which contains an entire table row of cells
type tableRow struct {
	Panel              // embedded panel
	cells []*tableCell // array of row cells
}

// tableCell is a panel which contains one cell (a label)
type tableCell struct {
	Panel       // embedded panel
	label Label // cell label
	value any   // cell current value
}

// NewTable creates and returns a pointer to a new Table with the
// specified width, height and columns
func NewTable(width, height float32, cols []TableColumn) (*Table, error) {
	t := new(Table)
	t.Panel.InitPanel(t, width, height)

	// Initialize table header
	t.header.InitPanel(&t.header, 0, 0)
	t.header.cmap = make(map[string]*tableColHeader)
	t.header.cols = make([]*tableColHeader, 0)

	// Create column header panels
	for ci := 0; ci < len(cols); ci++ {
		cdesc := cols[ci]
		// Column id must not be empty
		if cdesc.Id == "" {
			return nil, fmt.Errorf("column with empty id")
		}
		// Column id must be unique
		if t.header.cmap[cdesc.Id] != nil {
			return nil, fmt.Errorf("column with duplicate id")
		}
		// Creates a column header
		c := new(tableColHeader)
		c.InitPanel(c, 0, 0)
		c.label = NewLabel(cdesc.Header)
		c.Add(c.label)
		c.id = cdesc.Id
		c.minWidth = cdesc.Minwidth
		if c.minWidth < tableColMinWidth {
			c.minWidth = tableColMinWidth
		}
		c.width = cdesc.Width
		if c.width < c.minWidth {
			c.width = c.minWidth
		}
		c.align = cdesc.Align
		c.format = cdesc.Format
		c.formatFunc = cdesc.FormatFunc
		c.expand = cdesc.Expand
		// Sets default format and order
		if c.format == "" {
			c.format = "%v"
		}
		c.order = ci
		c.SetVisible(!cdesc.Hidden)
		t.header.cmap[c.id] = c
		// Sets column header width and height
		width := cdesc.Width
		if width < c.label.Width()+c.MinWidth() {
			width = c.label.Width() + c.MinWidth()
		}
		c.SetContentSize(width, c.label.Height())
		// Adds the column header to the header panel
		t.header.cols = append(t.header.cols, c)
		t.header.Panel.Add(c)
	}
	// Creates last header
	t.header.lastPan.InitPanel(&t.header, 0, 0)
	t.header.Panel.Add(&t.header.lastPan)

	// Add header panel to the table panel
	t.Panel.Add(&t.header)

	// Creates status panel
	t.statusPanel.InitPanel(&t.statusPanel, 0, 0)
	t.statusPanel.SetVisible(false)
	t.statusLabel = NewLabel("")
	t.statusPanel.Add(t.statusLabel)
	t.Panel.Add(&t.statusPanel)

	// Subscribe to events
	t.Panel.Subscribe(t.onGuiEvent)
	t.recalc()
	return t, nil
}

// ShowHeader shows or hides the table header
func (t *Table) ShowHeader(show bool) {
	if t.header.Visible() == show {
		return
	}
	t.header.SetVisible(show)
	t.recalc()
}

// ShowColumn sets the visibility of the column with the specified id
// If the column id does not exit the function panics.
func (t *Table) ShowColumn(col string, show bool) {
	c := t.header.cmap[col]
	if c == nil {
		panic(tableErrInvCol)
	}
	if c.Visible() == show {
		return
	}
	c.SetVisible(show)
	t.recalc()
}

// ShowAllColumns shows all the table columns
func (t *Table) ShowAllColumns() {
	recalc := false
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		if !c.Visible() {
			c.SetVisible(true)
			recalc = true
		}
	}
	if !recalc {
		return
	}
	t.recalc()
}

// RowCount returns the current number of rows in the table
func (t *Table) RowCount() int {
	return len(t.rows)
}

// SetRows clears all current rows of the table and
// sets new rows from the specifying parameter.
// Each row is a map keyed by the colum id.
// The map value currently can be a string or any number type
// If a row column is not found it is ignored
func (t *Table) SetRows(values []map[string]any) {
	// Add missing rows
	if len(values) > len(t.rows) {
		count := len(values) - len(t.rows)
		for row := 0; row < count; row++ {
			t.insertRow(len(t.rows), nil)
		}
		// Remove remaining rows
	} else if len(values) < len(t.rows) {
		for row := len(values); row < len(t.rows); row++ {
			t.removeRow(row)
		}
	}

	// Set rows values
	for row := 0; row < len(values); row++ {
		t.setRow(row, values[row])
	}
	t.firstRow = 0
	t.recalc()
}

// SetRow sets the value of all the cells of the specified row from
// the specified map indexed by column id.
func (t *Table) SetRow(row int, values map[string]any) {
	if row < 0 || row >= len(t.rows) {
		panic(tableErrInvRow)
	}
	t.setRow(row, values)
	t.recalc()
}

// SetCell sets the value of the cell specified by its row and column id
// The function panics if the passed row or column id is invalid
func (t *Table) SetCell(row int, colid string, value any) {
	if row < 0 || row >= len(t.rows) {
		panic(tableErrInvRow)
	}
	if t.header.cmap[colid] == nil {
		panic(tableErrInvCol)
	}
	t.setCell(row, colid, value)
	t.recalc()
}

// SetColFormat sets the formatting string (Printf) for the specified column
// Update must be called to update the table.
func (t *Table) SetColFormat(id, format string) {
	c := t.header.cmap[id]
	if c == nil {
		panic(tableErrInvCol)
	}
	c.format = format
}

// SetColOrder sets the exhibition order of the specified column.
// The previous column which has the specified order will have
// the original column order.
func (t *Table) SetColOrder(colid string, order int) {
	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	// Checks exhibition order
	if order < 0 || order > len(t.header.cols) {
		panic(tableErrInvRow)
	}
	// Find the exhibition order for the specified column
	for ci := 0; ci < len(t.header.cols); ci++ {
		if t.header.cols[ci] == c {
			// If the order of the specified column is the same, nothing to do
			if ci == order {
				return
			}
			// Swap column orders
			prev := t.header.cols[order]
			t.header.cols[order] = c
			t.header.cols[ci] = prev
			break
		}
	}

	// Recalculates the header and all rows
	t.recalc()
}

// SetColWidth sets the specified column width and may
// change the widths of the columns to the right
func (t *Table) SetColWidth(colid string, width float32) {
	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	t.setColWidth(c, width)
}

// SetColExpand sets the column expand factor.
// When the table width is increased the columns widths are
// increased proportionally to their expand factor.
// A column with expand factor = 0 is not increased.
func (t *Table) SetColExpand(colid string, expand float32) {
	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	if expand < 0 {
		c.expand = 0
	} else {
		c.expand = expand
	}
	t.recalc()
}

// AddRow adds a new row at the end of the table with the specified values
func (t *Table) AddRow(values map[string]any) {
	t.InsertRow(len(t.rows), values)
}

// InsertRow inserts the specified values in a new row at the specified index
func (t *Table) InsertRow(row int, values map[string]any) {
	// Checks row index
	if row < 0 || row > len(t.rows) {
		panic(tableErrInvRow)
	}
	t.insertRow(row, values)
	t.recalc()
}

// RemoveRow removes from the specified row from the table
func (t *Table) RemoveRow(row int) {
	// Checks row index
	if row < 0 || row >= len(t.rows) {
		panic(tableErrInvRow)
	}
	t.removeRow(row)
	maxFirst := t.calcMaxFirst()
	if t.firstRow > maxFirst {
		t.firstRow = maxFirst
	}
	t.recalc()
}

// Clear removes all rows from the table
func (t *Table) Clear() {
	for ri := 0; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		t.Panel.Remove(trow)
		trow.DisposeChildren(true)
		trow.Dispose()
	}
	t.rows = nil
	t.firstRow = 0
	t.recalc()
}

// ShowStatus sets the visibility of the status lines at the bottom of the table
func (t *Table) ShowStatus(show bool) {
	if t.statusPanel.Visible() == show {
		return
	}
	t.statusPanel.SetVisible(show)
	t.recalcStatus()
	t.recalc()
}

// SetStatusText sets the text of status line at the bottom of the table
// It does not change its current visibility
func (t *Table) SetStatusText(text string) {
	t.statusLabel.SetText(text)
}

// Rows returns a slice of maps with the contents of the table rows
// specified by the rows first and last index.
// To get all the table rows, use Rows(0, -1)
func (t *Table) Rows(fi, li int) []map[string]any {
	if fi < 0 || fi >= len(t.header.cols) {
		panic(tableErrInvRow)
	}
	if li < 0 {
		li = len(t.rows) - 1
	} else if li < 0 || li >= len(t.rows) {
		panic(tableErrInvRow)
	}
	if li < fi {
		panic("Last index less than first index")
	}
	res := make([]map[string]any, li-li+1)
	for ri := fi; ri <= li; ri++ {
		trow := t.rows[ri]
		rmap := make(map[string]any)
		for ci := 0; ci < len(t.header.cols); ci++ {
			c := t.header.cols[ci]
			rmap[c.id] = trow.cells[c.order].value
		}
		res = append(res, rmap)
	}
	return res
}

// Row returns a map with the current contents of the specified row index
func (t *Table) Row(ri int) map[string]any {
	if ri < 0 || ri > len(t.header.cols) {
		panic(tableErrInvRow)
	}
	res := make(map[string]any)
	trow := t.rows[ri]
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		res[c.id] = trow.cells[c.order].value
	}
	return res
}

// Cell returns the current content of the specified cell
func (t *Table) Cell(col string, ri int) any {
	c := t.header.cmap[col]
	if c == nil {
		panic(tableErrInvCol)
	}
	if ri < 0 || ri >= len(t.rows) {
		panic(tableErrInvRow)
	}
	trow := t.rows[ri]
	return trow.cells[c.order].value
}

// SortColumn sorts the specified column interpreting its values as strings or numbers
// and sorting in ascending or descending order.
// This sorting is independent of the sort configuration of column set when the table was created
func (t *Table) SortColumn(col string, asString bool, asc bool) {
	c := t.header.cmap[col]
	if c == nil {
		panic(tableErrInvCol)
	}
	if len(t.rows) < 2 {
		return
	}
	if asString {
		ts := tableSortString{rows: t.rows, col: c.order, asc: asc, format: c.format}
		sort.Sort(ts)
	} else {
		ts := tableSortNumber{rows: t.rows, col: c.order, asc: asc}
		sort.Sort(ts)
	}
	t.recalc()
}

// setRow sets the value of all the cells of the specified row from
// the specified map indexed by column id.
func (t *Table) setRow(row int, values map[string]any) {
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		cv, ok := values[c.id]
		if !ok {
			continue
		}
		t.setCell(row, c.id, cv)
	}
}

// setCell sets the value of the cell specified by its row and column id
func (t *Table) setCell(row int, colid string, value any) {
	c := t.header.cmap[colid]
	if c == nil {
		return
	}
	cell := t.rows[row].cells[c.order]
	cell.label.SetText(fmt.Sprintf(c.format, value))
	cell.value = value
}

// insertRow is the internal version of InsertRow which does not call recalc()
func (t *Table) insertRow(row int, values map[string]any) {
	// Creates tableRow panel
	trow := new(tableRow)
	trow.InitPanel(trow, 0, 0)
	trow.cells = make([]*tableCell, 0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// Creates tableRow cell panel
		cell := new(tableCell)
		cell.InitPanel(cell, 0, 0)
		cell.label.InitLabel("", StyleDefault().Font)
		cell.Add(&cell.label)
		trow.cells = append(trow.cells, cell)
		trow.Panel.Add(cell)
	}
	t.Panel.Add(trow)

	// Inserts tableRow in the table rows at the specified index
	t.rows = append(t.rows, nil)
	copy(t.rows[row+1:], t.rows[row:])
	t.rows[row] = trow

	// Sets the new row values from the specified map
	if values != nil {
		t.SetRow(row, values)
	}
	t.recalcRow(row)
}

// removeRow removes from the table the row specified its index
func (t *Table) removeRow(row int) {
	// Get row to be removed
	trow := t.rows[row]

	// Remove row from table children
	t.Panel.Remove(trow)

	// Remove row from rows array
	copy(t.rows[row:], t.rows[row+1:])
	t.rows[len(t.rows)-1] = nil
	t.rows = t.rows[:len(t.rows)-1]

	// Dispose row resources
	trow.DisposeChildren(true)
	trow.Dispose()
}

// onResize receives subscribed resize events for this table
func (t *Table) onResize(_ core.GuiResizeEvent) {
	t.recalc()
	t.recalcStatus()
}

// setColWidth sets the width of the specified column
func (t *Table) setColWidth(c *tableColHeader, width float32) {
	// Sets the column width
	if width < c.minWidth {
		width = c.minWidth
	}
	if c.Width() == width {
		return
	}
	dw := width - c.Width()
	c.SetWidth(width)

	// Find the column with expand != 0
	hasExpand := false
	ci := -1
	for i := 0; i < len(t.header.cols); i++ {
		current := t.header.cols[i]
		if current == c {
			ci = i
		}
		if current.expand > 0 && current.Visible() {
			hasExpand = true
		}
	}
	if ci >= len(t.header.cols) {
		panic("Internal: column not found")
	}
	// If no column is expandable, nothing more
	if !hasExpand {
		t.recalc()
		return
	}
	// Calculates the width of the columns at the right
	rwidth := float32(0)
	for i := ci + 1; i < len(t.header.cols); i++ {
		c := t.header.cols[i]
		if !c.Visible() {
			continue
		}
		rwidth += c.Width()
	}
	// Distributes the delta to the columns at the right
	for i := ci + 1; i < len(t.header.cols); i++ {
		c := t.header.cols[i]
		if !c.Visible() {
			continue
		}
		cdelta := -dw * (c.Width() / rwidth)
		newWidth := c.Width() + cdelta
		if newWidth < c.minWidth {
			newWidth = c.minWidth
		}
		c.SetWidth(newWidth)
	}
	t.recalc()
}

// recalcHeader recalculates and sets the position and size of the header panels
func (t *Table) recalcHeader() {
	// Calculates total width, height, expansion and available width space
	hwidth := float32(0)
	height := float32(0)
	wspace := float32(0)
	totalExpand := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is invisible, ignore
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		if c.Height() > height {
			height = c.Height()
		}
		if c.expand > 0 {
			totalExpand += c.expand
		}
		hwidth += c.Width()
	}
	// Total table width
	twidth := t.ContentWidth()
	// Available space for columns: may be negative
	wspace = twidth - hwidth

	// If no expandable column, keeps the columns widths
	if totalExpand == 0 {
	} else if wspace >= 0 {
		for ci := 0; ci < len(t.header.cols); ci++ {
			// If column is invisible, ignore
			c := t.header.cols[ci]
			if !c.Visible() {
				continue
			}
			// There is space available and if column is expandable,
			// expands it proportionally to the other expandable columns
			factor := c.expand / totalExpand
			w := factor * wspace
			c.SetWidth(c.Width() + w)
		}
	} else {
		acols := make([]*tableColHeader, 0)
		awidth := float32(0)
		widthAvail := twidth
		// Sets the widths of the columns
		for ci := 0; ci < len(t.header.cols); ci++ {
			// If column is invisible, ignore
			c := t.header.cols[ci]
			if !c.Visible() {
				continue
			}
			// The table was reduced so shrinks this column proportionally to its current width
			factor := c.Width() / hwidth
			newWidth := factor * twidth
			if newWidth < c.minWidth {
				newWidth = c.minWidth
				c.SetWidth(newWidth)
				widthAvail -= c.minWidth
			} else {
				acols = append(acols, c)
				awidth += c.Width()
			}
		}
		for ci := 0; ci < len(acols); ci++ {
			c := acols[ci]
			factor := c.Width() / awidth
			newWidth := factor * widthAvail
			c.SetWidth(newWidth)
		}
	}

	// Sets the header panel and its internal panels positions
	posx := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is invisible, ignore
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		// Sets the column header panel position
		c.SetPosition(posx, 0)
		c.SetVisible(true)
		c.xl = posx
		posx += c.Width()
		c.xr = posx
	}

	// Last header
	w := t.ContentWidth() - posx
	if w > 0 {
		t.header.lastPan.SetVisible(true)
		t.header.lastPan.SetSize(w, height)
		t.header.lastPan.SetPosition(posx, 0)
	} else {
		t.header.lastPan.SetVisible(false)
	}

	// Header container
	t.header.SetWidth(t.ContentWidth())
	t.header.SetContentHeight(height)
}

// recalcStatus recalculates and sets the position and size of the status panel and its label
func (t *Table) recalcStatus() {
	if !t.statusPanel.Visible() {
		return
	}
	t.statusPanel.SetContentHeight(t.statusLabel.Height())
	py := t.ContentHeight() - t.statusPanel.Height()
	t.statusPanel.SetPosition(0, py)
	t.statusPanel.SetWidth(t.ContentWidth())
}

// recalc calculates the visibility, positions and sizes of all row cells.
// should be called in the following situations:
// - the table is resized
// - row is added, inserted or removed
// - column alignment and expansion changed
// - column visibility is changed
// - horizontal or vertical scroll position changed
func (t *Table) recalc() {
	// Get available row height for rows
	starty, theight := t.rowsHeight()

	// Determines if it is necessary to show the scrollbar or not.
	py := starty
	for ri := 0; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		py += trow.height
	}
	// Recalculates the header
	t.recalcHeader()

	// Sets the position and sizes of all cells of the visible rows
	py = starty
	for ri := 0; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		// If row is before first row or its y coordinate is greater the table height,
		// sets it invisible
		if ri < t.firstRow || py > starty+theight {
			trow.SetVisible(false)
			continue
		}
		t.recalcRow(ri)
		// Set row y position and visible
		trow.SetPosition(0, py)
		trow.SetVisible(true)
		// Set the last completely visible row index
		if py+trow.Height() <= starty+theight {
			t.lastRow = ri
		}
		py += trow.height
	}
	// Status panel must be on top of all the row panels
	t.SetTopChild(&t.statusPanel)
}

// recalcRow recalculates the positions and sizes of all cells of the specified row
// Should be called when the row is created and column visibility or order is changed.
func (t *Table) recalcRow(ri int) {
	trow := t.rows[ri]
	// Calculates and sets row height
	maxheight := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is hidden, ignore
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		cell := trow.cells[c.order]
		cellHeight := cell.MinHeight() + cell.label.Height()
		if cellHeight > maxheight {
			maxheight = cellHeight
		}
	}
	trow.SetContentHeight(maxheight)

	// Sets row cells sizes and positions and sets row width
	px := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is hidden, ignore
		col := t.header.cols[ci]
		cell := trow.cells[col.order]
		if !col.Visible() {
			cell.SetVisible(false)
			continue
		}
		// Sets cell position and size
		cell.SetPosition(px, 0)
		cell.SetVisible(true)
		cell.SetSize(col.Width(), trow.ContentHeight())
		// Checks for format function
		if col.formatFunc != nil {
			text := col.formatFunc(TableCell{t, ri, col.id, cell.value})
			cell.label.SetText(text)
		}
		// Sets the cell label alignment inside the cell
		if col.align != AlignNone {
			lx, ly := col.align.CalculatePosition(cell.ContentWidth(), cell.ContentHeight(), cell.label.Width(), cell.label.Height())
			cell.label.SetPosition(lx, ly)
		}
		px += col.Width()
	}
	trow.SetContentWidth(px)
}

// rowsHeight returns the available start y coordinate and height in the table for rows,
// considering the visibility of the header and status panels.
func (t *Table) rowsHeight() (float32, float32) {
	start := float32(0)
	height := t.ContentHeight()
	if t.header.Visible() {
		height -= t.header.Height()
		start += t.header.Height()
	}
	if t.statusPanel.Visible() {
		height -= t.statusPanel.Height()
	}
	if height < 0 {
		return 0, 0
	}
	return start, height
}

// calcMaxFirst calculates the maximum index of the first visible row
// such as the remaining rows fits completely inside the table
// It is used when scrolling the table vertically
func (t *Table) calcMaxFirst() int {
	_, total := t.rowsHeight()
	ri := len(t.rows) - 1
	if ri < 0 {
		return 0
	}
	height := float32(0)
	for {
		trow := t.rows[ri]
		height += trow.height
		if height > total {
			break
		}
		ri--
		if ri < 0 {
			break
		}
	}
	return ri + 1
}

func (t *Table) onGuiEvent(ev core.GuiEvent) bool {
	switch ev := ev.(type) {
	case core.GuiResizeEvent:
		t.onResize(ev)
	default:
		return false
	}
	return true
}

// tableSortString is an internal type implementing the sort.Interface
// and is used to sort a table column interpreting its values as strings
type tableSortString struct {
	rows   []*tableRow
	col    int
	asc    bool
	format string
}

func (ts tableSortString) Len() int      { return len(ts.rows) }
func (ts tableSortString) Swap(i, j int) { ts.rows[i], ts.rows[j] = ts.rows[j], ts.rows[i] }
func (ts tableSortString) Less(i, j int) bool {
	vi := ts.rows[i].cells[ts.col].value
	vj := ts.rows[j].cells[ts.col].value
	si := fmt.Sprintf(ts.format, vi)
	sj := fmt.Sprintf(ts.format, vj)
	if ts.asc {
		return si < sj
	}
	return sj < si
}

// tableSortNumber is an internal type implementing the sort.Interface
// and is used to sort a table column interpreting its values as numbers
type tableSortNumber struct {
	rows []*tableRow
	col  int
	asc  bool
}

func (ts tableSortNumber) Len() int      { return len(ts.rows) }
func (ts tableSortNumber) Swap(i, j int) { ts.rows[i], ts.rows[j] = ts.rows[j], ts.rows[i] }
func (ts tableSortNumber) Less(i, j int) bool {
	vi := ts.rows[i].cells[ts.col].value
	vj := ts.rows[j].cells[ts.col].value
	ni := cv2f64(vi)
	nj := cv2f64(vj)
	if ts.asc {
		return ni < nj
	}
	return nj < ni
}

// Try to convert an interface value to a float64 number
func cv2f64(v any) float64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case uint:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case int:
		return float64(n)
	case string:
		sv, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return sv
		}
		return 0
	default:
		return 0
	}
}
