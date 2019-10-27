package pathx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_reverseSlice(t *testing.T) {
	tests := []struct {
		slice           []string
		expectedReverse []string
	}{
		{
			slice:           nil,
			expectedReverse: nil,
		}, {
			slice:           []string{},
			expectedReverse: []string{},
		}, {
			slice:           []string{"1", "2", "3"},
			expectedReverse: []string{"3", "2", "1"},
		}, {
			slice:           []string{"1"},
			expectedReverse: []string{"1"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(strings.Join(test.slice, ","), func(t *testing.T) {
			reversed := reverseStringSlice(test.slice)
			assert.Equal(t, test.expectedReverse, reversed)
		})
	}
}

// the duplication is due to the fact that both function are
// given the same test suite to be played with, but small
// amount of details changes depending of the function.
// nolint: dupl
func Test_splitFullPathBackward(t *testing.T) {
	tests := []struct {
		path          string
		expectedSplit []string
	}{
		{
			path:          ".",
			expectedSplit: []string{},
		}, {
			path:          ".././..",
			expectedSplit: []string{"../..", ".."},
		}, {
			path:          "/./a/b/c",
			expectedSplit: []string{"/a/b/c", "/a/b", "/a", "/"},
		}, {
			path:          "./a/b/c",
			expectedSplit: []string{"a/b/c", "a/b", "a"},
		}, {
			path:          "../a/b/c",
			expectedSplit: []string{"../a/b/c", "../a/b", "../a", ".."},
		}, {
			path:          "/a/b/c/",
			expectedSplit: []string{"/a/b/c", "/a/b", "/a", "/"},
		}, {
			path:          "a/a/a/",
			expectedSplit: []string{"a/a/a", "a/a", "a"},
		}, {
			path:          "a-a-a/",
			expectedSplit: []string{"a-a-a"},
		}, {
			path:          "/////",
			expectedSplit: []string{"/"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.path, func(t *testing.T) {
			split := splitFullPathBackward(test.path)
			assert.Equal(t, test.expectedSplit, split)
		})
	}
}

// the duplication is due to the fact that both function are
// given the same test suite to be played with, but small
// amount of details changes depending of the function.
// nolint: dupl
func Test_splitPath(t *testing.T) {
	tests := []struct {
		path          string
		expectedSplit []string
	}{
		{
			path:          ".",
			expectedSplit: []string{},
		}, {
			path:          ".././..",
			expectedSplit: []string{"..", ".."},
		}, {
			path:          "/./a/b/c",
			expectedSplit: []string{"/", "a", "b", "c"},
		}, {
			path:          "./a/b/c",
			expectedSplit: []string{"a", "b", "c"},
		}, {
			path:          "../a/b/c",
			expectedSplit: []string{"..", "a", "b", "c"},
		}, {
			path:          "/a/b/c/",
			expectedSplit: []string{"/", "a", "b", "c"},
		}, {
			path:          "a/a/a/",
			expectedSplit: []string{"a", "a", "a"},
		}, {
			path:          "a-a-a/",
			expectedSplit: []string{"a-a-a"},
		}, {
			path:          "/////",
			expectedSplit: []string{"/"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.path, func(t *testing.T) {
			split := SplitPath(test.path)
			assert.Equal(t, test.expectedSplit, split)
		})
	}
}
