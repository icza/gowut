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

// Session monitor component interface and implementation.

package gwu

import (
	"time"
)

// SessMonitor interface defines a component which monitors and displays
// the session timeout and network connectivity at client side without
// interacting with the session.
//
// Default style classes: "gwu-SessMonitor", "gwu-SessMonitor-Expired",
// ".gwu-SessMonitor-Error"
type SessMonitor interface {
	// SessMonitor is a Timer, but it does not generate Events!
	Timer

	// SetJsConverter sets the Javascript function name which converts
	// a float second time value to a displayable string.
	// The default value is "convertSessTimeout" whose implementation is:
	//     function convertSessTimeout(sec) {
	//         if (sec <= 0)
	//             return "Expired!";
	//         else if (sec < 60)
	//             return "<1 min";
	//         else
	//             return "~" + Math.round(sec / 60) + " min";
	//     }
	SetJsConverter(jsFuncName string)

	// JsConverter returns the name of the Javascript function which converts
	// float second time values to displayable strings.
	JsConverter() string
}

// SessMonitor implementation
type sessMonitorImpl struct {
	timerImpl // Timer implementation
}

// NewSessMonitor creates a new SessMonitor.
// By default it is active repeats with 1 minute timeout duration.
func NewSessMonitor() SessMonitor {
	c := &sessMonitorImpl{
		timerImpl{compImpl: newCompImpl(nil), timeout: time.Minute, active: true, repeat: true},
	}
	c.Style().AddClass("gwu-SessMonitor")
	c.SetJsConverter("convertSessTimeout")
	return c
}

func (c *sessMonitorImpl) SetJsConverter(jsFuncName string) {
	c.SetAttr("gwuJsFuncName", jsFuncName)
}

func (c *sessMonitorImpl) JsConverter() string {
	return c.Attr("gwuJsFuncName")
}

var (
	strEmptySpan     = []byte("<span></span>") // "<span></span>"
	strJsCheckSessOp = []byte("checkSession(") // "checkSession("
)

func (c *sessMonitorImpl) Render(w Writer) {
	w.Write(strSpanOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	w.Write(strEmptySpan) // Placeholder for session timeout value

	w.Write(strScriptOp)
	c.renderSetupTimerJs(w, strJsCheckSessOp, int(c.id), strParenCl)
	// Call sess check right away:
	w.Write(strJsCheckSessOp)
	w.Writev(int(c.id))
	w.Write(strJsFuncCl)
	w.Write(strScriptCl)

	w.Write(strSpanCl)
}
