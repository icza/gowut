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

// Implementation of the GUI server which handles sessions,
// renders the windows and handles event dispatching.

package gwu

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Internal path constants.
const (
	_PATH_STATIC      = "_gwu_static/" // App path-relative path for GWU static contents.
	_PATH_EVENT       = "e"            // Window-relative path for sending events 
	_PATH_RENDER_COMP = "rc"           // Window-relative path for rendering a component 
)

// Parameters passed between the browser and the server.
const (
	_PARAM_EVENT_TYPE      = "et"   // Event type parameter name
	_PARAM_COMP_ID         = "cid"  // Component id parameter name
	_PARAM_COMP_VALUE      = "cval" // Component value parameter name
	_PARAM_FOCUSED_COMP_ID = "fcid" // Focused component id parameter name
	_PARAM_MOUSE_WX        = "mwx"  // Mouse x pixel coordinate (inside window)
	_PARAM_MOUSE_WY        = "mwy"  // Mouse y pixel coordinate (inside window)
	_PARAM_MOUSE_X         = "mx"   // Mouse x pixel coordinate (relative to source component)
	_PARAM_MOUSE_Y         = "my"   // Mouse y pixel coordinate (relative to source component)
	_PARAM_MOUSE_BTN       = "mb"   // Mouse button
	_PARAM_MOD_KEYS        = "mk"   // Modifier key states
	_PARAM_KEY_CODE        = "kc"   // Key code
)

// Event response actions (client actions to take after processing an event).
const (
	_ERA_NO_ACTION   = iota // Event processing OK and no action required 
	_ERA_RELOAD_WIN         // Window name to be reloaded
	_ERA_DIRTY_COMPS        // There are dirty components which needs to be refreshed
	_ERA_FOCUS_COMP         // Focus a compnent 
)

// GWU session id cookie name
const _GWU_SESSID_COOKIE = "gwu-sessid"

// SessionHandler interface defines a callback to get notified
// for certain events related to session life-cycles.
type SessionHandler interface {
	// Created is called when a new session is created.
	// At this time the client does not yet know about the session.
	Created(sess Session)

	// Removed is called when a session is being removed
	// from the server. After removal, the session id will become
	// an invalid session id.
	Removed(sess Session)
}

// Server interface defines the GUI server which handles sessions,
// renders the windows, components and handles event dispatching.
type Server interface {
	// The Server implements the Session interface:
	// there is one public session which is shared between
	// the "sessionless" requests.
	// This is to maintain windows without a session.
	Session

	// A server has text which will be used as the title
	// of the server.
	HasText

	// Secure returns if the server is configured to run
	// in secure (HTTPS) mode or in HTTP mode.
	Secure() bool

	// AppUrl returns the application URL string.
	AppUrl() string

	// AppPath returns the application path string.
	AppPath() string

	// AddSessCreatorName registers a nonexistent window name
	// whose path auto-creates a new session.
	// 
	// Normally sessions are created from event handlers during
	// event dispatching by calling Event.NewSession(). This
	// requires a public window and an event source component
	// (e.g. a Button) to create a session.
	// With AddSessCreatorName you can register nonexistent (meaning
	// not-yet added) window names whose path will trigger an automatic
	// session creation (if the current session is not private), and
	// with a registered SessionHandler you can build the window and
	// add it to the auto-created new session prior to it being served.
	// 
	// The text linking to the name will be included in the window list
	// if text is a non-empty string. 
	// 
	// Tip: You can use this to pre-register a login window for example.
	// You can call
	// 		AddSessCreatorName("login", "Login Window")
	// and in the Created() method of a registered SessionHandler:
	//		func (h MySessHanlder) Created(s gwu.Session) {
	//			win := gwu.NewWindow("login", "Login Window")
	//			// ...add content to the login window...
	// 			s.AddWindow(win)
	// 		}
	AddSessCreatorName(name, text string)

	// AddSHandler adds a new session handler.
	AddSHandler(handler SessionHandler)

	// AddStaticDir registers a directory whose content (files) recursively
	// will be served by the server when requested.
	// path is an app-path relative path to address a file, dir is the root directory
	// to search in.
	// 
	// Example:
	//     AddStaticDir("img", "/tmp/myimg")
	// And then the request "/appname/img/faces/happy.gif" will serve "/tmp/myimg/faces/happy.gif".
	// Note that the app name must be included in the request path!
	AddStaticDir(path, dir string) error

	// Theme returns the default CSS theme of the server.
	Theme() string

	// SetTheme sets the default CSS theme of the server.
	SetTheme(theme string)

	// SetLogger sets the logger to be used
	// to log incoming requests.
	// Pass nil to disable logging. This is the default.
	SetLogger(logger *log.Logger)

	// Start starts the GUI server and waits for incoming connections.
	// 
	// Sessionless window names may be specified as optional parameters
	// that will be opened in the default browser.
	// Tip: Pass an empty string to open the window list.
	// Tip: Not passing any window names will start the server silently
	// without opening any windows.
	Start(openWins ...string) error
}

