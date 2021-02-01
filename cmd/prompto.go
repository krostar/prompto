package main

import (
	"os"
	"runtime/debug"
	"syscall"

	"github.com/krostar/cli"
	"github.com/krostar/cli/app"
	"github.com/krostar/cli/mapper/spf13/cobra"

	"github.com/krostar/prompto/internal/prompto/port/command"
)

func main() {
	// don't bother with garbage collection on a such short time running program
	debug.SetGCPercent(-1)

	ctx, cancel := cli.NewContextCancelableBySignal(syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cli.Exit(ctx, cobra.Execute(ctx, cli.NewCommand(app.Name(), command.Prompt()), os.Args))
}
