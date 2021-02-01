package segment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/krostar/prompto/pkg/color"
	"golang.org/x/sys/unix"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type readOnly struct {
	cfg readOnlyConfig
	cwd string
}

type readOnlyConfig struct {
	Content string       `yaml:"content"`
	Color   color.Config `yaml:"color"`
}

func segmentReadOnly(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(readOnlyConfig)
	if !isArgConfig {
		return nil, errors.New("segmentReadOnly expected 1 arg of type readOnlyConfig")
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get working directory: %w", err)
	}

	return &readOnly{
		cfg: cfg,
		cwd: filepath.Clean(wd),
	}, nil
}

func (s *readOnly) SegmentName() string { return "read-only" }

func (s *readOnly) ProvideSegments() (domain.Segments, error) {
	var segments domain.Segments

	if unix.Access(s.cwd, unix.W_OK) != nil {
		segments = domain.Segments{
			domain.
				NewSegment(s.cfg.Content).
				SetStyle(s.cfg.Color.ToStyle()),
		}
	}

	return segments, nil
}
