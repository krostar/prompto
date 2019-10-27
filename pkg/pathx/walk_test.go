package pathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WalkSplit(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected []string
	}{
		"empty path": {
			path:     "",
			expected: []string{},
		}, "one element": {
			path:     "a",
			expected: []string{"a"},
		}, "multiple elements": {
			path:     "a/b/c",
			expected: []string{"a", "b", "c"},
		}, "cleaned path": {
			path:     "///a",
			expected: []string{"/", "a"},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			var i int
			err := WalkSplit(test.path, func(p string) error {
				require.False(t, i >= len(test.expected))
				assert.Equal(t, test.expected[i], p)
				i++
				return nil
			})
			assert.Equal(t, len(test.expected), i)
			require.NoError(t, err)
		})
	}
}

func Test_Walk(t *testing.T) {
	t.FailNow()
}

func Test_WalkBackward(t *testing.T) {
	t.FailNow()
}
