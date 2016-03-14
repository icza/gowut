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
	StBackground    = "background"     // Background (color)
	StBorder        = "border"         // Border
	StBorderLeft    = "border-left"    // Left border
	StBorderRight   = "border-right"   // Right border
	StBorderTop     = "border-top"     // Top border
	StBorderBottom  = "border-bottom"  // Bottom border
	StColor         = "color"          // (Foreground) color
	StCursor        = "cursor"         // Cursor
	StDisplay       = "display"        // Display
	StFontSize      = "font-size"      // Font size
	StFontStyle     = "font-style"     // Font style
	StFontWeight    = "font-weight"    // Font weight
	StHeight        = "height"         // Height
	StMargin        = "margin"         // Margin
	StMarginLeft    = "margin-left"    // Left margin
	StMarginRight   = "margin-right"   // Right margin
	StMarginTop     = "margin-top"     // Top margin
	StMarginBottom  = "margin-bottom"  // Bottom margin
	StPadding       = "padding"        // Padding
	StPaddingLeft   = "padding-left"   // Left padding
	StPaddingRight  = "padding-right"  // Right padding
	StPaddingTop    = "padding-top"    // Top padding
	StPaddingBottom = "padding-bottom" // Bottom padding
	StWhiteSpace    = "white-space"    // White-space
	StWidth         = "width"          // Width
)

// The 17 standard color constants.
const (
	ClrAqua    = "Aqua"    // Aqua    (#00FFFF)
	ClrBlack   = "Black"   // Black   (#000000)
	ClrBlue    = "Blue"    // Blue    (#0000FF)
	ClrFuchsia = "Fuchsia" // Fuchsia (#FF00FF)
	ClrGray    = "Gray"    // Gray    (#808080)
	ClrGrey    = "Grey"    // Grey    (#808080)
	ClrGreen   = "Green"   // Green   (#008000)
	ClrLime    = "Lime"    // Lime    (#00FF00)
	ClrMaroon  = "Maroon"  // Maroon  (#800000)
	ClrNavy    = "Navy"    // Navy    (#000080)
	ClrOlive   = "Olive"   // Olive   (#808000)
	ClrPurple  = "Purple"  // Purple  (#800080)
	ClrRed     = "Red"     // Red     (#FF0000)
	ClrSilver  = "Silver"  // Silver  (#C0C0C0)
	ClrTeal    = "Teal"    // Teal    (#008080)
	ClrWhite   = "White"   // White   (#FFFFFF)
	ClrYellow  = "Yellow"  // Yellow  (#FFFF00)
)

// Border style constants.
const (
	BrdStyleSolid  = "solid"  // Solid
	BrdStyleDashed = "dashed" // Dashed
	BrdStyleDotted = "dotted" // Dotted
	BrdStyleDouble = "double" // Double
	BrdStyleGroove = "groove" // 3D grooved border
	BrdStyleRidge  = "ridge"  // 3D ridged border
	BrdStyleInset  = "inset"  // 3D inset border
	BrdStyleOutset = "outset" // 3D outset border
)

// Font weight constants.
const (
	FontWeightNormal  = "normal"  // Normal
	FontWeightBold    = "bold"    // Bold
	FontWeightBolder  = "bolder"  // Bolder
	FontWeightLighter = "lighter" // Lighter
)

// Font style constants.
const (
	FontStyleNormal = "normal" // Normal
	FontStyleItalic = "italic" // Italic
)

// Mouse cursor constants.
const (
	CursorAuto      = "auto"      // Default. Web browser sets the cursor.
	CursorCrosshair = "crosshair" // Crosshair
	CursorDefault   = "default"   // The default cursor.
	CursorHelp      = "help"      // Help
	CursorMove      = "move"      // Move
	CursorPointer   = "pointer"   // Pointer
	CursorProgress  = "progress"  // Progress
	CursorText      = "text"      // Text
	CursorWait      = "wait"      // Wait
	CursorInherit   = "inherit"   // The cursor should be inherited from the parent element.
)

// Display mode constants.
const (
	DisplayNone    = "none"    // The element will not be displayed.
	DisplayBlock   = "block"   // The element is displayed as a block.
	DisplayInline  = "inline"  // The element is displayed as an in-line element. This is the default.
	DisplayInherit = "inherit" // The display property value will be inherited from the parent element.
)

