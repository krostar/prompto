// Package cli defines the CLI port.
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/krostar/clix"
	"github.com/krostar/config/defaulter"
	"github.com/krostar/logger"
	"github.com/krostar/prompto/pkg/color"
	"github.com/spf13/cobra"

	"github.com/krostar/prompto/internal/pkg/app"
	"github.com/krostar/prompto/internal/prompto/adapter/segment"
	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/internal/prompto/domain/usecase"
)

// CommandWritePrompt creates and writes prompt to the standard output.
func CommandWritePrompt(ctx context.Context) (*cobra.Command, context.Context, error) {
	var cfg writePromptConfig
	if err := defaulter.SetDefault(&cfg); err != nil {
		return nil, nil, fmt.Errorf("unable to set config defaults: %w", err)
	}

	var cfgPath promptConfigFile

	cfgPath.SetDefault()

	cmd := writePromptCommandAndFlags(&cfg, &cfgPath)
	cmd.RunE = clix.ExecHandler(ctx, func(showHelp func()) (clix.Handler, error) {
		if err := cfgPath.load(&cfg.Prompt); err != nil {
			return nil, fmt.Errorf("unable to load config: %w", err)
		}

		return &writePromptCommand{
			showHelp:    showHelp,
			cfg:         cfg,
			log:         clix.LoggerFromContext(ctx),
			writePrompt: usecase.WritePrompts(os.Stdout),
		}, nil
	})

	return cmd, ctx, nil
}

func writePromptCommandAndFlags(cfg *writePromptConfig, cfgPath *promptConfigFile) *cobra.Command {
	cmd := &cobra.Command{
		Short:                 "Display prompt on standard output",
		Example:               app.Name() + "--left",
		SilenceErrors:         true,
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()
	flags.VarP(cfgPath,
		"config", "c",
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
	flags.UintVarP(&cfg.Prompt.Segments.LastCMDExecTime.DurationNS,
		"last-cmd-duration", "d", cfg.Prompt.Segments.LastCMDExecTime.DurationNS,
		"execution time of the last command, in nanoseconds",
	)
	flags.StringVar(&cfg.Shell,
		"shell", cfg.Shell,
		"escape sequence to use for colorization, based on shell name",
	)
	flags.UintVarP(&cfg.Prompt.Segments.LastCMDExecStatus.StatusCode,
		"last-cmd-status", "s", cfg.Prompt.Segments.LastCMDExecStatus.StatusCode,
		"exit code of the executed last command",
	)

	return cmd
}

type writePromptCommand struct {
	showHelp func()
	cfg      writePromptConfig
	log      logger.Logger

	writePrompt usecase.PromptWriterFunc
}

type writePromptConfig struct {
	LeftOnly  bool
	RightOnly bool
	Shell     string

	Prompt promptConfig
}

func (c writePromptConfig) SetDefault() {
	c.Shell = os.Getenv("SHELL")
}

func (c *writePromptCommand) Handle(ctx context.Context, args, dashed []string) error {
	var (
		creationRequests []usecase.PromptCreationRequest
		err              error
	)

	colorizer, err := color.NewColorizer(c.cfg.Shell)
	if err != nil {
		return fmt.Errorf("unable to create colorizer: %w", err)
	}

	if !c.cfg.LeftOnly && !c.cfg.RightOnly {
		c.cfg.LeftOnly = true
		c.cfg.RightOnly = true
	}

	if c.cfg.LeftOnly {
		var segmenters []domain.SegmentsProvider
		segmenters, err = segment.ProvideSegments(c.cfg.Prompt.LeftSegments, c.cfg.Prompt.Segments)
		creationRequests = append(creationRequests, usecase.PromptCreationRequest{
			Direction:        domain.DirectionLeft,
			Colorizer:        colorizer,
			SegmentsProvider: segmenters,
			SeparatorConfig:  c.cfg.Prompt.Separator,
		})
	}

	if c.cfg.RightOnly {
		var segmenters []domain.SegmentsProvider
		segmenters, err = segment.ProvideSegments(c.cfg.Prompt.RightSegments, c.cfg.Prompt.Segments)
		creationRequests = append(creationRequests, usecase.PromptCreationRequest{
			Direction:        domain.DirectionRight,
			Colorizer:        colorizer,
			SegmentsProvider: segmenters,
			SeparatorConfig:  c.cfg.Prompt.Separator,
		})
	}

	if err != nil {
		return err
	}

	if err := c.writePrompt(ctx, creationRequests...); err != nil {
		return fmt.Errorf("unable to write prompt to stdout: %w", err)
	}

	return nil
}