// Server implementation.
type serverImpl struct {
	sessionImpl // Single public session implementation
	hasTextImpl // Has text implementation

	appName           string             // Application name (part of the application path)
	addr              string             // Server address
	secure            bool               // Tells if the server is configured to run in secure (HTTPS) mode
	appPath           string             // Application path
	appUrl            string             // Application URL
	sessions          map[string]Session // Sessions
	certFile, keyFile string             // Certificate and key files for secure (HTTPS) mode
	sessCreatorNames  map[string]string  // Session creator names
	sessionHandlers   []SessionHandler   // Registered session handlers
	theme             string             // Default CSS theme of the server
	logger            *log.Logger        // Logger.
}

// NewServer creates a new GUI server in HTTP mode.
// The specified app name will be part of the application path (the first part).
// If addr is empty string, "localhost:3434" will be used.
// 
// Tip: Pass an empty string as appName to place the GUI server to the root path ("/").
func NewServer(appName, addr string) Server {
	return newServerImpl(appName, addr, "", "")
}

// NewServerTLS creates a new GUI server in secure (HTTPS) mode.
// The specified app name will be part of the application path (the first part).
// If addr is empty string, "localhost:3434" will be used.
// 
// Tip: Pass an empty string as appName to place the GUI server to the root path ("/").
// Tip: You can use generate_cert.go in crypto/tls to generate
// a test certificate and key file (cert.pem andkey.pem). 
func NewServerTLS(appName, addr, certFile, keyFile string) Server {
	return newServerImpl(appName, addr, certFile, keyFile)
}

// newServerImpl creates a new serverImpl.
func newServerImpl(appName, addr, certFile, keyFile string) *serverImpl {
	if len(addr) == 0 {
		addr = "localhost:3434"
	}

	s := &serverImpl{sessionImpl: newSessionImpl(false), appName: appName, addr: addr, sessions: make(map[string]Session),
		sessCreatorNames: make(map[string]string), theme: THEME_DEFAULT}

	if len(s.appName) == 0 {
		s.appPath = "/"
	} else {
		s.appPath = "/" + s.appName + "/"
	}

	if len(certFile) == 0 || len(keyFile) == 0 {
		s.secure = false
		s.appUrl = "http://" + addr + s.appPath
	} else {
		s.secure = true
		s.appUrl = "https://" + addr + s.appPath
		s.certFile = certFile
		s.keyFile = keyFile
	}

	return s
}

func (s *serverImpl) Secure() bool {
	return s.secure
}

func (s *serverImpl) AppUrl() string {
	return s.appUrl
}

func (s *serverImpl) AppPath() string {
	return s.appPath
}

func (s *serverImpl) AddSessCreatorName(name, text string) {
	if len(name) > 0 {
		s.sessCreatorNames[name] = text
	}
}

func (s *serverImpl) AddSHandler(handler SessionHandler) {
	s.sessionHandlers = append(s.sessionHandlers, handler)
}

// newSession creates a new (private) Session.
// The event is optional. If specified and the current session
// (as returned by Event.Session()) is private, it will be removed first.
// The new session is set to the event, and also returned.
func (s *serverImpl) newSession(e *eventImpl) Session {
	if e != nil {
		// First remove old session
		s.removeSess(e)
	}

	sessImpl := newSessionImpl(true)
	sess := &sessImpl
	if e != nil {
		e.shared.session = sess
	}
	// Store new session
	s.sessions[sess.Id()] = sess

	if s.logger != nil {
		s.logger.Println("SESSION created:", sess.Id())
	}

	// Notify session handlers
	for _, handler := range s.sessionHandlers {
		handler.Created(sess)
	}

	return sess
}

