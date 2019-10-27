// Package app exposes app-related information that are injected at compile time.
package app

import (
	"time"
	"unicode"

	"github.com/pkg/errors"
)

// build injection requires these global variables
// nolint: gochecknoglobals
var (
	name             = "app"                  // app name
	version          = "0.0.0-0-gmaster"      // see https://semver.org/ to have a description of the format
	buildAtRaw       = "1970-01-01T00:00:00Z" // build date in RFC3339 format
	builtAt          time.Time                // set based on buildAtRaw by init()
	alphaNumericName string                   // set based on name by init()
)

// this is required as it depends on build-time variable injection
// nolint: gochecknoinits
func init() {
	t, err := time.Parse(time.RFC3339, buildAtRaw)
	if err != nil {
		panic(errors.Wrap(err, "buildAtRaw must be injected with RFC3339 format"))
	}

	builtAt = t

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			alphaNumericName += string(r)
		}
	}
}

// Name returns the app name as set at build injection.
func Name() string { return name }

// AlphaNumericName returns the app name with only alphanumeric caracters.
func AlphaNumericName() string { return alphaNumericName }

// Version returns the app version as set at build injection.
func Version() string { return version }

// BuiltAt returns the app built time as set at build injection.
func BuiltAt() time.Time { return builtAt }
