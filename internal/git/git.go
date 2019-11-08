package git

import (
	"os"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func CloneRepo(path, url, secretPath string) (*git.Repository, error) {
	publicKey, err := ssh.NewPublicKeysFromFile("git", secretPath, "")
	if err != nil {
		return nil, err
	}

	return git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     publicKey,
	})
}
