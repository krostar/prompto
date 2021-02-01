package color

import (
	"fmt"
	"strconv"
	"strings"
)

// Style stores a foreground and background color.
type Style [2]Color

// NewStyle creates a new Style based on both foreground and background color.
// Color are applied to the style only if they have the right kind.
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

// ANSISprintf works the same as fmt.Sprintf but surround output with
// the ansi escape sequences the Style + reset, to apply colors.
func (s Style) ANSISprintf(format string, args ...interface{}) string {
	var code []string

	if s[0].kind == KindForeground {
		r, g, b := s[0].RGB()
		code = append(code, "38;2;"+strconv.Itoa(int(r))+";"+strconv.Itoa(int(g))+";"+strconv.Itoa(int(b)))
	}

	if s[1].kind == KindBackground {
		r, g, b := s[1].RGB()
		code = append(code, "48;2;"+strconv.Itoa(int(r))+";"+strconv.Itoa(int(g))+";"+strconv.Itoa(int(b)))
	}

	return fmt.Sprintf("\x1b["+strings.Join(code, ";")+"m"+format+"\x1b[0m", args...)
}

// Colors splits back a Style to its components: a foreground and a background color.
func (s Style) Colors() (Color, Color) {
	return s[0], s[1]
}
