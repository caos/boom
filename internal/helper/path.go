package helper

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func GetAbsPath(pathParts ...string) (string, error) {

	filePath := filepath.Join(pathParts...)
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "Error while getting absolute path for %s", filePath)
	}
	return absFilePath, nil
}
