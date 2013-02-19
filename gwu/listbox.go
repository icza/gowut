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
// Suggested event type to handle changes: ETYPE_CHANGE
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
	_STR_SELIDXS = []byte("selIdxs(this)") // "selIdxs(this)"
)

// NewListBox creates a new ListBox.
func NewListBox(values []string) ListBox {
	c := &listBoxImpl{newCompImpl(_STR_SELIDXS), newHasEnabledImpl(), values, false, make([]bool, len(values)), 1}
	c.AddSyncOnETypes(ETYPE_CHANGE)
	c.Style().AddClass("gwu-ListBox")
	return c
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
	for i, _ := range c.selected {
		c.selected[i] = false
	}

	// And now select that needs to be selected
	for _, idx := range indices {
		c.selected[idx] = true
	}
}

func (c *listBoxImpl) ClearSelected() {
	for i, _ := range c.selected {
		c.selected[i] = false
	}
}

func (c *listBoxImpl) preprocessEvent(event Event, r *http.Request) {
	value := r.FormValue(_PARAM_COMP_VALUE)
	if len(value) == 0 {
		return
	}

	// Set selected indices
	c.ClearSelected()
	for _, sidx := range strings.Split(value, ",") {
		if idx, err := strconv.Atoi(sidx); err == nil {
			c.selected[idx] = true
		}
	}
}

var (
	_STR_SELECT_OP     = []byte("<select")                        // "<select"
	_STR_MULTIPLE      = []byte(" multiple=\"multiple\"")         // " multiple=\"multiple\""
	_STR_OPTION_OP_SEL = []byte("<option selected=\"selected\">") // "<option selected=\"selected\">"
	_STR_OPTION_OP     = []byte("<option>")                       // "<option>"
	_STR_OPTION_CL     = []byte("</option>")                      // "</option>"
	_STR_SELECT_CL     = []byte("</select>")                      // "</select>"
)

func (c *listBoxImpl) Render(w writer) {
	w.Write(_STR_SELECT_OP)
	if c.multi {
		w.Write(_STR_MULTIPLE)
	}
	w.WriteAttr("size", strconv.Itoa(c.rows))
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	for i, value := range c.values {
		if c.selected[i] {
			w.Write(_STR_OPTION_OP_SEL)
		} else {
			w.Write(_STR_OPTION_OP)
		}
		w.Writees(value)
		w.Write(_STR_OPTION_CL)
	}

	w.Write(_STR_SELECT_CL)
}
