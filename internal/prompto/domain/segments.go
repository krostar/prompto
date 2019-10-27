package domain

// Segments stores multiple segments of a prompt.
type Segments []*Segment

// SetSeparators sets the separator given a direction and a separator config.
func (ss Segments) SetSeparators(d Direction, cfg SeparatorConfig) error {
	var err error

	switch d {
	case DirectionLeft:
		err = ss.setLeftSeparatorsToSegments(cfg)
	case DirectionRight:
		ss.InverseOrder()
		err = ss.setRightSeparatorsToSegments(cfg)
	}

	return err
}

func (ss Segments) setLeftSeparatorsToSegments(cfg SeparatorConfig) error {
	var (
		previous *Segment
		err      error
	)

	for _, s := range ss {
		if previous != nil {
			var separator *Separator

			separator, err = NewSeparator(
				DirectionLeft, cfg,
				previous.Style(), s.Style(),
				s.SeparatorColor(),
			)
			if err != nil {
				break
			}

			s.SetSeparator(separator)
		}

		previous = s
	}

	return err
}

func (ss Segments) setRightSeparatorsToSegments(cfg SeparatorConfig) error {
	var err error

	for i := len(ss) - 1; i >= 0; i-- {
		if i-1 >= 0 {
			var separator *Separator

			separator, err = NewSeparator(
				DirectionRight, cfg,
				ss[i].Style(), ss[i-1].Style(),
				ss[i].SeparatorColor(),
			)
			if err != nil {
				break
			}

			ss[i].SetSeparator(separator)
		}
	}

	return err
}

// InverseOrder reverses the order of the segments.
func (ss Segments) InverseOrder() {
	last := len(ss) - 1

	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

// SegmentsProvider defines how to provide segments.
type SegmentsProvider interface {
	ProvideSegments() (Segments, error)
}
