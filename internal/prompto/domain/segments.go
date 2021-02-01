package domain

import (
	"github.com/krostar/prompto/pkg/color"
)

// Segments stores multiple segments of a prompt.
type Segments []*Segment

func (ss Segments) applyDirectionAndSeparators(d Direction, cfg SeparatorConfig) {
	if d == DirectionRight {
		ss.InverseOrder()
	}

	var previous *Segment
	for _, s := range ss {
		s.direction = d

		if previous != nil && !previous.separatorDisabledForNextSegment {
			ss.setSegmentSeparator(s, d, cfg, previous.style)
		}

		previous = s
	}
}

func (ss Segments) setSegmentSeparator(s *Segment, d Direction, cfg SeparatorConfig, previousStyle color.Style) {
	var separator *Separator

	switch d {
	case DirectionLeft:
		separator = NewSeparator(d, cfg, previousStyle, s.style)
	case DirectionRight:
		separator = NewSeparator(d, cfg, s.style, previousStyle)
	}

	if separator != nil {
		s.separator = separator
	}
}

// InverseOrder reverses the order of the segments.
func (ss Segments) InverseOrder() {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
