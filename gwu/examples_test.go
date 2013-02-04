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

package gwu_test

import (
	"code/google/com/p/gowut/gwu"
)

// Example code determining which button was clicked. 
func ExampleButton() {
	b := gwu.NewButton("Click me")
	b.AddEHandlerFunc(func(e gwu.Event) {
		if e.MouseBtn() == gwu.MOUSE_BTN_MIDDLE {
			// Middle click
		}
	}, gwu.ETYPE_CLICK)
}

// Example code determining what kind of key is involved. 
func ExampleTextBox() {
	b := gwu.NewTextBox("")
	tb.AddSyncOnETypes(gwu.ETYPE_KEY_UP) // This is here so we will see up-to-date value in the event handler
	b.AddEHandlerFunc(func(e gwu.Event) {
		if e.ModKey(gwu.MOD_KEY_SHIFT) {
			// SHIFT is pressed
		}

		c := e.KeyCode()
		switch {
		case c == gwu.KEY_ENTER: // Enter
		case c >= gwu.KEY_0 && c <= gwu.KEY_9:
			fallthrough
		case c >= gwu.KEY_NUMPAD_0 && c <= gwuKEY_NUMPAD_9: // Number
		case c >= gwu.KEY_A && c <= gwu.KEY_Z: // Letter
		case c >= gwu.KEY_F1 && c <= gwu.KEY_F12: // Function key
		}
	}, gwu.ETYPE_KEY_UP)
}
