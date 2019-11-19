package domain

import (
	"fmt"
	"io"

	"github.com/krostar/prompto/pkg/color"
)

// Separator defines what a separator is.
type Separator struct {
	style   color.Style
	content string
}

// WriteTo implements io.WriterTo for Separator
// a segment's separator with style.
func (s *Separator) WriteTo(colorizer color.Colorizer, w io.Writer) (int64, error) {
	n, err := w.Write([]byte(colorizer.Colorize(s.style, s.content)))
	return int64(n), err
}

// SeparatorConfig defines how to configure separators.
type SeparatorConfig struct {
	Content     SeparatorContentConfig `yaml:"content"`
	ThinFGColor map[uint8]uint8        `yaml:"thin-fg-color"`
}

type SeparatorContentConfig struct {
	Left      string `yaml:"left"`
	LeftThin  string `yaml:"left-thin"`
	Right     string `yaml:"right"`
	RightThin string `yaml:"right-thin"`
}

// NewSeparator creates a new separator given the previous and current segment's style.
func NewSeparator(d Direction, cfg SeparatorConfig, previous, current color.Style) (*Separator, error) {
	_, prevBG, err := previous.SplitToColors()
	if err != nil {
		return nil, fmt.Errorf("unable to get colors from previous segment style: %w", err)
	}

	_, currBG, err := current.SplitToColors()
	if err != nil {
		return nil, fmt.Errorf("unable to get colors from current segment style: %w", err)
	}

	sameBackground := prevBG.Value() == currBG.Value()

	var fg color.Color

	if sameBackground {
		fgValue, exists := cfg.ThinFGColor[prevBG.Value()]
		if exists {
			fg = color.NewFGColor(fgValue)
		} else {
			fg = color.NewFGColor(prevBG.Value())
		}
	} else {
		fg = color.NewFGColor(prevBG.Value())
	}

	var content string

	switch d {
	case DirectionLeft:
		if sameBackground {
			content = cfg.Content.LeftThin
		} else {
			content = cfg.Content.Left
		}
	case DirectionRight:
		if sameBackground {
			content = cfg.Content.RightThin
		} else {
			content = cfg.Content.Right
		}
	}

	return &Separator{
		content: content,
		style:   color.NewStyle(fg, currBG),
	}, nil
}

// FinalSeparator returns the final separator given the last segment's style.
func FinalSeparator(d Direction, cfg SeparatorConfig, lastSegmentStyle color.Style) (*Separator, error) {
	var content string

	switch d {
	case DirectionLeft:
		content = cfg.Content.Left
	case DirectionRight:
		content = cfg.Content.Right
	}

	_, bg, err := lastSegmentStyle.SplitToColors()
	if err != nil {
		return nil, fmt.Errorf("unable to split colors from last segment style %q: %w", lastSegmentStyle, err)
	}

	return &Separator{
		content: content,
		style:   color.NewStyle(color.NewFGColor(bg.Value()), color.Color{}),
	}, err
}
