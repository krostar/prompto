package segment

import (
	"os"
	"path/filepath"
	"sort"
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
	path  string

	ColorForeground     uint8 `yaml:"fg"`
	ColorBackground     uint8 `yaml:"bg"`
	SeparatorForeground uint8 `yaml:"separator-fg"`

	ReplaceWith string `yaml:"replace-with"`
}

func (c *cwdConfig) getUsefulSpecial(cwd string) []cwdConfigSpecial {
	var specials []cwdConfigSpecial

	for path, special := range c.Special {
		if path == cwdSpecialLast { // this one's special, we handle it manually
			continue
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
			special.path = path
			specials = append(specials, special)
		}
	}

	// reorder specials by depth
	sort.Slice(specials, func(i, j int) bool {
		return specials[i].depth < specials[j].depth
	})

	return specials
}
