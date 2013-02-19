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

// State button interfaces and implementations
// (CheckBox, RadioButton, SwitchButton).

package gwu

import (
	"net/http"
	"strconv"
)

// StateButton interface defines a button which has a boolean state:
// true/false or selected/deselected.
type StateButton interface {
	// stateButton is a button
	Button

	// State returns the state of the button.
	State() bool

	// SetState sets the state of the button.
	// In case of RadioButton, the button's RadioGroup is managed
	// so that only one can be selected.
	SetState(state bool)
}

// CheckBox interface defines a check box, a button which has
// 2 states: selected/deselected.
// 
// Suggested event type to handle changes: ETYPE_CLICK
// 
// Default style classes: "gwu-CheckBox", "gwu-CheckBox-Disabled"
type CheckBox interface {
	// CheckBox is a StateButton.
	StateButton
}

// SwitchButton interface defines a button which can be switched
// ON and OFF.
// 
// Suggested event type to handle changes: ETYPE_CLICK
// 
// Default style classes: "gwu-SwitchButton", "gwu-SwitchButton-On-Active"
// "gwu-SwitchButton-On-Inactive", "gwu-SwitchButton-Off-Active",
// "gwu-SwitchButton-Off-Inactive"
type SwitchButton interface {
	// SwitchButton is a component.
	Comp

	// SwitchButton can be enabled/disabled.
	HasEnabled

	// State returns the state of the switch button.
	State() bool

	// SetState sets the state of the switch button.
	SetState(state bool)

	// On returns the text displayed for the ON side.
	On() string

	// Off returns the text displayed for the OFF side.
	Off() string

	// SetOnOff sets the texts of the ON and OFF sides.
	SetOnOff(on, off string)
}

// RadioGroup interface defines the group for grouping radio buttons.
type RadioGroup interface {
	// Name returns the name of the radio group.
	Name() string

	// Selected returns the selected radio button of the group.
	Selected() RadioButton

	// PrevSelected returns the radio button that was selected
	// before the current selected radio button.
	PrevSelected() RadioButton

	// setSelected sets the selected radio button of the group,
	// and before that sets the current selected as the prev selected
	setSelected(selected RadioButton)
}

// RadioButton interface defines a radio button, a button which has
// 2 states: selected/deselected.
// In addition to the state, radio buttons belong to a group,
// and in each group only one radio button can be selected.
// Selecting an unselected radio button deselects the selected
// radio button of the group, if there was one.
// 
// Suggested event type to handle changes: ETYPE_CLICK
// 
// Default style classes: "gwu-RadioButton", "gwu-RadioButton-Disabled"
type RadioButton interface {
	// RadioButton is a StateButton.
	StateButton

	// Group returns the group of the radio button.
	Group() RadioGroup

	// setStateProp sets the state of the button
	// without managing the group of the radio button.
	setStateProp(state bool)
}

// RadioGroup implementation.
type radioGroupImpl struct {
	name         string      // Name of the radio group
	selected     RadioButton // Selected radio button of the group
	prevSelected RadioButton // Previous selected radio button of the group
}

// StateButton implementation.
type stateButtonImpl struct {
	buttonImpl // Button implementation 

	state         bool       // State of the button
	inputType     []byte     // Type of the underlying input tag
	group         RadioGroup // Group of the button
	inputId       ID         // distinct id for the rendered input tag
	disabledClass string     // Disabled style class
}

// SwitchButton implementation.
type switchButtonImpl struct {
	compImpl // Component implementation

	onButton, offButton *buttonImpl // ON and OFF button implementations
	state               bool        // State of the switch
}

// NewRadioGroup creates a new RadioGroup.
func NewRadioGroup(name string) RadioGroup {
	return &radioGroupImpl{name: name}
}

var (
	_STR_CHECKBOX     = []byte("checkbox")     // "checkbox"
	_STR_RADIO        = []byte("radio")        // "radio"
	_STR_THIS_CHECKED = []byte("this.checked") // "this.checked"
)

// NewCheckBox creates a new CheckBox.
// The initial state is false.
func NewCheckBox(text string) CheckBox {
	c := newStateButtonImpl(text, _STR_CHECKBOX, nil, "gwu-CheckBox-Disabled")
	c.Style().AddClass("gwu-CheckBox")
	return c
}

