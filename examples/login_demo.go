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

// A GWU example application with login window and session management.

package main

import (
	"code.google.com/p/gowut/gwu"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type MyButtonHandler struct {
	counter int
	text    string
}

func (h *MyButtonHandler) HandleEvent(e gwu.Event) {
	// Check if event source is a Button, just to be sure...
	// We add this handler to a button only, so this'll be always false.
	if b, isButton := e.Src().(gwu.Button); isButton {
		b.SetText(b.Text() + h.text)
		h.counter++
		b.SetToolTip("You've clicked " + strconv.Itoa(h.counter) + " times!")
		e.MarkDirty(b)
	}
}

type GreenHandler int

func (h *GreenHandler) HandleEvent(e gwu.Event) {
	var state bool
	src := e.Src()

	switch c := src.(type) {
	case gwu.CheckBox:
		state = c.State()
	case gwu.RadioButton:
		state = c.State()
	}

	if state {
		src.Style().SetBackground(gwu.CLR_GREEN)
	} else {
		src.Style().SetBackground("")
	}
	e.MarkDirty(src)
}

var greenHandler_ = GreenHandler(0)
var greenHandler = &greenHandler_

func buildPrivateWins(s gwu.Session) {
	// Create and build a window
	win := gwu.NewWindow("main", "Main Window")
	win.Style().SetFullWidth()
	win.SetCellPadding(2)

	p := gwu.NewPanel()
	p.SetLayout(gwu.LAYOUT_HORIZONTAL)
	p.SetCellPadding(2)
	p.Add(gwu.NewLabel("I'm a label! Try clicking on the button=>"))
	p.Add(gwu.NewLink("Google Home", "https://google.com"))
	img := gwu.NewImage("", "https://www.google.com/images/srpr/logo3w.png")
	img.Style().SetSize("25%", "25%")
	p.Add(img)
	win.Add(p)
	button := gwu.NewButton("Click me")
	button.AddEHandler(&MyButtonHandler{text: ":-)"}, gwu.ETYPE_CLICK)
	win.Add(button)
	extraBtns := gwu.NewPanel()
	extraBtns.SetLayout(gwu.LAYOUT_NATURAL)
	button.AddEHandlerFunc(func(e gwu.Event) {
		extraBtn := gwu.NewButton("Extra #" + strconv.Itoa(extraBtns.CompsCount()))
		extraBtn.AddEHandlerFunc(func(e gwu.Event) {
			extraBtn.Parent().Remove(extraBtn)
			e.MarkDirty(extraBtns)
		}, gwu.ETYPE_CLICK)
		extraBtns.Insert(extraBtn, 0)
		e.MarkDirty(extraBtns)
	}, gwu.ETYPE_CLICK)
	win.Add(extraBtns)

	p = gwu.NewPanel()
	p.SetLayout(gwu.LAYOUT_HORIZONTAL)
	p.SetCellPadding(2)
	p.Style().SetBorder2(1, gwu.BRD_STYLE_SOLID, gwu.CLR_BLACK)
	p.Add(gwu.NewLabel("A drop-down list being"))
	wideListBox := gwu.NewListBox([]string{"50", "100", "150", "200", "250"})
	wideListBox.Style().SetWidth("50")
	wideListBox.AddEHandlerFunc(func(e gwu.Event) {
		wideListBox.Style().SetWidth(wideListBox.SelectedValue() + "px")
		e.MarkDirty(wideListBox)
	}, gwu.ETYPE_CHANGE)
	p.Add(wideListBox)
	p.Add(gwu.NewLabel("pixel wide. And a multi-select list:"))
	listBox := gwu.NewListBox([]string{"First", "Second", "Third", "Forth", "Fifth", "Sixth"})
	listBox.SetMulti(true)
	listBox.SetRows(4)
	p.Add(listBox)
	countLabel := gwu.NewLabel("Selected count: 0")
	listBox.AddEHandlerFunc(func(e gwu.Event) {
		selCount := len(listBox.SelectedIndices())
		countLabel.SetText("Selected count: " + strconv.Itoa(selCount))
		e.MarkDirty(countLabel)
	}, gwu.ETYPE_CHANGE)
	p.Add(countLabel)
	win.Add(p)

	greenCheckBox := gwu.NewCheckBox("I'm a check box. When checked, I'm green!")
	greenCheckBox.AddEHandlerFunc(func(e gwu.Event) {
		if greenCheckBox.State() {
			greenCheckBox.Style().SetBackground(gwu.CLR_GREEN)
		} else {
			greenCheckBox.Style().SetBackground("")
		}
		e.MarkDirty(greenCheckBox)
	}, gwu.ETYPE_CLICK)
	greenCheckBox.AddEHandler(greenHandler, gwu.ETYPE_CLICK)
	win.Add(greenCheckBox)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.Style().SetBorder2(1, gwu.BRD_STYLE_SOLID, gwu.CLR_BLACK)
	table.EnsureSize(2, 4)
	table.Add(gwu.NewLabel("TAB-"), 0, 0)
	table.Add(gwu.NewLabel("LE"), 0, 1)
	table.Add(gwu.NewLabel("DE-"), 0, 2)
	table.Add(gwu.NewLabel("MO"), 0, 3)
	table.Add(gwu.NewLabel("Enter your name:"), 1, 0)
	tb := gwu.NewTextBox("")
	tb.AddSyncOnETypes(gwu.ETYPE_KEY_UP)
	table.Add(tb, 1, 1)
	table.Add(gwu.NewLabel("You entered:"), 1, 2)
	nameLabel := gwu.NewLabel("")
	nameLabel.Style().SetColor(gwu.CLR_RED)
	tb.AddEHandlerFunc(func(e gwu.Event) {
		nameLabel.SetText(tb.Text())
		e.MarkDirty(nameLabel)
	}, gwu.ETYPE_CHANGE, gwu.ETYPE_KEY_UP)
	table.Add(nameLabel, 1, 3)
	win.Add(table)

	table = gwu.NewTable()
	table.Style().SetBorder2(1, gwu.BRD_STYLE_SOLID, gwu.CLR_BLACK)
	table.SetAlign(gwu.HA_RIGHT, gwu.VA_TOP)
	table.EnsureSize(5, 5)
	for row := 0; row < 5; row++ {
		group := gwu.NewRadioGroup(strconv.Itoa(row))
		for col := 0; col < 5; col++ {
			radio := gwu.NewRadioButton("= "+strconv.Itoa(col)+" =", group)
			radio.AddEHandlerFunc(func(e gwu.Event) {
				radios := []gwu.RadioButton{radio, radio.Group().PrevSelected()}
				for _, radio := range radios {
					if radio != nil {
						if radio.State() {
							radio.Style().SetBackground(gwu.CLR_GREEN)
						} else {
							radio.Style().SetBackground("")
						}
						e.MarkDirty(radio)
					}
				}
			}, gwu.ETYPE_CLICK)
			table.Add(radio, row, col)
		}
	}
	table.SetColSpan(2, 1, 2)
	table.SetRowSpan(3, 1, 2)
	table.CellFmt(2, 2).Style().SetSizePx(150, 80)
	table.CellFmt(2, 2).SetAlign(gwu.HA_RIGHT, gwu.VA_BOTTOM)
	table.RowFmt(2).Style().SetBackground("#808080")
	table.RowFmt(2).SetAlign(gwu.HA_DEFAULT, gwu.VA_MIDDLE)
	table.RowFmt(3).Style().SetBackground("#d0d0d0")
	table.RowFmt(4).Style().SetBackground("#b0b0b0")
	win.Add(table)

	tabPanel := gwu.NewTabPanel()
	tabPanel.SetTabBarPlacement(gwu.TB_PLACEMENT_TOP)
	for i := 0; i < 6; i++ {
		if i == 3 {
			img := gwu.NewImage("", "https://www.google.com/images/srpr/logo3w.png")
			img.Style().SetWidthPx(100)
			tabPanel.Add(img, gwu.NewLabel("This is some long content, random="+strconv.Itoa(rand.Int())))
			continue
		}
		tabPanel.AddString(strconv.Itoa(i)+". tab", gwu.NewLabel("This is some long content, random="+strconv.Itoa(rand.Int())))
	}
	win.Add(tabPanel)
	tabPanel = gwu.NewTabPanel()
	tabPanel.SetTabBarPlacement(gwu.TB_PLACEMENT_LEFT)
	tabPanel.TabBarFmt().SetVAlign(gwu.VA_BOTTOM)
	for i := 7; i < 11; i++ {
		l := gwu.NewLabel("This is some long content, random=" + strconv.Itoa(rand.Int()))
		if i == 9 {
			img := gwu.NewImage("", "https://www.google.com/images/srpr/logo3w.png")
			img.Style().SetWidthPx(100)
			tabPanel.Add(img, l)
			tabPanel.CellFmt(l).Style().SetSizePx(400, 400)
			continue
		}
		tabPanel.AddString(strconv.Itoa(i)+". tab", l)
		tabPanel.CellFmt(l).Style().SetSizePx(400, 400)
	}
	win.Add(tabPanel)
	s.AddWin(win)

	win2 := gwu.NewWindow("main2", "Main2 Window")
	win2.Add(gwu.NewLabel("This is just a test 2nd window."))
	back := gwu.NewButton("Back")
	back.AddEHandlerFunc(func(e gwu.Event) {
		e.ReloadWin(win.Name())
	}, gwu.ETYPE_CLICK)
	win2.Add(back)
	s.AddWin(win2)
}

func buildLoginWin(s gwu.Session) {
	win := gwu.NewWindow("login", "Login Window")
	win.Style().SetFullSize()
	win.SetAlign(gwu.HA_CENTER, gwu.VA_MIDDLE)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HA_CENTER)
	p.SetCellPadding(2)

	l := gwu.NewLabel("Test GUI Login Window")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("150%")
	p.Add(l)
	l = gwu.NewLabel("Login")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("130%")
	p.Add(l)
	p.CellFmt(l).Style().SetBorder2(1, gwu.BRD_STYLE_DASHED, gwu.CLR_NAVY)
	l = gwu.NewLabel("user/pass: admin/a")
	l.Style().SetFontSize("80%").SetFontStyle(gwu.FONT_STYLE_ITALIC)
	p.Add(l)

	errL := gwu.NewLabel("")
	errL.Style().SetColor(gwu.CLR_RED)
	p.Add(errL)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(2, 2)
	table.Add(gwu.NewLabel("User name:"), 0, 0)
	tb := gwu.NewTextBox("")
	tb.Style().SetWidthPx(160)
	table.Add(tb, 0, 1)
	table.Add(gwu.NewLabel("Password:"), 1, 0)
	pb := gwu.NewPasswBox("")
	pb.Style().SetWidthPx(160)
	table.Add(pb, 1, 1)
	p.Add(table)
	b := gwu.NewButton("OK")
	b.AddEHandlerFunc(func(e gwu.Event) {
		if tb.Text() == "admin" && pb.Text() == "a" {
			e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
			buildPrivateWins(e.Session())
			e.ReloadWin("main")
		} else {
			e.SetFocusedComp(tb)
			errL.SetText("Invalid user name or password!")
			e.MarkDirty(errL)
		}
	}, gwu.ETYPE_CLICK)
	p.Add(b)
	l = gwu.NewLabel("")
	p.Add(l)
	p.CellFmt(l).Style().SetHeightPx(200)

	win.Add(p)
	win.SetFocusedCompId(tb.Id())

	p = gwu.NewPanel()
	p.SetLayout(gwu.LAYOUT_HORIZONTAL)
	p.SetCellPadding(2)
	p.Add(gwu.NewLabel("Here's an ON/OFF switch which enables/disables the other one:"))
	sw := gwu.NewSwitchButton()
	sw.SetOnOff("ENB", "DISB")
	sw.SetState(true)
	p.Add(sw)
	p.Add(gwu.NewLabel("And the other one:"))
	sw2 := gwu.NewSwitchButton()
	sw2.SetEnabled(true)
	sw2.Style().SetWidthPx(100)
	p.Add(sw2)
	sw.AddEHandlerFunc(func(e gwu.Event) {
		sw2.SetEnabled(sw.State())
		e.MarkDirty(sw2)
	}, gwu.ETYPE_CLICK)
	win.Add(p)

	s.AddWin(win)
}

type SessHandler struct{}

func (h SessHandler) Created(s gwu.Session) {
	fmt.Println("SESSION created:", s.Id())
	buildLoginWin(s)
}

func (h SessHandler) Removed(s gwu.Session) {
	fmt.Println("SESSION removed:", s.Id())
}

func main() {
	// Create GUI server
	server := gwu.NewServer("guitest", "")
	//server := gwu.NewServerTLS("guitest", "", "test_tls/cert.pem", "test_tls/key.pem")
	server.SetText("Test GUI Application")

	server.AddSessCreatorName("login", "Login Window")
	server.AddSHandler(SessHandler{})

	win := gwu.NewWindow("home", "Home Window")
	l := gwu.NewLabel("Home, sweet home of " + server.Text())
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("130%")
	win.Add(l)
	win.Add(gwu.NewLabel("Click on the button to login:"))
	b := gwu.NewButton("Login")
	b.AddEHandlerFunc(func(e gwu.Event) {
		e.ReloadWin("login")
	}, gwu.ETYPE_CLICK)
	win.Add(b)

	server.AddWin(win)

	server.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	// Start GUI server
	if err := server.Start(); err != nil {
		fmt.Println("Error: Cound not start GUI server:", err)
		return
	}
}
