package segment

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/krostar/prompto/pkg/pathx"
)

type cwdConfig struct {
	ColorForeground          uint8 `yaml:"fg"`
	ColorBackground          uint8 `yaml:"bg"`
	SeparatorForegroundColor uint8 `yaml:"separator-fg"`

	Special map[string]cwdConfigSpecial `yaml:"special"`
}

type cwdConfigSpecial struct {
	depth int

	ColorForeground     uint8 `yaml:"fg"`
	ColorBackground     uint8 `yaml:"bg"`
	SeparatorForeground uint8 `yaml:"separator-fg"`

	ReplaceWith string `yaml:"replace-with"`
}

func (c *cwdConfig) keepUsefulSpecialOnly(cwd string) {
	var (
		specials        = make(map[string]cwdConfigSpecial)
		maxSpecialDepth int
	)

	for path, special := range c.Special {
		if path == cwdSpecialLast { // this one's special, keep it
			specials[path] = special
		}

		// replace all environment-based specials
		pathSplit := pathx.SplitPath(path)
		for i, split := range pathSplit {
			if strings.HasPrefix(split, "$") {
				if p, isset := os.LookupEnv(split[1:]); isset && p != "" {
					pathSplit[i] = p
				}
			}
		}

		path = filepath.Join(pathSplit...)

		// add only the special usable with the cwd
		if strings.HasPrefix(cwd, path) {
			special.depth = len(pathx.SplitPath(path))
			if special.depth > maxSpecialDepth {
				maxSpecialDepth = special.depth
				specials[path] = special
			}
		}
	}

	// remove special with smaller depth
	for path, special := range specials {
		if path == cwdSpecialLast {
			continue
		}

		if special.depth != maxSpecialDepth {
			delete(specials, path)
		}
	}

	c.Special = specials
}
