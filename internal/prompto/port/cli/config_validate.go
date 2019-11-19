// Package cli defines the CLI port.
package cli

import (
	"context"
	"fmt"

	"github.com/krostar/clix"
	"github.com/krostar/logger"
	"github.com/spf13/cobra"
)

func CommandConfigValidate(ctx context.Context) (*cobra.Command, context.Context, error) {
	var cfg promptConfigFile

	cfg.SetDefault()

	cmd := validateConfigCommandAndFlags(&cfg)
	cmd.RunE = clix.ExecHandler(ctx, func(showHelp func()) (clix.Handler, error) {
		return &validateConfigCommand{
			showHelp: showHelp,
			cfg:      cfg,
			log:      clix.LoggerFromContext(ctx),
		}, nil
	})

	return cmd, ctx, nil
}

func validateConfigCommandAndFlags(cfg *promptConfigFile) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "",
		Example: "",
	}
	flags := cmd.Flags()
	flags.VarP(cfg,
		"config", "c",
		"configuration file to use to configure the prompt",
	)

	return cmd
}

type validateConfigCommand struct {
	showHelp func()
	cfg      promptConfigFile
	log      logger.Logger
}

func (c *validateConfigCommand) Handle(ctx context.Context, args, dashed []string) error {
	var cfg promptConfig

	cfg.Segments.LastCMDExecStatus.StatusCode = 42

	fmt.Println(cfg.Segments.LastCMDExecStatus.StatusCode)
	fmt.Println(cfg)

	if err := c.cfg.loadBinary(&cfg); err != nil {
		panic(err)
	}

	fmt.Println(cfg.Segments.LastCMDExecStatus.StatusCode)
	fmt.Println(cfg)

	return nil
}
