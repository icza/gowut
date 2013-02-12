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

// Defines the Event type and event handling.

package gwu

import (
	"strconv"
)

// Event type (kind) type.
type EventType int

// Converts an Event type to a string.
func (etype EventType) String() string {
	return strconv.Itoa(int(etype))
}

// Event types.
const (
	ETYPE_CLICK      EventType = iota // Mouse click event
	ETYPE_DBL_CLICK                   // Mouse double click event
	ETYPE_MOUSE_DOWN                  // Mouse down event
	ETYPE_MOUSE_MOVE                  // Mouse move event
	ETYPE_MOUSE_OVER                  // Mouse over event
	ETYPE_MOUSE_OUT                   // Mouse out event
	ETYPE_MOUSE_UP                    // Mouse up event
	ETYPE_KEY_DOWN                    // Key down event
	ETYPE_KEY_PRESS                   // Key press event
	ETYPE_KEY_UP                      // Key up event
	ETYPE_BLUR                        // Blur event (component loses focus)
	ETYPE_CHANGE                      // Change event
	ETYPE_FOCUS                       // Focus event (component gains focus)
)

// Attribute names for the event types.
var etypeAttrs map[EventType][]byte = map[EventType][]byte{
	ETYPE_CLICK:      []byte("onclick"),
	ETYPE_DBL_CLICK:  []byte("ondblclick"),
	ETYPE_MOUSE_DOWN: []byte("onmousedown"),
	ETYPE_MOUSE_MOVE: []byte("onmousemove"),
	ETYPE_MOUSE_OVER: []byte("onmouseover"),
	ETYPE_MOUSE_OUT:  []byte("onmouseout"),
	ETYPE_MOUSE_UP:   []byte("onmouseup"),
	ETYPE_KEY_DOWN:   []byte("onkeydown"),
	ETYPE_KEY_PRESS:  []byte("onkeypress"),
	ETYPE_KEY_UP:     []byte("onkeyup"),
	ETYPE_BLUR:       []byte("onblur"),
	ETYPE_CHANGE:     []byte("onchange"),
	ETYPE_FOCUS:      []byte("onfocus")}

// Mouse button type.
type MouseBtn int

// Mouse buttons
const (
	MOUSE_BTN_UNKNOWN MouseBtn = -1 // Unknown mouse button (info not available)
	MOUSE_BTN_LEFT             = 0  // Left mouse button
	MOUSE_BTN_MIDDLE           = 1  // Middle mouse button
	MOUSE_BTN_RIGHT            = 2  // Right mouse button
)

// Modifier key type.
type ModKey int

// Modifier key masks.
const (
	MOD_KEY_ALT   ModKey = 1 << iota // Alt key
	MOD_KEY_CTRL                     // Control key
	MOD_KEY_META                     // Meta key
	MOD_KEY_SHIFT                    // Shift key
)

// Key (keyboard key) type.
type Key int

// Some key codes.
const (
	KEY_BACKSPACE  Key = 8
	KEY_ENTER      Key = 13
	KEY_SHIFT      Key = 16
	KEY_CTRL       Key = 17
	KEY_ALT        Key = 18
	KEY_CAPS_LOCK  Key = 20
	KEY_ESCAPE     Key = 27
	KEY_SPACE      Key = 32
	KEY_PG_UP      Key = 33
	KEY_PG_DOWN    Key = 34
	KEY_END        Key = 35
	KEY_HOME       Key = 36
	KEY_LEFT       Key = 37
	KEY_UP         Key = 38
	KEY_RIGHT      Key = 39
	KEY_DOWN       Key = 40
	KEY_PRINT_SCRN Key = 44
	KEY_INSERT     Key = 45
	KEY_DEL        Key = 46

	KEY_0 Key = 48
	KEY_9 Key = 57

	KEY_A Key = 65
	KEY_Z Key = 90

	KEY_WIN Key = 91

	KEY_NUMPAD_0     Key = 96
	KEY_NUMPAD_9     Key = 105
	KEY_NUMPAD_MUL   Key = 106
	KEY_NUMPAD_PLUS  Key = 107
	KEY_NUMPAD_MINUS Key = 109
	KEY_NUMPAD_DOT   Key = 110
	KEY_NUMPAD_DIV   Key = 111

	KEY_F1  Key = 112
	KEY_F2  Key = 113
	KEY_F3  Key = 114
	KEY_F4  Key = 115
	KEY_F5  Key = 116
	KEY_F6  Key = 117
	KEY_F7  Key = 118
	KEY_F8  Key = 119
	KEY_F9  Key = 120
	KEY_F10 Key = 121
	KEY_F11 Key = 122
	KEY_F12 Key = 123

	KEY_NUM_LOCK    Key = 144
	KEY_SCROLL_LOCK Key = 145
)

