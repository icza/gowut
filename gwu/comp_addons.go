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

// Defines optional, additional features components might have.
// These include features only some component has, so it cannot be
// defined in Comp, and not worth making an own component type for these...
// ...not to mention these can be combined arbitrary.

package gwu

import (
	"strconv"
)

// HasText interface defines a modifiable text property.
type HasText interface {
	// Text returns the text.
	Text() string

	// SetText sets the text.
	SetText(text string)
}

// newHasTextImpl creates a new hasTextImpl
func newHasTextImpl(text string) hasTextImpl {
	return hasTextImpl{text}
}

// HasText implementation.
type hasTextImpl struct {
	text string // The text
}

func (c *hasTextImpl) Text() string {
	return c.text
}

func (c *hasTextImpl) SetText(text string) {
	c.text = text
}

// renderText renders the text.
func (c *hasTextImpl) renderText(w writer) {
	w.Writees(c.text)
}

// HasEnabled interface defines an enabled property.
type HasEnabled interface {
	// Enabled returns the enabled property.
	Enabled() bool

	// SetEnabled sets the enabled property.
	SetEnabled(enabled bool)
}

// newHasEnabledImpl returns a new hasEnabledImpl.
func newHasEnabledImpl() hasEnabledImpl {
	return hasEnabledImpl{true} // Enabled by default
}

// HasEnabled implementation.
type hasEnabledImpl struct {
	enabled bool // The enabled property
}

func (c *hasEnabledImpl) Enabled() bool {
	return c.enabled
}

func (c *hasEnabledImpl) SetEnabled(enabled bool) {
	c.enabled = enabled
}

var _STR_DISABLED = []byte(" disabled=\"disabled\"") // " disabled=\"disabled\""

// renderEnabled renders the enabled attribute.
func (c *hasEnabledImpl) renderEnabled(w writer) {
	if !c.enabled {
		w.Write(_STR_DISABLED)
	}
}

// HasUrl interface defines a URL string property.
type HasUrl interface {
	// URL returns the URL string.
	Url() string

	// SetUrl sets the URL string.
	SetUrl(url string)
}

// newHasUrlImpl creates a new hasUrlImpl
func newHasUrlImpl(url string) hasUrlImpl {
	return hasUrlImpl{url}
}

// HasUrl implementation.
type hasUrlImpl struct {
	url string // The URL string
}

func (c *hasUrlImpl) Url() string {
	return c.url
}

func (c *hasUrlImpl) SetUrl(url string) {
	c.url = url
}

// renderUrl renders the URL string.
func (c *hasUrlImpl) renderUrl(attr string, w writer) {
	w.WriteAttr(attr, c.url)
}

// Horizontal alignment type.
type HAlign string

// Horizontal alignment constants.
const (
	HA_LEFT   HAlign = "left"   // Horizontal left alignment
	HA_CENTER HAlign = "center" // Horizontal center alignment
	HA_RIGHT  HAlign = "right"  // Horizontal right alignment

	HA_DEFAULT HAlign = "" // Browser default (or inherited) horizontal alignment
)

// Vertical alignment type.
type VAlign string

// Vertical alignment constants.
const (
	VA_TOP    VAlign = "top"    // Vertical top alignment
	VA_MIDDLE VAlign = "middle" // Vertical center alignment
	VA_BOTTOM VAlign = "bottom" // Vertical bottom alignment

	VA_DEFAULT VAlign = "" // Browser default (or inherited) vertical alignment
)

// HasHVAlign interfaces defines a horizontal and a vertical
// alignment property.
type HasHVAlign interface {
	// HAlign returns the horizontal alignment.
	HAlign() HAlign

	// SetHAlign sets the horizontal alignment.
	SetHAlign(halign HAlign)

	// VAlign returns the vertical alignment.
	VAlign() VAlign

	// SetVAlign sets the vertical alignment.
	SetVAlign(valign VAlign)

	// SetAlign sets both the horizontal and vertical alignments.
	SetAlign(halign HAlign, valign VAlign)
}

