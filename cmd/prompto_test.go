package main

// const promptoBinaryPath = "../../build/prompto"
//
// var (
// 	leftExpected = "\x1b[38;2;228;228;228;48;2;48;48;48m 14s \x1b[0m\x1b[38;2;48;48;48;48;2;135;0;0m\ue0b0\x1b[0m\x1b[38;2;228;228;228;48;2;135;0;0m 32 \x1b[0m\x1b[38;2;135;0;0;48;2;48;48;48m\ue0b0\x1b[0m\x1b[38;2;188;188;188;48;2;48;48;48m $ \x1b[0m\x1b[38;2;48;48;48m\ue0b0\x1b[0m "
// 	leftArgs     = []string{
// 		"--config", "./testdata/prompto.yml",
// 		"--left",
// 		"--last-cmd-duration", "14000000000",
// 		"--last-cmd-status", "32",
// 	}
//
// 	rightExpected = " \x1b[38;2;0;135;0m\ue0b2\x1b[0m\x1b[38;2;228;228;228;48;2;0;135;0m• master \ue0a0 \x1b[0m\x1b[38;2;48;48;48;48;2;0;135;0m\ue0b2\x1b[0m\x1b[38;2;228;228;228;48;2;48;48;48m prompto \x1b[0m\x1b[38;2;128;128;128;48;2;48;48;48m\ue0b3\x1b[0m\x1b[38;2;188;188;188;48;2;48;48;48m cmd \x1b[0m\x1b[38;2;128;128;128;48;2;48;48;48m\ue0b3\x1b[0m\x1b[38;2;188;188;188;48;2;48;48;48m prompto \x1b[0m\x1b[38;2;128;128;128;48;2;48;48;48m\ue0b3\x1b[0m\x1b[38;2;128;128;128;48;2;48;48;48m\ue0b3\x1b[0m\x1b[38;2;228;228;228;48;2;48;48;48m ~ \x1b[0m"
// 	rightArgs     = []string{
// 		"--config", "./testdata/prompto.yml",
// 		"--right",
// 	}
// )
//
// func TestMain(t *testing.T) {
// 	ci, _ := strconv.ParseBool(os.Getenv("GOTEST_CI")) // nolint: errcheck, gosec
// 	if !ci {
// 		t.SkipNow()
// 	}
//
// 	t.Run("left prompt", func(t *testing.T) {
// 		stdout, stderr, status, err := execBinary(
// 			strings.Join(append([]string{promptoBinaryPath}, leftArgs...), " "),
// 		)
// 		require.NoError(t, err)
// 		require.Zero(t, status)
// 		assert.Empty(t, stderr)
// 		assert.Equal(t, leftExpected, stdout)
// 	})
//
// 	t.Run("right prompt", func(t *testing.T) {
// 		stdout, stderr, status, err := execBinary(
// 			strings.Join(append([]string{promptoBinaryPath}, rightArgs...), " "),
// 		)
// 		require.NoError(t, err)
// 		require.Zero(t, status)
// 		assert.Empty(t, stderr)
// 		assert.Equal(t, rightExpected, stdout)
// 	})
// }
//
// func Benchmark_executePromptoCLI_left(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		err := executePromptoCLI(context.Background(), leftArgs)
// 		if err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
//
// func Benchmark_executePromptoCLI_right(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		err := executePromptoCLI(context.Background(), rightArgs)
// 		if err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
//
// func execBinary(command string) (string, string, uint8, error) {
// 	var (
// 		bufOut bytes.Buffer
// 		bufErr bytes.Buffer
// 		status uint8
// 	)
//
// 	cmd := exec.Command("bash", "-c", command) // nolint: gosec
// 	cmd.Stdout = &bufOut
// 	cmd.Stderr = &bufErr
//
// 	if err := cmd.Run(); err != nil {
// 		if exit, ok := err.(*exec.ExitError); ok {
// 			status = uint8(exit.ExitCode())
// 		} else {
// 			return "", "", 255, fmt.Errorf("unable to execute and get exit code of command %q: %w", command, err)
// 		}
// 	}
//
// 	return bufOut.String(), bufErr.String(), status, nil
// }