// Empty event handler which does nothing.
const EMPTY_EHANDLER emptyEventHandler = emptyEventHandler(0)

// EventHandler interface defines a handler capable of handling events.
type EventHandler interface {
	// Handles the event.
	// 
	// If components are modified in a way that their view changes,
	// these components must be marked dirty in the event object
	// (so the client will see up-to-date state).
	// 
	// If the component tree is modified (new component added
	// or removed for example), then the Panel whose structure was modified
	// has to be marked dirty.
	HandleEvent(e Event)
}

// Event interface defines the event originating from components.
type Event interface {
	// Type returns the type of the event.
	Type() EventType

	// Src returns the source of the event,
	// the component the event is originating from
	Src() Comp

	// Mouse returns the mouse x and y coordinates relative to the component.
	// If no mouse coordinate info is available, (-1, -1) is returned.
	Mouse() (x, y int)

	// MouseWin returns the mouse x and y coordinates inside the window.
	// If no mouse coordinate info is available, (-1, -1) is returned.
	MouseWin() (x, y int)

	// MouseBtn returns the mouse button.
	// If no mouse button info is available, MOUSE_BTN_UNKNOWN is returned.
	MouseBtn() MouseBtn

	// ModKeys returns the states of the modifier keys.
	// The returned value contains the states of all modifier keys,
	// constants of type ModKey can be used to test a specific modifier key,
	// or use the ModKey method.
	ModKeys() int

	// ModKey returns the state of the specified modifier key.
	ModKey(modKey ModKey) bool

	// Key code returns the key code.
	KeyCode() Key

	// Requests the specified window to be reloaded
	// after processing the current event.
	// Tip: pass an empty string to reload the current window.
	ReloadWin(name string)

	// MarkDirty marks components dirty,
	// causing them to be re-rendered after processing the current event.
	// Component re-rendering happens without page reload in the browser.
	// 
	// Note: the Window itself (which is a Comp) can also be marked dirty
	// causing the whole window content to be re-rendered without page reload!
	// 
	// Marking a component dirty also marks all of its decendants dirty, recursively.
	// 
	// Also note that components will not be re-rendered multiple times.
	// For example if a child component and its parent component are both
	// marked dirty, the child component will only be re-rendered once. 
	MarkDirty(comps ...Comp)

	// SetFocusedComp sets the component to be focused after processing
	// the current event.
	SetFocusedComp(comp Comp)

	// Session returns the current session.
	// The Private() method of the session can be used to tell if the session
	// is a private session or the public shared session.
	Session() Session

	// NewSession creates a new (private) session.
	// If the current session (as returned by Session()) is private,
	// it will be removed first.
	NewSession() Session

	// RemoveSess removes (invalidates) the current session.
	// Only private sessions can be removed, calling this
	// when the current session (as returned by Session()) is public is a no-op.
	// After this method Session() will return the shared public session.
	RemoveSess()
}

