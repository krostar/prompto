package domain

// Direction defines the possible direction for a prompt to point to.
type Direction uint

const (
	// DirectionUnknown is the zero value for Direction, which is an invalid value.
	DirectionUnknown Direction = iota
	// DirectionLeft defines the direction that goes from left to right.
	DirectionLeft
	// DirectionRight defines the direction that goes from right to left.
	DirectionRight
)

// String implements stringer.
func (d Direction) String() string {
	switch d {
	case DirectionLeft:
		return "left"
	case DirectionRight:
		return "right"
	default:
		return "unknown"
	}
}
