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

/*
Package gwu implements an easy to use, platform independent Web UI Toolkit
in pure Go.


For additional documentation, News and more please visit the home page:
https://sites.google.com/site/gowebuitoolkit/


Introduction

Gowut (Go Web UI Toolkit) is a full-featured, easy to use, platform independent
Web UI Toolkit written in pure Go, no platform dependent native code is linked
or called.

The usage of the Gowut is similar to Google's GWT and the Java Swing toolkit.
If you are familiar with those, you will get started very easily. The main
difference compared to GWT is that this solution does not compile into JavaScript
but remains and runs as Go code (on the server side). Remaining on the server
side means you don't have to hassle with asynchronous event handlers like in GWT,
you can write real synchronous event handlers (like in Java Swing).

You can use this toolkit and build user interfaces with writing
Go code only: you can assemble the client interface in Go, and write
event handlers in Go.
You may optionally spice it and make it more customized with some HTML
and CSS (also added from Go code), but that is not required.

The UI can be simply assembled hierarchically from containers
and components. Components can generate events which are dispatched
to event handlers - also written in pure Go.
If there is no component for an HTML tag you wish to use, you can
use the Html component to wrap your custom HTML code. Components also allow
you to specify custom HTML attributes that will be added for their
(wrapper) HTML tags.

Creating user interfaces using Gowut does not require you to think like that
the clients will view it and interact with it through a browser.
The "browser" layer is hidden by Gowut.
While styling the components is done through CSS (either by calling
the style builder's methods or passing direct CSS codes), think of it
like a way similar to formatting HTML tags with CSS.

The state of the components are stored on server side, in the memory.
This means that if a browser is closed and reopened, or you navigate
away and back, the same state will be rendered again. AJAX technology
is used to automatically synchronize component's state from browser
to server, and to dispatch events.
AJAX technology is used also to refresh some parts (components) that
change (during event handling) without having to reload the whole page
to see the changes.

To quickly test it and see it in action, run the "Showcase of Features"
application by typing: (assuming you're in the root of your GOPATH)

	go run src/code.google.com/p/gowut/examples/showcase.go


Features of Gowut

-A component library to assemble your user interfaces with

-A GUI server which serves browser clients

-Session management

-Automatic event handling and dispatching

-(CSS) Style builder to easily manipulate the style of components


Server and Events and Sessions

The package contains a GUI server which is responsible to serve GUI
clients which are standard browsers. The user interface can be viewed
from any browsers (including smart phones) which makes this a cross
platform solution.
Starting the GUI server with a non-local address gives you
the possibility to view the GUI from a remote computer.
The server can be configured to run in normal mode (HTTP) or in secure
mode (HTTPS).

The GUI server also has Session management. By default windows added to the
server are public windows, and shared between all users (clients). This
means if a user changes the content (e.g. enters a text into a text box),
that text will be visible to all other users. This is suitable for most
desktop applications.

Sessions can be created during event handling (by calling the
Event.NewSession() method), and windows added to the session will only be
visible to the client associated with the session. If other users request
the same window, a new instance of the window is to be created and added
to their sessions.

Event handling is possible via event handlers. An event handler is
an implementation of the EventHandler interface. Event handlers have to be
attached to the components which will be the source of the event. Event
handlers are registered to event types or kinds (EventType) such as click
event (ETYPE_CLICK), value change event (ETYPE_CHANGE), key up event
(ETYPE_KEY_UP) etc.

The HandleEvent method of an event handler gets an Event value which has
multiple purposes and functions. 1) The event contains the parameters
of the event (such as the event type, the event source component, mouse
position in case of a mouse event etc.). 2) The Event is an accessor to the
Session associated with the client the event is originating from. Through
the event an event handler may access the current Session, create a new
Session or may remove it (invalidate it). 3) The event is also used
to define actions to be executed (automatically by Gowut) after the event
handling (post-event actions). For example if the event handler changes
a component, the handler has to mark it dirty causing it to be re-rendered
in the client browser, or an event handler can change the focused component,
or reload another window.

Creating a session from an event handler during event dispatching requires
a public window and an event source component (e.g. a Button).
There is another handy way to create sessions. Sessions can also be created
automatically by requesting pre-registered paths, paths of not-yet existing
windows. When such a window is requested and no private session associated
with the client exists, a new session will be created. A registered
SessionHandler can be used then to create the window prior to it being served.
Here's an example how to do it:
	// A SessionHandler implementation:
	type MySessHandler struct {}
	func (h SessHandler) Created(s gwu.Session) {
		win := gwu.NewWindow("login", "Login Window")
		// ...add content to the login window...
		s.AddWindow(win)
	}
	func (h SessHandler) Removed(s gwu.Session) {}

	// And to auto-create sessions for the login window:
	server := gwu.NewServer("guitest","")
	server.AddSessCreatorName("login", "Login Window")
	server.AddSHandler(MySessHandler{})

Despite the use of sessions if you access the application remotely (e.g. not
from localhost), security is only guaranteed if you configure the server to run
in secure (HTTPS) mode.


Under the hood

User interfaces are generated HTML documents which communicate with the server
with AJAX calls. The GUI server is based on the web server integrated in Go.

When a Window is requested by its URL, the Window will render a complete HTML
document. The Window will recursively include its child components.
Components render themselves into HTML codes.
When a component generates an event, the page in the browser will make an
AJAX call sending the event to the server. The event will be passed to all the
appropriate event handlers. Event handlers can mark components dirty,
specifying that they may have changed and they must be re-rendered.
When all the event handlers are done, the ids of the dirty components are sent
back, and the browser will request only to render the dirty components,
with AJAX calls, and the results will replace the old component nodes in the
HTML DOM.

Since the clients are HTTP browsers, the GWU sessions are implemented and
function as HTTP sessions. Cookies are used to maintain the browser sessions.


Styling

Styling the components is done through CSS. You can do this from Go code by
calling the style builder's methods, or you can create external CSS files.

The Comp interface contains a Style() method which returns the style builder
of the component. The builder can be used to set/manipulate the style class names
of the component (e.g. SetClass(), AddClass(), RemoveClass() methods).
The builder also has get and set methods for the common CSS attributes, and the
GWU package contains many CSS constants for CSS attribute values. Many styling
can be achieved using the builder's built-in methods and constants resulting in
the Go code containing no direct CSS at all. You can use the general Get() and
Set() methods of the style builder to manipulate any style attributes which it
does not have predefined methods for.

Each Gowut component has its own CSS class derived from its name using the "gwu-"
prefix, for example the Button component has the default CSS class "gwu-Button".
Many components use multiple CSS classes for their internal structure. These
classes are listed in the documentation of the components.
Gowut has multiple built-in CSS themes. A CSS theme is basically the collection
of the style definitions of the style classes used by the components. You can
set the default theme with the Server.SetTheme() method. This will be used for
all windows. You can set themes individually for windows too, using the
Window.SetTheme() method.

You can create your own external CSS files where you can extend/override the
definitions of the built-in style classes. For example you can define the
"gwu-Button" style class to have red background, and the result will be that all
Buttons will have red background without having to change their style individually.


Component palette

Containers to group and lay out components:
	Expander  - shows and hides a content comp when clicking on the header comp
	(Link)    - allows only one optional child
	Panel     - it has configurable layout
	Table     - it is dynamic and flexible
	TabPanel  - for tabbed displaying components (only 1 is visible at a time)
	Window    - top of component hierarchy, it is an extension of the Panel

Input components to get data from users:
	CheckBox
	ListBox    (it's either a drop-down list or a multi-line/multi-select list box)
	TextBox    (it's either a one-line text box or a multi-line text area)
	PasswBox
	RadioButton
	SwitchButton

Other components:
	Button
	Html
	Image
	Label
	Link
	Timer


Full application example

Let a full example follow here which is a complete application.
It builds a simple window, adds components to it, registers event handlers which
modify the content and starts the GUI server.
Component modifications (including both individual components and component
structure) will be seen without page reload.
All written in Go.

Source of this application is available here:
http://code.google.com/p/gowut/source/browse/examples/simple_demo.go

	type MyButtonHandler struct {
		counter int
		text    string
	}

	func (h *MyButtonHandler) HandleEvent(e gwu.Event) {
		if b, isButton := e.Src().(gwu.Button); isButton {
			b.SetText(b.Text() + h.text)
			h.counter++
			b.SetToolTip("You've clicked " + strconv.Itoa(h.counter) + " times!")
			e.MarkDirty(b)
		}
	}

	func main() {
		// Create and build a window
		win := gwu.NewWindow("main", "Test GUI Window")
		win.Style().SetFullWidth()
		win.SetHAlign(gwu.HA_CENTER)
		win.SetCellPadding(2)

		// Button which changes window content
		win.Add(gwu.NewLabel("I'm a label! Try clicking on the button=>"))
		btn := gwu.NewButton("Click me")
		btn.AddEHandler(&MyButtonHandler{text: ":-)"}, gwu.ETYPE_CLICK)
		win.Add(btn)
		btnsPanel := gwu.NewNaturalPanel()
		btn.AddEHandlerFunc(func(e gwu.Event) {
			// Create and add a new button...
			newbtn := gwu.NewButton("Extra #" + strconv.Itoa(btnsPanel.CompsCount()))
			newbtn.AddEHandlerFunc(func(e gwu.Event) {
				btnsPanel.Remove(newbtn) // ...which removes itself when clicked
				e.MarkDirty(btnsPanel)
			}, gwu.ETYPE_CLICK)
			btnsPanel.Insert(newbtn, 0)
			e.MarkDirty(btnsPanel)
		}, gwu.ETYPE_CLICK)
		win.Add(btnsPanel)

		// ListBox examples
		p := gwu.NewHorizontalPanel()
		p.Style().SetBorder2(1, gwu.BRD_STYLE_SOLID, gwu.CLR_BLACK)
		p.SetCellPadding(2)
		p.Add(gwu.NewLabel("A drop-down list being"))
		widelb := gwu.NewListBox([]string{"50", "100", "150", "200", "250"})
		widelb.Style().SetWidth("50")
		widelb.AddEHandlerFunc(func(e gwu.Event) {
			widelb.Style().SetWidth(widelb.SelectedValue() + "px")
			e.MarkDirty(widelb)
		}, gwu.ETYPE_CHANGE)
		p.Add(widelb)
		p.Add(gwu.NewLabel("pixel wide. And a multi-select list:"))
		listBox := gwu.NewListBox([]string{"First", "Second", "Third", "Forth", "Fifth", "Sixth"})
		listBox.SetMulti(true)
		listBox.SetRows(4)
		p.Add(listBox)
		countLabel := gwu.NewLabel("Selected count: 0")
		listBox.AddEHandlerFunc(func(e gwu.Event) {
			countLabel.SetText("Selected count: " + strconv.Itoa(len(listBox.SelectedIndices())))
			e.MarkDirty(countLabel)
		}, gwu.ETYPE_CHANGE)
		p.Add(countLabel)
		win.Add(p)

		// Self-color changer check box
		greencb := gwu.NewCheckBox("I'm a check box. When checked, I'm green!")
		greencb.AddEHandlerFunc(func(e gwu.Event) {
			if greencb.State() {
				greencb.Style().SetBackground(gwu.CLR_GREEN)
			} else {
				greencb.Style().SetBackground("")
			}
			e.MarkDirty(greencb)
		}, gwu.ETYPE_CLICK)
		win.Add(greencb)

		// TextBox with echo
		p = gwu.NewHorizontalPanel()
		p.Add(gwu.NewLabel("Enter your name:"))
		tb := gwu.NewTextBox("")
		tb.AddSyncOnETypes(gwu.ETYPE_KEY_UP)
		p.Add(tb)
		p.Add(gwu.NewLabel("You entered:"))
		nameLabel := gwu.NewLabel("")
		nameLabel.Style().SetColor(gwu.CLR_RED)
		tb.AddEHandlerFunc(func(e gwu.Event) {
			nameLabel.SetText(tb.Text())
			e.MarkDirty(nameLabel)
		}, gwu.ETYPE_CHANGE, gwu.ETYPE_KEY_UP)
		p.Add(nameLabel)
		win.Add(p)

		// Create and start a GUI server (omitting error check)
		server := gwu.NewServer("guitest", "localhost:8081")
		server.SetText("Test GUI App")
		server.AddWin(win)
		server.Start("") // Also opens windows list in browser
	}

Now start the application and open the http://localhost:8081/guitest/main URL in your
browser to see the window. You can also try visiting http://localhost:8081/guitest/
which will render the available window list.
Test the components. Now close the browser and reopen the page. Gowut remembers
everything.


Limitations

1) Attaching onmouseover and onmouseout event handlers to a component and
changing (re-rendering) the same component causes some trouble (the browsers
generate multiple mouseover and mouseout events because the same HTML node is replaced
under the mouse cursor).

2) Attaching onmousedown and onmouseup event handlers to a check box and re-rendering it
prevents ETYPE_CHANGE handlers being called when clicking on it.


Closing

From the MVC point of view looking at a Go application using Gowut, the Go
components are the Model, the generated (and manipulated) HTML document in the
browser is the View and the Controller is integrated in both.

Gowut is ideal to create (cross platform) user interfaces for desktop
applications written in Go. It is also easy and handy to write the admin
and also client interfaces of your Go web application using Gowut.

Happy UI coding in Go :-)


Links

Author: András Belicza

Author email: gmail.com, user name: iczaaa

Home page: https://sites.google.com/site/gowebuitoolkit/

Source code: http://code.google.com/p/gowut/

Discussion forum: https://groups.google.com/d/forum/gowebuitoolkit

Live demo: Coming soon...


*/
package gwu

// Gowut version information.
const (
	GOWUT_VERSION         = "0.8.0"          // Gowut version (major.minor.maintenance)
	GOWUT_RELEASE_DATE    = "2013-02-19 CET" // Gowut release date
	GOWUT_REL_DATE_LAYOUT = "2006-01-02 MST" // Gowut release date layout (for time.Parse())
)
