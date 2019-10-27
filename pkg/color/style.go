package color

import (
	"fmt"
	"strings"
)

const (
	ansiEscapeSequenceReset         = "\x1b[0m"
	ansiEscapeSequenceStyleTemplate = "\x1b[%sm"
	ansiForegroundColorTemplate     = "38;2;%d;%d;%d"
	ansiBackgroundColorTemplate     = "48;2;%d;%d;%d"
)

// Style stores a foreground and background color.
type Style [2]Color

// NewStyle creates a new Style based on both foreground and background color.
// Color are applied to the style only if they have the right kind.
// See NewXColor for more details about colors.
func NewStyle(fg, bg Color) Style {
	var s Style

	if fg.kind == KindForeground {
		s[0] = fg
	}

	if bg.kind == KindBackground {
		s[1] = bg
	}

	return s
}

// EscapeSequence returns the ansi escape sequences of the Style.
func (s Style) EscapeSequence() string {
	var code []string

	if s[0].kind == KindForeground {
		r, g, b := s[0].RGB()
		code = append(code, fmt.Sprintf(ansiForegroundColorTemplate, r, g, b))
	}

	if s[1].kind == KindBackground {
		r, g, b := s[1].RGB()
		code = append(code, fmt.Sprintf(ansiBackgroundColorTemplate, r, g, b))
	}

	return strings.Join(code, ";")
}

// Colors splits back a Style to its components: a foreground and a background color.
// If the style didn't have a foreground or a background, an empty / invalid Color is
// returned for this component instead.
func (s Style) Colors() (Color, Color, error) {
	return s[0], s[1], nil
}
