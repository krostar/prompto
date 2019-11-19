package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func CommandConfig(ctx context.Context) (*cobra.Command, context.Context, error) {
	return configCommandAndFlags(), ctx, nil
}

func configCommandAndFlags() *cobra.Command {
	return &cobra.Command{Use: "config"}
}
