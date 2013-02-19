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

// Defines the TextBox component.

package gwu

import (
	"net/http"
	"strconv"
)

// TextBox interface defines a component for text input purpose.
// 
// Suggested event type to handle actions: ETYPE_CHANGE
//
// By default the value of the TextBox is synchronized with the server
// on ETYPE_CHANGE event which is when the TextBox loses focus
// or when the ENTER key is pressed.
// If you want a TextBox to synchronize values during editing
// (while you type in characters), add the ETYPE_KEY_UP event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETYPE_KEY_UP)
// 
// Default style class: "gwu-TextBox"
type TextBox interface {
	// TextBox is a component.
	Comp

	// TextBox has text.
	HasText

	// TextBox can be enabled/disabled.
	HasEnabled

	// ReadOnly returns if the text box is read-only.
	ReadOnly() bool

	// SetReadOnly sets if the text box is read-only.
	SetReadOnly(readOnly bool)

	// Rows returns the number of displayed rows.
	Rows() int

	// SetRows sets the number of displayed rows.
	// rows=1 will make this a simple, one-line input text box,
	// rows>1 will make this a text area.
	SetRows(rows int)

	// Cols returns the number of displayed columns.
	Cols() int

	// SetCols sets the number of displayed columns.
	SetCols(cols int)

	// MaxLength returns the maximum number of characters
	// allowed in the text box.
	// -1 is returned if there is no maximum length set.
	MaxLength() int

	// SetMaxLength sets the maximum number of characters
	// allowed in the text box.
	// Pass -1 to not limit the maximum length.
	SetMaxLength(maxLength int)
}

// PasswBox interface defines a text box for password input purpose.
// 
// Suggested event type to handle actions: ETYPE_CHANGE
//
// By default the value of the PasswBox is synchronized with the server
// on ETYPE_CHANGE event which is when the PasswBox loses focus
// or when the ENTER key is pressed.
// If you want a PasswBox to synchronize values during editing
// (while you type in characters), add the ETYPE_KEY_UP event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETYPE_KEY_UP)
// 
// Default style class: "gwu-PasswBox"
type PasswBox interface {
	// PasswBox is a TextBox.
	TextBox
}

// TextBox implementation.
type textBoxImpl struct {
	compImpl       // Component implementation
	hasTextImpl    // Has text implementation
	hasEnabledImpl // Has enabled implementation

	isPassw    bool // Tells if the text box is a password box
	rows, cols int  // Number of displayed rows and columns.
}

var (
	_STR_ENC_URI_THIS_V = []byte("encodeURIComponent(this.value)") // "encodeURIComponent(this.value)"
)

// NewTextBox creates a new TextBox.
func NewTextBox(text string) TextBox {
	c := newTextBoxImpl(_STR_ENC_URI_THIS_V, text, false)
	c.Style().AddClass("gwu-TextBox")
	return &c
}

// NewPasswBox creates a new PasswBox.
func NewPasswBox(text string) TextBox {
	c := newTextBoxImpl(_STR_ENC_URI_THIS_V, text, true)
	c.Style().AddClass("gwu-PasswBox")
	return &c
}

// newTextBoxImpl creates a new textBoxImpl.
func newTextBoxImpl(valueProviderJs []byte, text string, isPassw bool) textBoxImpl {
	c := textBoxImpl{newCompImpl(valueProviderJs), newHasTextImpl(text), newHasEnabledImpl(), isPassw, 1, 20}
	c.AddSyncOnETypes(ETYPE_CHANGE)
	return c
}

func (c *textBoxImpl) ReadOnly() bool {
	ro := c.Attr("readonly")
	return len(ro) > 0
}

func (c *textBoxImpl) SetReadOnly(readOnly bool) {
	if readOnly {
		c.SetAttr("readonly", "readonly")
	} else {
		c.SetAttr("readonly", "")
	}
}

func (c *textBoxImpl) Rows() int {
	return c.rows
}

func (c *textBoxImpl) SetRows(rows int) {
	c.rows = rows
}

func (c *textBoxImpl) Cols() int {
	return c.cols
}

func (c *textBoxImpl) SetCols(cols int) {
	c.cols = cols
}

func (c *textBoxImpl) MaxLength() int {
	if ml := c.Attr("maxlength"); len(ml) > 0 {
		if i, err := strconv.Atoi(ml); err == nil {
			return i
		}
	}
	return -1
}

func (c *textBoxImpl) SetMaxLength(maxLength int) {
	if maxLength < 0 {
		c.SetAttr("maxlength", "")
	} else {
		c.SetAttr("maxlength", strconv.Itoa(maxLength))
	}
}

func (c *textBoxImpl) preprocessEvent(event Event, r *http.Request) {
	// Empty string for text box is a valid value.
	// So we have to check whether it is supplied, not just whether its len() > 0 
	value := r.FormValue(_PARAM_COMP_VALUE)
	if len(value) > 0 {
		c.text = value
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[_PARAM_COMP_VALUE] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			c.text = values[0]
		}
	}
}

func (c *textBoxImpl) Render(w writer) {
	if c.rows <= 1 || c.isPassw {
		c.renderInput(w)
	} else {
		c.renderTextArea(w)
	}
}

var (
	_STR_INPUT_OP = []byte("<input type=\"") // "<input type=\""
	_STR_PASSWORD = []byte("password")       // "password"
	_STR_TEXT     = []byte("text")           // "text"
	_STR_SIZE     = []byte("\" size=\"")     // "\" size=\""
	_STR_VALUE    = []byte(" value=\"")      // " value=\""
	_STR_INPUT_CL = []byte("\"/>")           // "\"/>"
)

// renderInput renders the component as an input HTML tag.
func (c *textBoxImpl) renderInput(w writer) {
	w.Write(_STR_INPUT_OP)
	if c.isPassw {
		w.Write(_STR_PASSWORD)
	} else {
		w.Write(_STR_TEXT)
	}
	w.Write(_STR_SIZE)
	w.Writev(c.cols)
	w.Write(_STR_QUOTE)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	w.Write(_STR_VALUE)
	c.renderText(w)
	w.Write(_STR_INPUT_CL)
}

var (
	_STR_TEXTAREA_OP    = []byte("<textarea")   // "<textarea"
	_STR_ROWS           = []byte(" rows=\"")    // " rows=\""
	_STR_COLS           = []byte("\" cols=\"")  // "\" cols=\""
	_STR_TEXTAREA_OP_CL = []byte("\">\n")       // "\">\n"
	_STR_TEXTAREA_CL    = []byte("</textarea>") // "</textarea>"
)

// renderTextArea renders the component as an textarea HTML tag.
func (c *textBoxImpl) renderTextArea(w writer) {
	w.Write(_STR_TEXTAREA_OP)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	// New line char after the <textarea> tag is ignored.
	// So we must render a newline after textarea, else if text value
	// starts with a new line, it will be ommitted!
	w.Write(_STR_ROWS)
	w.Writev(c.rows)
	w.Write(_STR_COLS)
	w.Writev(c.cols)
	w.Write(_STR_TEXTAREA_OP_CL)

	c.renderText(w)
	w.Write(_STR_TEXTAREA_CL)
}
