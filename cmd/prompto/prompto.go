package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/krostar/clix"

	"github.com/krostar/prompto/internal/pkg/app"
	"github.com/krostar/prompto/internal/prompto/port/cli"
)

func main() {
	debug.SetGCPercent(-1)

	cmd := clix.
		Command(clix.WithLogger(
			cli.CommandWritePrompt,
			clix.LoggerWithAppName(app.Name()),
			clix.LoggerWithVersion(app.Version()),
		)).
		SubCommand(clix.
			Command(cli.CommandConfig).
			SubCommand(cli.CommandConfigCompile).
			SubCommand(cli.CommandConfigValidate).Build(),
		)

	ctx, cancel := clix.NewContextCancelableBySignal(os.Interrupt, os.Kill)
	defer cancel()

	if err := cmd.Exec(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
