package domain

import (
	"io"

	"github.com/krostar/prompto/pkg/color"
)

// Separator stores separator components.
type Separator struct {
	style   color.Style
	content string
}

// WriteTo writes to the provided writer the stylized separator.
func (s *Separator) WriteTo(colorizer color.Colorizer, w io.Writer) (int64, error) {
	n, err := w.Write([]byte(colorizer.Colorize(s.style, s.content)))
	return int64(n), err
}

// SeparatorConfig defines how to configure separators.
type SeparatorConfig struct {
	Content     SeparatorContentConfig `yaml:"content"`
	ThinFGColor map[string]string      `yaml:"thin-fg-color"`
}

// SeparatorContentConfig defines the separator content configuration.
type SeparatorContentConfig struct {
	Left      string `yaml:"left"`
	LeftThin  string `yaml:"left-thin"`
	Right     string `yaml:"right"`
	RightThin string `yaml:"right-thin"`
}

// NewSeparator creates a new separator given the previous and current segment's style.
func NewSeparator(d Direction, cfg SeparatorConfig, previous, current color.Style) *Separator {
	var (
		_, prevBG                 = previous.Colors()
		prevBGR, prevBGG, prevBGB = prevBG.RGB()

		_, currBG      = current.Colors()
		sameBackground = currBG.Equal(prevBG)

		fg color.Color
	)

	if sameBackground {
		fgValue, exists := cfg.ThinFGColor[prevBG.String()]
		if exists {
			fg = color.NewHexFGColor(fgValue)
		} else {
			fg = color.NewRGBFGColor(prevBGR, prevBGG, prevBGB)
		}
	} else {
		fg = color.NewRGBFGColor(prevBGR, prevBGG, prevBGB)
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
	}
}

// FinalSeparator returns the final separator given the last segment's style.
func FinalSeparator(d Direction, cfg SeparatorConfig, lastSegmentStyle color.Style) *Separator {
	var content string

	switch d {
	case DirectionLeft:
		content = cfg.Content.Left
	case DirectionRight:
		content = cfg.Content.Right
	}

	_, bg := lastSegmentStyle.Colors()
	r, g, b := bg.RGB()

	return &Separator{
		content: content,
		style:   color.NewStyle(color.NewRGBFGColor(r, g, b), color.Color{}),
	}
}
