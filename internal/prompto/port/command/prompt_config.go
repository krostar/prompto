package command

import (
	"context"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
)

// PromptConfig defines the meta-command that aggregates multiple config sub commands.
func PromptConfig() cli.Command { return &cmdConfig{} }

type cmdConfig struct{}

func (cmd cmdConfig) Description() string {
	return "PromptConfig regroups sub-commands that handle configuration."
}

func (cmd cmdConfig) Examples() []string {
	return []string{
		app.Name() + " help config",
		app.Name() + " config compile",
	}
}

func (cmd cmdConfig) Execute(_ context.Context, _, _ []string) error {
	return cli.ErrorWithExitStatus(cli.ErrorShowHelp(nil), 0)
}
