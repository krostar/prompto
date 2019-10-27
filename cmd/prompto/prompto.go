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
	debug.SetGCPercent(-1) // don't bother with garbage collection on a such short program

	ctx, cancel := clix.NewContextCancelableBySignal(os.Interrupt, os.Kill)
	defer cancel()

	err := clix.
		Command(clix.WithLogger(
			cli.CommandWritePrompt,
			clix.LoggerWithAppName(app.Name()),
			clix.LoggerWithVersion(app.Version()),
		)).
		SubCommand(clix.
			Command(cli.CommandConfig).
			SubCommand(cli.CommandConfigCompile).
			Build(),
		).
		Exec(ctx, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
