package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/krostar/prompto/internal/pkg/app"
)

// CommandConfig defines the meta-command config that aggregates multiple config sub commands.
// nolint: unparam
func CommandConfig(ctx context.Context) (*cobra.Command, context.Context, error) {
	return configCommandAndFlags(), ctx, nil
}

func configCommandAndFlags() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Config regroup sub-commands that handle somehow configuration-related stuff",
		Example: "\n\t" + app.Name() + "config compile",
	}
}
