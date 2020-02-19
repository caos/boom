package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caos/boom/internal/git"
	"github.com/caos/boom/internal/kubectl"
	"github.com/caos/orbiter/logging"
)

func UseFolder(logger logging.Logger, git *git.Client, deploy bool, tempDirectory, folderRelativePath string) error {

	command := kubectl.New("apply")
	if !deploy {
		command = kubectl.New("delete").AddFlag("--ignore-not-found")
	}

	if err := git.Clone(); err != nil {
		return err
	}

	folderPath := filepath.Join(tempDirectory, folderRelativePath)

	if err := RecreatePath(folderPath); err != nil {
		return err
	}

	files, err := git.ReadFolder(folderRelativePath)
	if err != nil {
		return err
	}

	for filename, file := range files {
		filePath := filepath.Join(folderPath, filename)
		err := ioutil.WriteFile(filePath, file, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return Run(logger, command.AddParameter("-f", folderPath).Build())
}
