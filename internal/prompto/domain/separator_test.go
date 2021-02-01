package domain

import (
	"bytes"
	"testing"

	"github.com/krostar/prompto/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewSeparator(t *testing.T) {
	cfg := SeparatorConfig{
		ThinFGColor: map[string]string{"#F0F0F0": "#F0F00F"},
		Content: SeparatorContentConfig{
			Left: "L", LeftThin: "l",
			Right: "R", RightThin: "r",
		},
	}
	tests := map[string]struct {
		d        Direction
		cfg      SeparatorConfig
		current  color.Style
		previous color.Style
		expected *Separator
	}{
		"left segments same background": {d: DirectionLeft, cfg: cfg,
			current:  color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#F0F0F0")),
			previous: color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#F0F0F0")),
			expected: &Separator{content: "l",
				style: color.NewStyle(color.NewHexFGColor("#F0F00F"), color.NewHexBGColor("#F0F0F0")),
			},
		},
		"left segments different background": {d: DirectionLeft, cfg: cfg,
			current:  color.NewStyle(color.NewHexFGColor("#F0F0F0"), color.NewHexBGColor("#F00FF0")),
			previous: color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#0FF00F")),
			expected: &Separator{content: "L",
				style: color.NewStyle(color.NewHexFGColor("#F00FF0"), color.NewHexBGColor("#0FF00F")),
			},
		},
		"right segments same background": {d: DirectionRight, cfg: cfg,
			current:  color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#F0F00F")),
			previous: color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#F0F00F")),
			expected: &Separator{content: "r",
				style: color.NewStyle(color.NewHexFGColor("#F0F00F"), color.NewHexBGColor("#F0F00F")),
			},
		},
		"right segments different background": {d: DirectionRight, cfg: cfg,
			current:  color.NewStyle(color.NewHexFGColor("#F0F0F0"), color.NewHexBGColor("#F00FF0")),
			previous: color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#0FF00F")),
			expected: &Separator{content: "R",
				style: color.NewStyle(color.NewHexFGColor("#F00FF0"), color.NewHexBGColor("#0FF00F")),
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			separator := NewSeparator(test.d, test.cfg, test.current, test.previous)
			assert.Equal(t, test.expected.style, separator.style)
			assert.Equal(t, test.expected.content, separator.content)
		})
	}
}

func Test_FinalSeparator(t *testing.T) {
	cfg := SeparatorConfig{
		ThinFGColor: map[string]string{"#F0F0F0": "#F0F00F"},
		Content: SeparatorContentConfig{
			Left: "L", LeftThin: "l",
			Right: "R", RightThin: "r",
		},
	}
	tests := map[string]struct {
		d        Direction
		cfg      SeparatorConfig
		style    color.Style
		expected *Separator
	}{
		"left direction": {d: DirectionLeft, cfg: cfg,
			style: color.NewStyle(color.NewHexFGColor("#F0F0F0"), color.NewHexBGColor("#F00FF0")),
			expected: &Separator{content: "L",
				style: color.NewStyle(color.NewHexFGColor("#F00FF0"), color.Color{}),
			},
		},
		"right direction": {d: DirectionRight, cfg: cfg,
			style: color.NewStyle(color.NewHexFGColor("#0F0F0F"), color.NewHexBGColor("#F0F00F")),
			expected: &Separator{content: "R",
				style: color.NewStyle(color.NewHexFGColor("#F0F00F"), color.Color{}),
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			separator := FinalSeparator(test.d, test.cfg, test.style)
			assert.Equal(t, test.expected.style, separator.style)
			assert.Equal(t, test.expected.content, separator.content)
		})
	}
}

func TestSeparator_WriteTo(t *testing.T) {
	var to bytes.Buffer
	_, err := (&Separator{
		content: "R",
		style:   color.NewStyle(color.NewHexFGColor("#F0F00F"), color.Color{}),
	}).WriteTo(&promptNoopColorizer{}, &to)
	require.NoError(t, err)
	assert.Equal(t, "R", to.String())
}
