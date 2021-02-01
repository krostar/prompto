package segment

import (
	"fmt"
	"os/exec"
	"strings"
)

type gitRepository struct {
	location string
	branch   string
	hasWIP   bool
	isSynced bool
}

type gitRepositoryGetter interface {
	get() (*gitRepository, error)
}

type gitCommandRepository struct{}

func (git *gitCommandRepository) get() (*gitRepository, error) {
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

	var repo gitRepository

	switch status {
	case 21:
		return nil, nil
	case 0, 31, 41:
		lines := strings.Split(string(out), "\n")

		repo.location = lines[0]
		repo.branch = lines[1]

		repo.hasWIP = status == 31
		repo.isSynced = status != 41 && !repo.hasWIP
	default:
		return nil, fmt.Errorf("unhandled status %d for command %q output %q", status, cmd, out)
	}

	return &repo, nil
}
