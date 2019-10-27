package segment

import (
	"github.com/krostar/prompto/internal/prompto/domain"
)

type newline struct{}

func segmentNewline(interface{}) (domain.SegmentsProvider, error) {
	return &newline{}, nil
}

func (s *newline) SegmentName() string { return "newline" }

func (s *newline) DisabledForRightPrompt() bool { return true }

func (s *newline) ProvideSegments() (domain.Segments, error) {
	return domain.Segments{domain.NewSegment("\n").DisableNextSegmentSeparator()}, nil
}
