package helm

import (
	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"
)

type TemplatorPreSteps interface {
	templator.HelmApplication
	HelmPreApplySteps(mntr.Monitor, *v1beta1.ToolsetSpec) ([]interface{}, error)
}

func (h *Helm) preApplySteps(app interface{}, spec *v1beta1.ToolsetSpec) error {

	pre, ok := app.(TemplatorPreSteps)
	if ok {

		logFields := map[string]interface{}{
			"application": pre.GetName().String(),
			"overlay":     h.overlay,
		}

		monitor := h.monitor.WithFields(logFields)
		monitor.Debug("Pre-steps")
		resources, err := pre.HelmPreApplySteps(monitor, spec)
		if err != nil {
			return errors.Wrapf(err, "Error while processing pre-steps for application %s", pre.GetName().String())
		}

		resultfilepath := h.GetResultsFilePath(pre.GetName(), h.overlay, h.templatorDirectoryPath)

		for i, resource := range resources {
			value, isString := resource.(string)

			if isString {
				err := helper.AddStringObjectToYaml(resultfilepath, value)
				if err != nil {
					return errors.Wrapf(err, "Error while adding element %d to result-file", i)
				}
			} else {
				err = helper.AddStructToYaml(resultfilepath, resource)
				if err != nil {
					return errors.Wrapf(err, "Error while adding element %d to result-file", i)
				}
			}
		}
	}
	return nil
}