// HasHVAlign implementation.
type hasHVAlignImpl struct {
	halign HAlign // Horizontal alignment
	valign VAlign // Vertical alignment
}

// newHasHVAlignImpl creates a new hasHVAlignImpl
func newHasHVAlignImpl(halign HAlign, valign VAlign) hasHVAlignImpl {
	return hasHVAlignImpl{halign, valign}
}

func (c *hasHVAlignImpl) HAlign() HAlign {
	return c.halign
}

func (c *hasHVAlignImpl) SetHAlign(halign HAlign) {
	c.halign = halign
}

func (c *hasHVAlignImpl) VAlign() VAlign {
	return c.valign
}

func (c *hasHVAlignImpl) SetVAlign(valign VAlign) {
	c.valign = valign
}

func (c *hasHVAlignImpl) SetAlign(halign HAlign, valign VAlign) {
	c.halign = halign
	c.valign = valign
}

// CellFmt interface defines a cell formatter which can be used to
// format and style the wrapper cells of individual components such as
// child components of a PanelView or a Table.
type CellFmt interface {
	// CellFmt allows overriding horizontal and vertical alignment.
	HasHVAlign

	// Style returns the Style builder of the wrapper cell.
	Style() Style

	// Attr returns the explicitly set value of the specified HTML attribute.
	attr(name string) string

	// SetAttr sets the value of the specified HTML attribute.
	// Pass an empty string value to delete the attribute.
	setAttr(name, value string)

	// iAttr returns the explicitly set value of the specified HTML attribute
	// as an int.
	// -1 is returned if the value is not set explicitly or is not an int.
	iAttr(name string) int

	// setIAttr sets the value of the specified HTML attribute as an int.
	setIAttr(name string, value int)
}

// CellFmt implementation
type cellFmtImpl struct {
	hasHVAlignImpl // Has horizontal and vertical alignment implementation

	styleImpl *styleImpl        // Style builder. Lazily initialized.
	attrs     map[string]string // Explicitly set HTML attributes for the cell. Lazily initalized.
}

// newCellFmtImpl creates a new cellFmtImpl.
// Default horizontal alignment is HA_DEFAULT,
// default vertical alignment is VA_DEFAULT.
func newCellFmtImpl() *cellFmtImpl {
	// Initialize hasHVAlignImpl with HA_DEFAULT and VA_DEFAULT
	// so if aligns are not changed, they will not be rendered =>
	// they will be inherited (from TR).
	return &cellFmtImpl{hasHVAlignImpl: newHasHVAlignImpl(HA_DEFAULT, VA_DEFAULT)}
}

func (c *cellFmtImpl) Style() Style {
	if c.styleImpl == nil {
		c.styleImpl = newStyleImpl()
	}
	return c.styleImpl
}

func (c *cellFmtImpl) attr(name string) string {
	return c.attrs[name]
}

func (c *cellFmtImpl) setAttr(name, value string) {
	if c.attrs == nil {
		c.attrs = make(map[string]string, 2)
	}
	if len(value) > 0 {
		c.attrs[name] = value
	} else {
		delete(c.attrs, name)
	}
}

func (c *cellFmtImpl) iAttr(name string) int {
	if value, err := strconv.Atoi(c.attr(name)); err == nil {
		return value
	}
	return -1
}

func (c *cellFmtImpl) setIAttr(name string, value int) {
	c.setAttr(name, strconv.Itoa(value))
}

// render renders the formatted HTML tag for the specified tag name.
// tag must start with a less than sign, e.g. "<td".
func (c *cellFmtImpl) render(tag []byte, w writer) {
	c.renderWithAligns(tag, c.halign, c.valign, w)
}

var _STR_VALIGN = []byte("vertical-align:") // "vertical-align:"

