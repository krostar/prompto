package domain

import (
	"bytes"
	"testing"

	"github.com/krostar/prompto/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewSegment(t *testing.T) {
	segment := NewSegment("a", "b", "c")
	assert.Equal(t, &Segment{
		contents:    []string{"a", "b", "c"},
		spaceBefore: true,
		spaceAfter:  true,
	}, segment)
}

func TestSegment_SetStyle(t *testing.T) {
	style := color.NewStyle(color.NewHexFGColor("#F0F0F0"), color.NewHexBGColor("#F0F0F0"))
	segment := NewSegment("segment").SetStyle(style)
	assert.Equal(t, style, segment.style)
}

func TestSegment_DisableSpaceAfter(t *testing.T) {
	segment := NewSegment("segment").DisableSpaceAfter()
	assert.False(t, segment.spaceAfter)
}

func TestSegment_DisableNextSegmentSeparator(t *testing.T) {
	segment := NewSegment("segment").DisableNextSegmentSeparator()
	assert.True(t, segment.separatorDisabledForNextSegment)
}

func TestSegment_contentWithSpace(t *testing.T) {
	tests := map[string]struct {
		d               Direction
		segment         *Segment
		expectedContent string
	}{
		"direction left with normal spaces": {d: DirectionLeft,
			segment:         NewSegment("a", "b", "c"),
			expectedContent: " a b c ",
		},
		"direction left with custom spaces": {d: DirectionLeft,
			segment:         NewSegment("a", "b", "c").DisableSpaceAfter(),
			expectedContent: " a b c",
		},
		"direction right with normal spaces": {d: DirectionRight,
			segment:         NewSegment("a", "b", "c"),
			expectedContent: " c b a ",
		},
		"direction right with custom spaces": {d: DirectionRight,
			segment:         NewSegment("a", "b", "c").DisableSpaceAfter(),
			expectedContent: "c b a ",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			test.segment.direction = test.d
			content := test.segment.contentWithSpace()
			assert.Equal(t, test.expectedContent, content)
		})
	}
}

func TestSegment_WriteTo(t *testing.T) {
	segment := *NewSegment("a", "b", "c")

	t.Run("without separator", func(t *testing.T) {
		var to bytes.Buffer
		_, err := segment.WriteTo(&promptNoopColorizer{}, &to)
		require.NoError(t, err)
		assert.Equal(t, " a b c ", to.String())
	})

	t.Run("with separator", func(t *testing.T) {
		var to bytes.Buffer
		segment := segment
		segment.separator = &Separator{content: "R",
			style: color.NewStyle(color.NewHexFGColor("#F00FF0"), color.NewHexBGColor("#0FF00F")),
		}
		_, err := segment.WriteTo(&promptNoopColorizer{}, &to)
		require.NoError(t, err)
		assert.Equal(t, "R a b c ", to.String())
	})
}
