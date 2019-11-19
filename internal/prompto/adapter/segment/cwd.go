package segment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/krostar/prompto/pkg/pathx"

	"github.com/krostar/prompto/internal/prompto/domain"
)

const cwdSpecialLast = "-last-"

type cwd struct {
	cfg      cwdConfig
	cwd      string
	cwdDepth int
	specials []*cwdConfigSpecial
}

func segmentCWD(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(cwdConfig)
	if !isArgConfig {
		return nil, errors.New("segmentCWD expected 1 arg of type cwd.Config")
	}

	cfg.setDefaultColorToSpecials()

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get current working directory: %w", err)
	}

	wd = filepath.Clean(wd)

	return &cwd{
		cwd:      wd,
		cwdDepth: len(pathx.SplitPath(wd)),
		cfg:      cfg,
		specials: cfg.getUsefulSpecial(wd),
	}, nil
}

func (s *cwd) ProvideSegments() (domain.Segments, error) {
	var segments domain.Segments

	if err := pathx.WalkBackward(s.cwd, func(path string) error {
		if ss := s.pathToSegments(path); ss != nil {
			segments = append(segments, ss...)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to walk through cwd directories: %w", err)
	}

	segments.InverseOrder()

	return segments, nil
}

func (s *cwd) specialSegmentConfig(path string) (*cwdConfigSpecial, bool) {
	var (
		mostInterestingSpecial *cwdConfigSpecial
		dropSegment            bool
	)

	if last, exists := s.cfg.Special[cwdSpecialLast]; exists && path == s.cwd {
		mostInterestingSpecial = last
	}

	// fmt.Println("===", path)
	for _, special := range s.specials {
		// fmt.Println(special.path, special.depth, s.cwdDepth, strings.HasPrefix(special.path, path))
		// special should be taken in consideration only if its depth
		// is smaller than cwd's depth and has the same prefix
		if special.depth <= s.cwdDepth && strings.HasPrefix(special.path, path) {
			// drop segment as it will be replaced
			if path != special.path && special.ReplaceWith != "" {
				dropSegment = true
				break
			}

			special := special
			mostInterestingSpecial = special

			break // we found one special, that's enough :)
		}
	}

	return mostInterestingSpecial, dropSegment
}

func (s *cwd) pathToSegments(path string) domain.Segments {
	// default style and content
	content, style := filepath.Base(path), s.cfg.Color.ToStyle()

	special, dropSegment := s.specialSegmentConfig(path)
	if dropSegment {
		return nil
	}

	if special != nil {
		style = special.Color.ToStyle()

		if special.ReplaceWith != "" {
			content = special.ReplaceWith
		}
	}

	segments := splitContentInSegments(content, func(content string) *domain.Segment {
		return domain.
			NewSegment(content).
			SetStyle(style).
			WithSpaceAround()
	})
	segments.InverseOrder()

	return segments
}
