// Package cli defines the CLI port.
package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/krostar/clix"
	"github.com/krostar/config"
	"github.com/krostar/config/defaulter"
	sourcefile "github.com/krostar/config/source/file"
	"github.com/krostar/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/krostar/prompto/internal/pkg/app"
	"github.com/krostar/prompto/internal/prompto/adapter/segment"
	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/internal/prompto/domain/usecase"
)

// CommandWritePrompt creates and writes prompt to the standard output.
func CommandWritePrompt(ctx context.Context) (*cobra.Command, context.Context, error) {
	var cfg writePromptCommandConfig
	if err := defaulter.SetDefault(&cfg); err != nil {
		return nil, nil, errors.Wrap(err, "unable to set config defaults")
	}

	var cfgPath = "prompto.yml"
	if home, isset := os.LookupEnv("HOME"); isset {
		cfgPath = filepath.Join(home, ".config", "prompto", cfgPath)
	}

	cmd := writePromptCommandAndFlags(&cfg, &cfgPath)
	cmd.RunE = clix.ExecHandler(ctx, func(showHelp func()) (clix.Handler, error) {
		if err := config.Load(&cfg, config.WithSources(sourcefile.New(
			cfgPath,
			sourcefile.FailOnUnknownFields(),
			sourcefile.MayNotExist(),
		))); err != nil {
			return nil, fmt.Errorf("unable to load config file: %w", err)
		}

		return &writePromptCommand{
			showHelp:    showHelp,
			cfg:         cfg,
			log:         clix.LoggerFromContext(ctx),
			writePrompt: usecase.WritePrompt(),
		}, nil
	})

	return cmd, ctx, nil
}

func writePromptCommandAndFlags(cfg *writePromptCommandConfig, cfgPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Short:                 "Display prompt on standard output",
		Example:               app.Name() + "--left",
		SilenceErrors:         true,
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
	}
	flags := cmd.Flags()
	flags.StringVarP(cfgPath,
		"config", "c", *cfgPath,
		"configuration file to use to configure the prompt",
	)
	flags.BoolVar(&cfg.LeftOnly,
		"left", cfg.LeftOnly,
		"display left prompt only",
	)
	flags.BoolVar(&cfg.RightOnly,
		"right", cfg.RightOnly,
		"display right prompt only",
	)
	flags.UintVarP(&cfg.Segments.LastCMDExecTime.DurationNS,
		"last-cmd-duration", "d", cfg.Segments.LastCMDExecTime.DurationNS,
		"execution time of the last command, in nanoseconds",
	)
	flags.UintVarP(&cfg.Segments.LastCMDExecStatus.StatusCode,
		"last-cmd-status", "s", cfg.Segments.LastCMDExecStatus.StatusCode,
		"exit code of the executed last command",
	)

	return cmd
}

type writePromptCommand struct {
	showHelp func()
	cfg      writePromptCommandConfig
	log      logger.Logger

	writePrompt usecase.PromptWriterFunc
}

type writePromptCommandConfig struct {
	LeftSegments  []string `yaml:"left-segments"`
	RightSegments []string `yaml:"right-segments"`

	LeftOnly  bool `yaml:"left-only"`
	RightOnly bool `yaml:"right-only"`

	Separator domain.SeparatorConfig `yaml:"separator"`
	Segments  segment.Config         `yaml:"segments"`
}

func (c *writePromptCommand) Handle(ctx context.Context, args, dashed []string) error {
	if !color.IsSupport256Color() {
		return fmt.Errorf("256-color no supported by terminal")
	}

	var (
		direction  domain.Direction
		segmenters []domain.SegmentsProvider
		err        error
	)

	switch {
	case c.cfg.LeftOnly:
		direction = domain.DirectionLeft
		segmenters, err = segment.ProvideSegments(c.cfg.LeftSegments, c.cfg.Segments)
	case c.cfg.RightOnly:
		direction = domain.DirectionRight
		segmenters, err = segment.ProvideSegments(c.cfg.RightSegments, c.cfg.Segments)
	default:
		err = errors.New("prompt direction must be chosen")
	}

	if err != nil {
		return err
	}

	if err := c.writePrompt(ctx, usecase.PromptCreationRequest{
		Direction:        direction,
		SegmentsProvider: segmenters,
		SeparatorConfig:  c.cfg.Separator,
	}, os.Stdout); err != nil {
		return fmt.Errorf("unable to write prompt to stdout: %w", err)
	}

	return nil
}
