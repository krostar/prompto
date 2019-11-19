// Package usecase contains the pure business-related methods.
package usecase

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/krostar/prompto/pkg/color"
	"go.uber.org/multierr"

	"github.com/krostar/prompto/internal/prompto/domain"
)

// WritePrompts usecase creates and writes prompt(s).
func WritePrompts(writeTo io.Writer) PromptWriterFunc {
	return (&promptWriter{
		writeTo: writeTo,
	}).WritePrompts
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
			multierr.Combine(errs, fmt.Errorf("unable to write %s prompt: %w", req.Direction, err))
		}
	}

	return errs
}

func (p *promptWriter) writePrompt(ctx context.Context, req PromptCreationRequest) error {
	var (
		parallelSegments = make([]domain.Segments, len(req.SegmentsProvider))

		segmentsErr error
		wg          sync.WaitGroup
		mutex       sync.Mutex
	)

	wg.Add(len(req.SegmentsProvider))

	for i, segmenter := range req.SegmentsProvider {
		segmenter := segmenter

		go func(index int) {
			defer wg.Done()
			s, err := segmenter.ProvideSegments()
			if err != nil {
				mutex.Lock()
				segmentsErr = multierr.Combine(segmentsErr, fmt.Errorf("unable to get prompt segment: %w", err))
				mutex.Unlock()
				return
			}
			parallelSegments[index] = s
		}(i)
	}

	wg.Wait()

	var segments domain.Segments
	for _, segment := range parallelSegments {
		segments = append(segments, segment...)
	}

	prompt, err := domain.NewPrompt(segments, req.Direction, req.SeparatorConfig)
	if err != nil {
		return fmt.Errorf("unable to create prompt: %w", err)
	}

	if _, err := prompt.WriteTo(req.Colorizer, p.writeTo); err != nil {
		return fmt.Errorf("unable to write prompt: %w", err)
	}

	return nil
}
