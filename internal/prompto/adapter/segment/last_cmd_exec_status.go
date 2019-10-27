package segment

import (
	"errors"
	"strconv"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type lastCmdExecStatus struct {
	cfg lastCmdExecStatusConfig
}

type lastCmdExecStatusConfig struct {
	StatusCode uint `yaml:"-"`

	Success lastCmdExecStatusStateConfig `yaml:"success"`
	Failure lastCmdExecStatusStateConfig `yaml:"failure"`
}

type lastCmdExecStatusStateConfig struct {
	ReplaceWith     string `yaml:"replace-with"`
	Hide            bool   `yaml:"hide"`
	ColorBackground uint8  `yaml:"bg"`
	ColorForeground uint8  `yaml:"fg"`
}

func segmentLastCmdExecStatus(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(lastCmdExecStatusConfig)
	if !isArgConfig {
		return nil, errors.New("segmentLastCmdExecStatus expected 1 arg of type lastCmdExecStatusConfig")
	}

	return &lastCmdExecStatus{
		cfg: cfg,
	}, nil
}

func (s *lastCmdExecStatus) ProvideSegments() (domain.Segments, error) {
	var cfg lastCmdExecStatusStateConfig
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
		domain.NewSegment(content).
			WithSpaceAround().
			SetStyle(domain.NewStyle(
				domain.NewFGColor(cfg.ColorForeground),
				domain.NewBGColor(cfg.ColorBackground),
			)),
	}, nil
}
