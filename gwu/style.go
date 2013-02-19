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

// Defines the style builder and style attribute constants.

package gwu

import (
	"strconv"
)

// Style attribute constants.
const (
	ST_BACKGROUND     = "background"     // Background (color)
	ST_BORDER         = "border"         // Border
	ST_BORDER_LEFT    = "border-left"    // Left border
	ST_BORDER_RIGHT   = "border-right"   // Right border
	ST_BORDER_TOP     = "border-top"     // Top border
	ST_BORDER_BOTTOM  = "border-bottom"  // Bottom border
	ST_COLOR          = "color"          // (Foreground) color
	ST_CURSOR         = "cursor"         // Cursor
	ST_DISPLAY        = "display"        // Display
	ST_FONT_SIZE      = "font-size"      // Font size
	ST_FONT_STYLE     = "font-style"     // Font style
	ST_FONT_WEIGHT    = "font-weight"    // Font weight
	ST_HEIGHT         = "height"         // Height
	ST_MARGIN         = "margin"         // Margin
	ST_MARGIN_LEFT    = "margin-left"    // Left margin
	ST_MARGIN_RIGHT   = "margin-right"   // Right margin
	ST_MARGIN_TOP     = "margin-top"     // Top margin
	ST_MARGIN_BOTTOM  = "margin-bottom"  // Bottom margin
	ST_PADDING        = "padding"        // Padding
	ST_PADDING_LEFT   = "padding-left"   // Left padding
	ST_PADDING_RIGHT  = "padding-right"  // Right padding
	ST_PADDING_TOP    = "padding-top"    // Top padding
	ST_PADDING_BOTTOM = "padding-bottom" // Bottom padding
	ST_WHITE_SPACE    = "white-space"    // White-space
	ST_WIDTH          = "width"          // Width
)

// The 17 standard color constants.
const (
	CLR_AQUA    = "Aqua"    // Aqua    (#00FFFF)
	CLR_BLACK   = "Black"   // Black   (#000000)
	CLR_BLUE    = "Blue"    // Blue    (#0000FF)
	CLR_FUCHSIA = "Fuchsia" // Fuchsia (#FF00FF)
	CLR_GRAY    = "Gray"    // Gray    (#808080)
	CLR_GREY    = "Grey"    // Grey    (#808080)
	CLR_GREEN   = "Green"   // Green   (#008000)
	CLR_LIME    = "Lime"    // Lime    (#00FF00)
	CLR_MAROON  = "Maroon"  // Maroon  (#800000)
	CLR_NAVY    = "Navy"    // Navy    (#000080)
	CLR_OLIVE   = "Olive"   // Olive   (#808000)
	CLR_PURPLE  = "Purple"  // Purple  (#800080)
	CLR_RED     = "Red"     // Red     (#FF0000)
	CLR_SILVER  = "Silver"  // Silver  (#C0C0C0)
	CLR_TEAL    = "Teal"    // Teal    (#008080)
	CLR_WHITE   = "White"   // White   (#FFFFFF)
	CLR_YELLOW  = "Yellow"  // Yellow  (#FFFF00)
)

// Border style constants.
const (
	BRD_STYLE_SOLID  = "solid"  // Solid
	BRD_STYLE_DASHED = "dashed" // Dashed
	BRD_STYLE_DOTTED = "dotted" // Dotted
	BRD_STYLE_DOUBLE = "double" // Double
	BRD_STYLE_GROOVE = "groove" // 3D grooved border 
	BRD_STYLE_RIDGE  = "ridge"  // 3D ridged border
	BRD_STYLE_INSET  = "inset"  // 3D inset border
	BRD_STYLE_OUTSET = "outset" // 3D outset border
)

// Font weight constants.
const (
	FONT_WEIGHT_NORMAL  = "normal"  // Normal
	FONT_WEIGHT_BOLD    = "bold"    // Bold
	FONT_WEIGHT_BOLDER  = "bolder"  // Bolder
	FONT_WEIGHT_LIGHTER = "lighter" // Lighter
)

// Font style constants.
const (
	FONT_STYLE_NORMAL = "normal" // Normal
	FONT_STYLE_ITALIC = "italic" // Italic
)

