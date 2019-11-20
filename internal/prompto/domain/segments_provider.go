package domain

// SegmentsProvider defines how to provide segments.
type SegmentsProvider interface {
	SegmentName() string
	ProvideSegments() (Segments, error)
}

type leftPromptSegmentDisabler interface {
	DisabledForLeftPrompt() bool
}

type rightPromptSegmentDisabler interface {
	DisabledForRightPrompt() bool
}

func IsSegmentsProviderDisabledForDirection(d Direction, sp SegmentsProvider) bool {
	switch d {
	case DirectionLeft:
		disabled, ok := (sp).(leftPromptSegmentDisabler)
		return ok && disabled.DisabledForLeftPrompt()
	case DirectionRight:
		disabled, ok := (sp).(rightPromptSegmentDisabler)
		return ok && disabled.DisabledForRightPrompt()
	default:
		return true
	}
}
