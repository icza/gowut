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
// Suggested event type to handle actions: ETypeChange
//
// By default the value of the TextBox is synchronized with the server
// on ETypeChange event which is when the TextBox loses focus
// or when the ENTER key is pressed.
// If you want a TextBox to synchronize values during editing
// (while you type in characters), add the ETypeKeyUp event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETypeKeyUp)
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
// Suggested event type to handle actions: ETypeChange
//
// By default the value of the PasswBox is synchronized with the server
// on ETypeChange event which is when the PasswBox loses focus
// or when the ENTER key is pressed.
// If you want a PasswBox to synchronize values during editing
// (while you type in characters), add the ETypeKeyUp event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETypeKeyUp)
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
	isFile    bool // Tells if the text box accepts a file path
	rows, cols int  // Number of displayed rows and columns.
}

var (
	strEncURIThisV = []byte("encodeURIComponent(this.value)") // "encodeURIComponent(this.value)"
)

// NewTextBox creates a new TextBox.
func NewTextBox(text string) TextBox {
	c := newTextBoxImpl(strEncURIThisV, text, false, false)
	c.Style().AddClass("gwu-TextBox")
	return &c
}

// NewFileInputBox creates a new TextBox that accept a file path.
func NewFileInputBox(text string) TextBox {
	c := newTextBoxImpl(strEncURIThisV, text, false, true)
	c.Style().AddClass("gwu-TextBox")
	return &c
}

// NewPasswBox creates a new PasswBox.
func NewPasswBox(text string) TextBox {
	c := newTextBoxImpl(strEncURIThisV, text, true, false)
	c.Style().AddClass("gwu-PasswBox")
	return &c
}

// newTextBoxImpl creates a new textBoxImpl.
func newTextBoxImpl(valueProviderJs []byte, text string, isPassw, isFile bool) textBoxImpl {
	c := textBoxImpl{newCompImpl(valueProviderJs), newHasTextImpl(text), newHasEnabledImpl(), isPassw, isFile, 1, 20}
	c.AddSyncOnETypes(ETypeChange)
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
	value := r.FormValue(paramCompValue)
	if len(value) > 0 {
		c.text = value
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[paramCompValue] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			c.text = values[0]
		}
	}
}

func (c *textBoxImpl) Render(w Writer) {
	if c.rows <= 1 || c.isPassw {
		c.renderInput(w)
	} else {
		c.renderTextArea(w)
	}
}

var (
	strInputOp  = []byte(`<input type="`) // `<input type="`
	strPassword = []byte("password")      // "password"
	strText     = []byte("text")          // "text"
	strFile     = []byte("file")          // "file"
	strSize     = []byte(`" size="`)      // `" size="`
	strValue    = []byte(` value="`)      // ` value="`
	strInputCl  = []byte(`"/>`)           // `"/>`
)

// renderInput renders the component as an input HTML tag.
func (c *textBoxImpl) renderInput(w Writer) {
	w.Write(strInputOp)
	if c.isPassw {
		w.Write(strPassword)
	} else if c.isFile{
		w.Write(strFile)
	} else {
		w.Write(strText)
	}
	w.Write(strSize)
	w.Writev(c.cols)
	w.Write(strQuote)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	w.Write(strValue)
	c.renderText(w)
	w.Write(strInputCl)
}

var (
	strTextareaOp   = []byte("<textarea")   // "<textarea"
	strRows         = []byte(` rows="`)     // ` rows="`
	strCols         = []byte(`" cols="`)    // `" cols="`
	strTextAreaOpCl = []byte("\">\n")       // "\">\n"
	strTextAreaCl   = []byte("</textarea>") // "</textarea>"
)

// renderTextArea renders the component as an textarea HTML tag.
func (c *textBoxImpl) renderTextArea(w Writer) {
	w.Write(strTextareaOp)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	// New line char after the <textarea> tag is ignored.
	// So we must render a newline after textarea, else if text value
	// starts with a new line, it will be omitted!
	w.Write(strRows)
	w.Writev(c.rows)
	w.Write(strCols)
	w.Writev(c.cols)
	w.Write(strTextAreaOpCl)

	c.renderText(w)
	w.Write(strTextAreaCl)
}