// Mouse cursor constants.
const (
	CURSOR_AUTO      = "auto"      // Default. Web browser sets the cursor.
	CURSOR_CROSSHAIR = "crosshair" // Crosshair
	CURSOR_DEFAULT   = "default"   // The default cursor.
	CURSOR_HELP      = "help"      // Help
	CURSOR_MOVE      = "move"      // Move
	CURSOR_POINTER   = "pointer"   // Pointer
	CURSOR_PROGRESS  = "progress"  // Progress
	CURSOR_TEXT      = "text"      // Text
	CURSOR_WAIT      = "wait"      // Wait
	CURSOR_INHERIT   = "inherit"   // The cursor should be inherited from the parent element.
)

// Display mode constants.
const (
	DISPLAY_NONE    = "none"    // The element will not be displayed.
	DISPLAY_BLOCK   = "block"   // The element is displayed as a block.
	DISPLAY_INLINE  = "inline"  // The element is displayed as an in-line element. This is the default.
	DISPLAY_INHERIT = "inherit" // The display property value will be inherited from the parent element.
)

// White space constants.
const (
	WHITE_SPACE_NORMAL   = "normal"   // Sequences of white spaces are collapsed into a single whitespace. Text will wrap when neccessary. This is the default.
	WHITE_SPACE_NOWRAP   = "nowrap"   // Sequences of whitespace will collapse into a single whitespace. Text will never wrap to the next line (the text is in one line).
	WHITE_SPACE_PRE      = "pre"      // Whitespace is preserved. Text will only wrap on line breaks.
	WHITE_SPACE_PRE_LINE = "pre-line" // Sequences of whitespace will collapse into a single whitespace. Text will wrap when necessary and on line breaks.
	WHITE_SPACE_PRE_WRAP = "pre-wrap" // Whitespace is preserved. Text will wrap when necessary, and on line breaks.
	WHITE_SPACE_INHERIT  = "inherit"  // Whitespace property will be inherited from the parent element.
)

