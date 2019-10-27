package color

// Kind defines where the color applies.
type Kind int

const (
	// KindUnknown is the zero value for ColorKind,
	// which is an invalid value.
	KindUnknown Kind = iota
	// KindForeground applies the color to the foreground.
	KindForeground
	// KindBackground applies the color to the background.
	KindBackground
)
