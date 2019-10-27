package segment

import (
	"sort"
	"strings"

	"github.com/krostar/prompto/pkg/color"
	"github.com/krostar/prompto/pkg/pathx"
)

type cwdConfig struct {
	Color   color.Config                 `yaml:"color"`
	Special map[string]*cwdConfigSpecial `yaml:"special"`
}

func (c *cwdConfig) setDefaultColorToSpecials() {
	for _, special := range c.Special {
		special.Color.SetDefaultColor(c.Color)
	}
}

type cwdConfigSpecial struct {
	depth int
	path  string

	Color       color.Config `yaml:"color"`
	ReplaceWith string       `yaml:"replace-with"`
}

func (c *cwdConfig) getUsefulSpecial(cwd string) []*cwdConfigSpecial {
	var specials []*cwdConfigSpecial

	for path, special := range c.Special {
		if path == cwdSpecialLast { // this one's special, we handle it manually
			continue
		}

		// replace all environment-based specials
		path = replaceEnvironmentInPath(path)

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
