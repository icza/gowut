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

// ListBox component interface and implementation.

package gwu

import (
	"net/http"
	"strconv"
	"strings"
)

// ListBox interface defines a component which allows selecting one or multiple values
// from a predefined list.
//
// Suggested event type to handle changes: ETypeChange
//
// Default style class: "gwu-ListBox"
type ListBox interface {
	// ListBox is a component
	Comp

	// ListBox can be enabled/disabled.
	HasEnabled

	// Multi tells if multiple selections are allowed.
	Multi() bool

	// SetMulti sets whether multiple selections are allowed.
	SetMulti(multi bool)

	// Rows returns the number of displayed rows.
	Rows() int

	// SetRows sets the number of displayed rows.
	// rows=1 will make this ListBox a dropdown list (if multi is false!).
	// Note that if rows is greater than 1, most browsers enforce a visual minimum size
	// (about 4 rows) even if rows is less than that.
	SetRows(rows int)

	// SelectedValue retruns the first selected value.
	// Empty string is returned if nothing is selected.
	SelectedValue() string

	// SelectedValues retruns all the selected values.
	SelectedValues() []string

	// Selected tells if the value at index i is selected.
	Selected(i int) bool

	// SelectedIdx returns the first selected index.
	// Returns -1 if nothing is selected.
	SelectedIdx() int

	// SelectedIndices returns a slice of the indices of the selected values.
	SelectedIndices() []int

	// SetSelected sets the selection state of the value at index i.
	SetSelected(i int, selected bool)

	// SetSelectedIndices sets the (only) selected values.
	// Only values will be selected that are contained in the specified indices slice.
	SetSelectedIndices(indices []int)

	// ClearSelected deselects all values.
	ClearSelected()
	
	// SetContent(content []string)
	SetContent(content []string)
}

// ListBox implementation.
type listBoxImpl struct {
	compImpl       // Component implementation
	hasEnabledImpl // Has enabled implementation

	values   []string // Values to choose from
	multi    bool     // Allow multiple selection
	selected []bool   // Array of selection state of the values
	rows     int      // Number of displayed rows
}

var (
	strSelidx = []byte("selIdxs(this)") // "selIdxs(this)"
)

// NewListBox creates a new ListBox.
func NewListBox(values []string) ListBox {
	c := &listBoxImpl{newCompImpl(strSelidx), newHasEnabledImpl(), values, false, make([]bool, len(values)), 1}
	c.AddSyncOnETypes(ETypeChange)
	c.Style().AddClass("gwu-ListBox")
	return c
}

func	(c *listBoxImpl)SetContent(newvalues [] string)	{
	c.values=newvalues
	c.selected=make([]bool,len(newvalues))
}

func (c *listBoxImpl) Multi() bool {
	return c.multi
}

func (c *listBoxImpl) SetMulti(multi bool) {
	c.multi = multi
}

func (c *listBoxImpl) Rows() int {
	return c.rows
}

func (c *listBoxImpl) SetRows(rows int) {
	c.rows = rows
}

func (c *listBoxImpl) SelectedValue() string {
	if i := c.SelectedIdx(); i >= 0 {
		return c.values[i]
	}

	return ""
}

func (c *listBoxImpl) SelectedValues() (sv []string) {
	for i, s := range c.selected {
		if s {
			sv = append(sv, c.values[i])
		}
	}
	return
}

func (c *listBoxImpl) Selected(i int) bool {
	return c.selected[i]
}

func (c *listBoxImpl) SelectedIdx() int {
	for i, s := range c.selected {
		if s {
			return i
		}
	}
	return -1
}

func (c *listBoxImpl) SelectedIndices() (si []int) {
	for i, s := range c.selected {
		if s {
			si = append(si, i)
		}
	}
	return
}

func (c *listBoxImpl) SetSelected(i int, selected bool) {
	c.selected[i] = selected
}

func (c *listBoxImpl) SetSelectedIndices(indices []int) {
	// First clear selected slice
	for i := range c.selected {
		c.selected[i] = false
	}

	// And now select that needs to be selected
	for _, idx := range indices {
		c.selected[idx] = true
	}
}

func (c *listBoxImpl) ClearSelected() {
	for i := range c.selected {
		c.selected[i] = false
	}
}

func (c *listBoxImpl) preprocessEvent(event Event, r *http.Request) {
	value := r.FormValue(paramCompValue)
	c.ClearSelected()
	if len(value) == 0 {
		return
	}

	// Set selected indices
	for _, sidx := range strings.Split(value, ",") {
		if idx, err := strconv.Atoi(sidx); err == nil {
			c.selected[idx] = true
		}
	}
}

var (
	strSelectOp    = []byte("<select")                      // "<select"
	strMultiple    = []byte(` multiple="multiple"`)         // ` multiple="multiple"`
	strOptionOpSel = []byte(`<option selected="selected">`) // `<option selected="selected">`
	strOptionOp    = []byte("<option>")                     // "<option>"
	strOptionCl    = []byte("</option>")                    // "</option>"
	strSelectCl    = []byte("</select>")                    // "</select>"
)

func (c *listBoxImpl) Render(w Writer) {
	w.Write(strSelectOp)
	if c.multi {
		w.Write(strMultiple)
	}
	w.WriteAttr("size", strconv.Itoa(c.rows))
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	for i, value := range c.values {
		if c.selected[i] {
			w.Write(strOptionOpSel)
		} else {
			w.Write(strOptionOp)
		}
		w.Writees(value)
		w.Write(strOptionCl)
	}

	w.Write(strSelectCl)
}
