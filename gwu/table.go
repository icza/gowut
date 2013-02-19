// Copyright (C) 2013 Andras Belicza. All rights reserved.
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Table component interface and implementation.

package gwu

// Table interface defines a container which lays out its children
// using a configurable, flexible table.
// The size of the table grows dynamically, on demand. However,
// if table size is known or can be guessed before/during building it,
// it is recommended to call EnsureSize to minimize reallocations
// in the background.
// 
// Default style class: "gwu-Table"
type Table interface {
	// Table is a TableView.
	TableView

	// EnsureSize ensures that the table will have at least the specified
	// rows, and at least the specified columns in rows whose index is < rows.
	EnsureSize(rows, cols int)

	// EnsureCols ensures that the table will have at least the specified
	// cols at the specified row.
	// This implicitly includes that the table must have at least (row+1) rows.
	// If the table have less than (row+1) rows, empty rows will be added first.
	EnsureCols(row, cols int)

	// CompsCount returns the number of components added to the table.
	CompsCount() int

	// CompAt returns the component at the specified row and column.
	// Returns nil if row or column are invalid.
	CompAt(row, col int) Comp

	// CompIdx returns the row and column of the specified component in the table.
	// (-1, -1) is returned if the component is not added to the table.
	CompIdx(c Comp) (row, col int)

	// RowFmt returns the row formatter of the specified table row.
	// If the table does not have a row specified by row, nil is returned.
	RowFmt(row int) CellFmt

	// CellFmt returns the cell formatter of the specified table cell.
	// If the table does not have a cell specified by row and col,
	// nil is returned.
	CellFmt(row, col int) CellFmt

	// Add adds a component to the table.
	// Return value indicates if the component was added successfully.
	// Returns false if row or col is negative.
	Add(c Comp, row, col int) bool

	// RowSpan returns the row span of the specified table cell.
	// -1 is returned if the table does not have a cell specified by row and col.
	RowSpan(row, col int) int

	// SetRowSpan sets the row span of the specified table cell.
	// If the table does not have a cell specified by row and col,
	// this is a no-op.
	SetRowSpan(row, col, rowSpan int)

	// ColSpan returns the col span of the specified table cell.
	// -1 is returned if the table does not have a cell specified by row and col.
	ColSpan(row, col int) int

	// SetColSpan sets the col span of the specified table cell.
	// If the table does not have a cell specified by row and col,
	// this is a no-op.
	SetColSpan(row, col, colSpan int)
}

// cellIdx type specifies a cell by its row and col indices.
type cellIdx struct {
	row, col int // Row and col indices of the cell.
}

// Table implementation.
type tableImpl struct {
	tableViewImpl // TableView implementation

	comps    [][]Comp                 // Components added to the table. Structure: comps[rowIdx][colIdx]
	rowFmts  map[int]*cellFmtImpl     // Lazily initialized row formatters of the rows
	cellFmts map[cellIdx]*cellFmtImpl // Lazily initialized cell formatters of the cells
}

// NewTable creates a new Table.
// Default horizontal alignment is HA_DEFAULT,
// default vertical alignment is VA_DEFAULT.
func NewTable() Table {
	c := &tableImpl{tableViewImpl: newTableViewImpl()}
	c.Style().AddClass("gwu-Table")
	c.SetCellSpacing(0)
	c.SetCellPadding(0)
	return c
}

func (c *tableImpl) Remove(c2 Comp) bool {
	row, col := c.CompIdx(c2)
	if row < 0 {
		return false
	}

	c2.setParent(nil)
	c.comps[row][col] = nil

	return true
}

func (c *tableImpl) ById(id ID) Comp {
	if c.id == id {
		return c
	}

	for _, rowComps := range c.comps {
		for _, c2 := range rowComps {
			if c2 == nil {
				continue
			}
			if c2.Id() == id {
				return c2
			}

			if c3, isContainer := c2.(Container); isContainer {
				if c4 := c3.ById(id); c4 != nil {
					return c4
				}
			}
		}
	}
	return nil
}

func (c *tableImpl) Clear() {
	// Clear row formatters
	if c.rowFmts != nil {
		c.rowFmts = nil
	}
	// Clear cell formatters
	if c.cellFmts != nil {
		c.cellFmts = nil
	}

	for _, rowComps := range c.comps {
		for _, c2 := range rowComps {
			c2.setParent(nil)
		}
	}
	c.comps = nil
}

func (c *tableImpl) EnsureSize(rows, cols int) {
	c.ensureRows(rows)

	// Ensure column count in each row
	for i := 0; i < rows; i++ {
		c.EnsureCols(i, cols)
	}
}

func (c *tableImpl) EnsureCols(row, cols int) {
	c.ensureRows(row + 1)

	rowComps := c.comps[row]
	if cols > len(rowComps) {
		c.comps[row] = append(rowComps, make([]Comp, cols-len(rowComps))...)
	}
}

