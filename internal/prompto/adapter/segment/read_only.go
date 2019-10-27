package segment

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/krostar/prompto/internal/prompto/domain"
	"golang.org/x/sys/unix"
)

type readOnly struct {
	cfg readOnlyConfig
	cwd string
}

type readOnlyConfig struct {
	Content string `yaml:"content"`

	ColorForeground uint8 `yaml:"fg"`
	ColorBackground uint8 `yaml:"bg"`
}

func segmentReadOnly(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(readOnlyConfig)
	if !isArgConfig {
		return nil, errors.New("segmentReadOnly expected 1 arg of type readOnlyConfig")
	}

	pwd, isset := os.LookupEnv("PWD")
	if !isset {
		return nil, errors.New("pwd environment variable is not set")
	}

	return &readOnly{
		cfg: cfg,
		cwd: filepath.Clean(pwd),
	}, nil
}

func (s *readOnly) ProvideSegments() (domain.Segments, error) {
	var segments domain.Segments
	if unix.Access(s.cwd, unix.W_OK) != nil {
		segments = domain.Segments{
			domain.NewSegment(s.cfg.Content).
				WithSpaceAround().
				SetStyle(domain.NewStyle(
					domain.NewFGColor(s.cfg.ColorForeground),
					domain.NewBGColor(s.cfg.ColorBackground),
				)),
		}
	}

	return segments, nil
}
