package domain

import (
	"errors"

	"github.com/krostar/prompto/pkg/color"
)

// Segments stores multiple segments of a prompt.
type Segments []*Segment

// ApplyDirectionAndSeparators sets the separators given a direction and a separator config.
func (ss Segments) ApplyDirectionAndSeparators(d Direction, cfg SeparatorConfig) error {
	if d == DirectionRight {
		ss.InverseOrder()
	}

	var (
		previous *Segment
		err      error
	)

	for _, s := range ss {
		s.setDirection(d)

		if previous != nil && !previous.separatorDisabledForNextSegment {
			if err = ss.setSegmentSeparator(s, d, cfg, previous.Style()); err != nil {
				break
			}
		}

		previous = s
	}

	return err
}

func (ss Segments) setSegmentSeparator(s *Segment, d Direction, cfg SeparatorConfig, previousStyle color.Style) error {
	var (
		separator *Separator
		err       error
	)

	switch d {
	case DirectionLeft:
		separator, err = NewSeparator(d, cfg, previousStyle, s.Style())
	case DirectionRight:
		separator, err = NewSeparator(d, cfg, s.Style(), previousStyle)
	default:
		err = errors.New("unknown direction")
	}

	if err != nil {
		return err
	}

	s.setSeparator(*separator)

	return nil
}

// InverseOrder reverses the order of the segments.
func (ss Segments) InverseOrder() {
	last := len(ss) - 1

	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
