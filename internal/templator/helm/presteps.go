package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
)

type TemplatorPreSteps interface {
	templator.HelmApplication
	HelmPreApplySteps(*v1beta1.ToolsetSpec) ([]interface{}, error)
}

func (h *Helm) preApplySteps(app interface{}, spec *v1beta1.ToolsetSpec) templator.Templator {
	if h.status != nil {
		return h
	}

	pre, ok := app.(TemplatorPreSteps)
	if ok {

		logFields := map[string]interface{}{
			"application": pre.GetName().String(),
			"overlay":     h.overlay,
		}
		logFields["logID"] = "HELM-36S6r895dyInePv"
		h.logger.WithFields(logFields).Info("Additional steps before apply")

		resources, err := pre.HelmPreApplySteps(spec)
		if err != nil {
			h.status = err
			return h
		}

		resultfilepath := h.GetResultsFilePath(pre.GetName(), h.overlay, h.templatorDirectoryPath)

		for _, resource := range resources {
			value, isString := resource.(string)
			if isString {
				h.status = helper.AddStringToYaml(resultfilepath, value)
			} else {
				h.status = helper.AddStructToYaml(resultfilepath, resource)
			}

			if h.status != nil {
				return h
			}
		}
	}
	return h
}
