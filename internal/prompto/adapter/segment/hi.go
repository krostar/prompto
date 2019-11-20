package segment

import (
	"errors"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type hi struct {
	cfg hiConfig
}

type hiConfig struct {
	Content string       `yaml:"content"`
	Color   color.Config `yaml:"color"`
}

func segmentHi(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(hiConfig)
	if !isArgConfig {
		return nil, errors.New("segmentHi expected 1 arg of type hiConfig")
	}

	return &hi{cfg: cfg}, nil
}

func (s *hi) SegmentName() string {
	return "hi"
}

func (s *hi) ProvideSegments() (domain.Segments, error) {
	return domain.Segments{
		domain.
			NewSegment(s.cfg.Content).
			WithSpaceAround().
			SetStyle(s.cfg.Color.ToStyle()),
	}, nil
}