// Style interface contains utility methods for manipulating
// the style of a component.
// You can think of it as the Style Builder.
// Set methods return the style reference so setting the values
// of multiple style attributes can be chained.
type Style interface {
	// AddClass adds a style class name to the class name list.
	AddClass(class string) Style

	// SetClass sets a style class name, removing all previously
	// added style class names.
	// Tip: set an empty string class name to remove all class names.
	SetClass(class string) Style

	// RemoveClass removes a style class name.
	// If the specified class is not found, this is a no-op.
	RemoveClass(class string) Style

	// Get returns the explicitly set value of the specified style attribute.
	// Explicitly set style attributes will be concatenated and rendered
	// as the "style" HTML attribute of the component.
	Get(name string) string

	// Set sets the value of the specified style attribute.
	// Pass an empty string value to delete the specified style attribute.
	Set(name, value string) Style

	// Size returns the size.
	Size() (width, height string)

	// SetSize sets the width and height.
	SetSize(width, height string) Style

	// SetSizePx sets the width and height, in pixels.
	SetSizePx(width, height int) Style

	// SetFullSize sets full width (100%) and height (100%).
	SetFullSize() Style

	// Padding returns the padding.
	// (The "padding" style attribute only.)
	Padding() string

	// SetPadding sets the padding.
	// (The "padding" style attribute only.)
	SetPadding(value string) Style

	// SetPadding2 sets the padding specified by parts.
	// (The "padding" style attribute only.)
	SetPadding2(top, right, bottom, left string) Style

	// SetPaddingPx sets the padding specified by parts, in pixels.
	// (The "padding" style attribute only.)
	SetPaddingPx(top, right, bottom, left int) Style

	// PaddingLeft returns the left padding.
	// (The "padding-left" style attribute only.)
	PaddingLeft() string

	// SetPaddingLeft sets the left padding.
	// (The "padding-left" style attribute only.)
	SetPaddingLeft(value string) Style

	// SetPaddingLeftPx sets the left padding, in pixels.
	// (The "padding-left" style attribute only.)
	SetPaddingLeftPx(width int) Style

	// PaddingRight returns the right padding.
	// (The "padding-right" style attribute only.)
	PaddingRight() string

	// SetPaddingRight sets the right padding.
	// (The "padding-right" style attribute only.)
	SetPaddingRight(value string) Style

	// SetPaddingRightPx sets the right padding, in pixels.
	// (The "padding-right" style attribute only.)
	SetPaddingRightPx(width int) Style

	// PaddingTop returns the top padding.
	// (The "padding-top" style attribute only.)
	PaddingTop() string

	// SetPaddingTop sets the top padding.
	// (The "padding-top" style attribute only.)
	SetPaddingTop(value string) Style

	// SetPaddingTopPx sets the top padding, in pixels.
	// (The "padding-top" style attribute only.)
	SetPaddingTopPx(height int) Style

	// PaddingBottom returns the bottom padding.
	// (The "padding-bottom" style attribute only.)
	PaddingBottom() string

	// SetPaddingBottom sets the bottom padding.
	// (The "padding-bottom" style attribute only.)
	SetPaddingBottom(value string) Style

	// SetPaddingBottomPx sets the bottom padding, in pixels.
	// (The "padding-bottom" style attribute only.)
	SetPaddingBottomPx(height int) Style

	// Margin returns the margin.
	// (The "margin" style attribute only.)
	Margin() string

	// SetMargin sets the margin.
	// (The "margin" style attribute only.)
	SetMargin(value string) Style

	// SetMargin2 sets the margin specified by parts.
	// (The "margin" style attribute only.)
	SetMargin2(top, right, bottom, left string) Style

	// SetMarginPx sets the margin specified by parts, in pixels.
	// (The "margin" style attribute only.)
	SetMarginPx(top, right, bottom, left int) Style

	// MarginLeft returns the left margin.
	// (The "margin-left" style attribute only.)
	MarginLeft() string

	// SetMarginLeft sets the left margin.
	// (The "margin-left" style attribute only.)
	SetMarginLeft(value string) Style

	// SetMarginLeftPx sets the left margin, in pixels.
	// (The "margin-left" style attribute only.)
	SetMarginLeftPx(width int) Style

	// MarginRight returns the right margin.
	// (The "margin-right" style attribute only.)
	MarginRight() string

	// SetMarginRight sets the right margin.
	// (The "margin-right" style attribute only.)
	SetMarginRight(value string) Style

	// SetMarginRightPx sets the right margin, in pixels.
	// (The "margin-right" style attribute only.)
	SetMarginRightPx(width int) Style

	// MarginTop returns the top margin.
	// (The "margin-top" style attribute only.)
	MarginTop() string

	// SetMarginTop sets the top margin.
	// (The "margin-top" style attribute only.)
	SetMarginTop(value string) Style

	// SetMarginTopPx sets the top margin, in pixels.
	// (The "margin-top" style attribute only.)
	SetMarginTopPx(height int) Style

	// MarginBottom returns the bottom margin.
	// (The "margin-bottom" style attribute only.)
	MarginBottom() string

	// SetMarginBottom sets the bottom margin.
	// (The "margin-bottom" style attribute only.)
	SetMarginBottom(value string) Style

	// SetMarginBottomPx sets the bottom margin, in pixels.
	// (The "margin-bottom" style attribute only.)
	SetMarginBottomPx(height int) Style

	// Background returns the background (color).
	Background() string

	// SetBackground sets the background (color).
	SetBackground(value string) Style

	// Border returns the border.
	Border() string

	// SetBorder sets the border.
	SetBorder(value string) Style

	// SetBorder2 sets the border specified by parts.
	// (The "border" style attribute only.)
	SetBorder2(width int, style, color string) Style

	// BorderLeft returns the left border.
	BorderLeft() string

	// SetBorderLeft sets the left border.
	SetBorderLeft(value string) Style

	// SetBorderLeft2 sets the left border specified by parts.
	// (The "border-left" style attribute only.)
	SetBorderLeft2(width int, style, color string) Style

	// BorderRight returns the right border.
	BorderRight() string

	// SetBorderRight sets the right border.
	SetBorderRight(value string) Style

	// SetBorderRight2 sets the right border specified by parts.
	// (The "border-right" style attribute only.)
	SetBorderRight2(width int, style, color string) Style

	// BorderTop returns the top border.
	BorderTop() string

	// SetBorderTop sets the top border.
	SetBorderTop(value string) Style

	// SetBorderTop2 sets the top border specified by parts.
	// (The "border-top" style attribute only.)
	SetBorderTop2(width int, style, color string) Style

	// BorderBottom returns the bottom border.
	BorderBottom() string

	// SetBorderBottom sets the bottom border.
	SetBorderBottom(value string) Style

	// SetBorderBottom2 sets the bottom border specified by parts.
	// (The "border-bottom" style attribute only.)
	SetBorderBottom2(width int, style, color string) Style

	// Color returns the (foreground) color.
	Color() string

	// SetColor sets the (foreground) color.
	SetColor(value string) Style

	// Cursor returns the (mouse) cursor.
	Cursor() string

	// SetCursor sets the (mouse) cursor.
	SetCursor(value string) Style

	// Display returns the display mode.
	Display() string

	// SetDisplay sets the display mode
	SetDisplay(value string) Style

	// FontSize returns the font size.
	FontSize() string

	// SetFontSize sets the font size.
	SetFontSize(value string) Style

	// FontStyle returns the font style.
	FontStyle() string

	// SetFontStyle sets the font style.
	SetFontStyle(value string) Style

	// FontWeight returns the font weight.
	FontWeight() string

	// SetFontWeight sets the font weight.
	SetFontWeight(value string) Style

	// Width returns the width.
	Width() string

	// SetWidth sets the width.
	SetWidth(value string) Style

	// SetWidthPx sets the width, in pixels.
	SetWidthPx(width int) Style

	// SetFullWidth sets full width (100%).
	SetFullWidth() Style

	// Height returns the height.
	Height() string

	// SetHeight sets the height.
	SetHeight(value string) Style

	// SetHeightPx sets the height.
	SetHeightPx(height int) Style

	// SetFullHeight sets full height (100%).
	SetFullHeight() Style

	// WhiteSpace returns the white space attribute value.
	WhiteSpace() string

	// SetWhiteSpace sets the white space attribute value.
	SetWhiteSpace(value string) Style

	// render renders all style information (style class names
	// and style attributes).
	render(w writer)

	// renderClasses renders the style class names.
	renderClasses(w writer)

	// renderAttrs renders the style attributes.
	renderAttrs(w writer)
}

