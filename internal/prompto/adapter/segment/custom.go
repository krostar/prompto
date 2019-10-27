package segment

import (
	"errors"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type custom struct {
	cfg customConfig
}

type customConfig struct {
	Exec     string                       `yaml:"exec"`
	Color    color.Config                 `yaml:"color"`
	Statuses map[uint8]customStatusConfig `yaml:"status"`
}

type customStatusConfig struct {
	Hide  bool         `yaml:"hide"`
	Color color.Config `yaml:"color"`
}

func segmentCustom(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(customConfig)
	if !isArgConfig {
		return nil, errors.New("segmentCustom expected 1 arg of type customConfig")
	}
	return &custom{cfg: cfg}, nil
}

func (s *custom) SegmentName() string { return "custom" }

func (s *custom) ProvideSegments() (domain.Segments, error) {
	stdout, _, status, err := execCommand(s.cfg.Exec)
	if err != nil {
		return nil, err
	}

	var style color.Style

	if statusCfg, isset := s.cfg.Statuses[status]; isset {
		if statusCfg.Hide {
			return nil, nil
		}

		statusCfg.Color.SetDefaultColor(s.cfg.Color)

		style = statusCfg.Color.ToStyle()
	} else {
		style = s.cfg.Color.ToStyle()
	}

	return splitContentInSegments(stdout, func(content string) *domain.Segment {
		return domain.
			NewSegment(content).
			SetStyle(style).
			WithSpaceAround()
	}), nil
}
