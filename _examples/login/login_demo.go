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
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/icza/gowut/gwu"
)

type myButtonHandler struct {
	counter int
	text    string
}

func (h *myButtonHandler) HandleEvent(e gwu.Event) {
	// Check if event source is a Button, just to be sure...
	// We add this handler to a button only, so this'll be always false.
	if b, isButton := e.Src().(gwu.Button); isButton {
		b.SetText(b.Text() + h.text)
		h.counter++
		b.SetToolTip(fmt.Sprintf("You've clicked %d times!", h.counter))
		e.MarkDirty(b)
	}
}

type greenHandler struct{}

func (greenHandler) HandleEvent(e gwu.Event) {
	var state bool
	src := e.Src()

	switch c := src.(type) {
	case gwu.CheckBox:
		state = c.State()
	case gwu.RadioButton:
		state = c.State()
	}

	if state {
		src.Style().SetBackground(gwu.ClrGreen)
	} else {
		src.Style().SetBackground("")
	}
	e.MarkDirty(src)
}

func buildPrivateWins(s gwu.Session) {
	// Create and build a window
	win := gwu.NewWindow("main", "Main Window")
	win.Style().SetFullWidth()
	win.SetCellPadding(2)

	p := gwu.NewPanel()
	p.SetLayout(gwu.LayoutHorizontal)
	p.SetCellPadding(2)
	p.Add(gwu.NewLabel("I'm a label! Try clicking on the button=>"))
	p.Add(gwu.NewLink("Google Home", "https://google.com"))
	img := gwu.NewImage("", "https://www.google.com/images/srpr/logo3w.png")
	img.Style().SetSize("25%", "25%")
	p.Add(img)
	win.Add(p)
	button := gwu.NewButton("Click me")
	button.AddEHandler(&myButtonHandler{text: ":-)"}, gwu.ETypeClick)
	win.Add(button)
	extraBtns := gwu.NewPanel()
	extraBtns.SetLayout(gwu.LayoutNatural)
	button.AddEHandlerFunc(func(e gwu.Event) {
		extraBtn := gwu.NewButton(fmt.Sprintf("Extra #%d", extraBtns.CompsCount()))
		extraBtn.AddEHandlerFunc(func(e gwu.Event) {
			extraBtn.Parent().Remove(extraBtn)
			e.MarkDirty(extraBtns)
		}, gwu.ETypeClick)
		extraBtns.Insert(extraBtn, 0)
		e.MarkDirty(extraBtns)
	}, gwu.ETypeClick)
	win.Add(extraBtns)

	p = gwu.NewPanel()
	p.SetLayout(gwu.LayoutHorizontal)
	p.SetCellPadding(2)
	p.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	p.Add(gwu.NewLabel("A drop-down list being"))
	wideListBox := gwu.NewListBox([]string{"50", "100", "150", "200", "250"})
	wideListBox.Style().SetWidth("50")
	wideListBox.AddEHandlerFunc(func(e gwu.Event) {
		wideListBox.Style().SetWidth(wideListBox.SelectedValue() + "px")
		e.MarkDirty(wideListBox)
	}, gwu.ETypeChange)
	p.Add(wideListBox)
	p.Add(gwu.NewLabel("pixel wide. And a multi-select list:"))
	listBox := gwu.NewListBox([]string{"First", "Second", "Third", "Forth", "Fifth", "Sixth"})
	listBox.SetMulti(true)
	listBox.SetRows(4)
	p.Add(listBox)
	countLabel := gwu.NewLabel("Selected count: 0")
	listBox.AddEHandlerFunc(func(e gwu.Event) {
		selCount := len(listBox.SelectedIndices())
		countLabel.SetText(fmt.Sprintf("Selected count: %d", selCount))
		e.MarkDirty(countLabel)
	}, gwu.ETypeChange)
	p.Add(countLabel)
	win.Add(p)

	greenCheckBox := gwu.NewCheckBox("I'm a check box. When checked, I'm green!")
	greenCheckBox.AddEHandlerFunc(func(e gwu.Event) {
		if greenCheckBox.State() {
			greenCheckBox.Style().SetBackground(gwu.ClrGreen)
		} else {
			greenCheckBox.Style().SetBackground("")
		}
		e.MarkDirty(greenCheckBox)
	}, gwu.ETypeClick)
	greenCheckBox.AddEHandler(greenHandler{}, gwu.ETypeClick)
	win.Add(greenCheckBox)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	table.EnsureSize(2, 4)
	table.Add(gwu.NewLabel("TAB-"), 0, 0)
	table.Add(gwu.NewLabel("LE"), 0, 1)
	table.Add(gwu.NewLabel("DE-"), 0, 2)
	table.Add(gwu.NewLabel("MO"), 0, 3)
	table.Add(gwu.NewLabel("Enter your name:"), 1, 0)
	tb := gwu.NewTextBox("")
	tb.AddSyncOnETypes(gwu.ETypeKeyUp)
	table.Add(tb, 1, 1)
	table.Add(gwu.NewLabel("You entered:"), 1, 2)
	nameLabel := gwu.NewLabel("")
	nameLabel.Style().SetColor(gwu.ClrRed)
	tb.AddEHandlerFunc(func(e gwu.Event) {
		nameLabel.SetText(tb.Text())
		e.MarkDirty(nameLabel)
	}, gwu.ETypeChange, gwu.ETypeKeyUp)
	table.Add(nameLabel, 1, 3)
	win.Add(table)

	table = gwu.NewTable()
	table.Style().SetBorder2(1, gwu.BrdStyleSolid, gwu.ClrBlack)
	table.SetAlign(gwu.HARight, gwu.VATop)
	table.EnsureSize(5, 5)
	for row := 0; row < 5; row++ {
		group := gwu.NewRadioGroup(strconv.Itoa(row))
		for col := 0; col < 5; col++ {
			radio := gwu.NewRadioButton(fmt.Sprintf("= %d =", col), group)
			radio.AddEHandlerFunc(func(e gwu.Event) {
				radios := []gwu.RadioButton{radio, radio.Group().PrevSelected()}
				for _, radio := range radios {
					if radio != nil {
						if radio.State() {
							radio.Style().SetBackground(gwu.ClrGreen)
						} else {
							radio.Style().SetBackground("")
						}
						e.MarkDirty(radio)
					}
				}
			}, gwu.ETypeClick)
			table.Add(radio, row, col)
		}
	}
	table.SetColSpan(2, 1, 2)
	table.SetRowSpan(3, 1, 2)
	table.CellFmt(2, 2).Style().SetSizePx(150, 80)
	table.CellFmt(2, 2).SetAlign(gwu.HARight, gwu.VABottom)
	table.RowFmt(2).Style().SetBackground("#808080")
	table.RowFmt(2).SetAlign(gwu.HADefault, gwu.VAMiddle)
	table.RowFmt(3).Style().SetBackground("#d0d0d0")
	table.RowFmt(4).Style().SetBackground("#b0b0b0")
	win.Add(table)

	tabPanel := gwu.NewTabPanel()
	tabPanel.SetTabBarPlacement(gwu.TbPlacementTop)
	for i := 0; i < 6; i++ {
		if i == 3 {
			img := gwu.NewImage("", "https://www.google.com/images/srpr/logo3w.png")
			img.Style().SetWidthPx(100)
			tabPanel.Add(img, gwu.NewLabel(fmt.Sprintf("This is some long content, random=%d", rand.Int())))
			continue
		}
		tabPanel.AddString(fmt.Sprintf("%d. tab", i), gwu.NewLabel(fmt.Sprintf("This is some long content, random=%d", rand.Int())))
	}
	win.Add(tabPanel)
	tabPanel = gwu.NewTabPanel()
	tabPanel.SetTabBarPlacement(gwu.TbPlacementLeft)
	tabPanel.TabBarFmt().SetVAlign(gwu.VABottom)
	for i := 7; i < 11; i++ {
		l := gwu.NewLabel(fmt.Sprintf("This is some long content, random=%d", rand.Int()))
		if i == 9 {
			img := gwu.NewImage("", "https://www.google.com/images/srpr/logo3w.png")
			img.Style().SetWidthPx(100)
			tabPanel.Add(img, l)
			tabPanel.CellFmt(l).Style().SetSizePx(400, 400)
			continue
		}
		tabPanel.AddString(fmt.Sprintf("%d. tab", i), l)
		tabPanel.CellFmt(l).Style().SetSizePx(400, 400)
	}
	win.Add(tabPanel)
	s.AddWin(win)

	win2 := gwu.NewWindow("main2", "Main2 Window")
	win2.Add(gwu.NewLabel("This is just a test 2nd window."))
	back := gwu.NewButton("Back")
	back.AddEHandlerFunc(func(e gwu.Event) {
		e.ReloadWin(win.Name())
	}, gwu.ETypeClick)
	win2.Add(back)
	s.AddWin(win2)
}

