package git

import (
	"os"

	"github.com/caos/utils/logging"
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
		logging.Log("GIT-5TP1NETBBdY2M4B").OnError(err).Debugf("Failed to clone repo %s", url)
		return nil, err
	}
	g.Repo = repo

	ref, err := g.Repo.Head()
	if err != nil {
		logging.Log("GIT-mMj1dZIWSoG2nZx").OnError(err).Debugf("Failed to get head of repo %s", url)
		return nil, err
	}

	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		logging.Log("GIT-juNgPH9agv09jNr").OnError(err).Debugf("Failed to get last commit of repo %s", url)
		return nil, err
	}
	prevTree, err := commit.Tree()
	if err != nil {
		logging.Log("GIT-wYeDNaCqmhn8x64").OnError(err).Debugf("Failed to get tree of last commit of repo %s", url)
		return nil, err
	}
	g.prevTree = prevTree

	return g, nil
}

func (g *Git) cloneRepo(localPath, url, secretPath string) (*git.Repository, error) {
	publicKey, err := ssh.NewPublicKeysFromFile("git", secretPath, "")
	if err != nil {
		logging.Log("GIT-ZImVXjm9lnrJwSu").OnError(err).Debugf("Failed to parse secret for repo %s", url)
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
		logging.Log("GIT-2PPaIdlguhB16n0").OnError(err).Debugf("Failed to get workingtree of repo in path %s", path)
		return false, err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err == git.NoErrAlreadyUpToDate {
		return false, nil
	}
	if err != nil {
		logging.Log("GIT-jJmpRiLFTHaDOXv").OnError(err).Debugf("Failed to pull of repo in path %s", path)
		return false, err
	}

	ref, err := g.Repo.Head()
	if err != nil {
		logging.Log("GIT-G2dOR3UlXsJ6wMU").OnError(err).Debugf("Failed to get the head of repo in path %s", path)
		return false, err
	}

	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		logging.Log("GIT-nP3h5VWLRAkEzVb").OnError(err).Debugf("Failed to get last commit of repo in path %s", path)
		return false, err
	}

	currentTree, err := commit.Tree()
	if err != nil {
		logging.Log("GIT-hRVoI2fcJ2EyJe3").OnError(err).Debugf("Failed to get tree of last commit of repo in path %s", path)
		return false, err
	}

	changes, err := currentTree.Diff(g.prevTree)
	if err != nil {
		logging.Log("GIT-Z90BMxv7QUhugWV").OnError(err).Debugf("Failed to diff changes of repo in path %s", path)
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
