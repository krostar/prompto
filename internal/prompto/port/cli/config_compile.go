// Package cli defines the CLI port.
package cli

import (
	"context"
	"fmt"

	"github.com/krostar/clix"
	"github.com/krostar/config/defaulter"
	"github.com/krostar/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/krostar/prompto/internal/pkg/app"
)

// CommandConfigCompile binary compiles the configuration to speed-up the loading time.
// nolint: unparam
func CommandConfigCompile(ctx context.Context) (*cobra.Command, context.Context, error) {
	var cfg promptConfigFile

	cfg.SetDefault()

	cmd := compileConfigCommandAndFlags(&cfg)
	cmd.RunE = clix.ExecHandler(ctx, func(showHelp func()) (clix.Handler, error) {
		return &compileConfigCommand{
			showHelp: showHelp,
			cfg:      cfg,
			log:      clix.LoggerFromContext(ctx),
		}, nil
	})

	return cmd, ctx, nil
}

func compileConfigCommandAndFlags(cfg pflag.Value) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile configuration file to speed up deserialization time",
		Example: "\n\t" +
			app.Name() + "config compile\n\t" +
			app.Name() + "config compile --config ~/.config/prompto/config.yaml",
	}

	flags := cmd.Flags()
	flags.VarP(cfg, "config", "c", "configuration file to compile")

	return cmd
}

type compileConfigCommand struct {
	showHelp func()
	cfg      promptConfigFile
	log      logger.Logger
}

func (c *compileConfigCommand) Handle(ctx context.Context, args, dashed []string) error {
	var cfg promptConfig

	if err := defaulter.SetDefault(&cfg); err != nil {
		return fmt.Errorf("unable to set config defaults: %w", err)
	}

	if err := c.cfg.loadOriginal(&cfg); err != nil {
		return fmt.Errorf("unable to load config file: %w", err)
	}

	if err := c.cfg.generateBinary(&cfg); err != nil {
		return fmt.Errorf("unable to generate binary config: %w", err)
	}

	return nil
}
