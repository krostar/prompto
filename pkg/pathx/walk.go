// Package pathx defines multiple helpers around paths.
package pathx

import (
	"fmt"
	"sort"
)

// WalkSplit walks recursively through all directories part of the
// provided path, beginning with the first directory.
// Example for /a/b/c/d, walk func will be called in this order:
// /
// a
// b
// c
// d
func WalkSplit(path string, walk func(path string) (err error)) error {
	directories := SplitPath(path)
	for _, dir := range directories {
		if err := walk(dir); err != nil {
			return fmt.Errorf("walking through %q failed: %w", path, err)
		}
	}

	return nil
}

// Walk walks recursively through all directories part of the
// provided path, beginning with the first directory.
// Example for /a/b/c/d, walk func will be called in this order:
// /
// /a
// /a/b
// /a/b/c
// /a/b/c/d
func Walk(path string, walk func(path string) (err error)) error {
	directories := reverseStringSlice(splitFullPathBackward(path))
	for _, dir := range directories {
		if err := walk(dir); err != nil {
			return fmt.Errorf("walking through %q failed: %w", path, err)
		}
	}

	return nil
}

// WalkBackward works the same as Walk, except it does it in reverse order.
// Example for /a/b/c/d, walk func will be called in this order:
// /a/b/c/d
// /a/b/c
// /a/b
// /a
// /
func WalkBackward(path string, walk func(path string) (err error)) error {
	directories := splitFullPathBackward(path)
	sort.Sort(sort.Reverse(sort.StringSlice(directories)))

	for _, dir := range directories {
		if err := walk(dir); err != nil {
			return fmt.Errorf("walking through %q failed: %w", dir, err)
		}
	}

	return nil
}
