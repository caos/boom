package git

import (
	"io/ioutil"
	"os"

	"github.com/caos/utils/logging"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	git "gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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

	logging.Log("GIT-4Sia0VjJ79gb7cw").Info("Cloned...")
	ref, err := g.Repo.Head()
	if err != nil {
		logging.Log("GIT-mMj1dZIWSoG2nZx").OnError(err).Debugf("Failed to get head of repo %s", url)
		return nil, err
	}

	logging.Log("GIT-4Sia0VjJ79gb7cw").Info("Get last commit...")
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
	logging.Log("GIT-pQnw5FfIqAk0eWc").Infof("Cloned new GitCRD %s%s", url, localPath)

	return g, nil
}

func (g *Git) cloneRepo(localPath, url, secretPath string) (*git.Repository, error) {

	logging.Log("GIT-vNU9maj2Rfo5rRU").Infof("SecretPath %s", secretPath)
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

	logging.Log("GIT-6ccjBaSlm0DwboL").Infof("PlainClone %s %s", localPath, url)
	return git.PlainClone(localPath, false, &git.CloneOptions{
		URL:          url,
		SingleBranch: true,
		Depth:        1,
		Progress:     os.Stdout,
		Auth:         auth,
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
