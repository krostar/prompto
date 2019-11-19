package color

// Color aggregates both 8-bit color and the kind where the
// color applies. Value must be a 8-bit 256 mode color.
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
// or https://jonasjacek.github.io/colors/ to get more details
// on the value possibilities.
type Color struct {
	value uint8
	kind  ColorKind
}

// Value returns the color's value.
func (c Color) Value() uint8 { return c.value }

// ColorKind defines where the color applies.
type ColorKind int

const (
	// ColorKindUnknown is the zero value for ColorKind,
	// which is an invalid value.
	ColorKindUnknown ColorKind = iota
	// ColorKindForeground applies the color to the foreground.
	ColorKindForeground
	// ColorKindBackground applies the color to the background.
	ColorKindBackground
)

// NewFGColor creates a new foreground color. See the Color
// type definition for more details about the possible values.
func NewFGColor(value uint8) Color {
	return Color{
		value: value,
		kind:  ColorKindForeground,
	}
}

// NewBGColor creates a new background color. See the Color
// type definition for more details about the possible values.
func NewBGColor(value uint8) Color {
	return Color{
		value: value,
		kind:  ColorKindBackground,
	}
}
