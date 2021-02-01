// Package command defines the CLI port.
package command

import (
	"context"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
)

func PromptConfigCompile() cli.Command { return &cmdConfigCompile{} }

type cmdConfigCompile struct {
	// cfgFile cmdPromptConfigFile
}

func (cmd *cmdConfigCompile) Description() string {
	return "Compile configuration file to speed up deserialization time"
}

func (cmd *cmdConfigCompile) Examples() []string {
	return []string{
		app.Name() + "prompt config compile",
		app.Name() + "prompt config compile --config ~/.config/prompto/config.yaml",
	}
}

func (cmd *cmdConfigCompile) Hooks() *cli.Hooks {
	return &cli.Hooks{
		BeforeFlagsDefinition: func(ctx context.Context) error {
			// if err := config.SetDefault(&cmd.cfgFile); err != nil {
			// 	return fmt.Errorf("unable to set config defaults: %v", err)
			// }
			return nil
		},
	}
}

func (cmd *cmdConfigCompile) Flags() []cli.Flag {
	return []cli.Flag{
		// cli.FlagCustom("config", "c", &cmd.cfgFile, "configuration file to compile"),
	}
}

func (cmd *cmdConfigCompile) Execute(_ context.Context, _, _ []string) error {
	// var cfg loadPromptConfig
	//
	// if err := config.SetDefault(&cfg); err != nil {
	// 	return fmt.Errorf("unable to set config defaults: %w", err)
	// }
	//
	// if err := cmd.cfgFile.loadNotCompiledFile(&cfg); err != nil {
	// 	return fmt.Errorf("unable to load config file: %w", err)
	// }
	//
	// if err := cmd.cfgFile.compile(&cfg); err != nil {
	// 	return fmt.Errorf("unable to generate binary config: %w", err)
	// }
	//
	return nil
}
