package helm

import (
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/boom/internal/templator/helm/helmcommand"
	"github.com/pkg/errors"
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

	monitor := h.monitor.WithFields(logFields)

	monitor.Debug("Deleting old results")
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

	if err := h.runHelmTemplate(h.overlay, app, spec, resultAbsFilePath); err != nil {
		h.status = err
		return h
	}

	deleteKind := "Namespace"
	h.status = helper.DeleteKindFromYaml(resultAbsFilePath, deleteKind)
	if h.status != nil {
		h.status = errors.Wrapf(h.status, "Error while trying to delete kind %s from results", deleteKind)
		return h
	}

	// mutate templated results
	if err := h.mutate(app, spec).GetStatus(); err != nil {
		h.status = err
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

func (h *Helm) runHelmTemplate(overlay string, app templator.HelmApplication, spec *v1beta1.ToolsetSpec, resultAbsFilePath string) error {
	if h.status != nil {
		return h.status
	}

	logFields := map[string]interface{}{
		"application": app.GetName().String(),
		"overlay":     overlay,
		"action":      "templating",
	}
	monitor := h.monitor.WithFields(logFields)

	monitor.Debug("Generate values with toolsetSpec")
	chartInfo := app.GetChartInfo()
	values := app.SpecToHelmValues(monitor, spec)

	valuesAbsFilePath, err := helper.GetAbsPath(h.templatorDirectoryPath, app.GetName().String(), overlay, "values.yaml")
	if err != nil {
		monitor.Error(err)
		return err
	}

	if err := helper.AddStructToYaml(valuesAbsFilePath, values); err != nil {
		monitor.Error(err)
		return err
	}

	monitor.Debug("Generate result through helm template")
	out, err := helmcommand.Template(&helmcommand.TemplateConfig{
		TempFolderPath:   h.templatorDirectoryPath,
		ChartName:        chartInfo.Name,
		ReleaseName:      app.GetName().String(),
		ReleaseNamespace: app.GetNamespace(),
		ValuesFilePath:   valuesAbsFilePath,
	})
	if err != nil {
		monitor.Error(err)
		return err
	}

	return helper.AddStringToYaml(resultAbsFilePath, string(out))
}
