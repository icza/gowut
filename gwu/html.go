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

// HTML interface defines a component which wraps an HTML text into a component.
//
// Default style class: "gwu-HTML"
type HTML interface {
	// HTML is a component.
	Comp

	// HTML returns the HTML text.
	HTML() string

	// SetHTML sets the HTML text.
	SetHTML(html string)
}

// HTML implementation
type htmlImpl struct {
	compImpl // Component implementation

	html string // HTML text
}

// NewHTML creates a new HTML.
func NewHTML(html string) HTML {
	c := &htmlImpl{newCompImpl(nil), html}
	c.Style().AddClass("gwu-Html")
	return c
}

func (c *htmlImpl) HTML() string {
	return c.html
}

func (c *htmlImpl) SetHTML(html string) {
	c.html = html
}

func (c *htmlImpl) Render(w Writer) {
	w.Write(strSpanOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	w.Writes(c.html)

	w.Write(strSpanCl)
}