// removeSess removes (invalidates) the current session of the specified event.
// Only private sessions can be removed, calling this
// when the current session (as returned by Event.Session()) is public is a no-op.
// After this method Event.Session() will return the shared public session.
func (s *serverImpl) removeSess(e *eventImpl) {
	if e.shared.session.Private() {
		s.removeSess2(e.shared.session)
		e.shared.session = &s.sessionImpl
	}
}

// removeSess2 removes (invalidates) the specified session.
// Only private sessions can be removed, calling this
// the public session is a no-op.
func (s *serverImpl) removeSess2(sess Session) {
	if sess.Private() {
		if s.logger != nil {
			s.logger.Println("SESSION removed:", sess.Id())
		}

		// Notify session handlers
		for _, handler := range s.sessionHandlers {
			handler.Removed(sess)
		}
		delete(s.sessions, sess.Id())
	}
}

// addSessCookie lets the client know about the specified (new) session
// by setting the GWU session id cookie.
// Also clears the new flag of the session.
func (s *serverImpl) addSessCookie(sess Session, w http.ResponseWriter) {
	// HttpOnly: do not allow non-HTTP access to it (like javascript) to prevent stealing it...
	// Secure: only send it over HTTPS
	// MaxAge: to specify the max age of the cookie in seconds, else it's a session cookie and gets deleted after the browser is closed.
	c := http.Cookie{Name: _GWU_SESSID_COOKIE, Value: sess.Id(), Path: s.appPath, HttpOnly: true, Secure: s.secure,
		MaxAge: 72 * 60 * 60} // 72 hours max age
	http.SetCookie(w, &c)

	sess.clearNew()
}

// sessCleaner periodically checks whether private sessions has timed out
// in an endless loop. If a session has timed out, removes it.
// This method is to start as a new go routine.
func (s *serverImpl) sessCleaner() {
	sleep := 10 * time.Second
	for {
		now := time.Now()

		// TODO synchronization?
		for _, sess := range s.sessions {
			if now.Sub(sess.Accessed()) > sess.Timeout() {
				s.removeSess2(sess)
			}
		}

		time.Sleep(sleep)
	}
}

func (s *serverImpl) AddStaticDir(path, dir string) error {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if len(path) == 0 {
		return errors.New("path cannot be empty string!")
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	path = s.appPath + path

	if path == s.appPath+_PATH_STATIC {
		return errors.New("path cannot be '" + _PATH_STATIC + "' (reserved)!")
	}

	http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))

	return nil
}

func (s *serverImpl) Theme() string {
	return s.theme
}

func (s *serverImpl) SetTheme(theme string) {
	s.theme = theme
}