// NewSwitchButton creates a new SwitchButton.
// Default texts for ON and OFF sides are: "ON" and "OFF".
// The initial state is false (OFF).
func NewSwitchButton() SwitchButton {
	onButton := newButtonImpl(nil, "ON")
	offButton := newButtonImpl(nil, "OFF")

	// We only want to switch the state if the opposite button is pressed
	// (e.g. OFF is pressed when switch is ON and vice versa;
	// if ON is pressed when switch is ON, do not switch to OFF):
	valueProviderJs := []byte("sbtnVal(event,'" + onButton.Id().String() + "','" + offButton.Id().String() + "')")

	c := &switchButtonImpl{newCompImpl(valueProviderJs), &onButton, &offButton, true} // Note the "true" state, so the following SetState(false) will be executed (different states)!
	c.AddSyncOnETypes(ETYPE_CLICK)
	c.SetAttr("cellspacing", "0")
	c.SetAttr("cellpadding", "0")
	c.Style().AddClass("gwu-SwitchButton")
	c.SetState(false)
	return c
}

// NewRadioButton creates a new radio button.
// The initial state is false.
func NewRadioButton(text string, group RadioGroup) RadioButton {
	c := newStateButtonImpl(text, _STR_RADIO, group, "gwu-RadioButton-Disabled")
	c.Style().AddClass("gwu-RadioButton")
	return c
}

// newStateButtonImpl creates a new stateButtonImpl.
func newStateButtonImpl(text string, inputType []byte, group RadioGroup, disabledClass string) *stateButtonImpl {
	c := &stateButtonImpl{newButtonImpl(_STR_THIS_CHECKED, text), false, inputType, group, nextCompId(), disabledClass}
	// Use ETYPE_CLICK because IE fires onchange only when focus is lost...
	c.AddSyncOnETypes(ETYPE_CLICK)
	return c
}

func (r *radioGroupImpl) Name() string {
	return r.name
}

func (r *radioGroupImpl) Selected() RadioButton {
	return r.selected
}

func (r *radioGroupImpl) PrevSelected() RadioButton {
	return r.prevSelected
}

func (r *radioGroupImpl) setSelected(selected RadioButton) {
	r.prevSelected = r.selected
	r.selected = selected
}

// SetEnabled sets the enabled property.
// We have some extra job to do when changing enabled status:
// we have to manage disabled class style.
func (c *stateButtonImpl) SetEnabled(enabled bool) {
	if enabled {
		c.Style().RemoveClass(c.disabledClass)
	} else {
		c.Style().AddClass(c.disabledClass)
	}

	c.hasEnabledImpl.SetEnabled(enabled)
}

func (c *stateButtonImpl) State() bool {
	return c.state
}

func (c *stateButtonImpl) SetState(state bool) {
	// Only continue if state changes:
	if c.state == state {
		return
	}

	if c.group != nil {
		// Group management: if a new radio button is selected, the old one must be deselected.
		sel := c.group.Selected()

		if sel == nil {
			// no prev selection
			if state {
				c.group.setSelected(c)
			}
		} else {
			// There is a prev selection
			if state {
				if !sel.Equals(c) {
					sel.setStateProp(false)
					c.group.setSelected(c)
				}
			} else {
				// There is prev selection, and our new state is false
				// (and our prev state was true => we are selected)
				c.group.setSelected(nil)
			}
		}
	}

	c.state = state
}

func (c *stateButtonImpl) Group() RadioGroup {
	return c.group
}

func (c *stateButtonImpl) setStateProp(state bool) {
	c.state = state
}

func (c *stateButtonImpl) preprocessEvent(event Event, r *http.Request) {
	value := r.FormValue(_PARAM_COMP_VALUE)
	if len(value) == 0 {
		return
	}

	if v, err := strconv.ParseBool(value); err == nil {
		// Call SetState instead of assigning to the state property
		// because SetState properly manages radio groups.
		c.SetState(v)
	}
}

