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

// Label component interface and implementation.

package gwu

// Label interface defines a component which wraps a text into a component.
// 
// Default style class: "gwu-Label"
type Label interface {
	// Label is a component.
	Comp

	// Label has text.
	HasText
}

// Label implementation
type labelImpl struct {
	compImpl    // Component implementation
	hasTextImpl // Has text implementation
}

// NewLabel creates a new Label.
func NewLabel(text string) Label {
	c := &labelImpl{newCompImpl(nil), newHasTextImpl(text)}
	c.Style().AddClass("gwu-Label")
	return c
}

func (c *labelImpl) Render(w writer) {
	w.Write(_STR_SPAN_OP)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	c.renderText(w)

	w.Write(_STR_SPAN_CL)
}
