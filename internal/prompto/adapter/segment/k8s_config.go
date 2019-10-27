package segment

import (
	"os"
	"path/filepath"

	"github.com/krostar/prompto/pkg/color"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sConfig struct {
	ConfigFile string                                  `yaml:"config-file"`
	Color      color.Config                            `yaml:"color"`
	DisplayIf  string                                  `yaml:"display-if"`
	Special    map[string]map[string]*k8sConfigSpecial `yaml:"special"`
}

func (c *k8sConfig) setDefaultColorToSpecials() {
	for _, cluster := range c.Special {
		for _, namespace := range cluster {
			namespace.Color.SetDefaultColor(c.Color)
		}
	}
}

type k8sConfigSpecial struct {
	Color       color.Config `yaml:"color"`
	ReplaceWith string       `yaml:"replace-with"`
}

func (c *k8sConfig) SetDefault() {
	if home, err := os.UserHomeDir(); err == nil {
		c.ConfigFile = filepath.Join(
			home,
			clientcmd.RecommendedHomeDir,
			clientcmd.RecommendedFileName,
		)
	}
}
