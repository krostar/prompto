package cli

// func Test_CommandWritePrompt(t *testing.T) {
// 	// test d'integration
// 	// t.FailNow()
// }

// func Test_writePromptCommandAndFlags(t *testing.T) {
// 	t.Run("command should have description/examples", func(t *testing.T) {
// 		var cfgPath promptConfigFile
// 		cmd := writePromptCommandAndFlags(&promptConfig{}, &cfgPath)
// 		assert.True(t, cmd.HasExample())
// 		assert.NotEmpty(t, cmd.Short)
// 	})
// 	t.Run("flags should be set and take cfg in consideration", func(t *testing.T) {
// 		var (
// 			segments segment.Config
// 			cfgPath  = promptConfigFile("./config.yml")
// 		)

// 		segments.LastCMDExecTime.DurationNS = 12345678
// 		segments.LastCMDExecStatus.StatusCode = 117

// 		cmd := writePromptCommandAndFlags(&promptConfig{
// 			LeftSegments:  []string{},
// 			RightSegments: []string{},
// 			LeftOnly:      true,
// 			RightOnly:     true,
// 			Segments:      segments,
// 		}, &cfgPath)

// 		for name, value := range map[string]string{
// 			"config":            "./config.yml",
// 			"left":              "true",
// 			"right":             "true",
// 			"last-cmd-duration": "12345678",
// 			"last-cmd-status":   "117",
// 		} {
// 			flag := cmd.Flag(name)
// 			require.NotNil(t, flag, name)
// 			assert.Equal(t, value, flag.Value.String(),
// 				fmt.Sprintf("%q flag should exists and equals %q", name, value),
// 			)
// 		}
// 	})
// }

// func TestWritePromptCommand_Handle(t *testing.T) {
// 	separatorConfig := domain.SeparatorConfig{
// 		Content: domain.SeparatorContentConfig{
// 			Left:      "c-left",
// 			LeftThin:  "c-left-thin",
// 			Right:     "c-right",
// 			RightThin: "c-right-thin",
// 		},
// 	}
// 	cmd := writePromptCommand{
// 		showHelp: func() {},
// 		log:      &logger.Noop{},
// 		cfg: promptConfig{
// 			Separator: separatorConfig,
// 		},
// 		writePrompt: func(context.Context, usecase.PromptCreationRequest) error { return nil },
// 	}
// 	term := os.Getenv("TERM")
// 	restoreTerm := func() { require.NoError(t, os.Setenv("TERM", term)) }
// 	setTermColorable := func() { require.NoError(t, os.Setenv("TERM", "xterm-256color")) }
// 	setTermNotColorable := func() { require.NoError(t, os.Setenv("TERM", "nope")) }
// 	usecaseStubCalled := func(called *bool) func(context.Context, usecase.PromptCreationRequest) error {
// 		return func(context.Context, usecase.PromptCreationRequest) error {
// 			*called = true
// 			return nil
// 		}
// 	}

// 	t.Run("color is not supported by terminal", func(t *testing.T) {
// 		setTermNotColorable()
// 		defer restoreTerm()

// 		called := false
// 		cmd := cmd

// 		cmd.writePrompt = usecaseStubCalled(&called)

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.Error(t, err)
// 		assert.False(t, called)
// 	})

// 	t.Run("left and right are chose", func(t *testing.T) {
// 		setTermColorable()
// 		defer restoreTerm()

// 		called := false
// 		cmd := cmd

// 		cmd.cfg.LeftOnly = true
// 		cmd.cfg.RightOnly = true
// 		cmd.writePrompt = usecaseStubCalled(&called)

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.Error(t, err)
// 		assert.False(t, called)
// 	})

// 	t.Run("neither left or right are chose", func(t *testing.T) {
// 		setTermColorable()
// 		defer restoreTerm()

// 		called := false
// 		cmd := cmd

// 		cmd.writePrompt = usecaseStubCalled(&called)

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.Error(t, err)
// 		assert.False(t, called)
// 	})

// 	t.Run("left or right with error", func(t *testing.T) {
// 		setTermColorable()
// 		defer restoreTerm()

// 		called := false
// 		cmd := cmd

// 		cmd.cfg.LeftOnly = true
// 		cmd.cfg.LeftSegments = []string{segment.SegmentNameUnknown}
// 		cmd.writePrompt = usecaseStubCalled(&called)

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.Error(t, err)
// 		assert.False(t, called)
// 	})

// 	// test is way too similar, but two different tests make sens
// 	// nolint: dupl
// 	t.Run("left prompt with valid config", func(t *testing.T) {
// 		setTermColorable()
// 		defer restoreTerm()

// 		cmd := cmd

// 		cmd.cfg.LeftOnly = true
// 		cmd.cfg.LeftSegments = []string{segment.SegmentNameStub}
// 		cmd.cfg.Segments.Stub = segment.StubConfig{
// 			Segments: domain.Segments{domain.NewSegment("hello")},
// 		}
// 		cmd.writePrompt = func(_ context.Context, pcr usecase.PromptCreationRequest) error {
// 			assert.Equal(t, domain.DirectionLeft, pcr.Direction)
// 			assert.Equal(t, separatorConfig, pcr.SeparatorConfig)
// 			assert.Len(t, pcr.SegmentsProvider, 1)
// 			return nil
// 		}

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.NoError(t, err)
// 	})

// 	// test is way too similar, but two different tests make sens
// 	// nolint: dupl
// 	t.Run("right prompt with valid config", func(t *testing.T) {
// 		setTermColorable()
// 		defer restoreTerm()

// 		cmd := cmd

// 		cmd.cfg.RightOnly = true
// 		cmd.cfg.RightSegments = []string{segment.SegmentNameStub}
// 		cmd.cfg.Segments.Stub = segment.StubConfig{
// 			Segments: domain.Segments{domain.NewSegment("hello")},
// 		}
// 		cmd.writePrompt = func(_ context.Context, pcr usecase.PromptCreationRequest) error {
// 			assert.Equal(t, domain.DirectionRight, pcr.Direction)
// 			assert.Equal(t, separatorConfig, pcr.SeparatorConfig)
// 			assert.Len(t, pcr.SegmentsProvider, 1)
// 			return nil
// 		}

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.NoError(t, err)
// 	})

// 	t.Run("usecase failed", func(t *testing.T) {
// 		setTermColorable()
// 		defer restoreTerm()

// 		cmd := cmd

// 		cmd.cfg.LeftOnly = true
// 		cmd.writePrompt = func(_ context.Context, pcr usecase.PromptCreationRequest) error {
// 			return errors.New("boum")
// 		}

// 		err := cmd.Handle(context.Background(), nil, nil)
// 		require.Error(t, err)
// 	})
// }