type styleImpl struct {
	classes []string          // Style classes.
	attrs   map[string]string // Explicitly set style attributes. Lazily initialized.
}

// newStyleImpl creates a new styleImpl.
func newStyleImpl() *styleImpl {
	return &styleImpl{}
}

func (s *styleImpl) AddClass(class string) Style {
	s.classes = append(s.classes, class)
	return s
}

func (s *styleImpl) SetClass(class string) Style {
	s.classes = s.classes[0:0]
	if len(class) > 0 {
		s.classes = append(s.classes, class)
	}
	return s
}

func (s *styleImpl) RemoveClass(class string) Style {
	for i, class_ := range s.classes {
		if class_ == class {
			oldClasses := s.classes
			s.classes = append(oldClasses[0:i], oldClasses[i+1:]...)
			oldClasses[len(oldClasses)-1] = ""
			break
		}
	}

	return s
}

func (s *styleImpl) Get(name string) string {
	return s.attrs[name]
}

func (s *styleImpl) Set(name, value string) Style {
	if s.attrs == nil {
		s.attrs = make(map[string]string)
	}

	if len(value) > 0 {
		s.attrs[name] = value
	} else {
		delete(s.attrs, name)
	}
	return s
}

func (s *styleImpl) Size() (width, height string) {
	return s.Get(ST_WIDTH), s.Get(ST_HEIGHT)
}

func (s *styleImpl) SetSize(width, height string) Style {
	s.Set(ST_WIDTH, width)
	s.Set(ST_HEIGHT, height)
	return s
}

func (s *styleImpl) SetSizePx(width, height int) Style {
	return s.SetSize(strconv.Itoa(width)+"px", strconv.Itoa(height)+"px")
}

func (s *styleImpl) SetFullSize() Style {
	return s.SetSize("100%", "100%")
}

func (s *styleImpl) Padding() string {
	return s.Get(ST_PADDING)
}

func (s *styleImpl) SetPadding(value string) Style {
	return s.Set(ST_PADDING, value)
}

func (s *styleImpl) SetPadding2(top, right, bottom, left string) Style {
	return s.SetPadding(top + " " + right + " " + bottom + " " + left)
}

func (s *styleImpl) SetPaddingPx(top, right, bottom, left int) Style {
	return s.SetPadding(strconv.Itoa(top) + "px " + strconv.Itoa(right) + "px " + strconv.Itoa(bottom) + "px " + strconv.Itoa(left) + "px")
}

