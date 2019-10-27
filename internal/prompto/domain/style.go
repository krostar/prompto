package domain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

// Style store a foreground and background color.
type Style [2]Color

// NewStyle creates a new Style based on both
// foreground and background color. Color are applied
// to the style only if they have the right kind.
// See NewXColor for more details about colors.
func NewStyle(fg, bg Color) Style {
	var s Style

	if fg.kind == ColorKindForeground {
		s[0] = fg
	}

	if bg.kind == ColorKindBackground {
		s[1] = bg
	}

	return s
}

// Colorize behave has Sprintf, but also adds the style to the output.
func (s Style) Colorize(format string, a ...interface{}) string {
	return color.NewPrinter(s.String()).Sprintf(format, a...)
}

// String implements Stringer, and return the style ansi code.
func (s Style) String() string {
	var colors []uint8

	if s[0].kind == ColorKindForeground {
		colors = append(colors, s[0].value)
	}

	if s[1].kind == ColorKindBackground {
		colors = append(colors, s[1].value)
	}

	return color.S256(colors...).Code()
}

// SplitToColors splits back a Style to its components: a
// foreground and a background color. If the style didn't have
// a foreground or a background, an empty / invalid Color is
// returned for this component instead.
func (s Style) SplitToColors() (Color, Color, error) {
	codes := strings.Split(s.String(), ";")
	toUint8 := func(s string) (uint8, error) {
		u, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("unable to parse ansi 8 bit color from %q: %w", s, err)
		}

		return uint8(u), nil
	}

	var (
		fg, bg           Color
		fgValue, bgValue uint8
		err              error
	)

	switch l := len(codes); {
	case l >= 6:
		bgValue, err = toUint8(codes[5])
		if err != nil {
			err = fmt.Errorf("unable to parse background color: %w", err)
			break
		}

		bg = NewBGColor(bgValue)

		fallthrough
	case l >= 3:
		fgValue, err = toUint8(codes[2])
		if err != nil {
			err = fmt.Errorf("unable to parse foreground color: %w", err)
			break
		}

		fg = NewFGColor(fgValue)
	}

	return fg, bg, err
}
