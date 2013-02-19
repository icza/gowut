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

// Our improved writer with helper methods to easier write data we need.

package gwu

import (
	"errors"
	"fmt"
	"html"
	"io"
	"strconv"
)

// Number of cached ints.
const _CACHED_INTS = 32

// Byte slice vars (constants) of frequently used strings.
// Render methods use these to avoid array allocations
// when converting strings to byte slices in order to write them.
var (
	_STR_SPACE    = []byte(" ")   // " " (space string)
	_STR_QUOTE    = []byte("\"")  // "\"" (quotation mark)
	_STR_EQ_QUOTE = []byte("=\"") // "=\"" (equal sign and a quotation mark)
	_STR_COMMA    = []byte(",")   // "," (comma string)
	_STR_COLON    = []byte(":")   // ":" (colon string)
	_STR_SEMICOL  = []byte(";")   // ";" (semicolon string)
	_STR_LT       = []byte("<")   // "<" (less than string)
	_STR_GT       = []byte(">")   // ">" (greater than string)

	_STR_SPAN_OP  = []byte("<span")    // "<span"
	_STR_SPAN_CL  = []byte("</span>")  // "</span>"
	_STR_TABLE_OP = []byte("<table")   // "<table"
	_STR_TABLE_CL = []byte("</table>") // "</table>"
	_STR_TD       = []byte("<td>")     // "<td>"
	_STR_TR       = []byte("<tr>")     // "<tr>"
	_STR_TD_OP    = []byte("<td")      // "<td"
	_STR_TR_OP    = []byte("<tr")      // "<tr"

	_STR_STYLE = []byte(" style=\"") // " style=\""
	_STR_CLASS = []byte(" class=\"") // " class=\""
	_STR_ALIGN = []byte(" align=\"") // " align=\""

	_STR_INTS  [_CACHED_INTS][]byte                                            // Numbers
	_STR_BOOLS = map[bool][]byte{false: []byte("false"), true: []byte("true")} // Bools
)

// init initializes the cached ints.
func init() {
	for i := 0; i < _CACHED_INTS; i++ {
		_STR_INTS[i] = []byte(strconv.Itoa(i))
	}
}

// writer is an improved writer with helper methods
// to easier write data we need
type writer struct {
	io.Writer // Writer implementation
}

// NewWriter returns an implementation of our writer.
func NewWriter(w io.Writer) writer {
	return writer{w}
}

// Writev writes a value.
func (w writer) Writev(v interface{}) (n int, err error) {
	switch v2 := v.(type) {
	case string:
		return w.Write([]byte(v2))
	case int:
		if v2 < _CACHED_INTS && v2 >= 0 {
			return w.Write(_STR_INTS[v2])
		}
		return w.Write([]byte(strconv.Itoa(v2)))
	case []byte:
		return w.Write(v2)
	case fmt.Stringer:
		return w.Write([]byte(v2.String()))
	case bool:
		return w.Write(_STR_BOOLS[v2])
	}

	fmt.Printf("Not supported type: %T\n", v)
	return 0, errors.New(fmt.Sprintf("Not supported type: %T", v))
}

// Writevs writes values.
func (w writer) Writevs(v ...interface{}) (n int, err error) {
	for _, v2 := range v {
		var m int
		m, err = w.Writev(v2)
		n += m
		if err != nil {
			return
		}
	}
	return
}

// Writes writes a string.
func (w writer) Writes(s string) (n int, err error) {
	return w.Write([]byte(s))
}

// Writess writes strings.
func (w writer) Writess(ss ...string) (n int, err error) {
	for _, s := range ss {
		var m int
		m, err = w.Write([]byte(s))
		n += m
		if err != nil {
			return
		}
	}
	return
}

// Writees writes a string after html-escaping it.
func (w writer) Writees(s string) (n int, err error) {
	return w.Write([]byte(html.EscapeString(s)))
}

// WriteAttr writes an attribute in the form of:
// " name=\"value\""
func (w writer) WriteAttr(name, value string) (n int, err error) {
	// Easiest implementation would be:
	// return w.Writevs(_STR_SPACE, name, _STR_EQ_QUOTE, value, _STR_QUOTE)

	// ...but since this is called very frequently, I allow some extra lines of code
	// for the sake of efficiency (also avoiding array allocation for the varargs...) 

	n, err = w.Write(_STR_SPACE)
	if err != nil {
		return
	}

	var m int
	m, err = w.Write([]byte(name))
	n += m
	if err != nil {
		return
	}

	m, err = w.Write(_STR_EQ_QUOTE)
	n += m
	if err != nil {
		return
	}

	m, err = w.Write([]byte(value))
	n += m
	if err != nil {
		return
	}

	m, err = w.Write(_STR_QUOTE)
	n += m

	return
}
