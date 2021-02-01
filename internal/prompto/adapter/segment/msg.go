package segment

import (
	"errors"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type msg struct {
	cfg msgConfig
}

type msgConfig struct {
	Content string       `yaml:"content"`
	Color   color.Config `yaml:"color"`
}

func segmentMsg(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(msgConfig)
	if !isArgConfig {
		return nil, errors.New("segmentMsg expected 1 arg of type msgConfig")
	}

	return &msg{cfg: cfg}, nil
}

func (s *msg) SegmentName() string { return "msg" }

func (s *msg) ProvideSegments() (domain.Segments, error) {
	return domain.Segments{
		domain.
			NewSegment(s.cfg.Content).
			SetStyle(s.cfg.Color.ToStyle()),
	}, nil
}
