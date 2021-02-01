package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetDefaultColor(t *testing.T) {
	var (
		white = "#FFFFFF"
		black = "#000000"
		red   = "#FF0000"

		def = Config{
			Foreground: &black,
			Background: &red,
		}
	)

	t.Run("all set", func(t *testing.T) {
		cfg := Config{
			Foreground: &white,
			Background: &black,
		}
		cfg.SetDefaultColor(def)

		assert.Equal(t, Config{
			Foreground: &white,
			Background: &black,
		}, cfg)
	})

	t.Run("bg unset", func(t *testing.T) {
		cfg := Config{Foreground: &white}
		cfg.SetDefaultColor(def)

		assert.Equal(t, Config{
			Foreground: &white,
			Background: &red,
		}, cfg)
	})

	t.Run("fg unset", func(t *testing.T) {
		cfg := Config{Background: &white}
		cfg.SetDefaultColor(def)

		assert.Equal(t, Config{
			Foreground: &black,
			Background: &white,
		}, cfg)
	})
}

func TestConfig_ToStyle(t *testing.T) {
	var (
		white = "#FFFFFF"
		black = "#000000"

		cfg = Config{
			Foreground: &black,
			Background: &white,
		}
	)

	fg, bg := cfg.ToStyle().Colors()
	assert.Equal(t, NewHexFGColor(black), fg)
	assert.Equal(t, NewHexBGColor(white), bg)
}
