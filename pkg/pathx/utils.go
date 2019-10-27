package pathx

import (
	"path/filepath"
	"strings"
)

func reverseStringSlice(array []string) []string {
	last := len(array) - 1

	for i := 0; i < len(array)/2; i++ {
		array[i], array[last-i] = array[last-i], array[i]
	}

	return array
}

func splitFullPathBackward(path string) []string {
	directories := []string{}

	for dir := filepath.Clean(path); ; dir = filepath.Dir(dir) {
		if dir == "." {
			break
		}

		directories = append(directories, dir)

		if dir == string(filepath.Separator) {
			break
		}
	}

	return directories
}

// SplitPath splits a path on the separator.
func SplitPath(path string) []string {
	path = filepath.Clean(path)

	switch path {
	case ".":
		return []string{}
	case string(filepath.Separator):
		return []string{string(filepath.Separator)}
	default:
		directories := strings.Split(path, string(filepath.Separator))
		for i, dir := range directories {
			if dir == "" {
				directories[i] = string(filepath.Separator)
			}
		}

		return directories
	}
}
