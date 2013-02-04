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

// Built-in static CSS themes of GWU.

package gwu

// Built-in CSS themes.
const (
	THEME_DEFAULT = "default" // Default CSS theme
	THEME_DEBUG   = "debug"   // Debug CSS theme, useful for developing/debugging purposes. 
)

// resNameStaticCss returns the CSS resource name
// for the specified CSS theme.
func resNameStaticCss(theme string) string {
	// E.g. "gowut-default-0.8.0.css"
	return "gowut-" + theme + "-" + GOWUT_VERSION + ".css"
}

var staticCss map[string][]byte = make(map[string][]byte)

func init() {
	staticCss[resNameStaticCss(THEME_DEFAULT)] = []byte("" +
		`
body {font-family:Arial}

.gwu-Window {}

.gwu-Panel {}

.gwu-Table {}

.gwu-Label {}

.gwu-Link {}

.gwu-Image {}

.gwu-Button {}

.gwu-CheckBox {}

.gwu-RadioButton {}

.gwu-ListBox {}

.gwu-TextBox {}

.gwu-PasswBox {}

.gwu-Html {}

.gwu-SwitchButton {}
.gwu-SwitchButton-On-Active {background:#00a000; color:#d0ffd0}
.gwu-SwitchButton-Off-Active {background:#d03030; color:#ffd0d0}
.gwu-SwitchButton-On-Inactive, .gwu-SwitchButton-Off-Inactive {background:#606060; color:#909090}
.gwu-SwitchButton-On-Inactive:enabled, .gwu-SwitchButton-Off-Inactive:enabled {cursor:pointer}
.gwu-SwitchButton-On-Active, .gwu-SwitchButton-Off-Active, .gwu-SwitchButton-On-Inactive, .gwu-SwitchButton-Off-Inactive {margin:0px;border: 0px; width:100%}
.gwu-SwitchButton-On-Active:disabled, .gwu-SwitchButton-Off-Active:disabled, .gwu-SwitchButton-On-Inactive:disabled, .gwu-SwitchButton-Off-Inactive:disabled {color:black}

.gwu-TabBar {}
.gwu-TabBar-Top {border-bottom:5px solid #8080f8}
.gwu-TabBar-Bottom {border-top:5px solid #8080f8}
.gwu-TabBar-Left {border-right:5px solid #8080f8; padding-top:1px}
.gwu-TabBar-Right {border-left:5px solid #8080f8; padding-top:1px}
.gwu-TabBar-NotSelected {padding-left:5px; padding-right:5px; border:1px solid white  ; background:#c0c0ff; cursor:default; width:100%}
.gwu-TabBar-Selected    {padding-left:5px; padding-right:5px; border:1px solid #8080f8; background:#8080f8; cursor:default; width:100%}
.gwu-TabPanel {}
.gwu-TabPanel-Content {border:1px solid #8080f8}
`)

	staticCss[resNameStaticCss(THEME_DEBUG)] = []byte(string(staticCss[resNameStaticCss(THEME_DEFAULT)]) +
		`
.gwu-Window td {border:1px solid black}
.gwu-Table td {border:1px solid black}
.gwu-Panel td {border:1px solid black}
.gwu-TabPanel td {border:1px solid black}
`)
}
