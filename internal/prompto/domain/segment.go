package domain

import (
	"fmt"
	"io"
)

// Segment is a part of a prompt. It is composed
// of a content, a style, and a separator.
type Segment struct {
	content string

	style              Style
	separator          *Separator
	thinSeparatorColor Color
}

// NewSegment creates a new segment.
func NewSegment(content string) *Segment { return &Segment{content: content} }

// Style returns the style of the segment.
func (s *Segment) Style() Style { return s.style }

// SetStyle sets the style of the segment.
func (s *Segment) SetStyle(style Style) *Segment {
	s.style = style
	return s
}

// Separator returns the separator of the segment.
func (s *Segment) Separator() *Separator { return s.separator }

// SetSeparator sets the separator of the segment.
func (s *Segment) SetSeparator(separator *Separator) *Segment {
	s.separator = separator
	return s
}

// SeparatorColor returns the separator color of the segment.
func (s *Segment) SeparatorColor() Color { return s.thinSeparatorColor }

// SetSeparatorColor sets the separator color of the segment.
func (s *Segment) SetSeparatorColor(color Color) *Segment {
	s.thinSeparatorColor = color
	return s
}

// WithSpaceAround adds some space around segment.
func (s *Segment) WithSpaceAround() *Segment {
	s.content = " " + s.content + " "
	return s
}

// WriteTo implements io.WriterTo for Segment
// to write segment's content with provided style.
func (s *Segment) WriteTo(w io.Writer) (int64, error) {
	var wrote int64

	if s.separator != nil {
		nSeparator, errSeparator := s.separator.WriteTo(w)
		wrote += nSeparator

		if errSeparator != nil {
			return wrote, fmt.Errorf("unable to write separator: %w", errSeparator)
		}
	}

	nSegment, errSegment := w.Write([]byte(s.style.Colorize(s.content)))
	wrote += int64(nSegment)

	if errSegment != nil {
		return wrote, fmt.Errorf("unable to write segment: %w", errSegment)
	}

	return wrote, nil
}
