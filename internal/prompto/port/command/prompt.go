package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/krostar/cli"
	"github.com/krostar/config"
	"github.com/krostar/logger"
	"github.com/krostar/logger/zap"

	"github.com/krostar/prompto/internal/prompto/adapter/segment"
	"github.com/krostar/prompto/internal/prompto/domain"
	"github.com/krostar/prompto/internal/prompto/domain/usecase"
	"github.com/krostar/prompto/pkg/color"
)

func Prompt() cli.Command { return &cmdPrompt{} }

type cmdPrompt struct {
	cfg cmdPromptConfig
	log logger.Logger

	loadPromptConfig usecase.LoadPromptConfigFunc
	writePrompt      usecase.PromptWriterFunc
}

type cmdPromptConfig struct {
	Logger logger.Config
	File   string

	Shell string

	LeftPromptOnly  bool
	RightPromptOnly bool

	LastCommandDurationNanoseconds uint
	LastCommandStatus              uint
}

func (cfg *cmdPromptConfig) SetDefault() {
	if shell := os.Getenv("SHELL"); shell != "" {
		cfg.Shell = filepath.Base(shell)
	}
	cfg.Logger = logger.Config{
		Verbosity: logger.LevelWarn.String(),
		Formatter: "console",
		WithColor: false,
		Output:    "stdout",
	}
	cfg.Logger.SetDefault()

	if home, err := os.UserHomeDir(); err == nil {
		cfg.File = filepath.Join(home, ".config", "prompto", "prompto.yml")
		cfg.Logger.Output = filepath.Join(home, ".config", "prompto", "prompto.log")
	}
}

func (cmd *cmdPrompt) Flags() []cli.Flag {
	return []cli.Flag{
		cli.FlagString("verbosity", "v", &cmd.cfg.Logger.Verbosity, fmt.Sprintf("verbosity of the logger (%s, %s, %s, %s)", logger.LevelDebug.String(), logger.LevelInfo.String(), logger.LevelWarn.String(), logger.LevelError.String())),
		cli.FlagString("shell", "", &cmd.cfg.Shell, "escape sequence to use for colorization, based on shell name"),
		cli.FlagString("file", "c", &cmd.cfg.File, "configuration file to use to configure the prompt"),
		cli.FlagBool("left", "", &cmd.cfg.LeftPromptOnly, "display left prompt only"),
		cli.FlagBool("right", "", &cmd.cfg.RightPromptOnly, "display right prompt only"),
		cli.FlagUint("last-cmd-duration", "d", &cmd.cfg.LastCommandDurationNanoseconds, "execution duration of the last command in nanosecond"),
		cli.FlagUint("last-cmd-status", "s", &cmd.cfg.LastCommandStatus, "exit status of the last executed command"),
	}
}

func (cmd *cmdPrompt) Hooks() *cli.Hooks {
	return &cli.Hooks{
		BeforeFlagsDefinition: func(_ context.Context) error {
			if err := config.SetDefault(&cmd.cfg); err != nil {
				return fmt.Errorf("unable to set config defaults: %v", err)
			}
			return nil
		},
		PersistentBeforeCommandExecution: func(ctx context.Context) error {
			log, flushLogs, err := zap.New(zap.WithConfig(cmd.cfg.Logger), zap.WithoutTime())
			if err != nil {
				return fmt.Errorf("unable to build logger: %v", err)
			}

			cmd.log = log

			setLoggerInContextContainer(ctx, cmd.log)
			cli.SetExitLogger(ctx, logger.WriteCloserLevel(cmd.log, flushLogs, logger.LevelError))

			return nil
		},
		BeforeCommandExecution: func(ctx context.Context) error {
			cmd.writePrompt = usecase.WritePrompts(os.Stdout)
			cmd.loadPromptConfig = usecase.LoadPromptConfig()
			return nil
		},
	}
}

func (cmd *cmdPrompt) Execute(ctx context.Context, _, _ []string) error {
	colorizer, err := color.NewColorizer(cmd.cfg.Shell)
	if err != nil {
		return fmt.Errorf("unable to create colorizer: %w", err)
	}

	promptConfig, err := cmd.loadPromptConfig(ctx, cmd.cfg.File, uint16(cmd.cfg.LastCommandStatus), time.Duration(cmd.cfg.LastCommandDurationNanoseconds))
	if err != nil {
		return fmt.Errorf("unable to load segments configuration: %v", err)
	}

	var creationRequests []usecase.PromptCreationRequest

	if !cmd.cfg.LeftPromptOnly && !cmd.cfg.RightPromptOnly {
		cmd.cfg.LeftPromptOnly = true
		cmd.cfg.RightPromptOnly = true
	}

	if cmd.cfg.LeftPromptOnly {
		var segmenters []domain.SegmentsProvider
		segmenters, err = segment.ProvideSegments(promptConfig.LeftSegments, promptConfig.Segments)
		creationRequests = append(creationRequests, usecase.PromptCreationRequest{
			Direction:        domain.DirectionLeft,
			Colorizer:        colorizer,
			SegmentsProvider: segmenters,
			SeparatorConfig:  promptConfig.Separators,
		})
	}

	if cmd.cfg.RightPromptOnly {
		var segmenters []domain.SegmentsProvider
		segmenters, err = segment.ProvideSegments(promptConfig.RightSegments, promptConfig.Segments)
		creationRequests = append(creationRequests, usecase.PromptCreationRequest{
			Direction:        domain.DirectionRight,
			Colorizer:        colorizer,
			SegmentsProvider: segmenters,
			SeparatorConfig:  promptConfig.Separators,
		})
	}

	if err != nil {
		return err
	}

	if err := cmd.writePrompt(ctx, creationRequests...); err != nil {
		return fmt.Errorf("unable to write prompt to stdout: %w", err)
	}

	return nil
}