// Event implementation.
type eventImpl struct {
	server *serverImpl // Server implementation

	etype EventType // Event type
	src   Comp      // Source of the event, the component the event is originating from

	x, y        int         // Mouse coordinates (relative to component)
	wx, wy      int         // Mouse coordinates (inside the window)
	mbtn        MouseBtn    // Mouse button
	modKeys     int         // State of the modifier keys
	keyCode     Key         // Key code
	reload      bool        // Tells if the window has to be reloaded
	reloadWin   string      // The name of the window to be reloaded
	dirtyComps  map[ID]Comp // The dirty components
	focusedComp Comp        // Component to be focused after the event processing
	session     Session     // Session
}

// newEventImpl creates a new eventImpl
func newEventImpl(server *serverImpl, etype EventType, src Comp, session Session) *eventImpl {
	e := eventImpl{server: server, etype: etype, src: src, dirtyComps: make(map[ID]Comp, 2), session: session}
	return &e
}

func (e *eventImpl) Type() EventType {
	return e.etype
}

func (e *eventImpl) Src() Comp {
	return e.src
}

func (e *eventImpl) Mouse() (x, y int) {
	return e.x, e.y
}

func (e *eventImpl) MouseWin() (x, y int) {
	return e.wx, e.wy
}

func (e *eventImpl) MouseBtn() MouseBtn {
	return e.mbtn
}

func (e *eventImpl) ModKeys() int {
	return e.modKeys
}

func (e *eventImpl) ModKey(modKey ModKey) bool {
	return e.modKeys&int(modKey) != 0
}

func (e *eventImpl) KeyCode() Key {
	return e.keyCode
}

func (e *eventImpl) ReloadWin(name string) {
	e.reload = true
	e.reloadWin = name
}

func (e *eventImpl) MarkDirty(comps ...Comp) {
	// We can optimize "on the run" (during dispatching) because we rely on the fact
	// that if the component tree is modified later by a handler, the Panel
	// whose structure was modified will also be marked dirty.
	//
	// So for example if a Panel (P) is already dirty, marking dirty one of its child (A) can be omitted
	// even if later the panel (P) is removed completely, and its child (A) is added to another Panel (P2).
	// In this case P2 will be (must be) marked dirty, and the child (A) will be re-rendered properly
	// along with P2.

	for _, comp := range comps {
		if !e.isDirty(comp) { // If not yet dirty
			// Before adding it, remove all components that are
			// descendants of comp, they will inherit the dirty mark from comp.
			for id, c := range e.dirtyComps {
				if c.DescendantOf(comp) {
					delete(e.dirtyComps, id)
				}
			}

			e.dirtyComps[comp.Id()] = comp
		}
	}
}

// IsDirty returns true if the specified component is already marked dirty.
// Note that a component being dirty makes all of its descendants dirty, recursively.
// 
// Also note that the "dirty" flag might change during the event dispatching
// because if a "clean" component is moved from a dirty parent to a clean parent,
// its inherited dirty flag changes from true to false.
func (e *eventImpl) isDirty(c2 Comp) bool {
	// First-class being dirty:
	if _, found := e.dirtyComps[c2.Id()]; found {
		return true
	}

	// Second-class being dirty:
	for _, c := range e.dirtyComps {
		if c2.DescendantOf(c) {
			return true
		}
	}

	return false
}

func (e *eventImpl) SetFocusedComp(comp Comp) {
	e.focusedComp = comp
}

func (e *eventImpl) Session() Session {
	return e.session
}

func (e *eventImpl) NewSession() Session {
	return e.server.newSession(e)
}

func (e *eventImpl) RemoveSess() {
	e.server.removeSess(e)
}

// Handler function wrapper
type handlerFuncWrapper struct {
	hf func(e Event) // The handler function to be called as part of implementing the EventHandler interface
}

// HandleEvent forwards the call to the handler function.
func (hfw handlerFuncWrapper) HandleEvent(e Event) {
	hfw.hf(e)
}

// Empty Event Handler type.
type emptyEventHandler int

// HandleEvent does nothing as to this is an empty event handler.
func (ee emptyEventHandler) HandleEvent(e Event) {
}
