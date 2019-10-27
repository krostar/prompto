// Package usecase contains the pure business-related methods.
package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/krostar/prompto/internal/prompto/domain"
)

// WritePrompt usecase creates and writes a prompt.
func WritePrompt() PromptWriterFunc {
	return (&promptWriter{}).WritePrompt
}

// PromptWriterFunc defines the function signature to write a prompt.
type PromptWriterFunc func(context.Context, PromptCreationRequest, io.Writer) error

// PromptCreationRequest defines how to create a prompt.
type PromptCreationRequest struct {
	SegmentsProvider []domain.SegmentsProvider
	Direction        domain.Direction
	SeparatorConfig  domain.SeparatorConfig
}

type promptWriter struct{}

func (p *promptWriter) WritePrompt(ctx context.Context, req PromptCreationRequest, to io.Writer) error {
	var segments domain.Segments

	for _, segmenter := range req.SegmentsProvider {
		s, err := segmenter.ProvideSegments()
		if err != nil {
			return fmt.Errorf("unable to get %s prompt segment: %w", "segment-name", err)
		}

		segments = append(segments, s...)
	}

	prompt, err := domain.NewPrompt(segments, req.Direction, req.SeparatorConfig)
	if err != nil {
		return fmt.Errorf("unable to create prompt: %w", err)
	}

	if _, err := prompt.WriteTo(to); err != nil {
		return fmt.Errorf("unable to write prompt: %w", err)
	}

	return nil
}
