package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/krostar/config"
	sourcefile "github.com/krostar/config/source/file"
	"github.com/vmihailenco/msgpack"

	"github.com/krostar/prompto/internal/prompto/adapter/segment"
	"github.com/krostar/prompto/internal/prompto/domain"
)

type promptConfig struct {
	LeftSegments  []string               `yaml:"left-segments"`
	RightSegments []string               `yaml:"right-segments"`
	Separator     domain.SeparatorConfig `yaml:"separator"`
	Segments      segment.Config         `yaml:"segments"`
}

type promptConfigFile string

// implements flag.Value
func (cf promptConfigFile) String() string { return string(cf) }
func (cf promptConfigFile) Type() string   { return "string" }
func (cf *promptConfigFile) Set(s string) error {
	*cf = promptConfigFile(s)
	return nil
}

func (cf *promptConfigFile) SetDefault() {
	dir := "."

	if home, err := os.UserHomeDir(); err == nil {
		dir = filepath.Join(home, ".config", "prompto")
	}

	*cf = promptConfigFile(filepath.Join(dir, "prompto.yml"))
}

func (cf promptConfigFile) binaryName() string {
	return strings.TrimSuffix(cf.String(), filepath.Ext(cf.String())) + ".bin"
}

func (cf promptConfigFile) canUseBinaryInstead() bool {
	if ext := filepath.Ext(cf.String()); ext != ".bin" {
		if fi, err := os.Stat(cf.binaryName()); err == nil {
			if !fi.IsDir() && fi.Mode().Perm()&(1<<2) == 0 {
				return true
			}
		}
	}

	return false
}

func (cf promptConfigFile) loadBinary(cfg *promptConfig) error {
	binFile, err := os.OpenFile(cf.binaryName(), os.O_RDONLY, 0400)
	if err != nil {
		return fmt.Errorf("unable to open binary config %q: %w", cf.String(), err)
	}
	defer binFile.Close() // nolint: errcheck, gosec

	if err := msgpack.NewDecoder(binFile).Decode(cfg); err != nil {
		return fmt.Errorf("unable to unmarshal binary config %q: %w", cf.String(), err)
	}

	return nil
}

func (cf promptConfigFile) generateBinary(cfg *promptConfig) error {
	binaryName := cf.binaryName()
	binFile, err := os.OpenFile(binaryName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to create binary file %q: %w", binaryName, err)
	}
	defer binFile.Close() // nolint: errcheck, gosec

	if err := msgpack.NewEncoder(binFile).UseJSONTag(true).Encode(&cfg); err != nil {
		return fmt.Errorf("unable to marshal config file: %w", err)
	}

	return nil
}

func (cf promptConfigFile) loadOriginal(cfg *promptConfig) error {
	return config.Load(&cfg, config.WithSources(sourcefile.New(
		cf.String(),
		sourcefile.FailOnUnknownFields(),
		sourcefile.MayNotExist(),
	)))
}

func (cf *promptConfigFile) load(cfg *promptConfig) error {
	var (
		err        error
		configFile = cf.String()
	)

	if cf.canUseBinaryInstead() {
		*cf = promptConfigFile(strings.TrimSuffix(configFile, filepath.Ext(configFile)) + ".bin")
		err = cf.loadBinary(cfg)
	} else {
		err = cf.loadOriginal(cfg)
	}

	return err
}
