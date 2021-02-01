package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirection_String(t *testing.T) {
	tests := []struct {
		d            Direction
		expectedRepr string
	}{
		{
			DirectionLeft,
			"left",
		}, {
			DirectionRight,
			"right",
		}, {
			Direction(42),
			"unknown",
		}, {
			Direction(0),
			"unknown",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedRepr, test.d.String())
	}
}
