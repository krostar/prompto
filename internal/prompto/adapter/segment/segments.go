// Package segment defines all segments and exposes an unique way to provide them.
package segment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/pkg/pathx"
)

const (
	// SegmentNameUnknown defines the segment name to use for failure.
	SegmentNameUnknown = ""
	// SegmentNameStub defines the segment name to use for stub.
	SegmentNameStub = "-stub-"
)

// this is easier to handle the mapping this way
// nolint: gochecknoglobals
var segmentsMapper = map[string]struct {
	create       func(interface{}) (domain.SegmentsProvider, error)
	configGetter func(Config) interface{}
}{
	SegmentNameStub: {
		create:       segmentStub,
		configGetter: func(cfg Config) interface{} { return cfg.Stub },
	},
	"cwd": {
		create:       segmentCWD,
		configGetter: func(cfg Config) interface{} { return cfg.CWD },
	},
	"last-cmd-exec-status": {
		create:       segmentLastCmdExecStatus,
		configGetter: func(cfg Config) interface{} { return cfg.LastCMDExecStatus },
	},
	"last-cmd-exec-time": {
		create:       segmentLastCmdExecTime,
		configGetter: func(cfg Config) interface{} { return cfg.LastCMDExecTime },
	},
	"read-only": {
		create:       segmentReadOnly,
		configGetter: func(cfg Config) interface{} { return cfg.ReadOnly },
	},
	"git": {
		create:       segmentGIT,
		configGetter: func(cfg Config) interface{} { return cfg.GIT },
	},
}

// Config stores the configuration for all segments provider.
type Config struct {
	Stub StubConfig `yaml:"-"`

	CWD               cwdConfig               `yaml:"cwd"`
	LastCMDExecStatus lastCmdExecStatusConfig `yaml:"last-cmd-exec-status"`
	LastCMDExecTime   lastCmdExecTimeConfig   `yaml:"last-cmd-exec-time"`
	ReadOnly          readOnlyConfig          `yaml:"read-only"`
	GIT               gitConfig               `yaml:"git"`
}

// ProvideSegments provides segments based on configuration.
func ProvideSegments(segments []string, cfg Config) ([]domain.SegmentsProvider, error) {
	var segmenters []domain.SegmentsProvider

	for _, segment := range segments {
		s, exists := segmentsMapper[segment]
		if !exists {
			return nil, fmt.Errorf("segmenter %q does not exists", segment)
		}

		segmenter, err := s.create(s.configGetter(cfg))
		if err != nil {
			return nil, fmt.Errorf("unable to create segmenter %q: %w", segment, err)
		}

		segmenters = append(segmenters, segmenter)
	}

	return segmenters, nil
}

func replaceEnvironmentInPath(path string) string {
	pathSplit := pathx.SplitPath(path)

	for i, split := range pathSplit {
		if strings.HasPrefix(split, "$") {
			if p, isset := os.LookupEnv(split[1:]); isset && p != "" {
				pathSplit[i] = p
			}
		}
	}

	return filepath.Join(pathSplit...)
}
