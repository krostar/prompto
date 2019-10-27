// Package gitx simplifies git related information retrieval.
package gitx

import (
	"fmt"
	"os/exec"
)

// IsGITBinaryInstalled checks whenever the git command is reach-able.
func IsGITBinaryInstalled() (bool, error) {
	if err := exec.Command("command", "-v", "git").Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			return false, nil
		}

		return false, fmt.Errorf("command execution failed %q: %w", "command -v git", err)
	}

	return true, nil
}
