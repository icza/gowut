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

// Implementation of the GWU session.

package gwu

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

// Session interface defines the session to the GWU users (clients).
type Session interface {
	// Id returns the id of the session.
	Id() string

	// New tells if the session is new meaning the client
	// does not (yet) know about it.
	New() bool

	// Private tells if the session is a private session.
	// There is only one public session, and it is shared
	// between the "sessionless" users.
	Private() bool

	// AddWin adds a window to the session.
	// Returns an error if window name is empty or
	// a window with the same name has already been added.
	AddWin(w Window) error

	// RemoveWin removes a window from the session.
	// Returns if the window was removed from the session.
	RemoveWin(w Window) bool

	// SortedWins returns a sorted slice of windows.
	// The slice is sorted by window text (title).
	SortedWins() []Window

	// WinByName returns a window specified by its name.
	WinByName(name string) Window

	// Attr returns the value of an attribute stored in the session.
	// TODO use an interface type something like "serializable".
	Attr(name string) interface{}

	// SetAttr sets the value of an attribute stored in the session.
	// Pass the nil value to delete the attribute.
	SetAttr(name string, value interface{})

	// Created returns the time when the session was created.
	Created() time.Time

	// Accessed returns the time when the session was last accessed.
	Accessed() time.Time

	// Timeout returns the session timeout.
	Timeout() time.Duration

	// SetTimeout sets the session timeout.
	SetTimeout(timeout time.Duration)

	// access registers an access to the session.
	access()

	// ClearNew clears the new flag.
	// After this New() will return false.
	clearNew()

	// rwMutex returns the RW mutex of the session.
	rwMutex() *sync.RWMutex
}

// Session implementation.
type sessionImpl struct {
	id       string                 // Id of the session
	isNew    bool                   // Tells if the session is new
	created  time.Time              // Creation time
	accessed time.Time              // Last accessed time
	windows  map[string]Window      // Windows of the session
	attrs    map[string]interface{} // Attributes stored in the session
	timeout  time.Duration          // Session timeout

	rwMutex_ *sync.RWMutex // RW mutex to synchronize session (and related Window and component) access
}

// newSessionImpl creates a new sessionImpl.
// The default timeout is 30 minutes.
func newSessionImpl(private bool) sessionImpl {
	var id string
	// The public session has an empty string id
	if private {
		id = genId()
	}

	now := time.Now()

	// Initialzie private sessions as new, but not the public session
	return sessionImpl{id: id, isNew: private, created: now, accessed: now, windows: make(map[string]Window),
		attrs: make(map[string]interface{}), timeout: 30 * time.Minute, rwMutex_: &sync.RWMutex{}}
}

// Number of valid id runes.
// Must be a power of 2!
const _ID_RUNES_COUNT = 64

// Mask to get an id rune idx from a random byte.
const _ID_RUNES_IDX_MASK = _ID_RUNES_COUNT - 1

// Valid runes to be used for session ids
// Its length must be _ID_RUNES_COUNT.
var _ID_RUNES = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_")

func init() {
	// Is _ID_RUNES_COUNT a power of 2?
	if _ID_RUNES_COUNT&(_ID_RUNES_COUNT-1) != 0 {
		panic(fmt.Sprint("_ID_RUNES_COUNT is not a power of 2: ", _ID_RUNES_COUNT))
	}
	if len(_ID_RUNES) != _ID_RUNES_COUNT {
		panic(fmt.Sprint("len(_ID_RUNES) != ", _ID_RUNES_COUNT))
	}
}

// Length of the session ids
const _ID_LENGTH = 22

// genId generates a session id.
func genId() string {
	r := make([]byte, _ID_LENGTH)
	io.ReadFull(rand.Reader, r)

	id := make([]rune, _ID_LENGTH)
	for i := 0; i < _ID_LENGTH; i++ {
		id[i] = _ID_RUNES[r[i]&_ID_RUNES_IDX_MASK]
	}

	return string(id)
}

func (s *sessionImpl) Id() string {
	return s.id
}

func (s *sessionImpl) New() bool {
	return s.isNew
}

func (s *sessionImpl) Private() bool {
	return len(s.id) > 0
}

func (s *sessionImpl) AddWin(w Window) error {
	if len(w.Name()) == 0 {
		return errors.New("Window name cannot be empty string!")
	}
	if _, exists := s.windows[w.Name()]; exists {
		return errors.New("A window with the same name has already been added: " + w.Name())
	}

	s.windows[w.Name()] = w

	return nil
}

func (s *sessionImpl) RemoveWin(w Window) bool {
	win := s.windows[w.Name()]
	if win != nil && win.Id() == w.Id() {
		delete(s.windows, w.Name())
		return true
	}
	return false
}

func (s *sessionImpl) SortedWins() []Window {
	wins := make(WinSlice, len(s.windows))

	i := 0
	for _, win := range s.windows {
		wins[i] = win
		i++
	}

	sort.Sort(wins)

	return wins
}

func (s *sessionImpl) WinByName(name string) Window {
	return s.windows[name]
}

func (s *sessionImpl) Attr(name string) interface{} {
	return s.attrs[name]
}

func (s *sessionImpl) SetAttr(name string, value interface{}) {
	if value == nil {
		delete(s.attrs, name)
	} else {
		s.attrs[name] = value
	}
}

func (s *sessionImpl) Created() time.Time {
	return s.created
}

func (s *sessionImpl) Accessed() time.Time {
	return s.accessed
}

func (s *sessionImpl) Timeout() time.Duration {
	return s.timeout
}

func (s *sessionImpl) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *sessionImpl) access() {
	s.accessed = time.Now()
}

func (s *sessionImpl) clearNew() {
	s.isNew = false
}

func (s *sessionImpl) rwMutex() *sync.RWMutex {
	return s.rwMutex_
}
