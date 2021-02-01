package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsSegmentsProviderDisabledForDirection(t *testing.T) {
	t.Run("segment provider no disabler", func(t *testing.T) {
		var sp segmentProviderNoop

		assert.False(t, IsSegmentsProviderDisabledForDirection(DirectionLeft, &sp))
		assert.False(t, IsSegmentsProviderDisabledForDirection(DirectionRight, &sp))
		assert.True(t, IsSegmentsProviderDisabledForDirection(DirectionUnknown, &sp))
	})

	t.Run("segment provider disables direction left", func(t *testing.T) {
		var sp segmentProviderLeftDisabled

		assert.True(t, IsSegmentsProviderDisabledForDirection(DirectionLeft, &sp))
		assert.False(t, IsSegmentsProviderDisabledForDirection(DirectionRight, &sp))
		assert.True(t, IsSegmentsProviderDisabledForDirection(DirectionUnknown, &sp))
	})

	t.Run("segment provider disables direction right", func(t *testing.T) {
		var sp segmentProviderRightDisabled

		assert.False(t, IsSegmentsProviderDisabledForDirection(DirectionLeft, &sp))
		assert.True(t, IsSegmentsProviderDisabledForDirection(DirectionRight, &sp))
		assert.True(t, IsSegmentsProviderDisabledForDirection(DirectionUnknown, &sp))
	})
}

type segmentProviderNoop struct{}

func (s segmentProviderNoop) SegmentName() string                { return "noop" }
func (s segmentProviderNoop) ProvideSegments() (Segments, error) { return nil, nil }

type segmentProviderLeftDisabled struct{ segmentProviderNoop }

func (s segmentProviderLeftDisabled) DisabledForLeftPrompt() bool { return true }

type segmentProviderRightDisabled struct{ segmentProviderNoop }

func (s segmentProviderRightDisabled) DisabledForRightPrompt() bool { return true }
