// Package usecase contains the pure business-related methods.
package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/krostar/prompto/internal/prompto/domain"
)

// WritePrompt usecase creates and writes a prompt.
func WritePrompt(writeTo io.Writer) PromptWriterFunc {
	return (&promptWriter{
		writeTo: writeTo,
	}).WritePrompt
}

// PromptWriterFunc defines the function signature to write a prompt.
type PromptWriterFunc func(context.Context, PromptCreationRequest) error

// PromptCreationRequest defines how to create a prompt.
type PromptCreationRequest struct {
	SegmentsProvider []domain.SegmentsProvider
	Direction        domain.Direction
	SeparatorConfig  domain.SeparatorConfig
}

type promptWriter struct {
	writeTo io.Writer
}

func (p *promptWriter) WritePrompt(ctx context.Context, req PromptCreationRequest) error {
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

	if _, err := prompt.WriteTo(p.writeTo); err != nil {
		return fmt.Errorf("unable to write prompt: %w", err)
	}

	return nil
}
