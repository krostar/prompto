package usecase

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/krostar/config"
	sourcefile "github.com/krostar/config/source/file"
	"github.com/vmihailenco/msgpack"

	"github.com/krostar/prompto/internal/prompto/adapter/segment"
	"github.com/krostar/prompto/internal/prompto/domain"
)

func LoadPromptConfig() LoadPromptConfigFunc { return (&loadPromptConfig{}).loadPromptConfig }

type PromptConfig struct {
	LeftSegments  []string `yaml:"left-segments"`
	RightSegments []string `yaml:"right-segments"`

	Separators domain.SeparatorConfig `yaml:"separator"`
	Segments   segment.Config         `yaml:"segments"`
}

type LoadPromptConfigFunc func(ctx context.Context, configFile string, lastCommandStatus uint16, lastCommandExecutionDuration time.Duration) (*PromptConfig, error)

type loadPromptConfig struct{}

func (uc *loadPromptConfig) loadPromptConfig(_ context.Context, configFile string, lastCommandStatus uint16, lastCommandExecutionDuration time.Duration) (*PromptConfig, error) {
	cfg := new(PromptConfig)
	cfg.Segments.LastCMDExecStatus.StatusCode = lastCommandStatus
	cfg.Segments.LastCMDExecTime.DurationNS = uint64(lastCommandExecutionDuration.Nanoseconds())

	binaryConfigFile := strings.TrimSuffix(configFile, filepath.Ext(configFile)) + ".bin"

	if err := uc.loadCompiledFile(binaryConfigFile, cfg); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := uc.loadNotCompiledFile(configFile, cfg); err != nil {
				return nil, fmt.Errorf("unable to load regular configuration file %q: %v", configFile, err)
			}
		} else {
			return nil, fmt.Errorf("unable to load binary configuration file %q: %v", binaryConfigFile, err)
		}
	}

	return cfg, nil
}

func (uc loadPromptConfig) loadCompiledFile(compiledConfigFilename string, cfg *PromptConfig) error {
	binary, err := ioutil.ReadFile(compiledConfigFilename)
	if err != nil {
		return fmt.Errorf("unable to read binary config %q: %w", compiledConfigFilename, err)
	}

	if err := msgpack.Unmarshal(binary, cfg); err != nil {
		return fmt.Errorf("unable to unmarshal binary config %q: %w", compiledConfigFilename, err)
	}

	return nil
}

func (uc loadPromptConfig) loadNotCompiledFile(configFilename string, cfg *PromptConfig) error {
	return config.Load(&cfg, config.WithSources(sourcefile.New(configFilename, sourcefile.FailOnUnknownFields())))
}
