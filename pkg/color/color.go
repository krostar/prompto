// Package color handles ANSI colors and colorization.
package color

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/multierr"
)

// Color stores a color definition: RGB composition and kind.
type Color struct {
	kind Kind
	r    uint8
	g    uint8
	b    uint8
}

// String returns the hex representation of a color.
func (c Color) String() string {
	return strings.ToUpper(fmt.Sprintf("#%02x%02x%02x", c.r, c.g, c.b))
}

// RGB returns the color composites.
func (c Color) RGB() (uint8, uint8, uint8) {
	return c.r, c.g, c.b
}

// Equal checks whenever two colors are the same.
func (c Color) Equal(cc Color) bool {
	return c.kind == cc.kind && c.r == cc.r && c.g == cc.g && c.b == cc.b
}

// NewRGBFGColor creates a new foreground color from its RGB components.
func NewRGBFGColor(r, g, b uint8) Color {
	return Color{kind: KindForeground, r: r, g: g, b: b}
}

// NewRGBBGColor creates a new background color from its RGB components.
func NewRGBBGColor(r, g, b uint8) Color {
	return Color{kind: KindBackground, r: r, g: g, b: b}
}

func newHexColor(hex string) (Color, error) {
	if hex[0] != '#' {
		return Color{}, errors.New("invalid format")
	}

	var err error

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}

		err = multierr.Append(err, errors.New("invalid byte"))

		return 0
	}

	var c Color

	switch len(hex) {
	case 7:
		c.r = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		c.g = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		c.b = hexToByte(hex[5])<<4 + hexToByte(hex[6])
	case 4:
		c.r = hexToByte(hex[1]) * 17
		c.g = hexToByte(hex[2]) * 17
		c.b = hexToByte(hex[3]) * 17
	default:
		return Color{}, errors.New("invalid length")
	}

	if err != nil {
		return Color{}, err
	}

	return c, nil
}

// NewHexFGColor parses foreground hex color.
// Invalid provided values return KindUnknown color kind.
func NewHexFGColor(value string) Color {
	color, err := newHexColor(value)
	if err != nil {
		color = Color{}
	} else {
		color.kind = KindForeground
	}

	return color
}

// NewHexBGColor parses background hex color.
// Invalid provided values return KindUnknown color kind.
func NewHexBGColor(value string) Color {
	color, err := newHexColor(value)
	if err != nil {
		color = Color{}
	} else {
		color.kind = KindBackground
	}

	return color
}
