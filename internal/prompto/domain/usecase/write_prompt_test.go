package usecase

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/krostar/prompto/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/prompto/internal/prompto/adapter/segment"
	"github.com/krostar/prompto/internal/prompto/domain"
)

func TestPromptWriter_WritePrompts_ok(t *testing.T) {
	provider, err := segment.ProvideSegments([]string{segment.SegmentNameStub},
		segment.Config{Stub: segment.StubConfig{
			Segments: domain.Segments{
				domain.NewSegment("hello"),
				domain.NewSegment("world"),
			},
		}},
	)
	require.NoError(t, err)

	var to bytes.Buffer
	err = WritePrompts(&to)(context.Background(), PromptCreationRequest{
		Direction:        domain.DirectionLeft,
		Colorizer:        &color.NoopColorizer{},
		SegmentsProvider: provider,
		SeparatorConfig: domain.SeparatorConfig{
			Content: domain.SeparatorContentConfig{
				Left:     "/",
				LeftThin: ">",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, " hello > world / ", to.String())
}

func TestPromptWriter_WritePrompts_ko_segment_is_disabled_for_direction(t *testing.T) {
	provider, err := segment.ProvideSegments([]string{"newline"}, segment.Config{})
	require.NoError(t, err)

	err = WritePrompts(nil)(context.Background(), PromptCreationRequest{
		Direction:        domain.DirectionRight,
		SegmentsProvider: provider,
	})
	assert.Error(t, err)
}

func TestPromptWriter_WritePrompts_ko_provider_failed(t *testing.T) {
	provider, err := segment.ProvideSegments([]string{segment.SegmentNameStub},
		segment.Config{Stub: segment.StubConfig{Error: errors.New("boum")}},
	)
	require.NoError(t, err)

	err = WritePrompts(nil)(context.Background(), PromptCreationRequest{
		Direction:        domain.DirectionRight,
		SegmentsProvider: provider,
	})
	assert.Error(t, err)
}

func TestPromptWriter_WritePrompts_ko_number_of_prompts(t *testing.T) {
	err := WritePrompts(nil)(context.Background())
	assert.Error(t, err)

	err = WritePrompts(nil)(context.Background(), PromptCreationRequest{}, PromptCreationRequest{}, PromptCreationRequest{})
	assert.Error(t, err)
}

func TestPromptWriter_WritePrompts_multithread_providers(t *testing.T) {
	t.FailNow()
}
