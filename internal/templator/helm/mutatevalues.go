package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/orbiter/mntr"
)

type TemplatorMutateValues interface {
	templator.HelmApplication
	HelmMutateValues(mntr.Monitor, *v1beta1.ToolsetSpec, string) error
}

func (h *Helm) mutateValue(app interface{}, spec *v1beta1.ToolsetSpec, valuesAbsFilePath string) templator.Templator {
	if h.status != nil {
		return h
	}

	mutate, ok := app.(TemplatorMutateValues)
	if ok {

		logFields := map[string]interface{}{
			"application": mutate.GetName().String(),
			"overlay":     h.overlay,
		}
		mutateMonitor := h.monitor.WithFields(logFields)

		mutateMonitor.Debug("Mutate values")

		if err := mutate.HelmMutateValues(mutateMonitor, spec, valuesAbsFilePath); err != nil {
			h.status = err
			return h
		}
	}

	return h
}
