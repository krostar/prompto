package segment

import (
	"errors"
	"time"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type lastCMDExecTime struct {
	cfg lastCMDExecTimeConfig
}

type lastCMDExecTimeConfig struct {
	DurationNS uint64 `json:"-" yaml:"-"`

	Color            color.Config                    `yaml:"color"`
	TresholdDisplay  time.Duration                   `yaml:"treshold-display"`
	TresholdTruncate map[time.Duration]time.Duration `yaml:"treshold-truncate"`
}

func segmentLastCMDExecTime(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(lastCMDExecTimeConfig)
	if !isArgConfig {
		return nil, errors.New("segmentLastCMDExecTime expected 1 arg of type lastCMDExecTimeConfig")
	}

	return &lastCMDExecTime{cfg: cfg}, nil
}

func (s *lastCMDExecTime) SegmentName() string { return "last command exec time" }

func (s *lastCMDExecTime) ProvideSegments() (domain.Segments, error) {
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
		domain.
			NewSegment(d.String()).
			SetStyle(s.cfg.Color.ToStyle()),
	}, nil
}
