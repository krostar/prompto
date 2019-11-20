// Package segment defines all segments and exposes an unique way to provide them.
package segment

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/pkg/pathx"
)

const (
	// SegmentNameUnknown defines the segment name to use for failure.
	SegmentNameUnknown = ""
	SegmentNameCustom  = "custom-"
	SegmentNameNewline = "newline"
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
	SegmentNameNewline: {
		create:       segmentNewline,
		configGetter: func(cfg Config) interface{} { return nil },
	},
	"cwd": {
		create:       segmentCWD,
		configGetter: func(cfg Config) interface{} { return cfg.CWD },
	},
	"git": {
		create:       segmentGIT,
		configGetter: func(cfg Config) interface{} { return cfg.GIT },
	},
	"hi": {
		create:       segmentHi,
		configGetter: func(cfg Config) interface{} { return cfg.Hi },
	},
	"k8s": {
		create:       segmentK8S,
		configGetter: func(cfg Config) interface{} { return cfg.K8S },
	},
	"last-cmd-exec-status": {
		create:       segmentLastCMDExecStatus,
		configGetter: func(cfg Config) interface{} { return cfg.LastCMDExecStatus },
	},
	"last-cmd-exec-time": {
		create:       segmentLastCMDExecTime,
		configGetter: func(cfg Config) interface{} { return cfg.LastCMDExecTime },
	},
	"read-only": {
		create:       segmentReadOnly,
		configGetter: func(cfg Config) interface{} { return cfg.ReadOnly },
	},
}

// Config stores the configuration for all segments provider.
type Config struct {
	Stub   StubConfig              `yaml:"-"`
	Custom map[string]customConfig `yaml:"custom"`

	CWD               cwdConfig               `yaml:"cwd"`
	GIT               gitConfig               `yaml:"git"`
	Hi                hiConfig                `yaml:"hi"`
	K8S               k8sConfig               `yaml:"k8s"`
	LastCMDExecStatus lastCMDExecStatusConfig `yaml:"last-cmd-exec-status"`
	LastCMDExecTime   lastCMDExecTimeConfig   `yaml:"last-cmd-exec-time"`
	ReadOnly          readOnlyConfig          `yaml:"read-only"`
}

// ProvideSegments provides segments based on configuration.
func ProvideSegments(segments []string, cfg Config) ([]domain.SegmentsProvider, error) {
	var segmenters []domain.SegmentsProvider

	for _, segment := range segments {
		var (
			segmentProviderFunc func(interface{}) (domain.SegmentsProvider, error)
			segmentProviderCfg  interface{}
		)

		switch {
		case strings.HasPrefix(segment, SegmentNameCustom):
			segmentProviderFunc = segmentCustom
			segmentProviderCfg = cfg.Custom[segment[len(SegmentNameCustom):]]
		default:
			s, exists := segmentsMapper[segment]
			if !exists {
				return nil, fmt.Errorf("segmenter %q does not exists", segment)
			}

			segmentProviderFunc = s.create
			segmentProviderCfg = s.configGetter(cfg)
		}

		segmenter, err := segmentProviderFunc(segmentProviderCfg)
		if err != nil {
			return nil, fmt.Errorf("unable to create segmenter %q: %w", segment, err)
		}

		segmenters = append(segmenters, segmenter)
	}

	return segmenters, nil
}

func execCommand(command string) (string, string, uint8, error) {
	var (
		bufOut bytes.Buffer
		bufErr bytes.Buffer
		status uint8
	)

	cmd := exec.Command("bash", "-c", command) // nolint: gosec
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr

	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			status = uint8(exit.ExitCode())
		} else {
			return "", "", 0, fmt.Errorf("unable to execute and get exit code of command %q: %w", command, err)
		}
	}

	return bufOut.String(), bufErr.String(), status, nil
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

func createSubstituteTemplate(template, value string, m map[string]string) (bool, error) {
	var matched bool

	re, err := regexp.Compile(template)
	if err != nil {
		return matched, fmt.Errorf("unable to compile template %q: %w", template, err)
	}

	match := re.FindStringSubmatch(value)

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" && i < len(match) {
			m[name] = match[i]
			matched = true
		}
	}

	return matched, nil
}

func substituteWithTemplate(value *string, m map[string]string) {
	for key, replaceWith := range m {
		*value = strings.Replace(*value, "{"+key+"}", replaceWith, -1)
	}
}

func splitContentInSegments(content string, forEach func(content string) *domain.Segment) domain.Segments {
	var segments domain.Segments

	for _, path := range pathx.SplitPath(content) {
		segment := forEach(path)
		if segment != nil {
			segments = append(segments, segment)
		}
	}

	return segments
}
