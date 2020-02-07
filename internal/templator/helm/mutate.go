package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/orbiter/logging"
)

type TemplatorMutate interface {
	templator.HelmApplication
	HelmMutate(logging.Logger, *v1beta1.ToolsetSpec, string) error
}

func (h *Helm) mutate(app interface{}, spec *v1beta1.ToolsetSpec) templator.Templator {
	if h.status != nil {
		return h
	}

	mutate, ok := app.(TemplatorMutate)
	if ok {

		logFields := map[string]interface{}{
			"application": mutate.GetName().String(),
			"overlay":     h.overlay,
		}
		mutateLogger := h.logger.WithFields(logFields)

		logFields["logID"] = "HELM-wm1sjJ2Bmr0VMYV"
		mutateLogger.WithFields(logFields).Debug("Mutate before apply")

		resultfilepath := h.GetResultsFilePath(mutate.GetName(), h.overlay, h.templatorDirectoryPath)

		if err := mutate.HelmMutate(mutateLogger, spec, resultfilepath); err != nil {
			h.status = err
			return h
		}
	}
	return h
}
