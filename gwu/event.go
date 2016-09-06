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
	"net/http"
	"strconv"
)

// EventType is the event type (kind) type.
type EventType int

// Converts an Event type to a string.
func (etype EventType) String() string {
	return strconv.Itoa(int(etype))
}

// Event types.
const (
	// General events for all components
	ETypeClick     EventType = iota // Mouse click event
	ETypeDblClick                   // Mouse double click event
	ETypeMousedown                  // Mouse down event
	ETypeMouseMove                  // Mouse move event
	ETypeMouseOver                  // Mouse over event
	ETypeMouseOut                   // Mouse out event
	ETypeMouseUp                    // Mouse up event
	ETypeKeyDown                    // Key down event
	ETypeKeyPress                   // Key press event
	ETypeKeyUp                      // Key up event
	ETypeBlur                       // Blur event (component loses focus)
	ETypeChange                     // Change event (value change)
	ETypeFocus                      // Focus event (component gains focus)

	// Window events (for Window only)
	ETypeWinLoad   // Window load event
	ETypeWinUnload // Window unload event

	// Internal events, generated and dispatched internally while processing another event
	ETypeStateChange // State change
)

// EventCategory is the event type category.
type EventCategory int

// Event type categories.
const (
	ECatGeneral  EventCategory = iota // General event type for all components
	ECatWindow                        // Window event type for Window only
	ECatInternal                      // Internal event generated and dispatched internally while processing another event

	ECatUnknown EventCategory = -1 // Unknown event category
)

// Category returns the event type category.
func (etype EventType) Category() EventCategory {
	switch {
	case etype >= ETypeClick && etype <= ETypeFocus:
		return ECatGeneral
	case etype >= ETypeWinLoad && etype <= ETypeWinUnload:
		return ECatWindow
	case etype >= ETypeStateChange && etype <= ETypeStateChange:
		return ECatInternal
	}

	return ECatUnknown
}

// Attribute names for the general event types; only for the general event types.
var etypeAttrs map[EventType][]byte = map[EventType][]byte{
	ETypeClick:     []byte("onclick"),
	ETypeDblClick:  []byte("ondblclick"),
	ETypeMousedown: []byte("onmousedown"),
	ETypeMouseMove: []byte("onmousemove"),
	ETypeMouseOver: []byte("onmouseover"),
	ETypeMouseOut:  []byte("onmouseout"),
	ETypeMouseUp:   []byte("onmouseup"),
	ETypeKeyDown:   []byte("onkeydown"),
	ETypeKeyPress:  []byte("onkeypress"),
	ETypeKeyUp:     []byte("onkeyup"),
	ETypeBlur:      []byte("onblur"),
	ETypeChange:    []byte("onchange"),
	ETypeFocus:     []byte("onfocus")}

// Function names for window event types.
var etypeFuncs map[EventType][]byte = map[EventType][]byte{
	ETypeWinLoad:   []byte("onload"),
	ETypeWinUnload: []byte("onbeforeunload")} // Bind it to onbeforeunload (instead of onunload) for several reasons (onunload might cause trouble for AJAX; onunload is not called in IE if page is just refreshed...)

// MouseBtn is the mouse button type.
type MouseBtn int

// Mouse buttons
const (
	MouseBtnUnknown MouseBtn = -1 // Unknown mouse button (info not available)
	MouseBtnLeft             = 0  // Left mouse button
	MouseBtnMiddle           = 1  // Middle mouse button
	MouseBtnRight            = 2  // Right mouse button
)

// ModKey is the modifier key type.
type ModKey int

// Modifier key masks.
const (
	ModKeyAlt   ModKey = 1 << iota // Alt key
	ModKeyCtrl                     // Control key
	ModKeyMeta                     // Meta key
	ModKeyShift                    // Shift key
)

// Key (keyboard key) type.
type Key int

