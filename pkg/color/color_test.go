package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor_String(t *testing.T) {
	assert.Equal(t, "#A70C67", Color{r: 0xa7, g: 0x0c, b: 0x67}.String())
}

func TestColor_RGB(t *testing.T) {
	r, g, b := Color{r: 167, g: 12, b: 103}.RGB()
	assert.Equal(t, uint8(167), r)
	assert.Equal(t, uint8(12), g)
	assert.Equal(t, uint8(103), b)
}

func TestColor_Equal(t *testing.T) {
	c := Color{r: 167, g: 12, b: 103}

	cc := Color{r: 167, g: 12, b: 103}
	assert.True(t, c.Equal(cc))

	cc.b = 104
	assert.False(t, c.Equal(cc))
}

func Test_NewRGBFGColor(t *testing.T) {
	color := NewRGBFGColor(12, 34, 56)
	assert.Equal(t, Color{
		kind: KindForeground,
		r:    12, g: 34, b: 56,
	}, color)
}

func Test_NewRGBBGColor(t *testing.T) {
	color := NewRGBBGColor(12, 34, 56)
	assert.Equal(t, Color{
		kind: KindBackground,
		r:    12, g: 34, b: 56,
	}, color)
}

func Test_newHexColor(t *testing.T) {
	tests := map[string]struct {
		expectedColor   Color
		expectedFailure bool
	}{
		"#123456": {
			expectedColor: Color{r: 0x12, g: 0x34, b: 0x56},
		}, "#ffffff": {
			expectedColor: Color{r: 0xff, g: 0xff, b: 0xff},
		}, "#000000": {
			expectedColor: Color{r: 0x00, g: 0x00, b: 0x00},
		}, "#ABC": {
			expectedColor: Color{r: 0xaa, g: 0xbb, b: 0xcc},
		}, "#ABCD": {
			expectedFailure: true,
		}, "42": {
			expectedFailure: true,
		}, "helloworld": {
			expectedFailure: true,
		}, "#BCDEFG": {
			expectedFailure: true,
		},
	}

	for hex, test := range tests {
		hex, test := hex, test
		t.Run(hex, func(t *testing.T) {
			color, err := newHexColor(hex)
			assert.Equal(t, test.expectedFailure, err != nil)
			assert.Equal(t, test.expectedColor, color)
		})
	}
}

func Test_NewHexFGColor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		color := NewHexFGColor("#123456")
		assert.Equal(t, Color{
			kind: KindForeground,
			r:    0x12, g: 0x34, b: 0x56,
		}, color)
	})

	t.Run("failed to parse color", func(t *testing.T) {
		color := NewHexFGColor("invalid")
		assert.Equal(t, Color{}, color)
	})
}

func Test_NewHexBGColor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		color := NewHexBGColor("#123456")
		assert.Equal(t, Color{
			kind: KindBackground,
			r:    0x12, g: 0x34, b: 0x56,
		}, color)
	})

	t.Run("failed to parse color", func(t *testing.T) {
		color := NewHexBGColor("invalid")
		assert.Equal(t, Color{}, color)
	})
}
