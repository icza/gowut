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

// Window component interface and implementation.

package gwu

// The Window interface is the top of the component hierarchy.
// A Window defines the content seen in the browser window.
// Multiple windows can be created, but only one is visible
// at a time in the browser. The Window interface is the
// equivalent of the browser page.
//
// Default style class: "gwu-Window"
type Window interface {
	// Window is a Panel, child components can be added to it.
	Panel

	// A window has text which will be used as the title
	// of the browser window.
	HasText

	// Name returns the name of the window.
	// The name appears in the URL.
	Name() string

	// SetName sets the name of the window.
	SetName(name string)

	// AddHeadHtml adds an HTML text which will be included
	// in the HTML <head> section.
	AddHeadHtml(html string)

	// RemoveHeadHtml removes an HTML head text
	// that was previously added with AddHeadHtml().
	RemoveHeadHtml(html string)

	// SetFocusedCompId sets the id of the currently focused component.
	SetFocusedCompId(id ID)

	// Theme returns the CSS theme of the window.
	// If an empty string is returned, the server's theme will be used.
	Theme() string

	// SetTheme sets the default CSS theme of the window.
	// If an empty string is set, the server's theme will be used.
	SetTheme(theme string)

	// RenderWin renders the window as a complete HTML document.
	RenderWin(w Writer, s Server)
}

// WinSlice is a slice of windows which implements sort.Interface so it
// can be sorted by window text (title).
type WinSlice []Window

func (w WinSlice) Len() int {
	return len(w)
}

func (w WinSlice) Less(i, j int) bool {
	return w[i].Text() < w[j].Text()
}

func (w WinSlice) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Window implementation
type windowImpl struct {
	panelImpl   // Panel implementation
	hasTextImpl // Has text implementation

	name          string   // Window name
	heads         []string // Additional head HTML texts
	focusedCompId ID       // Id of the last reported focused component
	theme         string   // CSS theme of the window
}

// NewWindow creates a new window.
// The default layout strategy is LayoutVertical.
func NewWindow(name, text string) Window {
	c := &windowImpl{panelImpl: newPanelImpl(), hasTextImpl: newHasTextImpl(text), name: name}
	c.Style().AddClass("gwu-Window")
	return c
}

func (w *windowImpl) Name() string {
	return w.name
}

func (w *windowImpl) SetName(name string) {
	w.name = name
}

func (w *windowImpl) AddHeadHtml(html string) {
	w.heads = append(w.heads, html)
}

func (w *windowImpl) RemoveHeadHtml(html string) {
	for i, v := range w.heads {
		if v == html {
			old := w.heads
			w.heads = append(w.heads[:i], w.heads[i+1:]...)
			old[len(old)-1] = ""
			return
		}
	}
}

func (w *windowImpl) SetFocusedCompId(id ID) {
	w.focusedCompId = id
}

func (w *windowImpl) Theme() string {
	return w.theme
}

func (w *windowImpl) SetTheme(theme string) {
	w.theme = theme
}

func (w *windowImpl) Render(wr Writer) {
	// Attaching window events is outside of the HTML tag denoted by the window's id.
	// This means if the window is re-rendered (not reloaded), changed window event handlers
	// will not be reflected.
	// This also avoids the effect of registering the event sender functions multiple times.

	// First render window event handlers as window functions.
	found := false
	for etype := range w.handlers {
		if etype.Category() != ECatWindow {
			continue
		}

		if !found {
			found = true
			wr.Write(strScriptOp)
		}
		// To render       : add<etypeFunc>(function(){se(null,etype,id);});
		// Example (onload): addonload(function(){se(null,13,4327);});
		wr.Writevs("add", etypeFuncs[etype], "(function(){se(null,", int(etype), ",", int(w.id), ");});")
	}
	if found {
		wr.Write(strScriptCl)
	}

	// And now call panelImpl's Render()
	w.panelImpl.Render(wr)
}

func (w *windowImpl) RenderWin(wr Writer, s Server) {
	// We could optimize this (store byte slices of static strings)
	// but windows are rendered "so rarely"...
	wr.Writes(`<html><head><meta http-equiv="content-type" content="text/html; charset=UTF-8"><title>`)
	wr.Writees(w.text)
	wr.Writess(`</title><link href="`, s.AppPath(), pathStatic)
	if w.theme == "" {
		wr.Writes(resNameStaticCss(s.Theme()))
	} else {
		wr.Writes(resNameStaticCss(w.theme))
	}
	wr.Writes(`" rel="stylesheet" type="text/css">`)
	w.renderDynJs(wr, s)
	wr.Writess(`<script src="`, s.AppPath(), pathStatic, resNameStaticJs, `"></script>`)
	wr.Writess(w.heads...)
	wr.Writes("</head><body>")

	w.Render(wr)

	wr.Writes("</body></html>")
}

// renderDynJs renders the dynamic JavaScript codes of Gowut.
func (w *windowImpl) renderDynJs(wr Writer, s Server) {
	wr.Write(strScriptOp)
	wr.Writess("var _pathApp='", s.AppPath(), "';")
	wr.Writess("var _pathSessCheck=_pathApp+'", pathSessCheck, "';")
	wr.Writess("var _pathWin='", s.AppPath(), w.name, "/';")
	wr.Writess("var _pathEvent=_pathWin+'", pathEvent, "';")
	wr.Writess("var _pathRenderComp=_pathWin+'", pathRenderComp, "';")
	wr.Writess("var _focCompId='", w.focusedCompId.String(), "';")
	wr.Write(strScriptCl)
}
