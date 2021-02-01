package command

import (
	"context"

	"github.com/krostar/cli"
	"github.com/krostar/logger"
)

type contextContainerKey uint16

const contextContainerKeyLogger contextContainerKey = iota + 1

type contextContainer struct {
	log logger.Logger
}

func setLoggerInContextContainer(ctx context.Context, log logger.Logger) {
	cli.SetMetadata(ctx, contextContainerKeyLogger, log)
}

func getLoggerFromContextContainer(ctx context.Context) logger.Logger {
	if log, ok := cli.GetMetadata(ctx, contextContainerKeyLogger).(logger.Logger); ok {
		return log
	}
	return nil
}
