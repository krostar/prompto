package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewStyle(t *testing.T) {
	fg := NewHexFGColor("#0ABCD0")
	bg := NewHexBGColor("#0DCBA0")

	t.Run("both colors ok", func(t *testing.T) {
		style := NewStyle(fg, bg)
		assert.Equal(t, fg, style[0])
		assert.Equal(t, bg, style[1])
	})

	t.Run("both colors ko", func(t *testing.T) {
		style := NewStyle(bg, fg)
		assert.Equal(t, Color{}, style[0])
		assert.Equal(t, Color{}, style[1])
	})

	t.Run("fg color ko", func(t *testing.T) {
		style := NewStyle(Color{kind: KindUnknown, r: 255}, bg)
		assert.Equal(t, Color{}, style[0])
		assert.Equal(t, bg, style[1])
	})

	t.Run("bg color ko", func(t *testing.T) {
		style := NewStyle(fg, Color{kind: KindUnknown, r: 255})
		assert.Equal(t, fg, style[0])
		assert.Equal(t, Color{}, style[1])
	})
}

func TestStyle_ANSISprintf(t *testing.T) {
	t.Run("both colors ok", func(t *testing.T) {
		out := NewStyle(NewHexFGColor("#0ABCD0"), NewHexBGColor("#0DCBA0")).
			ANSISprintf("hello %s: %d", "world", 42)
		assert.Equal(t, "\x1b[38;2;10;188;208;48;2;13;203;160mhello world: 42\x1b[0m", out)
	})

	t.Run("both colors ko", func(t *testing.T) {
		out := NewStyle(NewHexBGColor("#0ABCD0"), NewHexFGColor("#0DCBA0")).
			ANSISprintf("hello %s: %d", "world", 42)
		assert.Equal(t, "\x1b[mhello world: 42\x1b[0m", out)
	})

	t.Run("fg color ko", func(t *testing.T) {
		out := NewStyle(NewHexBGColor("#0ABCD0"), NewHexBGColor("#0DCBA0")).
			ANSISprintf("hello %s: %d", "world", 42)
		assert.Equal(t, "\x1b[48;2;13;203;160mhello world: 42\x1b[0m", out)
	})

	t.Run("bg color ko", func(t *testing.T) {
		out := NewStyle(NewHexFGColor("#0ABCD0"), NewHexFGColor("#0DCBA0")).
			ANSISprintf("hello %s: %d", "world", 42)
		assert.Equal(t, "\x1b[38;2;10;188;208mhello world: 42\x1b[0m", out)
	})
}

func TestStyle_Colors(t *testing.T) {
	expectedFG, expectedBG := NewHexFGColor("#0ABCD0"), NewHexBGColor("#0DCBA0")

	fg, bg := NewStyle(expectedFG, expectedBG).Colors()
	assert.Equal(t, expectedFG, fg)
	assert.Equal(t, expectedBG, bg)
}