var (
	_STR_INPUT     = []byte("<input type=\"")       // "<input type=\""
	_STR_ID        = []byte("\" id=\"")             // "\" id=\""
	_STR_NAME      = []byte(" name=\"")             // " name=\""
	_STR_CHECKED   = []byte(" checked=\"checked\"") // " checked=\"checked\""
	_STR_LABEL_FOR = []byte("><label for=\"")       // "><label for=\""
	_STR_LABEL_CL  = []byte("</label>")             // "</label>"
)

func (c *stateButtonImpl) Render(w writer) {
	// Proper state button consists of multiple HTML tags (input and label), so render a wrapper tag for them:
	w.Write(_STR_SPAN_OP)
	c.renderAttrsAndStyle(w)
	w.Write(_STR_GT)

	w.Write(_STR_INPUT)
	w.Write(c.inputType)
	w.Write(_STR_ID)
	w.Writev(int(c.inputId))
	w.Write(_STR_QUOTE)
	if c.group != nil {
		w.Write(_STR_NAME)
		w.Writes(c.group.Name())
		w.Write(_STR_QUOTE)
	}
	if c.state {
		w.Write(_STR_CHECKED)
	}
	c.renderEnabled(w)
	c.renderEHandlers(w)

	w.Write(_STR_LABEL_FOR)
	w.Writev(int(c.inputId))
	w.Write(_STR_QUOTE)
	// TODO readding click handler here causes double event sending...
	// But we might add mouseover and other handlers still...
	//c.renderEHandlers(w)
	w.Write(_STR_GT)
	c.renderText(w)
	w.Write(_STR_LABEL_CL)
	w.Write(_STR_SPAN_CL)
}

func (c *switchButtonImpl) Enabled() bool {
	return c.onButton.Enabled()
}

func (c *switchButtonImpl) SetEnabled(enabled bool) {
	c.onButton.SetEnabled(enabled)
	c.offButton.SetEnabled(enabled)
}

func (c *switchButtonImpl) State() bool {
	return c.state
}

func (c *switchButtonImpl) SetState(state bool) {
	// Only continue if state changes:
	if c.state == state {
		return
	}

	c.state = state

	if c.state {
		c.onButton.Style().SetClass("gwu-SwitchButton-On-Active")
		c.offButton.Style().SetClass("gwu-SwitchButton-Off-Inactive")
	} else {
		c.onButton.Style().SetClass("gwu-SwitchButton-On-Inactive")
		c.offButton.Style().SetClass("gwu-SwitchButton-Off-Active")
	}
}

func (c *switchButtonImpl) On() string {
	return c.onButton.Text()
}
func (c *switchButtonImpl) Off() string {
	return c.offButton.Text()
}

func (c *switchButtonImpl) SetOnOff(on, off string) {
	c.onButton.SetText(on)
	c.offButton.SetText(off)
}

func (c *switchButtonImpl) preprocessEvent(event Event, r *http.Request) {
	value := r.FormValue(_PARAM_COMP_VALUE)
	if len(value) == 0 {
		return
	}

	if v, err := strconv.ParseBool(value); err == nil {
		// Call SetState instead of assigning to the state property
		// because SetState properly changes style classes.
		c.SetState(v)
		// SwitchButtons' client code properly updates internal buttons' style,
		// so we're good not to mark the switch button dirty if state changes.
	}
}

var (
	_STR_CL_TR = []byte("><tr>")              // "><tr>"
	_STR_TD_50 = []byte("<td width=\"50%\">") // "<td width=\"50%\">"
)

func (c *switchButtonImpl) Render(w writer) {
	w.Write(_STR_TABLE_OP)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	// For Internet Explorer only:
	// Since state synchronization is done on ETYPE_CLICK, which will add a click handler
	// to the wrapper tag and not to the on/off buttons, the wrapper tag itself must be
	// disabled (must have a 'disabled' attribute) if the switch button is disabled in order
	// for clicks really be disabled.
	c.onButton.renderEnabled(w)
	w.Write(_STR_CL_TR)

	w.Write(_STR_TD_50)
	c.onButton.Render(w)

	w.Write(_STR_TD_50)
	c.offButton.Render(w)

	w.Write(_STR_TABLE_CL)
}
