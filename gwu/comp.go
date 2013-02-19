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

// Comp component interface and implementation.

package gwu

import (
	"html"
	"net/http"
	"strconv"
)

// Container interface defines a component that can contain other components.
// Since a Container is a component itself, it can be added to
// other containers as well. The contained components are called
// the child components.
type Container interface {
	// Container is a component.
	Comp

	// Remove removes a component from this container.
	// Return value indicates if the specified component was a child
	// and was removed successfully.
	// After a successful Remove the specified component's
	// Parent() method will return nil.
	Remove(c Comp) bool

	// ById finds a component (recursively) by its ID and returns it.
	// nil is returned if no child component is found (recursively)
	// with the specified ID.
	ById(id ID) Comp

	// Clear clears the container, removes all child components.
	Clear()
}

// Comp interface: the base of all UI components.
type Comp interface {
	// Id returns the unique id of the component
	Id() ID

	// Equals tells if this component is equal to the specified another component.
	Equals(c2 Comp) bool

	// Parent returns the component's parent container.
	Parent() Container

	// setParent sets the component's parent container.
	setParent(parent Container)

	// makeOrphan makes this component orphan: if the component
	// has a parent, the component will be removed from the parent.
	// Return value indicates if the component was a child
	// and was removed successfully.
	makeOrphan() bool

	// Attr returns the explicitly set value of the specified HTML attribute.
	Attr(name string) string

	// SetAttr sets the value of the specified HTML attribute.
	// Pass an empty string value to delete the attribute.
	SetAttr(name, value string)

	// IAttr returns the explicitly set value of the specified HTML attribute
	// as an int.
	// -1 is returned if the value is not set explicitly or is not an int.
	IAttr(name string) int

	// SetAttr sets the value of the specified HTML attribute as an int.
	SetIAttr(name string, value int)

	// ToolTip returns the tool tip of the component.
	ToolTip() string

	// SetToolTip sets the tool tip of the component.
	SetToolTip(toolTip string)

	// Style returns the Style builder of the component.
	Style() Style

	// DescendantOf tells if this component is a descendant of the specified another component.
	DescendantOf(c2 Comp) bool

	// AddEHandler adds a new event handler.
	AddEHandler(handler EventHandler, etypes ...EventType)

	// AddEHandlerFunc adds a new event handler generated from a handler function.
	AddEHandlerFunc(hf func(e Event), etypes ...EventType)

	// HandlersCount returns the number of added handlers.
	HandlersCount(etype EventType) int

	// SyncOnETypes returns the event types on which to synchronize component value
	// from browser to the server.
	SyncOnETypes() []EventType

	// AddSyncOnETypes adds additional event types on which to synchronize
	// component value from browser to the server.
	AddSyncOnETypes(etypes ...EventType)

	// PreprocessEvent preprocesses an incoming event before it is dispatched.
	// This gives the opportunity for components to update their new value
	// before event handlers are called for example.
	preprocessEvent(event Event, r *http.Request)

	// DispatchEvent dispatches the event to all registered event handlers.
	dispatchEvent(e Event)

	// Render renders the component (as HTML code).
	Render(w writer)
}

// Comp implementation.
type compImpl struct {
	id     ID        // The component id
	parent Container // Parent container

	attrs     map[string]string // Explicitly set HTML attributes for the component's wrapper tag.
	styleImpl *styleImpl        // Style builder.

	handlers        map[EventType][]EventHandler // Event handlers mapped from event type. Lazily initialized.
	valueProviderJs []byte                       // If the HTML representation of the component has a value, this JavaScript code code must provide it. It will be automatically sent as the PARAM_COMP_ID parameter.
	syncOnETypes    map[EventType]bool           // Tells on which event types should comp value sync happen.
}

// newCompImpl creates a new compImpl.
// If the component has a value, the valueProviderJs must be a
// JavaScript code which when evaluated provides the component's
// value. Pass an empty string if the component does not have a value.
func newCompImpl(valueProviderJs []byte) compImpl {
	id := nextCompId()
	return compImpl{id: id, attrs: map[string]string{"id": id.String()}, styleImpl: newStyleImpl(), valueProviderJs: valueProviderJs}
}

func (c *compImpl) Id() ID {
	return c.id
}

func (c *compImpl) Equals(c2 Comp) bool {
	return c.id == c2.Id()
}

func (c *compImpl) Parent() Container {
	return c.parent
}

func (c *compImpl) setParent(parent Container) {
	c.parent = parent
}

func (c *compImpl) makeOrphan() bool {
	if c.parent == nil {
		return false
	}

	return c.parent.Remove(c)
}

