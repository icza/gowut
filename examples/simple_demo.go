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

// A GWU example application with a single public window (no sessions).

package main

import (
	"code.google.com/p/gowut/gwu"
	"strconv"
)

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
