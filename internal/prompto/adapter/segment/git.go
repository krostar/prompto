package segment

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/krostar/prompto/pkg/gitx"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type git struct {
	cfg  gitConfig
	repo *gitx.Repository
}

type gitConfig struct {
	Ignore  []string       `yaml:"ignore"`
	Clean   gitStateConfig `yaml:"clean"`
	Changes gitStateConfig `yaml:"changes"`
}

type gitStateConfig struct {
	ColorForeground uint8 `yaml:"fg"`
	ColorBackground uint8 `yaml:"bg"`
}

func segmentGIT(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(gitConfig)
	if !isArgConfig {
		return nil, errors.New("segmentGIT expected 1 arg of type gitConfig")
	}

	repo, err := gitx.LocalRepository()
	if err != nil {
		if !errors.Is(err, gitx.ErrRepositoryDoesNotExists) {
			return nil, fmt.Errorf("unable to open local repository: %w", err)
		}
	}

	for i, ignore := range cfg.Ignore {
		cfg.Ignore[i] = filepath.Clean(replaceEnvironmentInPath(ignore))
	}

	return &git{
		cfg:  cfg,
		repo: repo,
	}, nil
}

func (s *git) ProvideSegments() (domain.Segments, error) {
	useless, err := s.isRepositoryMissingOrIgnored()
	if err != nil {
		return nil, fmt.Errorf("unable to check if repository should be looked at: %w", err)
	}
	if useless {
		return nil, nil
	}

	fg, bg, err := s.getColorDependingOnRepositoryState()
	if err != nil {
		return nil, fmt.Errorf("unable to get segment colors: %w", err)
	}

	headReference, err := s.repo.HeadReference()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve actual branch name: %w", err)
	}

	content := headReference + ""

	syncedWithRemote, err := s.repo.IsBranchSyncedWithRemote(headReference)
	if err != nil {
		return nil, fmt.Errorf("unable to check if synced with remote: %w", err)
	}
	if !syncedWithRemote {
		content += "•"
	}

	return domain.Segments{domain.
		NewSegment(content).
		SetStyle(domain.NewStyle(fg, bg)).WithSpaceAround(),
	}, nil
}

func (s *git) isRepositoryMissingOrIgnored() (bool, error) {
	if s.repo == nil {
		return true, nil
	}
	repoLocation, err := s.repo.AbsoluteLocation()
	if err != nil {
		return false, fmt.Errorf("unable to get absolute repository location: %w", err)
	}
	for _, ignore := range s.cfg.Ignore {
		if repoLocation == ignore {
			return true, nil
		}
	}
	return false, nil
}

func (s *git) getColorDependingOnRepositoryState() (domain.Color, domain.Color, error) {
	hasWIP, err := s.repo.HasWIP()
	if err != nil {
		return domain.Color{}, domain.Color{}, fmt.Errorf("unable to check work in progress status: %w", err)
	}

	var cfg gitStateConfig
	if hasWIP {
		cfg = s.cfg.Changes
	} else {
		cfg = s.cfg.Clean
	}

	return domain.NewFGColor(cfg.ColorForeground), domain.NewBGColor(cfg.ColorBackground), nil

}
