package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/orbiter/mntr"
)

type TemplatorMutate interface {
	templator.HelmApplication
	HelmMutate(mntr.Monitor, *v1beta1.ToolsetSpec, string) error
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
		mutateMonitor := h.monitor.WithFields(logFields)

		mutateMonitor.WithFields(logFields).Debug("Mutate before apply")

		resultfilepath := h.GetResultsFilePath(mutate.GetName(), h.overlay, h.templatorDirectoryPath)

		if err := mutate.HelmMutate(mutateMonitor, spec, resultfilepath); err != nil {
			h.status = err
			return h
		}
	}
	return h
}
