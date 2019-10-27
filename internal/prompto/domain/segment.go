package domain

import (
	"fmt"
	"io"
	"strings"

	"github.com/krostar/prompto/pkg/color"
)

// Segment is a part of a prompt.
type Segment struct {
	direction Direction
	contents  []string
	style     color.Style

	spaceBefore bool
	spaceAfter  bool

	separatorDisabledForNextSegment bool
	separator                       *Separator
}

// NewSegment creates a new segment.
// Multiple contents can be provided, they will be printed in direction order.
func NewSegment(contents ...string) *Segment {
	return &Segment{contents: contents}
}

// Style returns the style of the segment.
func (s *Segment) Style() color.Style { return s.style }

// SetStyle sets the style of the segment.
func (s *Segment) SetStyle(style color.Style) *Segment {
	s.style = style
	return s
}

// WithSpaceAround adds some space around segment.
func (s *Segment) WithSpaceAround() *Segment {
	s.spaceBefore = true
	s.spaceAfter = true
	return s
}

// WithSpaceBefore adds a space before the segment content.
func (s *Segment) WithSpaceBefore() *Segment {
	s.spaceBefore = true
	return s
}

// WithSpaceAfter adds a space after the segment content.
func (s *Segment) WithSpaceAfter() *Segment {
	s.spaceAfter = true
	return s
}

// DisableNextSegmentSeparator prevents the next segments from writing its separator.
func (s *Segment) DisableNextSegmentSeparator() *Segment {
	s.separatorDisabledForNextSegment = true
	return s
}

func (s *Segment) setDirection(d Direction) {
	s.direction = d
}

func (s *Segment) setSeparator(sep Separator) {
	s.separator = &sep
}

func (s *Segment) contentWithSpace() string {
	if s.direction == DirectionRight {
		// swap content order
		last := len(s.contents) - 1
		for i := 0; i < len(s.contents)/2; i++ {
			s.contents[i], s.contents[last-i] = s.contents[last-i], s.contents[i]
		}

		if !s.spaceBefore || !s.spaceAfter {
			s.spaceBefore = !s.spaceBefore
			s.spaceAfter = !s.spaceAfter
		}
	}

	content := strings.Join(s.contents, " ")

	if s.spaceBefore {
		content = " " + content
	}

	if s.spaceAfter {
		content += " "
	}

	return content
}

// WriteTo implements io.WriterTo for Segment to write segment's content with provided style.
func (s *Segment) WriteTo(colorizer color.Colorizer, w io.Writer) (int64, error) {
	var wrote int64

	if s.separator != nil {
		nSeparator, errSeparator := s.separator.WriteTo(colorizer, w)
		wrote += nSeparator

		if errSeparator != nil {
			return wrote, fmt.Errorf("unable to write separator: %w", errSeparator)
		}
	}

	nSegment, errSegment := w.Write([]byte(colorizer.Colorize(s.style, s.contentWithSpace())))
	wrote += int64(nSegment)

	if errSegment != nil {
		return wrote, fmt.Errorf("unable to write segment: %w", errSegment)
	}

	return wrote, nil
}
