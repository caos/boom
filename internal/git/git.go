package git

import (
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Git struct {
	Repo     *git.Repository
	prevTree *object.Tree
}

func New(localPath, url, secretPath string) (*Git, error) {
	g := &Git{}

	repo, err := g.cloneRepo(localPath, url, secretPath)
	if err != nil {
		return nil, err
	}
	g.Repo = repo

	ref, err := g.Repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	prevTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	g.prevTree = prevTree

	return g, nil
}

func (g *Git) cloneRepo(localPath, url, secretPath string) (*git.Repository, error) {
	publicKey, err := ssh.NewPublicKeysFromFile("git", secretPath, "")
	if err != nil {
		return nil, err
	}

	return git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     publicKey,
	})
}

func (g *Git) IsFileChanged(path string) (bool, error) {

	w, err := g.Repo.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err == git.NoErrAlreadyUpToDate {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	ref, err := g.Repo.Head()
	if err != nil {
		return false, err
	}

	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		return false, err
	}

	currentTree, err := commit.Tree()
	if err != nil {
		return false, err
	}

	changes, err := currentTree.Diff(g.prevTree)
	if err != nil {
		return false, err
	}
	g.prevTree = currentTree

	for _, c := range changes {
		if c.To.Name == path {
			return true, nil
		}
	}
	return false, nil

}