// render renders the formatted HTML tag for the specified tag name
// using the specified alignments instead of ours.
// tag must start with a less than sign, e.g. "<td".
func (c *cellFmtImpl) renderWithAligns(tag []byte, halign HAlign, valign VAlign, w writer) {
	w.Write(tag)

	for name, value := range c.attrs {
		w.WriteAttr(name, value)
	}

	if halign != HA_DEFAULT {
		w.Write(_STR_ALIGN)
		w.Writes(string(halign))
		w.Write(_STR_QUOTE)
	}

	if c.styleImpl != nil {
		c.styleImpl.renderClasses(w)
	}

	if valign != VA_DEFAULT || c.styleImpl != nil {
		w.Write(_STR_STYLE)
		if valign != VA_DEFAULT {
			w.Write(_STR_VALIGN)
			w.Writes(string(valign))
			w.Write(_STR_SEMICOL)
		}
		if c.styleImpl != nil {
			c.styleImpl.renderAttrs(w)
		}
		w.Write(_STR_QUOTE)
	}

	w.Write(_STR_GT)
}

// TableView interface defines a component which is rendered into a table.
type TableView interface {
	// TableView is a Container.
	Container

	// Border returns the border width of the table.
	Border() int

	// SetBorder sets the border width of the table.
	SetBorder(width int)

	// TableView has horizontal and vertical alignment.
	// This is the default horizontal and vertical alignment for
	// all children inside their enclosing cells.
	HasHVAlign

	// CellSpacing returns the cell spacing.
	CellSpacing() int

	// SetCellSpacing sets the cell spacing.
	// Has no effect if layout is LAYOUT_NATURAL.
	SetCellSpacing(spacing int)

	// CellPadding returns the cell spacing.
	CellPadding() int

	// SetCellPadding sets the cell padding.
	// Has no effect if layout is LAYOUT_NATURAL.
	SetCellPadding(padding int)
}

// TableView implementation.
type tableViewImpl struct {
	compImpl       // component implementation
	hasHVAlignImpl // Has horizontal and vertical alignment implementation 
}

// newTableViewImpl creates a new tableViewImpl.
// Default horizontal alignment is HA_DEFAULT,
// default vertical alignment is VA_DEFAULT.
func newTableViewImpl() tableViewImpl {
	// Initialize hasHVAlignImpl with HA_DEFAULT and VA_DEFAULT
	// so if aligns are not changed, they will not be rendered =>
	// they will be inherited (from TR).
	c := tableViewImpl{compImpl: newCompImpl(nil), hasHVAlignImpl: newHasHVAlignImpl(HA_DEFAULT, VA_DEFAULT)}
	c.SetCellSpacing(0)
	c.SetCellPadding(0)
	return c
}

func (c *tableViewImpl) Border() int {
	return c.IAttr("border")
}

func (c *tableViewImpl) SetBorder(width int) {
	c.SetIAttr("border", width)
}

func (c *tableViewImpl) CellSpacing() int {
	return c.IAttr("cellspacing")
}

func (c *tableViewImpl) SetCellSpacing(spacing int) {
	c.SetIAttr("cellspacing", spacing)
}

func (c *tableViewImpl) CellPadding() int {
	return c.IAttr("cellpadding")
}

func (c *tableViewImpl) SetCellPadding(padding int) {
	c.SetIAttr("cellpadding", padding)
}

var _STR_ST_VALIGN = []byte(" style=\"vertical-align:") // " style=\"vertical-align:"

// renderTr renders an HTML TR tag with horizontal and vertical
// alignment info included. 
func (c *tableViewImpl) renderTr(w writer) {
	w.Write(_STR_TR_OP)
	if c.halign != HA_DEFAULT {
		w.Write(_STR_ALIGN)
		w.Writes(string(c.halign))
		w.Write(_STR_QUOTE)
	}
	if c.valign != VA_DEFAULT {
		w.Write(_STR_ST_VALIGN)
		w.Writes(string(c.valign))
		w.Write(_STR_QUOTE)
	}
	w.Write(_STR_GT)
}
