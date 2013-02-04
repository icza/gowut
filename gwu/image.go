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

// Image component interface and implementation.

package gwu

// Image interface defines an image.
// 
// Default style class: "gwu-Image"
type Image interface {
	// Image is a component.
	Comp

	// Image has text which is its description (alternate text).
	HasText

	// Image has URL string.
	HasUrl
}

// Image implementation
type imageImpl struct {
	compImpl    // Component implementation
	hasTextImpl // Has text implementation
	hasUrlImpl  // Has text implementation
}

// NewImage creates a new Image component.
func NewImage(text, url string) Image {
	c := &imageImpl{newCompImpl(""), newHasTextImpl(text), newHasUrlImpl(url)}
	c.Style().AddClass("gwu-Image")
	return c
}

func (c *imageImpl) Render(w writer) {
	w.Writes("<img")
	c.renderUrl("src", w)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Writes(" alt=\"")
	c.renderText(w)
	w.Writes("\">")
}