func (s *styleImpl) PaddingLeft() string {
	return s.Get(ST_PADDING_LEFT)
}

func (s *styleImpl) SetPaddingLeft(value string) Style {
	return s.Set(ST_PADDING_LEFT, value)
}

func (s *styleImpl) SetPaddingLeftPx(width int) Style {
	return s.SetPaddingLeft(strconv.Itoa(width) + "px")
}

func (s *styleImpl) PaddingRight() string {
	return s.Get(ST_PADDING_RIGHT)
}

func (s *styleImpl) SetPaddingRight(value string) Style {
	return s.Set(ST_PADDING_RIGHT, value)
}

func (s *styleImpl) SetPaddingRightPx(width int) Style {
	return s.SetPaddingRight(strconv.Itoa(width) + "px")
}

func (s *styleImpl) PaddingTop() string {
	return s.Get(ST_PADDING_TOP)
}

func (s *styleImpl) SetPaddingTop(value string) Style {
	return s.Set(ST_PADDING_TOP, value)
}

func (s *styleImpl) SetPaddingTopPx(height int) Style {
	return s.SetPaddingTop(strconv.Itoa(height) + "px")
}

func (s *styleImpl) PaddingBottom() string {
	return s.Get(ST_PADDING_BOTTOM)
}

func (s *styleImpl) SetPaddingBottom(value string) Style {
	return s.Set(ST_PADDING_BOTTOM, value)
}

func (s *styleImpl) SetPaddingBottomPx(height int) Style {
	return s.SetPaddingBottom(strconv.Itoa(height) + "px")
}

func (s *styleImpl) Margin() string {
	return s.Get(ST_MARGIN)
}

func (s *styleImpl) SetMargin(value string) Style {
	return s.Set(ST_MARGIN, value)
}

func (s *styleImpl) SetMargin2(top, right, bottom, left string) Style {
	return s.SetMargin(top + " " + right + " " + bottom + " " + left)
}

func (s *styleImpl) SetMarginPx(top, right, bottom, left int) Style {
	return s.SetMargin(strconv.Itoa(top) + "px " + strconv.Itoa(right) + "px " + strconv.Itoa(bottom) + "px " + strconv.Itoa(left) + "px ")
}

func (s *styleImpl) MarginLeft() string {
	return s.Get(ST_MARGIN_LEFT)
}

func (s *styleImpl) SetMarginLeft(value string) Style {
	return s.Set(ST_MARGIN_LEFT, value)
}

func (s *styleImpl) SetMarginLeftPx(width int) Style {
	return s.SetMarginLeft(strconv.Itoa(width) + "px")
}

func (s *styleImpl) MarginRight() string {
	return s.Get(ST_MARGIN_RIGHT)
}

func (s *styleImpl) SetMarginRight(value string) Style {
	return s.Set(ST_MARGIN_RIGHT, value)
}

func (s *styleImpl) SetMarginRightPx(width int) Style {
	return s.SetMarginRight(strconv.Itoa(width) + "px")
}

func (s *styleImpl) MarginTop() string {
	return s.Get(ST_MARGIN_TOP)
}

func (s *styleImpl) SetMarginTop(value string) Style {
	return s.Set(ST_MARGIN_TOP, value)
}

func (s *styleImpl) SetMarginTopPx(height int) Style {
	return s.SetMarginTop(strconv.Itoa(height) + "px")
}

func (s *styleImpl) MarginBottom() string {
	return s.Get(ST_MARGIN_BOTTOM)
}

func (s *styleImpl) SetMarginBottom(value string) Style {
	return s.Set(ST_MARGIN_BOTTOM, value)
}

func (s *styleImpl) SetMarginBottomPx(height int) Style {
	return s.SetMarginBottom(strconv.Itoa(height) + "px")
}

func (s *styleImpl) Background() string {
	return s.Get(ST_BACKGROUND)
}

func (s *styleImpl) SetBackground(value string) Style {
	return s.Set(ST_BACKGROUND, value)
}

func (s *styleImpl) Border() string {
	return s.Get(ST_BORDER)
}

func (s *styleImpl) SetBorder(value string) Style {
	return s.Set(ST_BORDER, value)
}

