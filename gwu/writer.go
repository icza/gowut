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

// Gowut's improved and optimized writer with helper methods to easier write data we need.

package gwu

import (
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"strconv"
)

// Number of cached ints.
const cachedInts = 64

// Byte slice vars (constants) of frequently used strings.
// Render methods use these to avoid array allocations
// when converting strings to byte slices in order to write them.
var (
	strSpace    = []byte(" ")  // " " (space string)
	strQuote    = []byte(`"`)  // `"` (quotation mark)
	strEqQuote  = []byte(`="`) // `="` (equal sign and a quotation mark)
	strComma    = []byte(",")  // "," (comma string)
	strColon    = []byte(":")  // ":" (colon string)
	strSemicol  = []byte(";")  // ";" (semicolon string)
	strLT       = []byte("<")  // "<" (less than string)
	strGT       = []byte(">")  // ">" (greater than string)
	strParenCl  = []byte(")")  // ")" (closing parenthesis)
	strJsFuncCl = []byte(");") // ");" (closing parenthesis and a semicolon)

	strSpanOp   = []byte("<span")     // "<span"
	strSpanCl   = []byte("</span>")   // "</span>"
	strTableOp  = []byte("<table")    // "<table"
	strTableCl  = []byte("</table>")  // "</table>"
	strTD       = []byte("<td>")      // "<td>"
	strTR       = []byte("<tr>")      // "<tr>"
	strTDOp     = []byte("<td")       // "<td"
	strTROp     = []byte("<tr")       // "<tr"
	strScriptOp = []byte("<script>")  // "<script>"
	strScriptCl = []byte("</script>") // "</script>"

	strStyle = []byte(` style="`) // ` style="`
	strClass = []byte(` class="`) // ` class="`
	strAlign = []byte(` align="`) // ` align="`

	strInts  [cachedInts][]byte                                              // Numbers
	strBools = map[bool][]byte{false: []byte("false"), true: []byte("true")} // Bools
)

// init initializes the cached ints.
func init() {
	for i := 0; i < cachedInts; i++ {
		strInts[i] = []byte(strconv.Itoa(i))
	}
}

// Writer is an improved and optimized io.Writer with additionial helper methods
// to easier write data we need to render components.
type Writer interface {
	io.Writer // Writer is an io.Writer

	// Writev writes a value. It is highly optimized for certain values/types.
	// Supported value types are string, int, []byte, bool.
	Writev(v interface{}) (n int, err error)

	// Writevs writes values. It is highly optimized for certain values/types.
	// For supported value types see Writev().
	Writevs(v ...interface{}) (n int, err error)

	// Writes writes a string.
	Writes(s string) (n int, err error)

	// Writess writes strings.
	Writess(ss ...string) (n int, err error)

	// Writees writes a string after html-escaping it.
	Writees(s string) (n int, err error)

	// WriteAttr writes an attribute in the form of:
	// ` name="value"`
	WriteAttr(name, value string) (n int, err error)
}

// writerImpl is the implementation of our Writer.
type writerImpl struct {
	io.Writer // Writer implementation
}

// NewWriter returns a new Writer, wrapping the specified io.Writer.
func NewWriter(w io.Writer) Writer {
	return writerImpl{w}
}

func (w writerImpl) Writev(v interface{}) (n int, err error) {
	switch v2 := v.(type) {
	case string:
		return w.Write([]byte(v2))
	case int:
		if v2 < cachedInts && v2 >= 0 {
			return w.Write(strInts[v2])
		}
		return w.Write([]byte(strconv.Itoa(v2)))
	case []byte:
		return w.Write(v2)
	case fmt.Stringer:
		return w.Write([]byte(v2.String()))
	case bool:
		return w.Write(strBools[v2])
	}

	log.Printf("Not supported type: %T\n", v)
	return 0, errors.New(fmt.Sprintf("Not supported type: %T", v))
}

func (w writerImpl) Writevs(v ...interface{}) (n int, err error) {
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

func (w writerImpl) Writes(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func (w writerImpl) Writess(ss ...string) (n int, err error) {
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

func (w writerImpl) Writees(s string) (n int, err error) {
	return w.Write([]byte(html.EscapeString(s)))
}

func (w writerImpl) WriteAttr(name, value string) (n int, err error) {
	// Easiest implementation would be:
	// return w.Writevs(strSpace, name, strEqQuote, value, strQuote)

	// ...but since this is called very frequently, I allow some extra lines of code
	// for the sake of efficiency (also avoiding array allocation for the varargs...)

	n, err = w.Write(strSpace)
	if err != nil {
		return
	}

	var m int
	m, err = w.Write([]byte(name))
	n += m
	if err != nil {
		return
	}

	m, err = w.Write(strEqQuote)
	n += m
	if err != nil {
		return
	}

	m, err = w.Write([]byte(value))
	n += m
	if err != nil {
		return
	}

	m, err = w.Write(strQuote)
	n += m

	return
}
