package pathx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// the duplication is due to the fact that both function are
// given the same test suite to be played with, but small
// amount of details changes depending of the function.
// nolint: dupl
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
		}, "other multiple elements": {
			path:     "./a/b/c",
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

func Test_WalkSplit_error(t *testing.T) {
	err := WalkSplit("/", func(p string) error { return errors.New("boum") })
	require.Error(t, err)
}

// the duplication is due to the fact that both function are
// given the same test suite to be played with, but small
// amount of details changes depending of the function.
// nolint: dupl
func Test_Walk(t *testing.T) {
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
			expected: []string{"a", "a/b", "a/b/c"},
		}, "other multiple elements": {
			path:     "./a/b/c",
			expected: []string{"a", "a/b", "a/b/c"},
		}, "cleaned path": {
			path:     "///a",
			expected: []string{"/", "/a"},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			var i int
			err := Walk(test.path, func(p string) error {
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

func Test_Walk_error(t *testing.T) {
	err := Walk("/", func(p string) error { return errors.New("boum") })
	require.Error(t, err)
}

// the duplication is due to the fact that both function are
// given the same test suite to be played with, but small
// amount of details changes depending of the function.
// nolint: dupl
func Test_WalkBackward(t *testing.T) {
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
			expected: []string{"a/b/c", "a/b", "a"},
		}, "other multiple elements": {
			path:     "./a/b/c",
			expected: []string{"a/b/c", "a/b", "a"},
		}, "cleaned path": {
			path:     "///a",
			expected: []string{"/a", "/"},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			var i int
			err := WalkBackward(test.path, func(p string) error {
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

func Test_WalkBackward_error(t *testing.T) {
	err := WalkBackward("/", func(p string) error { return errors.New("boum") })
	require.Error(t, err)
}
