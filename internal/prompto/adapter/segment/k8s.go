package segment

import (
	"errors"
	"fmt"
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

func (s *k8s) SegmentName() string { return "k8s" }

func (s *k8s) ProvideSegments() (domain.Segments, error) {
	current, exists := s.k8sConfig.Contexts[s.k8sConfig.CurrentContext]
	if !exists {
		return nil, nil
	}

	tpl, content, style, err := s.content(current.Cluster, current.Namespace)
	if err != nil {
		return nil, fmt.Errorf("unable to get content and style: %w", err)
	}

	if !s.shouldDisplay(tpl) {
		return nil, nil
	}

	return splitContentInSegments(content, func(content string) *domain.Segment {
		return domain.
			NewSegment(content).
			SetStyle(style)
	}), nil
}

func (s *k8s) shouldDisplay(tpl map[string]string) bool {
	if s.cfg.DisplayIf != "" {
		cmd := s.cfg.DisplayIf
		substituteWithTemplate(&cmd, tpl)

		if _, _, status, err := execCommand(cmd); err != nil || status != 0 {
			return false
		}
	}

	return true
}

func (s *k8s) content(cluster, namespace string) (map[string]string, string, color.Style, error) {
	content := filepath.Join(cluster, namespace)
	style := s.cfg.Color.ToStyle()
	tpl := make(map[string]string)

	for c, nm := range s.cfg.Special {
		match, err := createSubstituteTemplate(c, cluster, tpl)
		if err != nil {
			return nil, "", color.Style{}, fmt.Errorf("unable to create template for cluster: %w", err)
		}

		if match {
			for n, spe := range nm {
				match, err := createSubstituteTemplate(n, namespace, tpl)
				if err != nil {
					return nil, "", color.Style{}, fmt.Errorf("unable to create template for namespace: %w", err)
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

	return tpl, content, style, nil
}
