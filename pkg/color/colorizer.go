package color

import (
	"fmt"
	"unicode/utf8"
)

type Colorizer interface {
	Colorize(style Style, format string, args ...interface{}) string
}

func NewColorizer(shell string) (Colorizer, error) {
	switch shell {
	case "zsh":
		return &zshColorizer{}, nil
	default:
		return &defaultANSIColorizer{}, nil
	}
}

type zshColorizer struct{}

func (zshColorizer) Colorize(style Style, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)

	return fmt.Sprintf(
		"%%{"+ansiEscapeSequenceStyleTemplate+"%s"+ansiEscapeSequenceReset+"%%%dG%%}",
		style.Code(),
		message,
		utf8.RuneCountInString(message),
	)
}

type defaultANSIColorizer struct{}

func (defaultANSIColorizer) Colorize(style Style, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf(
		ansiEscapeSequenceStyleTemplate+"%s"+ansiEscapeSequenceReset,
		style.Code(),
		message,
	)
}