func (s *serverImpl) SetLogger(logger *log.Logger) {
	s.logger = logger
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xgd-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func (s *serverImpl) Start(openWins ...string) error {
	http.HandleFunc(s.appPath, func(w http.ResponseWriter, r *http.Request) {
		s.serveHTTP(w, r)
	})

	http.HandleFunc(s.appPath+_PATH_STATIC, func(w http.ResponseWriter, r *http.Request) {
		s.serveStatic(w, r)
	})

	fmt.Println("Starting GUI server on:", s.appUrl)
	if s.logger != nil {
		s.logger.Println("Starting GUI server on:", s.appUrl)
	}

	for _, winName := range openWins {
		open(s.appUrl + winName)
	}

	go s.sessCleaner()

	var err error
	if s.secure {
		err = http.ListenAndServeTLS(s.addr, s.certFile, s.keyFile, nil)
	} else {
		err = http.ListenAndServe(s.addr, nil)
	}

	if err != nil {
		return err
	}
	return nil
}

// serveStatic handles the static contents of GWU.
func (s *serverImpl) serveStatic(w http.ResponseWriter, r *http.Request) {
	// Parts example: "/appname/_gwu_static/gwu-0.8.0.js" => {"", "appname", "_gwu_static", "gwu-0.8.0.js"}
	parts := strings.Split(r.URL.Path, "/")

	if len(s.appName) == 0 {
		// No app name, gui server resides in root
		if len(parts) < 2 {
			// This should never happen. Path is always at least a slash ("/").
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string and _PATH_STATIC
		parts = parts[2:]
	} else {
		// We have app name
		if len(parts) < 3 {
			// Missing app name from path
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string, app name and _PATH_STATIC
		parts = parts[3:]
	}

	res := parts[0]
	if res == _RES_NAME_STATIC_JS {
		w.Header().Set("Expires", time.Now().Add(72*time.Hour).Format(http.TimeFormat)) // Set 72 hours caching
		w.Header().Set("Content-Type", "application/x-javascript; charset=utf-8")
		w.Write(staticJs)
		return
	}
	if strings.HasSuffix(res, ".css") {
		cssCode := staticCss[res]
		if cssCode != nil {
			w.Header().Set("Expires", time.Now().Add(72*time.Hour).Format(http.TimeFormat)) // Set 72 hours caching
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			w.Write(cssCode)
			return
		}
	}

	http.NotFound(w, r)
}

// serveHTTP handles the incoming requests.
// Renders of the URL-selected window,
// and also handles event dispatching.
func (s *serverImpl) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if s.logger != nil {
		s.logger.Println("Incoming: ", r.URL.Path)
	}

	// Check session
	var sess Session
	c, err := r.Cookie(_GWU_SESSID_COOKIE)
	if err == nil {
		sess = s.sessions[c.Value]
	}
	if sess == nil {
		sess = &s.sessionImpl
	}
	sess.access()

	// Parts example: "/appname/winname/e?et=0&cid=1" => {"", "appname", "winname", "e"}
	parts := strings.Split(r.URL.Path, "/")

	if len(s.appName) == 0 {
		// No app name, gui server resides in root
		if len(parts) < 1 {
			// This should never happen. Path is always at least a slash ("/").
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string
		parts = parts[1:]
	} else {
		// We have app name
		if len(parts) < 2 {
			// Missing app name from path
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string and the app name
		parts = parts[2:]
	}

	if len(parts) < 1 || len(parts[0]) == 0 {
		// Missing window name, render window list
		s.renderWinList(sess, w, r)
		return
	}

	winName := parts[0]

	win := sess.WinByName(winName)
	// If not found and we're on an authenticated session, try the public window list
	if win == nil && sess.Private() {
		win = s.WinByName(winName) // Server is a Session, the public session
		if win != nil {
			s.access()
		}
	}
	// If still not found and no private session, try the session creator names
	if win == nil && !sess.Private() {
		_, found := s.sessCreatorNames[winName]
		if found {
			sess = s.newSession(nil)
			s.addSessCookie(sess, w)
			// Search again in the new session as SessionHandlers may have added windows.
			win = sess.WinByName(winName)
		}
	}

	if win == nil {
		// Invalid window name, render an error message with a link to the window list
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		NewWriter(w).Writess("<html><body>Window for name <b>'", winName, "'</b> not found. See the <a href=\"", s.appPath, "\">Window list</a>.</body></html>")
		return
	}

	var path string
	if len(parts) >= 2 {
		path = parts[1]
	}

	rwMutex := sess.rwMutex()

	switch path {
	case _PATH_EVENT:
		rwMutex.Lock()
		defer rwMutex.Unlock()

		s.handleEvent(sess, win, w, r)
	case _PATH_RENDER_COMP:
		rwMutex.RLock()
		defer rwMutex.RUnlock()

		// Render just a component
		s.renderComp(win, w, r)
	default:
		rwMutex.RLock()
		defer rwMutex.RUnlock()

		// Render the whole window
		win.RenderWin(NewWriter(w), s)
	}
}

// renderWinList renders the window list of a session as HTML document with clickable links. 
func (s *serverImpl) renderWinList(sess Session, wr http.ResponseWriter, r *http.Request) {
	if s.logger != nil {
		s.logger.Println("\tRending windows list.")
	}
	wr.Header().Set("Content-Type", "text/html; charset=utf-8")

	w := NewWriter(wr)

	w.Writes("<html><head><meta http-equiv=\"content-type\" content=\"text/html; charset=UTF-8\"><title>")
	w.Writees(s.text)
	w.Writess(" - Window list</title></head><body><h2>")
	w.Writees(s.text)
	w.Writes(" - Window list</h2>")

	// Render both private and public session windows
	sessions := make([]Session, 1, 2)
	sessions[0] = sess
	if sess.Private() {
		sessions = append(sessions, &s.sessionImpl)
	} else {
		// No private session yet, render session creators:
		if len(s.sessCreatorNames) > 0 {
			w.Writes("Session creators:<ul>") // TODO needs a better name
			for name, text := range s.sessCreatorNames {
				w.Writess("<li><a href=\"", s.appPath, name, "\">", text, "</a>")
			}
			w.Writes("</ul>")
		}
	}

	for _, session := range sessions {
		if session.Private() {
			w.Writes("Authenticated windows:")
		} else {
			w.Writes("Public windows:")
		}
		w.Writes("<ul>")
		for _, win := range session.SortedWins() {
			w.Writess("<li><a href=\"", s.appPath, win.Name(), "\">", win.Text(), "</a>")
		}
		w.Writes("</ul>")
	}

	w.Writes("</body></html>")
}

// renderComp renders just a component. 
func (s *serverImpl) renderComp(win Window, w http.ResponseWriter, r *http.Request) {
	id, err := AtoID(r.FormValue(_PARAM_COMP_ID))
	if err != nil {
		http.Error(w, "Invalid component id!", http.StatusBadRequest)
		return
	}

	if s.logger != nil {
		s.logger.Println("\tRendering comp:", id)
	}

	comp := win.ById(id)
	if comp == nil {
		http.Error(w, fmt.Sprint("Component not found: ", id), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // We send it as text!
	comp.Render(NewWriter(w))
}

// handleEvent handles the event dispatching.
func (s *serverImpl) handleEvent(sess Session, win Window, wr http.ResponseWriter, r *http.Request) {
	focCompId, err := AtoID(r.FormValue(_PARAM_FOCUSED_COMP_ID))
	if err == nil {
		win.SetFocusedCompId(focCompId)
	}

	id, err := AtoID(r.FormValue(_PARAM_COMP_ID))
	if err != nil {
		http.Error(wr, "Invalid component id!", http.StatusBadRequest)
		return
	}

	comp := win.ById(id)
	if comp == nil {
		if s.logger != nil {
			s.logger.Println("\tComp not found:", id)
		}
		http.Error(wr, fmt.Sprint("Component not found: ", id), http.StatusBadRequest)
		return
	}

	etype := parseIntParam(r, _PARAM_EVENT_TYPE)
	if etype < 0 {
		http.Error(wr, "Invalid event type!", http.StatusBadRequest)
		return
	}
	if s.logger != nil {
		s.logger.Println("\tEvent from comp:", id, " event:", etype)
	}

	event := newEventImpl(EventType(etype), comp, s, sess)
	shared := event.shared

	event.x = parseIntParam(r, _PARAM_MOUSE_X)
	if event.x >= 0 {
		event.y = parseIntParam(r, _PARAM_MOUSE_Y)
		shared.wx = parseIntParam(r, _PARAM_MOUSE_WX)
		shared.wy = parseIntParam(r, _PARAM_MOUSE_WY)
		shared.mbtn = MouseBtn(parseIntParam(r, _PARAM_MOUSE_BTN))
	} else {
		event.y, shared.wx, shared.wy, shared.mbtn = -1, -1, -1, -1
	}

	shared.modKeys = parseIntParam(r, _PARAM_MOD_KEYS)
	shared.keyCode = Key(parseIntParam(r, _PARAM_KEY_CODE))

	comp.preprocessEvent(event, r)

	// Dispatch event...
	comp.dispatchEvent(event)

	// Check if a new session was created during event dispatching
	if shared.session.New() {
		s.addSessCookie(shared.session, wr)
	}

	// ...and send back the result
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8") // We send it as text
	w := NewWriter(wr)
	hasAction := false
	// If we reload, nothing else matters
	if shared.reload {
		hasAction = true
		w.Writevs(_ERA_RELOAD_WIN, _STR_COMMA, shared.reloadWin)
	} else {
		if len(shared.dirtyComps) > 0 {
			hasAction = true
			w.Writev(_ERA_DIRTY_COMPS)
			for id, _ := range shared.dirtyComps {
				w.Write(_STR_COMMA)
				w.Writev(int(id))
			}
		}
		if shared.focusedComp != nil {
			if hasAction {
				w.Write(_STR_SEMICOL)
			} else {
				hasAction = true
			}
			w.Writevs(_ERA_FOCUS_COMP, _STR_COMMA, int(shared.focusedComp.Id()))
			// Also register focusable comp at window
			win.SetFocusedCompId(shared.focusedComp.Id())
		}
	}
	if !hasAction {
		w.Writev(_ERA_NO_ACTION)
	}
}

// parseIntParam parses an int param.
// If error occurs, -1 will be returned. 
func parseIntParam(r *http.Request, paramName string) int {
	if num, err := strconv.Atoi(r.FormValue(paramName)); err == nil {
		return num
	}
	return -1
}
