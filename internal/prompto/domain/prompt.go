// Package domain defines all models.
package domain

import (
	"fmt"
	"io"

	"github.com/krostar/prompto/pkg/color"
)

// Prompt stores and manipulate prompt attributes.
type Prompt struct {
	direction      Direction
	segments       Segments
	finalSeparator Separator
}

// NewPrompt creates a new prompt given the segments,
// direction, and separator configuration.
func NewPrompt(segments Segments, d Direction, separatorConfig SeparatorConfig) (*Prompt, error) {
	if len(segments) == 0 {
		return &Prompt{direction: d}, nil
	}

	finalSegment := segments[len(segments)-1]

	if err := segments.ApplyDirectionAndSeparators(d, separatorConfig); err != nil {
		return nil, fmt.Errorf("unable to set segments separators: %w", err)
	}

	finalSeparator, err := FinalSeparator(d, separatorConfig, finalSegment.Style())
	if err != nil {
		return nil, fmt.Errorf("unable to get final separator: %w", err)
	}

	return &Prompt{
		direction:      d,
		segments:       segments,
		finalSeparator: *finalSeparator,
	}, nil
}

// WriteTo implements io.WriterTo for Segments
// to write all segment's content with style.
func (p *Prompt) WriteTo(colorizer color.Colorizer, w io.Writer) (int64, error) {
	var wrote int64

	if p.direction == DirectionRight && len(p.segments) > 0 {
		w.Write([]byte(" ")) // nolint: errcheck, gosec
		fwrote, err := p.finalSeparator.WriteTo(colorizer, w)
		wrote += fwrote + 1

		if err != nil {
			return wrote, fmt.Errorf("unable to write first separator: %w", err)
		}
	}

	for _, s := range p.segments {
		sWrote, err := s.WriteTo(colorizer, w)
		wrote += sWrote

		if err != nil {
			return wrote, fmt.Errorf("unable to write part of segments: %w", err)
		}
	}

	if p.direction == DirectionLeft && len(p.segments) > 0 {
		fwrote, err := p.finalSeparator.WriteTo(colorizer, w)
		wrote += fwrote + 1

		if err != nil {
			return wrote, fmt.Errorf("unable to write last separator: %w", err)
		}

		w.Write([]byte(" ")) // nolint: errcheck, gosec
	}

	return wrote, nil
}
