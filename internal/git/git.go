package git

import (
	"io/ioutil"
	"os"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	git "gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Git struct {
	Repo     *git.Repository
	prevTree *object.Tree
	logger   logging.Logger
}

func New(logger logging.Logger, localPath, url, secretPath string) (*Git, error) {
	g := &Git{logger: logger}

	repoLogger := g.logger.WithFields(map[string]interface{}{
		"repo": url,
	})

	repo, err := g.cloneRepo(localPath, url, secretPath)
	if err != nil {
		repoLogger.WithFields(map[string]interface{}{
			"logID": "GIT-5TP1NETBBdY2M4B",
			"err":   err.Error(),
		}).Debug("Cloning failed")
		return nil, err
	}
	g.Repo = repo

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-4Sia0VjJ79gb7cw",
	}).Info("Cloned...")
	ref, err := g.Repo.Head()
	if err != nil {
		repoLogger.WithFields(map[string]interface{}{
			"logID": "GIT-mMj1dZIWSoG2nZx",
			"err":   err.Error(),
		}).Debug("Failed to get head")
		return nil, err
	}

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-4Sia0VjJ79gb7cw",
	}).Info("Get last commit...")
	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		repoLogger.WithFields(map[string]interface{}{
			"logID": "GIT-juNgPH9agv09jNr",
			"err":   err.Error(),
		}).Debug("Failed to get last commit")
		return nil, err
	}
	prevTree, err := commit.Tree()
	if err != nil {
		repoLogger.WithFields(map[string]interface{}{
			"logID": "GIT-wYeDNaCqmhn8x64",
			"err":   err.Error(),
		}).Debug("Failed to get tree of last commit")
		return nil, err
	}
	g.prevTree = prevTree
	repoLogger.WithFields(map[string]interface{}{
		"logID": "GIT-pQnw5FfIqAk0eWc",
		"path":  localPath,
	}).Info("Cloned new GitCRD")

	return g, nil
}

func (g *Git) cloneRepo(localPath, url, secretPath string) (*git.Repository, error) {

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-vNU9maj2Rfo5rRU",
		"path":  secretPath,
	}).Info("Using secret")
	sshKey, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	if err != nil {
		return nil, err
	}
	auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	// auth, err := ssh.NewPublicKeysFromFile("git", secretPath, "")
	// if err != nil {
	// 	logging.Log("GIT-ZImVXjm9lnrJwSu").OnError(err).Debugf("Failed to parse secret for repo %s", url)
	// 	return nil, err
	// }

	g.logger.WithFields(map[string]interface{}{
		"logID": "GIT-vNU9maj2Rfo5rRU",
		"repo":  url,
		"to":    localPath,
	}).Info("Cloning plain")
	return git.PlainClone(localPath, false, &git.CloneOptions{
		URL:          url,
		SingleBranch: true,
		Depth:        1,
		Progress:     os.Stdout,
		Auth:         auth,
	})
}

func (g *Git) IsFileChanged(path string) (changed bool, err error) {

	var action string
	defer func() {
		if err != nil {
			g.logger.WithFields(map[string]interface{}{
				"logID": "GIT-2PPaIdlguhB16n0",
				"path":  path,
			}).Error(errors.Wrapf(err, "Failed to %s of repo", action))
		}
	}()

	w, err := g.Repo.Worktree()
	if err != nil {
		action = "get workingtree"
		return false, err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err == git.NoErrAlreadyUpToDate {
		return false, nil
	}
	if err != nil {
		action = "pull"
		return false, err
	}

	ref, err := g.Repo.Head()
	if err != nil {
		action = "get the head"
		return false, err
	}

	commit, err := g.Repo.CommitObject(ref.Hash())
	if err != nil {
		action = "get last commit"
		return false, err
	}

	currentTree, err := commit.Tree()
	if err != nil {
		action = "get tree of last commit"
		return false, err
	}

	changes, err := currentTree.Diff(g.prevTree)
	if err != nil {
		action = "diff changes"
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
