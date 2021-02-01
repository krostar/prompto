package domain

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/krostar/prompto/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Prompt_no_segments(t *testing.T) {
	prompt := NewPrompt(nil, DirectionRight, SeparatorConfig{})
	assert.Equal(t, &Prompt{direction: DirectionRight}, prompt)

	var to bytes.Buffer
	_, err := prompt.WriteTo(&promptNoopColorizer{}, &to)
	require.NoError(t, err)
	assert.Equal(t, "", to.String())
}

// nolint: dupl
func Test_Prompt_left(t *testing.T) {
	prompt := NewPrompt(Segments{
		NewSegment("1").
			SetStyle(color.NewStyle(
				color.NewHexFGColor("#654321"), color.NewHexBGColor("#123456"),
			)),
		NewSegment("2").
			SetStyle(color.NewStyle(
				color.NewHexFGColor("#123456"), color.NewHexBGColor("#654321"),
			)),
		NewSegment("3").
			SetStyle(color.NewStyle(
				color.NewHexFGColor("#123456"), color.NewHexBGColor("#654321"),
			)),
	}, DirectionLeft, SeparatorConfig{
		Content: SeparatorContentConfig{
			Left:     "L",
			LeftThin: "l",
		},
		ThinFGColor: map[string]string{"#654321": "#FF00FF"},
	})

	var to bytes.Buffer
	n, err := prompt.WriteTo(&promptNoopColorizer{}, &to)
	require.NoError(t, err)
	assert.Equal(t, int64(13), n)
	assert.Equal(t, " 1 L 2 l 3 L ", to.String())
}

// nolint: dupl
func Test_Prompt_right(t *testing.T) {
	prompt := NewPrompt(Segments{
		NewSegment("1").
			SetStyle(color.NewStyle(
				color.NewHexFGColor("#654321"), color.NewHexBGColor("#123456"),
			)),
		NewSegment("2").
			SetStyle(color.NewStyle(
				color.NewHexFGColor("#123456"), color.NewHexBGColor("#654321"),
			)),
		NewSegment("3").
			SetStyle(color.NewStyle(
				color.NewHexFGColor("#123456"), color.NewHexBGColor("#654321"),
			)),
	}, DirectionRight, SeparatorConfig{
		Content: SeparatorContentConfig{
			Right:     "R",
			RightThin: "r",
		},
		ThinFGColor: map[string]string{"#654321": "#FF00FF"},
	})

	var to bytes.Buffer
	n, err := prompt.WriteTo(&promptNoopColorizer{}, &to)
	require.NoError(t, err)
	assert.Equal(t, int64(13), n)
	assert.Equal(t, " R 3 r 2 R 1 ", to.String())
}

type promptNoopColorizer struct{}

func (promptNoopColorizer) Colorize(style color.Style, format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
