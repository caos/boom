package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
	"github.com/pkg/errors"
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
		resources, err := pre.HelmPreApplySteps(spec)
		if err != nil {
			h.status = errors.Wrapf(err, "Error while processing pre-steps for application %s", pre.GetName().String())
			return h
		}

		resultfilepath := h.GetResultsFilePath(pre.GetName(), h.overlay, h.templatorDirectoryPath)

		for i, resource := range resources {
			value, isString := resource.(string)
			if isString {
				h.status = helper.AddStringObjectToYaml(resultfilepath, value)
			} else {
				h.status = helper.AddStructToYaml(resultfilepath, resource)
			}

			if h.status != nil {
				h.status = errors.Wrapf(err, "Error while adding element %i to result-file", i)
				return h
			}
		}
	}
	return h
}
