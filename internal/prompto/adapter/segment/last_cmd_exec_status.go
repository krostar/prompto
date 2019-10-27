package segment

import (
	"errors"
	"strconv"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type lastCMDExecStatus struct {
	cfg lastCMDExecStatusConfig
}

type lastCMDExecStatusConfig struct {
	StatusCode uint `json:"-" yaml:"-"`

	Success lastCMDExecStatusStateConfig `yaml:"success"`
	Failure lastCMDExecStatusStateConfig `yaml:"failure"`
}

type lastCMDExecStatusStateConfig struct {
	ReplaceWith string       `yaml:"replace-with"`
	Hide        bool         `yaml:"hide"`
	Color       color.Config `yaml:"color"`
}

func segmentLastCMDExecStatus(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(lastCMDExecStatusConfig)
	if !isArgConfig {
		return nil, errors.New("segmentLastCMDExecStatus expected 1 arg of type lastCMDExecStatusConfig")
	}
	return &lastCMDExecStatus{cfg: cfg}, nil
}

func (s *lastCMDExecStatus) SegmentName() string { return "last command exec status" }

func (s *lastCMDExecStatus) ProvideSegments() (domain.Segments, error) {
	var cfg lastCMDExecStatusStateConfig
	if s.cfg.StatusCode == 0 {
		cfg = s.cfg.Success
	} else {
		cfg = s.cfg.Failure
	}

	if cfg.Hide {
		return nil, nil
	}

	content := strconv.FormatUint(uint64(s.cfg.StatusCode), 10)
	if cfg.ReplaceWith != "" {
		content = cfg.ReplaceWith
	}

	return domain.Segments{
		domain.
			NewSegment(content).
			SetStyle(cfg.Color.ToStyle()).
			WithSpaceAround(),
	}, nil
}
