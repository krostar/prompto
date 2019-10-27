// Package usecase contains the pure business-related methods.
package usecase

import (
	"context"
	"fmt"
	"io"
	"sync"

	"go.uber.org/multierr"

	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/pkg/color"
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
			errs = multierr.Combine(errs, fmt.Errorf("unable to write %s prompt: %w", req.Direction, err))
		}
	}

	return errs
}

func (p *promptWriter) writePrompt(_ context.Context, req PromptCreationRequest) error {
	var (
		parallelSegments = make([]domain.Segments, len(req.SegmentsProvider))

		segmentsErr error
		wg          sync.WaitGroup
		mutex       sync.Mutex
	)

	for i, segmenter := range req.SegmentsProvider {
		segmenter := segmenter

		if domain.IsSegmentsProviderDisabledForDirection(req.Direction, segmenter) {
			segmentsErr = p.appendError(&mutex, segmentsErr, fmt.Errorf(
				"segment %s is not available for direction %s",
				segmenter.SegmentName(), req.Direction,
			))
			break
		}

		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			s, err := segmenter.ProvideSegments()
			if err != nil {
				segmentsErr = p.appendError(&mutex, segmentsErr, fmt.Errorf(
					"unable to get prompt segment %s: %w",
					segmenter.SegmentName(), err,
				))
				return
			}

			parallelSegments[index] = s
		}(i)
	}

	wg.Wait()

	if segmentsErr != nil {
		return segmentsErr
	}

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

func (p *promptWriter) appendError(mutex sync.Locker, errs, err error) error {
	mutex.Lock()
	errs = multierr.Combine(errs, err)
	mutex.Unlock()

	return errs
}
