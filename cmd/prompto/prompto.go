package main

import (
	"fmt"
	"os"

	"github.com/krostar/clix"

	"github.com/krostar/prompto/internal/pkg/app"
	"github.com/krostar/prompto/internal/prompto/port/cli"
)

func main() {
	cmd := clix.Command(clix.WithLogger(
		cli.CommandWritePrompt,
		clix.LoggerWithAppName(app.Name()),
		clix.LoggerWithVersion(app.Version()),
	))

	ctx, cancel := clix.NewContextCancelableBySignal(os.Interrupt, os.Kill)
	defer cancel()

	if err := cmd.Exec(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
