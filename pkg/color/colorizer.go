package color

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Colorizer defines how to colorize something with style.
type Colorizer interface {
	Colorize(style Style, format string, args ...interface{}) string
}

// NewColorizer returns a new colorizer based on the provided shell.
func NewColorizer(shell string) (Colorizer, error) {
	switch shell {
	case "zsh":
		return &zshANSIColorizer{}, nil
	case "fish", "bash":
		fallthrough
	case "":
		return &defaultANSIColorizer{}, nil
	default:
		return nil, fmt.Errorf("unhandled shell %q", shell)
	}
}

// NoopColorizer implements Colorizer without doing anything.
type NoopColorizer struct{}

// Colorize does nothing.
func (NoopColorizer) Colorize(style Style, format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

type defaultANSIColorizer struct{}

func (defaultANSIColorizer) Colorize(style Style, format string, args ...interface{}) string {
	return style.ANSISprintf(format, args...)
}

type zshANSIColorizer struct{}

func (zshANSIColorizer) Colorize(style Style, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	return "%{" + style.ANSISprintf(message) + "%" + strconv.Itoa(utf8.RuneCountInString(message)) + "G%}"
}
