package color

import (
	"fmt"
	"unicode/utf8"
)

// Colorizer defines how to colorize color.
type Colorizer interface {
	Colorize(style Style, format string, args ...interface{}) string
}

// NewColorizer returns a new colorizer based on the provided shell.
// nolint: unparam
func NewColorizer(shell string) (Colorizer, error) {
	switch shell {
	case "zsh":
		return &zshColorizer{}, nil
	default:
		return &defaultANSIColorizer{}, nil
	}
}

type defaultANSIColorizer struct{}

func (defaultANSIColorizer) Colorize(style Style, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf(
		ansiEscapeSequenceStyleTemplate+"%s"+ansiEscapeSequenceReset,
		style.EscapeSequence(),
		message,
	)
}

type zshColorizer struct{}

func (zshColorizer) Colorize(style Style, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf(
		"%%{"+ansiEscapeSequenceStyleTemplate+"%s"+ansiEscapeSequenceReset+"%%%dG%%}",
		style.EscapeSequence(),
		message,
		utf8.RuneCountInString(message),
	)
}
