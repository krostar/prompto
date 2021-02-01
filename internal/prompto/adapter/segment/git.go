package segment

import (
	"errors"
	"path/filepath"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type git struct {
	cfg gitConfig

	repository gitRepositoryGetter
}

type gitConfig struct {
	Ignore  []string       `yaml:"ignore"`
	Clean   gitStateConfig `yaml:"clean"`
	Changes gitStateConfig `yaml:"changes"`
}

type gitStateConfig struct {
	Color color.Config `yaml:"color"`
}

func segmentGIT(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(gitConfig)
	if !isArgConfig {
		return nil, errors.New("segmentGIT expected 1 arg of type gitConfig")
	}

	for i, ignore := range cfg.Ignore {
		cfg.Ignore[i] = filepath.Clean(replaceEnvironmentInPath(ignore))
	}

	return &git{
		cfg:        cfg,
		repository: &gitCommandRepository{},
	}, nil
}

func (s *git) SegmentName() string { return "git" }

func (s *git) ProvideSegments() (domain.Segments, error) {
	repo, err := s.repository.get()
	if err != nil {
		return nil, err
	}

	if repo == nil {
		return nil, nil
	}

	segment := domain.
		NewSegment(s.segmentContent(repo)...).
		SetStyle(s.segmentStyle(repo))

	if !repo.isSynced {
		segment.DisableSpaceAfter()
	}

	return domain.Segments{segment}, nil
}

func (s *git) segmentContent(repo *gitRepository) []string {
	content := []string{"", repo.branch}

	if !repo.isSynced {
		content = append(content, "•")
	}

	return content
}

func (s *git) segmentStyle(repo *gitRepository) color.Style {
	var cfg gitStateConfig

	if repo.hasWIP {
		cfg = s.cfg.Changes
	} else {
		cfg = s.cfg.Clean
	}

	return cfg.Color.ToStyle()
}