func (c *compImpl) Attr(name string) string {
	return c.attrs[name]
}

func (c *compImpl) SetAttr(name, value string) {
	if len(value) > 0 {
		c.attrs[name] = value
	} else {
		delete(c.attrs, name)
	}
}

func (c *compImpl) IAttr(name string) int {
	if value, err := strconv.Atoi(c.Attr(name)); err == nil {
		return value
	}
	return -1
}

func (c *compImpl) SetIAttr(name string, value int) {
	c.SetAttr(name, strconv.Itoa(value))
}

func (c *compImpl) ToolTip() string {
	return html.UnescapeString(c.Attr("title"))
}

func (c *compImpl) SetToolTip(toolTip string) {
	c.SetAttr("title", html.EscapeString(toolTip))
}

func (c *compImpl) Style() Style {
	return c.styleImpl
}

func (c *compImpl) DescendantOf(c2 Comp) bool {
	for parent := c.parent; parent != nil; parent = parent.Parent() {
		// Always compare components by id, because Comp.Parent()
		// only returns Parent and not the components real type (e.g. windowImpl)!
		if parent.Equals(c2) {
			return true
		}
	}

	return false
}

// renderAttrs renders the explicitly set attributes and styles.
func (c *compImpl) renderAttrsAndStyle(w writer) {
	for name, value := range c.attrs {
		w.WriteAttr(name, value)
	}

	c.styleImpl.render(w)
}

func (c *compImpl) AddEHandler(handler EventHandler, etypes ...EventType) {
	if c.handlers == nil {
		c.handlers = make(map[EventType][]EventHandler)
	}
	for _, etype := range etypes {
		c.handlers[etype] = append(c.handlers[etype], handler)
	}
}

func (c *compImpl) AddEHandlerFunc(hf func(e Event), etypes ...EventType) {
	c.AddEHandler(handlerFuncWrapper{hf}, etypes...)
}

func (c *compImpl) HandlersCount(etype EventType) int {
	return len(c.handlers[etype])
}

func (c *compImpl) SyncOnETypes() []EventType {
	if c.syncOnETypes == nil {
		return nil
	}

	etypes := make([]EventType, len(c.syncOnETypes))
	i := 0
	for etype := range c.syncOnETypes {
		etypes[i] = etype
		i++
	}
	return etypes
}

func (c *compImpl) AddSyncOnETypes(etypes ...EventType) {
	if c.syncOnETypes == nil {
		c.syncOnETypes = make(map[EventType]bool, len(etypes))
	}
	for _, etype := range etypes {
		if !c.syncOnETypes[etype] { // If not yet synced...
			c.syncOnETypes[etype] = true
			c.AddEHandler(EMPTY_EHANDLER, etype)
		}
	}
}

var (
	_STR_SE_PREFIX = []byte("=\"se(event,") // "=\"se(event,"
	_STR_SE_SUFFIX = []byte(")\"")          // ")\""
)

// rendrenderEventHandlers renders the event handlers as attributes.
func (c *compImpl) renderEHandlers(w writer) {
	for etype, _ := range c.handlers {
		etypeAttr := etypeAttrs[etype]
		if len(etypeAttr) == 0 { // Only general events are added to the etypeAttrs map
			continue
		}

		// To render                 : " <etypeAttr>=\"se(event,etype,compId,value)\""
		// Example (checkbox onclick): " onclick=\"se(event,0,4327,this.checked)\""
		w.Write(_STR_SPACE)
		w.Write(etypeAttr)
		w.Write(_STR_SE_PREFIX)
		w.Writev(int(etype))
		w.Write(_STR_COMMA)
		w.Writev(int(c.id))
		if len(c.valueProviderJs) > 0 && c.syncOnETypes != nil && c.syncOnETypes[etype] {
			w.Write(_STR_COMMA)
			w.Write(c.valueProviderJs)
		}
		w.Write(_STR_SE_SUFFIX)
	}
}

// THIS IS AN EMPTY IMPLEMENTATION AS NOT ALL COMPONENTS NEED THIS.
// THOSE WHO DO SHOULD DEFINE THEIR OWN.
func (b *compImpl) preprocessEvent(event Event, r *http.Request) {
}

func (c *compImpl) dispatchEvent(e Event) {
	for _, handler := range c.handlers[e.Type()] {
		handler.HandleEvent(e)
	}
}

// THIS IS AN EMPTY IMPLEMENTATION.
// ALL COMPONENTS SHOULD DEFINE THEIR OWN
func (c *compImpl) Render(w writer) {
}
