package gitx

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"gopkg.in/src-d/go-billy.v4/helper/chroot"
	gitlib "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type Repository struct {
	repo *gitlib.Repository
}

func LocalRepository() (*Repository, error) {
	gitInstalled, err := IsGITBinaryInstalled()
	if err != nil {
		return nil, fmt.Errorf("unable to check if git is installed: %w", err)
	}
	if !gitInstalled {
		return nil, fmt.Errorf("git is not installed but is required")
	}

	repo, err := gitlib.PlainOpenWithOptions(".", &gitlib.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err == gitlib.ErrRepositoryNotExists {
			return nil, ErrRepositoryDoesNotExists
		}
		return nil, fmt.Errorf("unable to open repository: %w", err)
	}

	return &Repository{
		repo: repo,
	}, nil
}

func (git *Repository) HeadReference() (string, error) {
	head, err := git.repo.Head()
	if err != nil {
		return "", fmt.Errorf("unable to get head: %w", err)
	}
	return head.Name().Short(), nil
}

func (git *Repository) AbsoluteLocation() (string, error) {
	rawFS, ok := git.repo.Storer.(*filesystem.Storage)
	if !ok {
		return "", fmt.Errorf("repository storage is not filesystem.Storage")
	}

	fs, ok := rawFS.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		return "", fmt.Errorf("filesystem is not chroot.ChrootHelper")
	}

	repoPath, err := filepath.Abs(fs.Root())
	if err != nil {
		return "", fmt.Errorf("unable to get absolute path: %w", err)
	}

	return filepath.Dir(repoPath), nil
}

func (git *Repository) IsBranchSyncedWithRemote(branchName string) (bool, error) {
	ref, err := git.repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return false, fmt.Errorf("unable to reference branch %q: %w", branchName, err)
	}

	remotes, err := git.getRemotesForRef(ref.Name())
	if err != nil {
		return false, fmt.Errorf("unable to get remote reference for %q: %w", ref.String(), err)
	}

	if len(remotes) == 0 {
		return false, nil
	}

	refHash := ref.Hash()
	for _, remote := range remotes {
		if remote != refHash {
			return false, nil
		}
	}
	return true, nil
}

func (git *Repository) HasWIP() (bool, error) {
	if err := exec.Command("git", "diff-index", "--quiet", "HEAD").Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (git *Repository) getRemotesForRef(refName plumbing.ReferenceName) ([]plumbing.Hash, error) {
	var remoteRefs []plumbing.Hash

	remotes, err := git.repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("unable to list remotes: %w", err)
	}

	for _, remote := range remotes {
		remoteRefName := plumbing.Revision(remote.Config().Name + "/" + refName.Short())
		remoteRef, err := git.repo.ResolveRevision(remoteRefName)
		if err != nil && err != plumbing.ErrReferenceNotFound {
			return nil, fmt.Errorf("unable to get remote revision named %q: %w", remoteRefName, err)
		}
		remoteRefs = append(remoteRefs, *remoteRef)
		break
	}

	return remoteRefs, nil
}