func buildLoginWin(s gwu.Session) {
	win := gwu.NewWindow("login", "Login Window")
	win.Style().SetFullSize()
	win.SetAlign(gwu.HACenter, gwu.VAMiddle)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(2)

	l := gwu.NewLabel("Test GUI Login Window")
	l.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("150%")
	p.Add(l)
	l = gwu.NewLabel("Login")
	l.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("130%")
	p.Add(l)
	p.CellFmt(l).Style().SetBorder2(1, gwu.BrdStyleDashed, gwu.ClrNavy)
	l = gwu.NewLabel("user/pass: admin/a")
	l.Style().SetFontSize("80%").SetFontStyle(gwu.FontStyleItalic)
	p.Add(l)

	errL := gwu.NewLabel("")
	errL.Style().SetColor(gwu.ClrRed)
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
	}, gwu.ETypeClick)
	p.Add(b)
	l = gwu.NewLabel("")
	p.Add(l)
	p.CellFmt(l).Style().SetHeightPx(200)

	win.Add(p)
	win.SetFocusedCompID(tb.ID())

	p = gwu.NewPanel()
	p.SetLayout(gwu.LayoutHorizontal)
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
	}, gwu.ETypeClick)
	win.Add(p)

	s.AddWin(win)
}

type sessHandler struct{}

func (h sessHandler) Created(s gwu.Session) {
	fmt.Println("SESSION created:", s.ID())
	buildLoginWin(s)
}

func (h sessHandler) Removed(s gwu.Session) {
	fmt.Println("SESSION removed:", s.ID())
}

func main() {
	// Create GUI server
	//server := gwu.NewServer("guitest", "")
	folder := "test_tls/"
	server := gwu.NewServerTLS("guitest", "", folder+"cert.pem", folder+"key.pem")
	server.SetText("Test GUI Application")

	server.AddSessCreatorName("login", "Login Window")
	server.AddSHandler(sessHandler{})

	win := gwu.NewWindow("home", "Home Window")
	l := gwu.NewLabel("Home, sweet home of " + server.Text())
	l.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("130%")
	win.Add(l)
	win.Add(gwu.NewLabel("Click on the button to login:"))
	b := gwu.NewButton("Login")
	b.AddEHandlerFunc(func(e gwu.Event) {
		e.ReloadWin("login")
	}, gwu.ETypeClick)
	win.Add(b)

	server.AddWin(win)

	server.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	// Start GUI server
	if err := server.Start(); err != nil {
		fmt.Println("Error: Cound not start GUI server:", err)
		return
	}
}
