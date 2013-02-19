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

// Built-in static JavaScript codes of Gowut.

package gwu

import (
	"strconv"
)

// Static JavaScript resource name
const _RES_NAME_STATIC_JS = "gowut-" + GOWUT_VERSION + ".js"

// Static javascript code
var staticJs []byte

func init() {
	// Init staticJs
	staticJs = []byte("" +
		// Param consts
		"var _pEventType='" + _PARAM_EVENT_TYPE +
		"',_pCompId='" + _PARAM_COMP_ID +
		"',_pCompValue='" + _PARAM_COMP_VALUE +
		"',_pFocCompId='" + _PARAM_FOCUSED_COMP_ID +
		"',_pMouseWX='" + _PARAM_MOUSE_WX +
		"',_pMouseWY='" + _PARAM_MOUSE_WY +
		"',_pMouseX='" + _PARAM_MOUSE_X +
		"',_pMouseY='" + _PARAM_MOUSE_Y +
		"',_pMouseBtn='" + _PARAM_MOUSE_BTN +
		"',_pModKeys='" + _PARAM_MOD_KEYS +
		"',_pKeyCode='" + _PARAM_KEY_CODE +
		"';\n" +
		// Modifier key masks
		"var _modKeyAlt=" + strconv.Itoa(int(MOD_KEY_ALT)) +
		",_modKeyCtlr=" + strconv.Itoa(int(MOD_KEY_CTRL)) +
		",_modKeyMeta=" + strconv.Itoa(int(MOD_KEY_META)) +
		",_modKeyShift=" + strconv.Itoa(int(MOD_KEY_SHIFT)) +
		";\n" +
		// Event response action consts
		"var _eraNoAction=" + strconv.Itoa(_ERA_NO_ACTION) +
		",_eraReloadWin=" + strconv.Itoa(_ERA_RELOAD_WIN) +
		",_eraDirtyComps=" + strconv.Itoa(_ERA_DIRTY_COMPS) +
		",_eraFocusComp=" + strconv.Itoa(_ERA_FOCUS_COMP) +
		";" +
		`

function createXmlHttp() {
	if (window.XMLHttpRequest) // IE7+, Firefox, Chrome, Opera, Safari
		return xmlhttp=new XMLHttpRequest();
	else // IE6, IE5
		return xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
}

// Send event
function se(event, etype, compId, compValue) {
	var xmlhttp = createXmlHttp();
	
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == 4 && xmlhttp.status == 200)
			procEresp(xmlhttp);
	}
	
	xmlhttp.open("POST", _pathEvent, true); // asynch call
	xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	
	var data="";
	
	if (etype != null)
		data += "&" + _pEventType + "=" + etype;
	if (compId != null)
		data += "&" + _pCompId + "=" + compId;
	if (compValue != null)
		data += "&" + _pCompValue + "=" + compValue;
	if (document.activeElement.id != null)
		data += "&" + _pFocCompId + "=" + document.activeElement.id;
	
	if (event != null) {
		if (event.clientX != null) {
			// Mouse data
			var x = event.clientX, y = event.clientY;
			data += "&" + _pMouseWX + "=" + x;
			data += "&" + _pMouseWY + "=" + y;
			var parent = document.getElementById(compId);
			do {
				x -= parent.offsetLeft;
				y -= parent.offsetTop;
			} while (parent = parent.offsetParent);
			data += "&" + _pMouseX + "=" + x;
			data += "&" + _pMouseY + "=" + y;
			data += "&" + _pMouseBtn + "=" + (event.button < 4 ? event.button : 1); // IE8 and below uses 4 for middle btn
		}
		
		var modKeys;
		modKeys += event.altKey ? _modKeyAlt : 0;
		modKeys += event.ctlrKey ? _modKeyCtlr : 0;
		modKeys += event.metaKey ? _modKeyMeta : 0;
		modKeys += event.shiftKey ? _modKeyShift : 0;
		data += "&" + _pModKeys + "=" + modKeys;
		data += "&" + _pKeyCode + "=" + (event.which ? event.which : event.keyCode);
	}
	
	xmlhttp.send(data);
}

function procEresp(xmlhttp) {
	var actions = xmlhttp.responseText.split(";");
	
	if (actions.length == 0) {
		window.alert("No response received!");
		return;
	}
	for (var i = 0; i < actions.length; i++) {
		var n = actions[i].split(",");
		
		switch (parseInt(n[0])) {
		case _eraDirtyComps:
			for (var j = 1; j < n.length; j++)
				rerenderComp(n[j]);
			break;
		case _eraFocusComp:
			if (n.length > 1)
				focusComp(parseInt(n[1]))
			break;
		case _eraNoAction:
			break;
		case _eraReloadWin:
			if (n.length > 1 && n[1].length > 0)
				window.location.href = _pathApp + n[1];
			else
				window.location.reload(true); // force reload
			break;
		default:
			window.alert("Unknown response code:" + n[0]);
			break;
		}
	}
}

function rerenderComp(compId) {
	var e = document.getElementById(compId);
	if (!e) // Component removed or not visible (e.g. on inactive tab of TabPanel)
		return;
	
	var xmlhttp = createXmlHttp();
	
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
			// Remember focused comp which might be replaced here:
			var focusedCompId = document.activeElement.id;
			e.outerHTML = xmlhttp.responseText;
			focusComp(focusedCompId);
			
			// Inserted JS code is not executed automatically, do it manually:
			// Have to "re-get" element by compId!
			var scripts = document.getElementById(compId).getElementsByTagName("script");
			for (var i = 0; i < scripts.length; i++) {
				eval(scripts[i].innerText);
			}
		}
	}
	
	xmlhttp.open("POST", _pathRenderComp, false); // synch call (if async, browser specific DOM rendering errors may arise)
	xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	
	xmlhttp.send(_pCompId + "=" + compId);
}

// Get selected indices (of an HTML select)
function selIdxs(select) {
	var selected = "";
	
	for (var i = 0; i < select.options.length; i++)
		if(select.options[i].selected)
			selected += i + ",";
	
	return selected;
}

// Get and update switch button value
function sbtnVal(event, onBtnId, offBtnId) {
	var onBtn = document.getElementById(onBtnId);
	var offBtn = document.getElementById(offBtnId);
	
	if (onBtn == null)
		return false;
	
	var value = onBtn == document.elementFromPoint(event.clientX, event.clientY);
	if (value) {
		onBtn.className = "gwu-SwitchButton-On-Active";
		offBtn.className = "gwu-SwitchButton-Off-Inactive";
	} else {
		onBtn.className = "gwu-SwitchButton-On-Inactive";
		offBtn.className = "gwu-SwitchButton-Off-Active";
	}
	
	return value;
}

function focusComp(compId) {
	if (compId != null) {
		var e = document.getElementById(compId);
		if (e) // Else component removed or not visible (e.g. on inactive tab of TabPanel)
			e.focus();
	}
}

function addonload(func) {
	var oldonload = window.onload;
	if (typeof window.onload != 'function') {
		window.onload = func;
	} else {
		window.onload = function() {
			if (oldonload)
				oldonload();
			func();
		}
	}
}

function addonbeforeunload(func) {
	var oldonbeforeunload = window.onbeforeunload;
	if (typeof window.onbeforeunload != 'function') {
		window.onbeforeunload = func;
	} else {
		window.onbeforeunload = function() {
			if (oldonbeforeunload)
				oldonbeforeunload();
			func();
		}
	}
}

var timers = new Object();

function setupTimer(compId, etype, timeout, repeat, active, reset) {
	var timer = timers[compId];
	
	if (timer != null) {
		var changed = timer.timeout != timeout || timer.repeat != repeat || timer.reset != reset;
		if (!active || changed) {
			if (timer.repeat)
				clearInterval(timer.id);
			else
				clearTimeout(timer.id);
			timers[compId] = null;
		}
		if (!changed)
			return;
	}
	if (!active)
		return;
	
	// Create new timer
	timers[compId] = timer = new Object();
	timer.timeout = timeout;
	timer.repeat = repeat;
	timer.reset = reset;
	
	// Start the timer
	var js = "se(null," + etype + "," + compId + ");";
	if (timer.repeat)
		timer.id = setInterval(js, timeout);
	else
		timer.id = setTimeout(js, timeout);
}

// INITIALIZATION

addonload(function() {
	focusComp(_focCompId);
});
`)
}
