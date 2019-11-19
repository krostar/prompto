package segment

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/krostar/prompto/pkg/color"
	"github.com/krostar/prompto/pkg/gitx"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type git struct {
	cfg gitConfig
}

type gitConfig struct {
	Ignore  []string       `yaml:"ignore"`
	Lite    []string       `yaml:"lite"`
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

	return &git{cfg: cfg}, nil
}

func (s *git) ProvideSegments() (domain.Segments, error) {
	repo, lite, err := s.getRepository()
	if err != nil {
		return nil, err
	}

	if repo == nil {
		return nil, nil
	}

	content, style, err := s.getDefaultComponents(repo)
	if err != nil {
		return nil, err
	}

	var syncedWithRemote bool
	if !lite {
		syncedWithRemote, err = repo.IsHeadSyncedWithRemote()
		if err != nil {
			return nil, fmt.Errorf("unable to check if synced with remote: %w", err)
		}

		if !syncedWithRemote {
			content = append(content, "•")
		}
	}

	segment := domain.
		NewSegment(content...).
		SetStyle(style).
		WithSpaceBefore()

	if syncedWithRemote {
		segment.WithSpaceAfter()
	}

	return domain.Segments{segment}, nil
}

func (s *git) getRepository() (*gitx.Repository, bool, error) {
	repo, err := gitx.LocalRepository()
	if err != nil {
		if errors.Is(err, gitx.ErrRepositoryDoesNotExists) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("unable to open local repository: %w", err)
	}

	repoLocation, err := repo.AbsoluteLocation()
	if err != nil {
		return nil, false, fmt.Errorf("unable to get absolute repository location: %w", err)
	}

	for _, i := range s.cfg.Ignore {
		if repoLocation == i {
			return nil, false, nil
		}
	}

	var lite bool

	for _, l := range s.cfg.Lite {
		if repoLocation == l {
			lite = true
			break
		}
	}

	return repo, lite, nil
}

func (s *git) getDefaultComponents(repo *gitx.Repository) ([]string, color.Style, error) {
	style, err := s.getStyleDependingOnRepositoryState(repo)
	if err != nil {
		return nil, color.Style{}, fmt.Errorf("unable to get segment style: %w", err)
	}

	headReference, err := repo.HeadReference()
	if err != nil {
		return nil, color.Style{}, fmt.Errorf("unable to retrieve actual branch name: %w", err)
	}

	return []string{"", headReference}, style, nil
}

func (s *git) getStyleDependingOnRepositoryState(repo *gitx.Repository) (color.Style, error) {
	hasWIP, err := repo.HasWIP()
	if err != nil {
		return color.Style{}, fmt.Errorf("unable to check work in progress status: %w", err)
	}

	var cfg gitStateConfig
	if hasWIP {
		cfg = s.cfg.Changes
	} else {
		cfg = s.cfg.Clean
	}

	return cfg.Color.ToStyle(), nil
}
