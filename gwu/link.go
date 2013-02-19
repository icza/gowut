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

// Link component interface and implementation.

package gwu

// Link interface defines a clickable link pointing to a URL.
// Links are usually used with a text, although Link is a
// container, and allows to set a child component
// which if set will also be a part of the clickable link.
// 
// Default style class: "gwu-Link"
type Link interface {
	// Link is a Container.
	Container

	// Link has text.
	HasText

	// Link has URL string.
	HasUrl

	// Target returns the target of the link.
	Target() string

	// SetTarget sets the target of the link.
	// Tip: pass "_blank" if you want the URL to open in a new window
	// (this is the default).
	SetTarget(target string)

	// Comp returns the optional child component, if set.
	Comp() Comp

	// SetComp sets the only child component
	// (which can be a Container of course).
	SetComp(c Comp)
}

// Link implementation.
type linkImpl struct {
	compImpl    // Component implementation
	hasTextImpl // Has text implementation
	hasUrlImpl  // Has text implementation

	comp Comp // Optional child component
}

// NewLink creates a new Link.
// By default links open in a new window (tab)
// because their target is set to "_blank".
func NewLink(text, url string) Link {
	c := &linkImpl{newCompImpl(nil), newHasTextImpl(text), newHasUrlImpl(url), nil}
	c.SetTarget("_blank")
	c.Style().AddClass("gwu-Link")
	return c
}

func (c *linkImpl) Remove(c2 Comp) bool {
	if c.comp == nil || !c.comp.Equals(c2) {
		return false
	}

	c2.setParent(nil)
	c.comp = nil

	return true
}

func (c *linkImpl) ById(id ID) Comp {
	if c.id == id {
		return c
	}

	if c.comp != nil {
		if c.comp.Id() == id {
			return c.comp
		}
		if c2, isContainer := c.comp.(Container); isContainer {
			if c3 := c2.ById(id); c3 != nil {
				return c3
			}
		}

	}

	return nil
}

func (c *linkImpl) Clear() {
	if c.comp != nil {
		c.comp.setParent(nil)
		c.comp = nil
	}
}

func (c *linkImpl) Target() string {
	return c.attrs["target"]
}

func (c *linkImpl) SetTarget(target string) {
	if len(target) == 0 {
		delete(c.attrs, "target")
	} else {
		c.attrs["target"] = target
	}
}

func (c *linkImpl) Comp() Comp {
	return c.comp
}

func (c *linkImpl) SetComp(c2 Comp) {
	c.comp = c2
}

var (
	_STR_A_OP = []byte("<a")   // "<a"
	_STR_A_CL = []byte("</a>") // "</a>"
)

func (c *linkImpl) Render(w writer) {
	w.Write(_STR_A_OP)
	c.renderUrl("href", w)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	c.renderText(w)

	if c.comp != nil {
		c.comp.Render(w)
	}

	w.Write(_STR_A_CL)
}
