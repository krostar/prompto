// Package color handles ANSI true color.
package color

import (
	"errors"
	"fmt"
)

// Color defines RGB true colors.
type Color struct {
	r    uint8
	g    uint8
	b    uint8
	kind Kind
}

// String returns the hex representation of an RGB color.
func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.r, c.g, c.b)
}

// RGB returns the color composites.
func (c Color) RGB() (uint8, uint8, uint8) {
	return c.r, c.g, c.b
}

// Equal checks whenever two colors are the same.
func (c Color) Equal(cc Color) bool {
	return c.kind == cc.kind && c.r == cc.r && c.g == cc.g && c.b == cc.b
}

// HexColor parses an hexadecimal like color (#FF00FF).
func HexColor(hex string) (Color, error) {
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

		err = errors.New("invalid byte")

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

	return c, err
}

// HexFGColor parses foreground hex color.
// Invalid provided values return KindUnknown color kind.
func HexFGColor(value string) Color {
	color, err := HexColor(value)
	if err == nil {
		color.kind = KindForeground
	}

	return color
}

// HexBGColor parses background hex color.
// Invalid provided values return KindUnknown color kind.
func HexBGColor(value string) Color {
	color, err := HexColor(value)
	if err == nil {
		color.kind = KindBackground
	}

	return color
}
