// Package usecase contains the pure business-related methods.
package usecase

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"

	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/pkg/color"
)

// WritePrompts usecase creates and writes prompt(s).
func WritePrompts(writeTo io.Writer) PromptWriterFunc {
	return (&promptWriter{writeTo: writeTo}).WritePrompts
}

// PromptWriterFunc defines the function signature to write a prompt.
type PromptWriterFunc func(context.Context, ...PromptCreationRequest) error

// PromptCreationRequest defines how to create a prompt.
type PromptCreationRequest struct {
	Direction        domain.Direction
	Colorizer        color.Colorizer
	SegmentsProvider []domain.SegmentsProvider
	SeparatorConfig  domain.SeparatorConfig
}

type promptWriter struct {
	writeTo io.Writer
}

func (p *promptWriter) WritePrompts(ctx context.Context, reqs ...PromptCreationRequest) error {
	var errs error

	if len(reqs) == 0 || len(reqs) > 2 {
		return fmt.Errorf("wrong number of prompts to print (%d)", len(reqs))
	}

	for _, req := range reqs {
		if err := p.writePrompt(ctx, req); err != nil {
			errs = multierr.Combine(errs, fmt.Errorf("unable to write %s prompt: %w", req.Direction, err))
		}
	}

	return errs
}

func (p *promptWriter) writePrompt(ctx context.Context, req PromptCreationRequest) error {
	parallelSegments := make([]domain.Segments, len(req.SegmentsProvider))

	wg, ctx := errgroup.WithContext(ctx)

	for index, segmenter := range req.SegmentsProvider {
		index, segmenter := index, segmenter

		wg.Go(func() error {
			if domain.IsSegmentsProviderDisabledForDirection(req.Direction, segmenter) {
				return fmt.Errorf("segment %s is not available for direction %s", segmenter.SegmentName(), req.Direction)
			}

			segments, err := segmenter.ProvideSegments()
			if err != nil {
				return fmt.Errorf("unable to get prompt segment %s: %v", segmenter.SegmentName(), err)
			}

			parallelSegments[index] = segments

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil
	}

	var segments domain.Segments
	for _, segment := range parallelSegments {
		segments = append(segments, segment...)
	}

	if _, err := domain.
		NewPrompt(segments, req.Direction, req.SeparatorConfig).
		WriteTo(req.Colorizer, p.writeTo); err != nil {
		return fmt.Errorf("unable to write prompt: %w", err)
	}

	return nil
}