// White space constants.
const (
	WhiteSpaceNormal  = "normal"   // Sequences of white spaces are collapsed into a single whitespace. Text will wrap when neccessary. This is the default.
	WhiteSpaceNowrap  = "nowrap"   // Sequences of whitespace will collapse into a single whitespace. Text will never wrap to the next line (the text is in one line).
	WhiteSpacePre     = "pre"      // Whitespace is preserved. Text will only wrap on line breaks.
	WhiteSpacePreLine = "pre-line" // Sequences of whitespace will collapse into a single whitespace. Text will wrap when necessary and on line breaks.
	WhiteSpacePreWrap = "pre-wrap" // Whitespace is preserved. Text will wrap when necessary, and on line breaks.
	WhiteSpaceInherit = "inherit"  // Whitespace property will be inherited from the parent element.
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
	render(w Writer)

	// renderClasses renders the style class names.
	renderClasses(w Writer)

	// renderAttrs renders the style attributes.
	renderAttrs(w Writer)
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
	return s.Get(StWidth), s.Get(StHeight)
}

func (s *styleImpl) SetSize(width, height string) Style {
	s.Set(StWidth, width)
	s.Set(StHeight, height)
	return s
}

func (s *styleImpl) SetSizePx(width, height int) Style {
	return s.SetSize(strconv.Itoa(width)+"px", strconv.Itoa(height)+"px")
}

func (s *styleImpl) SetFullSize() Style {
	return s.SetSize("100%", "100%")
}

func (s *styleImpl) Padding() string {
	return s.Get(StPadding)
}

func (s *styleImpl) SetPadding(value string) Style {
	return s.Set(StPadding, value)
}

func (s *styleImpl) SetPadding2(top, right, bottom, left string) Style {
	return s.SetPadding(top + " " + right + " " + bottom + " " + left)
}

func (s *styleImpl) SetPaddingPx(top, right, bottom, left int) Style {
	return s.SetPadding(strconv.Itoa(top) + "px " + strconv.Itoa(right) + "px " + strconv.Itoa(bottom) + "px " + strconv.Itoa(left) + "px")
}

func (s *styleImpl) PaddingLeft() string {
	return s.Get(StPaddingLeft)
}

func (s *styleImpl) SetPaddingLeft(value string) Style {
	return s.Set(StPaddingLeft, value)
}

func (s *styleImpl) SetPaddingLeftPx(width int) Style {
	return s.SetPaddingLeft(strconv.Itoa(width) + "px")
}

func (s *styleImpl) PaddingRight() string {
	return s.Get(StPaddingRight)
}

func (s *styleImpl) SetPaddingRight(value string) Style {
	return s.Set(StPaddingRight, value)
}

func (s *styleImpl) SetPaddingRightPx(width int) Style {
	return s.SetPaddingRight(strconv.Itoa(width) + "px")
}

func (s *styleImpl) PaddingTop() string {
	return s.Get(StPaddingTop)
}

func (s *styleImpl) SetPaddingTop(value string) Style {
	return s.Set(StPaddingTop, value)
}

func (s *styleImpl) SetPaddingTopPx(height int) Style {
	return s.SetPaddingTop(strconv.Itoa(height) + "px")
}

func (s *styleImpl) PaddingBottom() string {
	return s.Get(StPaddingBottom)
}

func (s *styleImpl) SetPaddingBottom(value string) Style {
	return s.Set(StPaddingBottom, value)
}

func (s *styleImpl) SetPaddingBottomPx(height int) Style {
	return s.SetPaddingBottom(strconv.Itoa(height) + "px")
}

func (s *styleImpl) Margin() string {
	return s.Get(StMargin)
}

func (s *styleImpl) SetMargin(value string) Style {
	return s.Set(StMargin, value)
}

func (s *styleImpl) SetMargin2(top, right, bottom, left string) Style {
	return s.SetMargin(top + " " + right + " " + bottom + " " + left)
}

func (s *styleImpl) SetMarginPx(top, right, bottom, left int) Style {
	return s.SetMargin(strconv.Itoa(top) + "px " + strconv.Itoa(right) + "px " + strconv.Itoa(bottom) + "px " + strconv.Itoa(left) + "px ")
}

func (s *styleImpl) MarginLeft() string {
	return s.Get(StMarginLeft)
}

func (s *styleImpl) SetMarginLeft(value string) Style {
	return s.Set(StMarginLeft, value)
}

func (s *styleImpl) SetMarginLeftPx(width int) Style {
	return s.SetMarginLeft(strconv.Itoa(width) + "px")
}

func (s *styleImpl) MarginRight() string {
	return s.Get(StMarginRight)
}

func (s *styleImpl) SetMarginRight(value string) Style {
	return s.Set(StMarginRight, value)
}

func (s *styleImpl) SetMarginRightPx(width int) Style {
	return s.SetMarginRight(strconv.Itoa(width) + "px")
}

func (s *styleImpl) MarginTop() string {
	return s.Get(StMarginTop)
}

func (s *styleImpl) SetMarginTop(value string) Style {
	return s.Set(StMarginTop, value)
}

func (s *styleImpl) SetMarginTopPx(height int) Style {
	return s.SetMarginTop(strconv.Itoa(height) + "px")
}

func (s *styleImpl) MarginBottom() string {
	return s.Get(StMarginBottom)
}