// Some key codes.
const (
	KeyBackspace Key = 8
	KeyEnter         = 13
	KeyShift         = 16
	KeyCtrl          = 17
	KeyAlt           = 18
	KeyCapsLock      = 20
	KeyEscape        = 27
	KeySpace         = 32
	KeyPgUp          = 33
	KeyPgDown        = 34
	KeyEnd           = 35
	KeyHome          = 36
	KeyLeft          = 37
	KeyUp            = 38
	KeyRight         = 39
	KeyDown          = 40
	KeyPrintScrn     = 44
	KeyInsert        = 45
	KeyDel           = 46

	Key0 = 48
	Key9 = 57

	KeyA = 65
	KeyZ = 90

	KeyWin = 91

	KeyNumpad0     = 96
	KeyNumpad9     = 105
	KeyNumpadMul   = 106
	KeyNumpadPlus  = 107
	KeyNumpadMinus = 109
	KeyNumpadDot   = 110
	KeyNumpadDiv   = 111

	KeyF1  = 112
	KeyF2  = 113
	KeyF3  = 114
	KeyF4  = 115
	KeyF5  = 116
	KeyF6  = 117
	KeyF7  = 118
	KeyF8  = 119
	KeyF9  = 120
	KeyF10 = 121
	KeyF11 = 122
	KeyF12 = 123

	KeyNumLock    = 144
	KeyScrollLock = 145
)

// EmptyEHandler is the empty event handler which does nothing.
const EmptyEHandler emptyEventHandler = 0

// EventHandler interface defines a handler capable of handling events.
type EventHandler interface {
	// Handles the event.
	//
	// If components are modified in a way that their view changes,
	// these components must be marked dirty in the event object
	// (so the client will see up-to-date state).
	//
	// If the component tree is modified (new component added
	// or removed for example), then the Container whose structure
	// was modified has to be marked dirty.
	HandleEvent(e Event)
}

// Event interface defines the event originating from components.
type Event interface {
	// Type returns the type of the event.
	Type() EventType

	// Src returns the source of the event,
	// the component the event is originating from
	Src() Comp

	// Parent returns the parent event if there's one.
	// Usually internal events have parent event for which the internal
	// event was created and dispatched.
	// The parent event can be used to identify the original source and event type.
	Parent() Event

	// Mouse returns the mouse x and y coordinates relative to the component.
	// If no mouse coordinate info is available, (-1, -1) is returned.
	Mouse() (x, y int)

	// MouseWin returns the mouse x and y coordinates inside the window.
	// If no mouse coordinate info is available, (-1, -1) is returned.
	MouseWin() (x, y int)

	// MouseBtn returns the mouse button.
	// If no mouse button info is available, MouseBtnUnknown is returned.
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
	// Marking a component dirty also marks all of its descendants dirty, recursively.
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

	// forkEvent forks a new Event from this one.
	// The new event will have a parent pointing to us.
	// Accessing/changing the session and defining post-event actions in the forked
	// event works as if they would be done on this event.
	forkEvent(etype EventType, src Comp) Event
}

// HasRequestResponse defines methods to acquire / access
// http.ResponseWriter and http.Request from something that supports this.
//
// The concrete type that implements Event does implement this too,
// but this is not added to the Event interface intentionally to not urge the use of this.
// Users should not rely on this as in a future implementation there might not be
// a response and request associated with an event.
// But this may be useful in certain scenarios, such as you need to know the client IP address,
// or you want to use custom authentication that needs the request/response.
//
// To get access to these methods, simply use a type assertion, asserting that the event value
// implements this interface. For example:
//
//     someButton.AddEHandlerFunc(func(e gwu.Event) {
//         if hrr, ok := e.(gwu.HasRequestResponse); ok {
//             req := hrr.Request()
//             log.Println("Client addr:", req.RemoteAddr)
//         }
//     }, gwu.ETypeClick)
type HasRequestResponse interface {
	// ResponseWriter returns the associated HTTP response writer.
	ResponseWriter() http.ResponseWriter

	// Request returns the associated HTTP request.
	Request() *http.Request
}

// Event implementation.
type eventImpl struct {
	etype  EventType  // Event type
	src    Comp       // Source of the event, the component the event is originating from
	parent *eventImpl // Optional parent event

	x, y int // Mouse coordinates (relative to component); not part of shared data because they component-relative

	shared *sharedEvtData // Shared event data
}