// EnsureRows ensures that the table will have at least the specified rows.
func (c *tableImpl) ensureRows(rows int) {
	if rows > len(c.comps) {
		c.comps = append(c.comps, make([][]Comp, rows-len(c.comps))...)
	}
}

func (c *tableImpl) CompsCount() (count int) {
	for _, rowComps := range c.comps {
		for _, c2 := range rowComps {
			if c2 != nil {
				count++
			}
		}
	}
	return
}

func (c *tableImpl) CompAt(row, col int) Comp {
	if row < 0 || col < 0 || row >= len(c.comps) {
		return nil
	}

	rowComps := c.comps[row]
	if col >= len(rowComps) {
		return nil
	}

	return rowComps[col]
}

func (c *tableImpl) CompIdx(c2 Comp) (int, int) {
	for row, rowComps := range c.comps {
		for col, c3 := range rowComps {
			if c3 == nil {
				continue
			}
			if c2.Equals(c3) {
				return row, col
			}
		}
	}

	return -1, -1
}

func (c *tableImpl) RowFmt(row int) CellFmt {
	if row < 0 || row >= len(c.comps) {
		return nil
	}

	if c.rowFmts == nil {
		c.rowFmts = make(map[int]*cellFmtImpl)
	}

	rf := c.rowFmts[row]
	if rf == nil {
		rf = newCellFmtImpl()
		c.rowFmts[row] = rf
	}

	return rf
}

func (c *tableImpl) CellFmt(row, col int) CellFmt {
	if row < 0 || col < 0 || row >= len(c.comps) || col >= len(c.comps[row]) {
		return nil
	}

	if c.cellFmts == nil {
		c.cellFmts = make(map[cellIdx]*cellFmtImpl)
	}

	ci := cellIdx{row, col}

	cf := c.cellFmts[ci]
	if cf == nil {
		cf = newCellFmtImpl()
		c.cellFmts[ci] = cf
	}

	return cf
}

func (c *tableImpl) Add(c2 Comp, row, col int) bool {
	c2.makeOrphan()

	// Quick check of row and col
	if row < 0 || col < 0 {
		return false
	}
	if row >= len(c.comps) || col >= len(c.comps[row]) {
		c.EnsureSize(row+1, col+1)
	}

	rowComps := c.comps[row]

	// Remove component if there is already one at the specified row and column:
	if rowComps[col] != nil {
		rowComps[col].setParent(nil)
	}

	rowComps[col] = c2
	c2.setParent(c)

	return true
}

func (c *tableImpl) RowSpan(row, col int) int {
	cf := c.CellFmt(row, col)
	if cf == nil {
		return -1
	}

	return cf.iAttr("rowspan")
}

func (c *tableImpl) SetRowSpan(row, col, rowSpan int) {
	cf := c.CellFmt(row, col)
	if cf == nil {
		return
	}

	if rowSpan < 2 {
		cf.setAttr("rowspan", "") // Delete attribute
	} else {
		cf.setIAttr("rowspan", rowSpan)
	}
}

func (c *tableImpl) ColSpan(row, col int) int {
	cf := c.CellFmt(row, col)
	if cf == nil {
		return -1
	}

	return cf.iAttr("colspan")
}

func (c *tableImpl) SetColSpan(row, col, colSpan int) {
	cf := c.CellFmt(row, col)
	if cf == nil {
		return
	}

	if colSpan < 2 {
		cf.setAttr("colspan", "") // Delete attribute
	} else {
		cf.setIAttr("colspan", colSpan)
	}
}

func (c *tableImpl) Render(w writer) {
	w.Write(_STR_TABLE_OP)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	// Create a reusable cell index
	ci := cellIdx{}

	for row, rowComps := range c.comps {
		c.renderRowTr(row, w)
		for col, c2 := range rowComps {
			ci.row, ci.col = row, col
			c.renderTd(ci, w)
			if c2 != nil {
				c2.Render(w)
			}
		}
	}

	w.Write(_STR_TABLE_CL)
}

// renderRowTr renders the formatted HTML TR tag for the specified row.
func (c *tableImpl) renderRowTr(row int, w writer) {
	var defha HAlign = c.halign // default halign of the table
	var defva VAlign = c.valign // default valign of the table

	if rf := c.rowFmts[row]; rf == nil {
		c.renderTr(w)
	} else {
		// If rf does not specify alignments, it means alignments must not be overriden,
		// default alignments of the table must be used!
		ha, va := rf.halign, rf.valign
		if ha == HA_DEFAULT {
			ha = defha
		}
		if va == VA_DEFAULT {
			va = defva
		}

		rf.renderWithAligns(_STR_TR_OP, ha, va, w)
	}
}

// renderTd renders the formatted HTML TD tag for the specified cell.
func (c *tableImpl) renderTd(ci cellIdx, w writer) {
	if cf := c.cellFmts[ci]; cf == nil {
		w.Write(_STR_TD)
	} else {
		cf.render(_STR_TD_OP, w)
	}
}
