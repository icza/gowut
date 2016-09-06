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

// Timer component interface and implementation.

package gwu

import (
	"time"
)

// Timer interface defines a component which can generate a timed event
// or a series of timed events periodically.
//
// Timers don't have a visual part, they are used only to generate events.
// The generated events are of type ETypeStateChange.
//
// Note that receiving an event from a Timer (like from any other components)
// updates the last accessed property of the associated session, causing
// a session never to expire if there are active timers on repeat at the
// client side.
//
// Also note that the Timer component operates at the client side meaning
// if the client is closed (or navigates away), events will not be generated.
// (This can also be used to detect if a Window is still open.)
type Timer interface {
	// Timer is a component.
	Comp

	// Timeout returns the timeout duration of the timer.
	Timeout() time.Duration

	// SetTimeout sets the timeout duration of the timer.
	// Event will be generated after the timeout period. If timer is on repeat,
	// events will be generated periodically after each timeout.
	//
	// Note: while this method allows you to pass an arbitrary time.Duration,
	// implementation might be using less precision (most likely millisecond).
	// Durations less than 1 ms might be rounded up to 1 ms.
	SetTimeout(timeout time.Duration)

	// Repeat tells if the timer is on repeat.
	Repeat() bool

	// SetRepeat sets if the timer is on repeat.
	// If timer is on repeat, events will be generated periodically after
	// each timeout.
	SetRepeat(repeat bool)

	// Active tells if the timer is active.
	Active() bool

	// SetActive sets if the timer is active.
	// If the timer is not active, events will not be generated.
	// If a timer is deactivated and activated again, its countdown is reset.
	SetActive(active bool)

	// Reset will cause the timer to restart/reschedule.
	// A Timer does not reset the countdown when it is re-rendered,
	// only if the timer config is changed (e.g. timeout or repeat).
	// By calling Reset() the countdown will reset when the timer is
	// re-rendered.
	Reset()
}

// Timer implementation
type timerImpl struct {
	compImpl // Component implementation

	timeout time.Duration // Timeout of the timer
	repeat  bool          // Tells if timer is on repeat
	active  bool          // Tells if the timer is active
	reset   int           // Reset counter
}

// NewTimer creates a new Timer.
// By default the timer is active and does not repeat.
func NewTimer(timeout time.Duration) Timer {
	return &timerImpl{compImpl: newCompImpl(nil), timeout: timeout, active: true}
}

func (c *timerImpl) Timeout() time.Duration {
	return c.timeout
}

func (c *timerImpl) SetTimeout(timeout time.Duration) {
	if timeout < time.Millisecond {
		timeout = time.Millisecond
	}
	c.timeout = timeout
}

func (c *timerImpl) Repeat() bool {
	return c.repeat
}

func (c *timerImpl) SetRepeat(repeat bool) {
	c.repeat = repeat
}

func (c *timerImpl) Active() bool {
	return c.active
}

func (c *timerImpl) SetActive(active bool) {
	c.active = active
}

func (c *timerImpl) Reset() {
	c.reset++
}

var (
	strSetupTimerOp = []byte("setupTimer(") // "setupTimer("
	strJsSendEvtOp  = []byte("se(null,")    // "se(null,"
)

// renderSetupTimerJs renders the Javascript code which sets up the timer.
// jsVs param holds the values which render Javascript code to be scheduled:
//     setupTimer(compId,"jscode",timeout,repeat,active,reset);
func (c *timerImpl) renderSetupTimerJs(w Writer, jsVs ...interface{}) {
	w.Write(strSetupTimerOp)
	w.Writev(int(c.id))
	w.Write(strComma)
	// js param
	w.Write(strQuote)
	w.Writevs(jsVs...)
	w.Write(strQuote)
	// end of js param
	w.Write(strComma)
	w.Writev(int(c.timeout / time.Millisecond))
	w.Write(strComma)
	w.Writev(c.repeat)
	w.Write(strComma)
	w.Writev(c.active)
	w.Write(strComma)
	w.Writev(c.reset)
	w.Write(strJsFuncCl)
}

func (c *timerImpl) Render(w Writer) {
	w.Write(strSpanOp)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(strGT)

	w.Write(strScriptOp)
	c.renderSetupTimerJs(w, strJsSendEvtOp, int(ETypeStateChange), strComma, int(c.id), strJsFuncCl)
	w.Write(strScriptCl)

	w.Write(strSpanCl)
}
