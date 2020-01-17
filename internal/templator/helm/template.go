package helm

import (
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/boom/internal/templator/helm/helmcommand"
)

func (h *Helm) Template(appInterface interface{}, spec *v1beta1.ToolsetSpec, resultFunc func(resultFilePath, namespace string) error) templator.Templator {
	if h.status != nil {
		return h
	}

	app, err := checkTemplatorInterface(appInterface)
	if err != nil {
		h.status = err
		return h
	}

	logFields := map[string]interface{}{
		"application": app.GetName().String(),
		"overlay":     h.overlay,
	}

	logFields["logID"] = "HELM-zqYnVAzGXHqBbhu"
	h.logger.WithFields(logFields).Info("Deleting old results")
	h.status = h.deleteResults(app)
	if h.status != nil {
		return h
	}

	var resultAbsFilePath string
	resultfilepath := h.GetResultsFilePath(app.GetName(), h.overlay, h.templatorDirectoryPath)

	resultAbsFilePath, h.status = filepath.Abs(resultfilepath)
	if h.status != nil {
		return h
	}

	out, err := h.runHelmTemplate(h.overlay, app, spec)
	if err != nil {
		h.status = err
		return h
	}

	h.status = helper.AddStringToYaml(resultAbsFilePath, out)
	if h.status != nil {
		return h
	}

	h.status = helper.DeleteKindFromYaml(resultAbsFilePath, "Namespace")
	if h.status != nil {
		return h
	}

	// pre apply steps
	if err := h.preApplySteps(app, spec).GetStatus(); err != nil {
		h.status = err
		return h
	}

	// func to apply
	h.status = resultFunc(resultAbsFilePath, app.GetNamespace())
	return h
}

func (h *Helm) runHelmTemplate(overlay string, app templator.HelmApplication, spec *v1beta1.ToolsetSpec) (string, error) {
	logFields := map[string]interface{}{
		"application": app.GetName().String(),
		"overlay":     overlay,
	}

	logFields["logID"] = "HELM-D3KJ1Qv5F8fxVA5"
	h.logger.WithFields(logFields).Info("Generate values with toolsetSpec")
	chartInfo := app.GetChartInfo()
	values := app.SpecToHelmValues(spec)

	valuesFilePath := filepath.Join(h.templatorDirectoryPath, app.GetName().String(), overlay, "values.yaml")
	valuesAbsFilePath, err := filepath.Abs(valuesFilePath)
	if err != nil {
		return "", err
	}

	if err := helper.StructToYaml(values, valuesAbsFilePath); err != nil {
		return "", err
	}

	logFields["logID"] = "HELM-siNod1Y2nYCVW0r"
	h.logger.WithFields(logFields).Info("Generate templator files")

	out, err := helmcommand.Template(h.templatorDirectoryPath, chartInfo.Name, app.GetName().String(), app.GetNamespace(), valuesAbsFilePath)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
