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
	cfg.keepUsefulSpecialOnly(pwd)

	return &cwd{
		cwd:      pwd,
		cwdDepth: len(pathx.SplitPath(pwd)),
		cfg:      cfg,
	}, nil
}

func (s *cwd) ProvideSegments() (domain.Segments, error) {
	var segments domain.Segments

	if err := pathx.WalkBackward(s.cwd, func(path string) error {
		if segment := s.pathToSegment(path); segment != nil {
			segments = append(segments, segment)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to walk through cwd directories: %w", err)
	}

	segments.InverseOrder()

	return segments, nil
}

func (s *cwd) pathToSegment(path string) *domain.Segment {
	// default style and content
	style := domain.NewStyle(
		domain.NewFGColor(s.cfg.ColorForeground),
		domain.NewBGColor(s.cfg.ColorBackground),
	)
	separatorColorFG := s.cfg.SeparatorForegroundColor
	content := filepath.Base(path)

	if last, exists := s.cfg.Special[cwdSpecialLast]; exists && path == s.cwd {
		style = domain.NewStyle(
			domain.NewFGColor(last.ColorForeground),
			domain.NewBGColor(last.ColorBackground),
		)
	}

	for p, special := range s.cfg.Special {
		// special should be taken in consideration only if
		//	its depth is smaller than cwd's depth
		//	it has the same prefix
		if special.depth <= s.cwdDepth && strings.HasPrefix(p, path) {
			// drop segment that will be replaced
			if path != p && special.ReplaceWith != "" {
				return nil
			}

			// otherwise apply the desired style
			style = domain.NewStyle(
				domain.NewFGColor(special.ColorForeground),
				domain.NewBGColor(special.ColorBackground),
			)
			separatorColorFG = special.SeparatorForeground

			if path == p && special.ReplaceWith != "" {
				content = special.ReplaceWith
			}

			break // we found one special, that's enough :)
		}
	}

	return domain.NewSegment(content).
		WithSpaceAround().
		SetStyle(style).
		SetSeparatorColor(domain.NewFGColor(separatorColorFG))
}
