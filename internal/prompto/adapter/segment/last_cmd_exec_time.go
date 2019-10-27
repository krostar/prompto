package segment

import (
	"errors"
	"time"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type lastCmdExecTime struct {
	cfg lastCmdExecTimeConfig
}

type lastCmdExecTimeConfig struct {
	DurationNS uint `yaml:"-"`

	ColorForeground uint8 `yaml:"fg"`
	ColorBackground uint8 `yaml:"bg"`

	TresholdDisplay  time.Duration                   `yaml:"treshold-display"`
	TresholdTruncate map[time.Duration]time.Duration `yaml:"treshold-truncate"`
}

func segmentLastCmdExecTime(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(lastCmdExecTimeConfig)
	if !isArgConfig {
		return nil, errors.New("segmentLastCmdExecTime expected 1 arg of type lastCmdExecTimeConfig")
	}

	return &lastCmdExecTime{
		cfg: cfg,
	}, nil
}

func (s *lastCmdExecTime) ProvideSegments() (domain.Segments, error) {
	d := time.Duration(s.cfg.DurationNS)

	if d < s.cfg.TresholdDisplay {
		return nil, nil
	}

	for treshold, truncate := range s.cfg.TresholdTruncate {
		if d > treshold {
			d = d.Truncate(truncate)
		}
	}

	return domain.Segments{
		domain.NewSegment(d.String()).
			WithSpaceAround().
			SetStyle(domain.NewStyle(
				domain.NewFGColor(s.cfg.ColorForeground),
				domain.NewBGColor(s.cfg.ColorBackground),
			)),
	}, nil
}
