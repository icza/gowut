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

// Defines the Html component.

package gwu

// Html interface defines a component which wraps an HTML text into a component.
// 
// Default style class: "gwu-Html"
type Html interface {
	// Html is a component.
	Comp

	// Html returns the HTML text.
	Html() string

	// SetHtml sets the HTML text.
	SetHtml(html string)
}

// Html implementation
type htmlImpl struct {
	compImpl // Component implementation

	html string // HTML text
}

// NewHtml creates a new Html.
func NewHtml(html string) Html {
	c := &htmlImpl{newCompImpl(nil), html}
	c.Style().AddClass("gwu-Html")
	return c
}

func (c *htmlImpl) Html() string {
	return c.html
}

func (c *htmlImpl) SetHtml(html string) {
	c.html = html
}

func (c *htmlImpl) Render(w writer) {
	w.Write(_STR_SPAN_OP)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	w.Writes(c.html)

	w.Write(_STR_SPAN_CL)
}
