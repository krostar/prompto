package domain

import (
	"fmt"
	"io"
)

// Separator defines what a separator is.
type Separator struct {
	style   Style
	content string
}

// WriteTo implements io.WriterTo for Separator
// a segment's separator with style.
func (s *Separator) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(s.style.Colorize(s.content)))
	return int64(n), err
}

// SeparatorConfig defines how to configure separators.
type SeparatorConfig struct {
	ContentLeft      string `yaml:"content-left"`
	ContentLeftThin  string `yaml:"content-left-thin"`
	ContentRight     string `yaml:"content-right"`
	ContentRightThin string `yaml:"content-right-thin"`
}

// NewSeparator creates a new separator given the previous and current segment's style.
func NewSeparator(d Direction, cfg SeparatorConfig, previous, current Style, thinSepColor Color) (*Separator, error) {
	_, prevBG, err := previous.SplitToColors()
	if err != nil {
		return nil, fmt.Errorf("unable to get colors from previous segment style: %w", err)
	}

	_, currBG, err := current.SplitToColors()
	if err != nil {
		return nil, fmt.Errorf("unable to get colors from current segment style: %w", err)
	}

	sameBackground := prevBG.Value() == currBG.Value()

	var fg Color
	if sameBackground && thinSepColor.kind == ColorKindForeground {
		fg = thinSepColor
	} else {
		fg = NewFGColor(prevBG.Value())
	}

	var content string

	switch d {
	case DirectionLeft:
		if sameBackground {
			content = cfg.ContentLeftThin
		} else {
			content = cfg.ContentLeft
		}
	case DirectionRight:
		if sameBackground {
			content = cfg.ContentRightThin
		} else {
			content = cfg.ContentRight
		}
	}

	return &Separator{
		content: content,
		style:   NewStyle(fg, currBG),
	}, nil
}

// FinalSeparator returns the final separator given the last segment's style.
func FinalSeparator(lastSegmentStyle Style, d Direction, cfg SeparatorConfig) (*Separator, error) {
	var content string

	switch d {
	case DirectionLeft:
		content = cfg.ContentLeft
	case DirectionRight:
		content = cfg.ContentRight
	}

	_, bg, err := lastSegmentStyle.SplitToColors()
	if err != nil {
		return nil, fmt.Errorf("unable to split colors from last segment style %q: %w", lastSegmentStyle, err)
	}

	return &Separator{
		content: content,
		style:   NewStyle(NewFGColor(bg.Value()), Color{}),
	}, err
}
