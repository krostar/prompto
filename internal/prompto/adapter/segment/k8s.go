package segment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/krostar/prompto/pkg/color"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type k8s struct {
	cfg       k8sConfig
	k8sConfig *api.Config
}

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

func segmentK8S(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(k8sConfig)
	if !isArgConfig {
		return nil, errors.New("segmentK8S expected 1 arg of type k8sConfig")
	}

	k8sConfig, err := clientcmd.LoadFromFile(cfg.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load k8s config file: %w", err)
	}

	cfg.setDefaultColorToSpecials()

	return &k8s{
		cfg:       cfg,
		k8sConfig: k8sConfig,
	}, nil
}

func (s *k8s) ProvideSegments() (domain.Segments, error) {
	current, exists := s.k8sConfig.Contexts[s.k8sConfig.CurrentContext]
	if !exists {
		return nil, nil
	}

	if !s.shouldDisplay() {
		return nil, nil
	}

	content, style, err := s.content(current.Cluster, current.Namespace)
	if err != nil {
		return nil, fmt.Errorf("unable to get content and style: %w", err)
	}

	return splitContentInSegments(content, func(content string) *domain.Segment {
		return domain.
			NewSegment(content).
			SetStyle(style).
			WithSpaceAround()
	}), nil
}

func (s *k8s) shouldDisplay() bool {
	if s.cfg.DisplayIf != "" {
		_, _, status, err := execCommand(s.cfg.DisplayIf)
		if err != nil || status != 0 {
			return false
		}
	}

	return true
}

func (s *k8s) content(cluster, namespace string) (string, color.Style, error) {
	content := filepath.Join(cluster, namespace)
	style := s.cfg.Color.ToStyle()
	tpl := make(map[string]string)

	for c, nm := range s.cfg.Special {
		match, err := createSubstituteTemplate(c, cluster, tpl)
		if err != nil {
			return "", color.Style{}, fmt.Errorf("unable to create template for cluster: %w", err)
		}

		if match {
			for n, spe := range nm {
				match, err := createSubstituteTemplate(n, namespace, tpl)
				if err != nil {
					return "", color.Style{}, fmt.Errorf("unable to create template for namespace: %w", err)
				}

				if match {
					if spe.ReplaceWith != "" {
						content = spe.ReplaceWith
					}

					style = spe.Color.ToStyle()

					break
				}
			}

			break
		}
	}

	substituteWithTemplate(&content, tpl)

	return content, style, nil
}