func (s *styleImpl) SetBorder2(width int, style, color string) Style {
	return s.SetBorder(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderLeft() string {
	return s.Get(ST_BORDER_LEFT)
}

func (s *styleImpl) SetBorderLeft(value string) Style {
	return s.Set(ST_BORDER_LEFT, value)
}

func (s *styleImpl) SetBorderLeft2(width int, style, color string) Style {
	return s.SetBorderLeft(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderRight() string {
	return s.Get(ST_BORDER_RIGHT)
}

func (s *styleImpl) SetBorderRight(value string) Style {
	return s.Set(ST_BORDER_RIGHT, value)
}

func (s *styleImpl) SetBorderRight2(width int, style, color string) Style {
	return s.SetBorderRight(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderTop() string {
	return s.Get(ST_BORDER_TOP)
}

func (s *styleImpl) SetBorderTop(value string) Style {
	return s.Set(ST_BORDER_TOP, value)
}

func (s *styleImpl) SetBorderTop2(width int, style, color string) Style {
	return s.SetBorderTop(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderBottom() string {
	return s.Get(ST_BORDER_BOTTOM)
}

func (s *styleImpl) SetBorderBottom(value string) Style {
	return s.Set(ST_BORDER_BOTTOM, value)
}

func (s *styleImpl) SetBorderBottom2(width int, style, color string) Style {
	return s.SetBorderBottom(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) Color() string {
	return s.Get(ST_COLOR)
}

func (s *styleImpl) SetColor(value string) Style {
	return s.Set(ST_COLOR, value)
}

func (s *styleImpl) Cursor() string {
	return s.Get(ST_CURSOR)
}

func (s *styleImpl) SetCursor(value string) Style {
	return s.Set(ST_CURSOR, value)
}

func (s *styleImpl) Display() string {
	return s.Get(ST_DISPLAY)
}

func (s *styleImpl) SetDisplay(value string) Style {
	return s.Set(ST_DISPLAY, value)
}

func (s *styleImpl) FontSize() string {
	return s.Get(ST_FONT_SIZE)
}

func (s *styleImpl) SetFontSize(value string) Style {
	return s.Set(ST_FONT_SIZE, value)
}

func (s *styleImpl) FontStyle() string {
	return s.Get(ST_FONT_STYLE)
}

func (s *styleImpl) SetFontStyle(value string) Style {
	return s.Set(ST_FONT_STYLE, value)
}

func (s *styleImpl) FontWeight() string {
	return s.Get(ST_FONT_WEIGHT)
}

func (s *styleImpl) SetFontWeight(value string) Style {
	return s.Set(ST_FONT_WEIGHT, value)
}

func (s *styleImpl) Height() string {
	return s.Get(ST_HEIGHT)
}

func (s *styleImpl) SetHeight(value string) Style {
	return s.Set(ST_HEIGHT, value)
}
func (s *styleImpl) SetHeightPx(height int) Style {
	return s.SetHeight(strconv.Itoa(height) + "px")
}

func (s *styleImpl) SetFullHeight() Style {
	return s.SetHeight("100%")
}

func (s *styleImpl) Width() string {
	return s.Get(ST_WIDTH)
}

func (s *styleImpl) SetWidth(value string) Style {
	return s.Set(ST_WIDTH, value)
}

func (s *styleImpl) SetWidthPx(width int) Style {
	return s.SetWidth(strconv.Itoa(width) + "px")
}

func (s *styleImpl) SetFullWidth() Style {
	return s.SetWidth("100%")
}

func (s *styleImpl) WhiteSpace() string {
	return s.Get(ST_WHITE_SPACE)
}

func (s *styleImpl) SetWhiteSpace(value string) Style {
	return s.Set(ST_WHITE_SPACE, value)
}

func (s *styleImpl) render(w writer) {
	s.renderClasses(w)

	if s.attrs != nil {
		w.Write(_STR_STYLE)
		s.renderAttrs(w)
		w.Write(_STR_QUOTE)
	}
}

func (s *styleImpl) renderClasses(w writer) {
	if len(s.classes) > 0 {
		w.Write(_STR_CLASS)
		for i, class := range s.classes {
			if i > 0 {
				w.Write(_STR_SPACE)
			}
			w.Writes(class)
		}
		w.Write(_STR_QUOTE)
	}
}

func (s *styleImpl) renderAttrs(w writer) {
	for name, value := range s.attrs {
		w.Writes(name)
		w.Write(_STR_COLON)
		w.Writes(value)
		w.Write(_STR_SEMICOL)
	}
}