// Event data shared between an event and its child events (forks).
type sharedEvtData struct {
	server *serverImpl // Server implementation

	wx, wy  int      // Mouse coordinates (inside the window)
	mbtn    MouseBtn // Mouse button
	modKeys int      // State of the modifier keys
	keyCode Key      // Key code

	reload      bool        // Tells if the window has to be reloaded
	reloadWin   string      // The name of the window to be reloaded
	dirtyComps  map[ID]Comp // The dirty components
	focusedComp Comp        // Component to be focused after the event processing
	session     Session     // Session

	rw  http.ResponseWriter // ResponseWriter of the HTTP request the event was created from
	req *http.Request       // Request of the HTTP request the event was created from
}

// newEventImpl creates a new eventImpl
func newEventImpl(etype EventType, src Comp, server *serverImpl, session Session,
	rw http.ResponseWriter, req *http.Request) *eventImpl {
	e := eventImpl{etype: etype, src: src,
		shared: &sharedEvtData{server: server, dirtyComps: make(map[ID]Comp, 2), session: session, rw: rw, req: req}}
	return &e
}

func (e *eventImpl) Type() EventType {
	return e.etype
}

func (e *eventImpl) Src() Comp {
	return e.src
}

func (e *eventImpl) Parent() Event {
	return e.parent
}

func (e *eventImpl) Mouse() (x, y int) {
	return e.x, e.y
}

func (e *eventImpl) MouseWin() (x, y int) {
	return e.shared.wx, e.shared.wy
}

func (e *eventImpl) MouseBtn() MouseBtn {
	return e.shared.mbtn
}

func (e *eventImpl) ModKeys() int {
	return e.shared.modKeys
}

func (e *eventImpl) ModKey(modKey ModKey) bool {
	return e.shared.modKeys&int(modKey) != 0
}

func (e *eventImpl) KeyCode() Key {
	return e.shared.keyCode
}

func (e *eventImpl) ReloadWin(name string) {
	e.shared.reload = true
	e.shared.reloadWin = name
}

func (e *eventImpl) MarkDirty(comps ...Comp) {
	// We can optimize "on the run" (during dispatching) because we rely on the fact
	// that if the component tree is modified later by a handler, the Container
	// whose structure was modified will also be marked dirty.
	//
	// So for example if a Panel (P) is already dirty, marking dirty one of its child (A) can be omitted
	// even if later the panel (P) is removed completely, and its child (A) is added to another Panel (P2).
	// In this case P2 will be (must be) marked dirty, and the child (A) will be re-rendered properly
	// along with P2.

	shared := e.shared

	for _, comp := range comps {
		if !shared.dirty(comp) { // If not yet dirty
			// Before adding it, remove all components that are
			// descendants of comp, they will inherit the dirty mark from comp.
			for id, c := range shared.dirtyComps {
				if c.DescendantOf(comp) {
					delete(shared.dirtyComps, id)
				}
			}

			shared.dirtyComps[comp.Id()] = comp
		}
	}
}

// dirty returns true if the specified component is already marked dirty.
// Note that a component being dirty makes all of its descendants dirty, recursively.
//
// Also note that the "dirty" flag might change during the event dispatching
// because if a "clean" component is moved from a dirty parent to a clean parent,
// its inherited dirty flag changes from true to false.
func (s *sharedEvtData) dirty(c2 Comp) bool {
	// First-class being dirty:
	if _, found := s.dirtyComps[c2.Id()]; found {
		return true
	}

	// Second-class being dirty:
	for _, c := range s.dirtyComps {
		if c2.DescendantOf(c) {
			return true
		}
	}

	return false
}

func (e *eventImpl) SetFocusedComp(comp Comp) {
	e.shared.focusedComp = comp
}

func (e *eventImpl) Session() Session {
	return e.shared.session
}

func (e *eventImpl) NewSession() Session {
	return e.shared.server.newSession(e)
}

func (e *eventImpl) RemoveSess() {
	e.shared.server.removeSess(e)
}

func (e *eventImpl) forkEvent(etype EventType, src Comp) Event {
	return &eventImpl{etype: etype, src: src, parent: e,
		x: -1, y: -1, // Mouse coordinates are unknown in the new source component...
		shared: e.shared}
}

func (e *eventImpl) ResponseWriter() http.ResponseWriter {
	return e.shared.rw
}

func (e *eventImpl) Request() *http.Request {
	return e.shared.req
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
