package gitx

import (
	"fmt"
	"os/exec"
	"strings"
)

type Repository struct {
	path   string
	branch string
	remote string
	hasWip bool
}

func LocalRepository() (*Repository, error) {
	out, err := exec.Command("bash", "-c", `
git rev-parse --show-toplevel ;
git status --porcelain=v2 --ignore-submodules --branch -z;
`).CombinedOutput()
	if err != nil {
		return nil, ErrRepositoryDoesNotExists
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) != 2 {
		return nil, fmt.Errorf("unable to parse git command output %q: %w", out, err)
	}

	var repo Repository

	repo.path = lines[0]
	status := strings.Split(lines[1], "\000")

	for _, data := range status[:len(status)-1] {
		if data[0] != '#' {
			repo.hasWip = true
			break
		}

		if strings.HasPrefix(data, "# branch.head ") {
			repo.branch = data[14:]
			continue
		}

		if strings.HasPrefix(data, "# branch.upstream ") {
			repo.remote = data[18:]
			continue
		}
	}

	return &repo, nil
}

func (git *Repository) HeadReference() (string, error) {
	return git.branch, nil
}

func (git *Repository) AbsoluteLocation() (string, error) {
	return git.path, nil
}

func (git *Repository) IsHeadSyncedWithRemote() (bool, error) {
	diff := fmt.Sprintf("git diff --quiet %s..%s --", git.branch, git.remote)

	if err := exec.Command("bash", "-c", diff).Run(); err != nil { // nolint: gosec
		return false, nil
	}

	return !git.hasWip, nil
}

func (git *Repository) HasWIP() (bool, error) {
	return git.hasWip, nil
}
