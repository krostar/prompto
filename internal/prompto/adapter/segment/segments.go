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
	// SegmentNameStub defines the segment name to use for stub.
	SegmentNameStub = "-stub-"

	segmentNameCustom = "custom-"
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
	"git": {
		create:       segmentGIT,
		configGetter: func(cfg Config) interface{} { return cfg.GIT },
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
	"msg": {
		create:       segmentMsg,
		configGetter: func(cfg Config) interface{} { return cfg.Msg },
	},
	"newline": {
		create:       segmentNewline,
		configGetter: func(cfg Config) interface{} { return nil },
	},
	"read-only": {
		create:       segmentReadOnly,
		configGetter: func(cfg Config) interface{} { return cfg.ReadOnly },
	},
	"user": {
		create:       segmentUser,
		configGetter: func(cfg Config) interface{} { return cfg.User },
	},
}

// Config stores the configuration for all segments provider.
type Config struct {
	Stub   StubConfig              `yaml:"-"`
	Custom map[string]customConfig `yaml:"custom"`

	CWD               cwdConfig               `yaml:"cwd"`
	GIT               gitConfig               `yaml:"git"`
	Msg               msgConfig               `yaml:"msg"`
	K8S               k8sConfig               `yaml:"k8s"`
	LastCMDExecStatus lastCMDExecStatusConfig `yaml:"last-cmd-exec-status"`
	LastCMDExecTime   lastCMDExecTimeConfig   `yaml:"last-cmd-exec-time"`
	ReadOnly          readOnlyConfig          `yaml:"read-only"`
	User              userConfig              `yaml:"user"`
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
		case strings.HasPrefix(segment, segmentNameCustom):
			segmentProviderFunc = segmentCustom
			segmentProviderCfg = cfg.Custom[segment[len(segmentNameCustom):]]
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

// nolint: unparam
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
	re, err := regexp.Compile(template)
	if err != nil {
		return false, fmt.Errorf("unable to compile template %q: %w", template, err)
	}

	match := re.FindStringSubmatch(value)
	matched := len(match) > 0

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" && i < len(match) {
			m[name] = match[i]
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
