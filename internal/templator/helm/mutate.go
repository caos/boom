package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/templator"
)

type TemplatorMutate interface {
	templator.HelmApplication
	HelmMutate(*v1beta1.ToolsetSpec, string) error
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
		logFields["logID"] = "HELM-wm1sjJ2Bmr0VMYV"
		h.logger.WithFields(logFields).Info("Mutate before apply")

		resultfilepath := h.GetResultsFilePath(mutate.GetName(), h.overlay, h.templatorDirectoryPath)

		if err := mutate.HelmMutate(spec, resultfilepath); err != nil {
			h.status = err
			return h
		}
	}
	return h
}
