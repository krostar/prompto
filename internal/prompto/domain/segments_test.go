package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegments_InverseOrder(t *testing.T) {
	s1, s2, s3 := NewSegment("1"), NewSegment("2"), NewSegment("3")
	ss := Segments{s1, s2, s3}
	ss.InverseOrder()

	assert.Equal(t, Segments{s3, s2, s1}, ss)
}
