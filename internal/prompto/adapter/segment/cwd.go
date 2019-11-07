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
	specials []cwdConfigSpecial
}

func segmentCWD(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(cwdConfig)
	if !isArgConfig {
		return nil, errors.New("segmentCWD expected 1 arg of type cwd.Config")
	}

	pwd, isset := os.LookupEnv("PWD")
	if !isset {
		return nil, errors.New("pwd environment variable is not set")
	}

	pwd = filepath.Clean(pwd)

	return &cwd{
		cwd:      pwd,
		cwdDepth: len(pathx.SplitPath(pwd)),
		cfg:      cfg,
		specials: cfg.getUsefulSpecial(pwd),
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

func (s *cwd) defaultSegmentComponents(path string) (string, domain.Style, uint8) {
	return filepath.Base(path),
		domain.NewStyle(
			domain.NewFGColor(s.cfg.ColorForeground),
			domain.NewBGColor(s.cfg.ColorBackground),
		),
		s.cfg.SeparatorForegroundColor
}

func (s *cwd) specialSegmentConfig(path string) (*cwdConfigSpecial, bool) {
	var (
		mostInterestingSpecial *cwdConfigSpecial
		dropSegment            bool
	)

	if last, exists := s.cfg.Special[cwdSpecialLast]; exists && path == s.cwd {
		mostInterestingSpecial = &last
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
			mostInterestingSpecial = &special

			break // we found one special, that's enough :)
		}
	}

	return mostInterestingSpecial, dropSegment
}

func (s *cwd) pathToSegments(path string) domain.Segments {
	// default style and content
	content, style, sepFGColor := s.defaultSegmentComponents(path)

	special, dropSegment := s.specialSegmentConfig(path)
	if dropSegment {
		return nil
	}

	if special != nil {
		style = domain.NewStyle(
			domain.NewFGColor(special.ColorForeground),
			domain.NewBGColor(special.ColorBackground),
		)
		sepFGColor = special.SeparatorForeground

		if special.ReplaceWith != "" {
			content = special.ReplaceWith
		}
	}

	var segments domain.Segments

	// in case content is separated with slashes, create multiple segments
	paths := pathx.SplitPath(content)
	for i := range paths {
		segments = append(segments, domain.
			NewSegment(paths[len(paths)-i-1]).
			WithSpaceAround().
			SetStyle(style).
			SetSeparatorColor(domain.NewFGColor(sepFGColor)),
		)
	}

	return segments
}
