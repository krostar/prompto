package gitx

import (
	"fmt"
	"os/exec"
	"strings"
)

// Repository helps dealing with git repositories.
type Repository struct {
	path     string
	branch   string
	hasWip   bool
	isSynced bool
}

// LocalRepository returns the first found git repository in the
// working directory tree, or an error if not found.
func LocalRepository() (*Repository, error) {
	var (
		status uint8
		cmd    = `
			git rev-parse --show-toplevel --abbrev-ref --symbolic-full-name @ || exit 21 && (
				test -z "$(git ls-files . --exclude-standard --others --)" &&
				git diff --quiet -- &&
				git diff --quiet --staged --
			) || exit 31 &&
			git diff --quiet @ @{upstream} -- || exit 41
		`
	)

	out, err := exec.Command("bash", "-c", cmd).CombinedOutput() // nolint: errcheck, gosec
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			status = uint8(exit.ExitCode())
		} else {
			return nil, fmt.Errorf("unable to execute and get exit code of git command %q: %w", cmd, err)
		}
	}

	var repo Repository

	switch status {
	case 21:
		return nil, ErrRepositoryDoesNotExists
	case 0, 31, 41:
		lines := strings.Split(string(out), "\n")

		repo.path = lines[0]
		repo.branch = lines[1]

		repo.hasWip = status == 31
		repo.isSynced = status != 41 && !repo.hasWip
	default:
		return nil, fmt.Errorf("unhandled status %d for command %q output %q", status, cmd, out)
	}

	return &repo, nil
}

// HeadReference returns the reference of HEAD.
// nolint: unparam
func (git *Repository) HeadReference() (string, error) { return git.branch, nil }

// AbsoluteLocation returns the absolute location of the repository.
// nolint: unparam
func (git *Repository) AbsoluteLocation() (string, error) { return git.path, nil }

// IsHeadSyncedWithRemote returns whenever the repository is synced with remote.
// nolint: unparam
func (git *Repository) IsHeadSyncedWithRemote() (bool, error) { return git.isSynced, nil }

// HasWIP returns whenever the repository has some work in progress.
// nolint: unparam
func (git *Repository) HasWIP() (bool, error) { return git.hasWip, nil }
