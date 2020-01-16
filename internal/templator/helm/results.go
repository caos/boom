package helm

import (
	"os"

	"github.com/caos/boom/internal/templator"
)

func (h *Helm) deleteResults(app templator.HelmApplication) error {
	resultsFileDirectory := h.getResultsFileDirectory(app.GetName(), h.overlay, h.templatorDirectoryPath)
	if err := os.RemoveAll(resultsFileDirectory); err != nil {
		return err
	}

	if err := os.MkdirAll(resultsFileDirectory, os.ModePerm); err != nil {
		return err
	}

	return nil
}