func (s *styleImpl) SetMarginBottom(value string) Style {
	return s.Set(StMarginBottom, value)
}

func (s *styleImpl) SetMarginBottomPx(height int) Style {
	return s.SetMarginBottom(strconv.Itoa(height) + "px")
}

func (s *styleImpl) Background() string {
	return s.Get(StBackground)
}

func (s *styleImpl) SetBackground(value string) Style {
	return s.Set(StBackground, value)
}

func (s *styleImpl) Border() string {
	return s.Get(StBorder)
}

func (s *styleImpl) SetBorder(value string) Style {
	return s.Set(StBorder, value)
}

func (s *styleImpl) SetBorder2(width int, style, color string) Style {
	return s.SetBorder(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderLeft() string {
	return s.Get(StBorderLeft)
}

func (s *styleImpl) SetBorderLeft(value string) Style {
	return s.Set(StBorderLeft, value)
}

func (s *styleImpl) SetBorderLeft2(width int, style, color string) Style {
	return s.SetBorderLeft(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderRight() string {
	return s.Get(StBorderRight)
}

func (s *styleImpl) SetBorderRight(value string) Style {
	return s.Set(StBorderRight, value)
}

func (s *styleImpl) SetBorderRight2(width int, style, color string) Style {
	return s.SetBorderRight(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderTop() string {
	return s.Get(StBorderTop)
}

func (s *styleImpl) SetBorderTop(value string) Style {
	return s.Set(StBorderTop, value)
}

func (s *styleImpl) SetBorderTop2(width int, style, color string) Style {
	return s.SetBorderTop(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) BorderBottom() string {
	return s.Get(StBorderBottom)
}

func (s *styleImpl) SetBorderBottom(value string) Style {
	return s.Set(StBorderBottom, value)
}

func (s *styleImpl) SetBorderBottom2(width int, style, color string) Style {
	return s.SetBorderBottom(strconv.Itoa(width) + "px " + style + " " + color)
}

func (s *styleImpl) Color() string {
	return s.Get(StColor)
}

func (s *styleImpl) SetColor(value string) Style {
	return s.Set(StColor, value)
}

func (s *styleImpl) Cursor() string {
	return s.Get(StCursor)
}

func (s *styleImpl) SetCursor(value string) Style {
	return s.Set(StCursor, value)
}

func (s *styleImpl) Display() string {
	return s.Get(StDisplay)
}

func (s *styleImpl) SetDisplay(value string) Style {
	return s.Set(StDisplay, value)
}

func (s *styleImpl) FontSize() string {
	return s.Get(StFontSize)
}

func (s *styleImpl) SetFontSize(value string) Style {
	return s.Set(StFontSize, value)
}

func (s *styleImpl) FontStyle() string {
	return s.Get(StFontStyle)
}

func (s *styleImpl) SetFontStyle(value string) Style {
	return s.Set(StFontStyle, value)
}

func (s *styleImpl) FontWeight() string {
	return s.Get(StFontWeight)
}

func (s *styleImpl) SetFontWeight(value string) Style {
	return s.Set(StFontWeight, value)
}

func (s *styleImpl) Height() string {
	return s.Get(StHeight)
}

func (s *styleImpl) SetHeight(value string) Style {
	return s.Set(StHeight, value)
}
func (s *styleImpl) SetHeightPx(height int) Style {
	return s.SetHeight(strconv.Itoa(height) + "px")
}

func (s *styleImpl) SetFullHeight() Style {
	return s.SetHeight("100%")
}

func (s *styleImpl) Width() string {
	return s.Get(StWidth)
}

func (s *styleImpl) SetWidth(value string) Style {
	return s.Set(StWidth, value)
}

func (s *styleImpl) SetWidthPx(width int) Style {
	return s.SetWidth(strconv.Itoa(width) + "px")
}

func (s *styleImpl) SetFullWidth() Style {
	return s.SetWidth("100%")
}

func (s *styleImpl) WhiteSpace() string {
	return s.Get(StWhiteSpace)
}

func (s *styleImpl) SetWhiteSpace(value string) Style {
	return s.Set(StWhiteSpace, value)
}

func (s *styleImpl) render(w Writer) {
	s.renderClasses(w)

	if s.attrs != nil {
		w.Write(strStyle)
		s.renderAttrs(w)
		w.Write(strQuote)
	}
}

func (s *styleImpl) renderClasses(w Writer) {
	if len(s.classes) > 0 {
		w.Write(strClass)
		for i, class := range s.classes {
			if i > 0 {
				w.Write(strSpace)
			}
			w.Writes(class)
		}
		w.Write(strQuote)
	}
}

func (s *styleImpl) renderAttrs(w Writer) {
	for name, value := range s.attrs {
		w.Writes(name)
		w.Write(strColon)
		w.Writes(value)
		w.Write(strSemicol)
	}
}
