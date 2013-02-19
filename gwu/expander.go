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

// Expander component interface and implementation.

package gwu

// Expander interface defines a component which can show and hide
// another component when clicked on the header.
// 
// You can register ETYPE_STATE_CHANGE event handlers which will be called when the user
// expands or collapses the expander by clicking on the header. The event source will be
// the expander. The event will have a parent event whose source will be the clicked
// header component and will contain the mouse coordinates.
// 
// Default style classes: "gwu-Expander", "gwu-Expander-Header",
// "gwuimg-collapsed", "gwu-Expander-Header-Expanded", "gwuimg-expanded",
// "gwu-Expander-Content"
type Expander interface {
	// Expander is a TableView.
	TableView

	// Header returns the header component of the expander.
	Header() Comp

	// SetHeader sets the header component of the expander.
	SetHeader(h Comp)

	// Content returns the content component of the expander.
	Content() Comp

	// SetContent sets the content component of the expander.
	SetContent(c Comp)

	// Expanded returns whether the expander is expanded.
	Expanded() bool

	// SetExpanded sets whether the expander is expanded.
	SetExpanded(expanded bool)

	// HeaderFmt returns the cell formatter of the header.
	HeaderFmt() CellFmt

	// ContentFmt returns the cell formatter of the content.
	ContentFmt() CellFmt
}

// Expander implementation.
type expanderImpl struct {
	tableViewImpl // TableView implementation

	header   Comp // Header component
	content  Comp // Content component
	expanded bool // Tells whether the expander is expanded

	headerFmt  *cellFmtImpl // Header cell formatter
	contentFmt *cellFmtImpl // Content cell formatter
}

// NewExpander creates a new Expander.
// By default expanders are collapsed.
func NewExpander() Expander {
	c := &expanderImpl{tableViewImpl: newTableViewImpl(), expanded: true, headerFmt: newCellFmtImpl(), contentFmt: newCellFmtImpl()}
	c.headerFmt.SetAlign(HA_LEFT, VA_MIDDLE)
	c.contentFmt.SetAlign(HA_LEFT, VA_TOP)
	c.Style().AddClass("gwu-Expander")
	// Init styles by changing expanded state, to the default value.
	c.SetExpanded(false)
	return c
}

func (c *expanderImpl) Remove(c2 Comp) bool {
	if c.content.Equals(c2) {
		c2.setParent(nil)
		c.content = nil
		return true
	}

	if c.header.Equals(c2) {
		c2.setParent(nil)
		c.header = nil
		return true
	}

	return false
}

func (c *expanderImpl) ById(id ID) Comp {
	if c.id == id {
		return c
	}

	if c.header != nil {
		if c.header.Id() == id {
			return c.header
		}
		if c2, isContainer := c.header.(Container); isContainer {
			if c3 := c2.ById(id); c3 != nil {
				return c3
			}
		}
	}

	if c.content != nil {
		if c.content.Id() == id {
			return c.content
		}
		if c2, isContainer := c.content.(Container); isContainer {
			if c3 := c2.ById(id); c3 != nil {
				return c3
			}
		}
	}

	return nil
}

func (c *expanderImpl) Clear() {
	if c.header != nil {
		c.header.setParent(nil)
		c.header = nil
	}
	if c.content != nil {
		c.content.setParent(nil)
		c.content = nil
	}
}

func (c *expanderImpl) Header() Comp {
	return c.header
}

func (c *expanderImpl) SetHeader(header Comp) {
	header.makeOrphan()
	c.header = header
	header.setParent(c)

	// TODO would be nice to remove this internal handler func when the header is removed!
	header.AddEHandlerFunc(func(e Event) {
		c.SetExpanded(!c.expanded)
		e.MarkDirty(c)
		if c.handlers[ETYPE_STATE_CHANGE] != nil {
			c.dispatchEvent(e.forkEvent(ETYPE_STATE_CHANGE, c))
		}
	}, ETYPE_CLICK)
}

func (c *expanderImpl) Content() Comp {
	return c.header
}

func (c *expanderImpl) SetContent(content Comp) {
	content.makeOrphan()
	c.content = content
	content.setParent(c)

	c.contentFmt.Style().AddClass("gwu-Expander-Content").SetFullSize()
}

func (c *expanderImpl) Expanded() bool {
	return c.expanded
}

func (c *expanderImpl) SetExpanded(expanded bool) {
	if c.expanded == expanded {
		return
	}

	style := c.headerFmt.Style()
	if c.expanded {
		style.RemoveClass("gwu-Expander-Header-Expanded")
		style.RemoveClass("gwuimg-expanded")
		style.AddClass("gwu-Expander-Header")
		style.AddClass("gwuimg-collapsed")
	} else {
		style.RemoveClass("gwu-Expander-Header")
		style.RemoveClass("gwuimg-collapsed")
		style.AddClass("gwu-Expander-Header-Expanded")
		style.AddClass("gwuimg-expanded")
	}

	c.expanded = expanded
}

func (c *expanderImpl) HeaderFmt() CellFmt {
	return c.headerFmt
}

func (c *expanderImpl) ContentFmt() CellFmt {
	return c.contentFmt
}

func (c *expanderImpl) Render(w writer) {
	w.Write(_STR_TABLE_OP)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	if c.header != nil {
		c.renderTr(w)
		c.headerFmt.render(_STR_TD_OP, w)
		c.header.Render(w)
	}

	if c.expanded && c.content != nil {
		c.renderTr(w)
		c.contentFmt.render(_STR_TD_OP, w)
		c.content.Render(w)
	}

	w.Write(_STR_TABLE_CL)
}
