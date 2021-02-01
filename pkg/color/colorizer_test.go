package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewColorizer(t *testing.T) {
	t.Run("unhandled value", func(t *testing.T) {
		colorizer, err := NewColorizer("unhandled")
		assert.Error(t, err)
		assert.Nil(t, colorizer)
	})

	t.Run("no shell", func(t *testing.T) {
		colorizer, err := NewColorizer("")
		assert.NoError(t, err)
		assert.IsType(t, &defaultANSIColorizer{}, colorizer)
	})

	t.Run("shell bash", func(t *testing.T) {
		colorizer, err := NewColorizer("bash")
		assert.NoError(t, err)
		assert.IsType(t, &defaultANSIColorizer{}, colorizer)
	})

	t.Run("shell fish", func(t *testing.T) {
		colorizer, err := NewColorizer("fish")
		assert.NoError(t, err)
		assert.IsType(t, &defaultANSIColorizer{}, colorizer)
	})

	t.Run("shell zsh", func(t *testing.T) {
		colorizer, err := NewColorizer("zsh")
		assert.NoError(t, err)
		assert.IsType(t, &zshANSIColorizer{}, colorizer)
	})
}

func TestDefaultANSIColorizer_Colorize(t *testing.T) {
	var colorizer defaultANSIColorizer

	out := colorizer.Colorize(
		NewStyle(NewHexFGColor("#F0F0F0"), NewHexBGColor("#0F0F0F")),
		"%d:%d",
		34, 35,
	)
	assert.Equal(t, "\x1b[38;2;240;240;240;48;2;15;15;15m34:35\x1b[0m", out)
}

func TestZSHANSIColorizer_Colorize(t *testing.T) {
	var colorizer zshANSIColorizer

	out := colorizer.Colorize(
		NewStyle(NewHexFGColor("#F0F0F0"), NewHexBGColor("#0F0F0F")),
		"%d:%d",
		34, 35,
	)
	assert.Equal(t, "%{\x1b[38;2;240;240;240;48;2;15;15;15m34:35\x1b[0m%5G%}", out)
}
