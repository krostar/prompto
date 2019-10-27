package segment

import (
	"errors"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type stub struct {
	cfg StubConfig
}

// StubConfig stores the results to return when calling the ProvideSegments method.
type StubConfig struct {
	Segments domain.Segments
	Error    error
}

func segmentStub(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(StubConfig)
	if !isArgConfig {
		return nil, errors.New("segmentstub expected 1 arg of type stubConfig")
	}
	return &stub{cfg: cfg}, nil
}

func (s *stub) SegmentName() string { return "stub" }

func (s *stub) ProvideSegments() (domain.Segments, error) { return s.cfg.Segments, s.cfg.Error }
