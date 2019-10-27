package domain

import (
	"fmt"
	"io"

	"github.com/krostar/prompto/pkg/color"
)

// Separator stores separator components.
type Separator struct {
	style   color.Style
	content string
}

// WriteTo implements io.WriterTo for Separator a segment's separator with style.
func (s *Separator) WriteTo(colorizer color.Colorizer, w io.Writer) (int64, error) {
	n, err := w.Write([]byte(colorizer.Colorize(s.style, s.content)))
	return int64(n), err
}

// SeparatorConfig defines how to configure separators.
type SeparatorConfig struct {
	Content     SeparatorContentConfig `yaml:"content"`
	ThinFGColor map[string]string      `yaml:"thin-fg-color"`
}

// SeparatorContentConfig defines the separator configuration.
type SeparatorContentConfig struct {
	Left      string `yaml:"left"`
	LeftThin  string `yaml:"left-thin"`
	Right     string `yaml:"right"`
	RightThin string `yaml:"right-thin"`
}

// NewSeparator creates a new separator given the previous and current segment's style.
func NewSeparator(d Direction, cfg SeparatorConfig, previous, current color.Style) (*Separator, error) {
	_, prevBG, err := previous.Colors()
	if err != nil {
		return nil, fmt.Errorf("unable to get colors from previous segment style: %w", err)
	}

	_, currBG, err := current.Colors()
	if err != nil {
		return nil, fmt.Errorf("unable to get colors from current segment style: %w", err)
	}

	var (
		sameBackground = currBG.Equal(prevBG)
		fg             color.Color
	)

	if sameBackground {
		fgValue, exists := cfg.ThinFGColor[prevBG.String()]
		if exists {
			fg = color.HexFGColor(fgValue)
		} else {
			fg = color.HexFGColor(prevBG.String())
		}
	} else {
		fg = color.HexFGColor(prevBG.String())
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

	_, bg, err := lastSegmentStyle.Colors()
	if err != nil {
		return nil, fmt.Errorf("unable to split colors from last segment style %q: %w", lastSegmentStyle, err)
	}

	return &Separator{
		content: content,
		style:   color.NewStyle(color.HexFGColor(bg.String()), color.Color{}),
	}, err
}
