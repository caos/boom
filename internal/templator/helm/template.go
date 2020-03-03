package helm

import (
	"os"
	"path/filepath"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/templator"
	"github.com/caos/boom/internal/templator/helm/helmcommand"
	"github.com/pkg/errors"
)

func (h *Helm) Template(appInterface interface{}, spec *v1beta1.ToolsetSpec, resultFunc func(resultFilePath, namespace string) error) error {
	app, err := checkTemplatorInterface(appInterface)
	if err != nil {
		return err
	}

	logFields := map[string]interface{}{
		"application": app.GetName().String(),
		"overlay":     h.overlay,
	}

	monitor := h.monitor.WithFields(logFields)

	monitor.Debug("Deleting old results")
	err = h.deleteResults(app)
	if err != nil {
		return err
	}

	var resultAbsFilePath string
	resultfilepath := h.GetResultsFilePath(app.GetName(), h.overlay, h.templatorDirectoryPath)

	resultAbsFilePath, err = filepath.Abs(resultfilepath)
	if err != nil {
		return err
	}

	valuesAbsFilePath, err := helper.GetAbsPath(h.templatorDirectoryPath, app.GetName().String(), h.overlay, "values.yaml")
	if err != nil {
		monitor.Error(err)
		return err
	}

	if err := h.prepareHelmTemplate(h.overlay, app, spec, valuesAbsFilePath); err != nil {
		monitor.Error(err)
		return err
	}

	if err := h.mutateValue(app, spec, valuesAbsFilePath); err != nil {
		monitor.Error(err)
		return err
	}

	if err := h.runHelmTemplate(h.overlay, app, valuesAbsFilePath, resultAbsFilePath); err != nil {
		monitor.Error(err)
		return err
	}

	deleteKind := "Namespace"
	err = helper.DeleteKindFromYaml(resultAbsFilePath, deleteKind)
	if err != nil {
		return errors.Wrapf(err, "Error while trying to delete kind %s from results", deleteKind)
	}

	// mutate templated results
	if err := h.mutate(app, spec); err != nil {
		return err
	}

	// pre apply steps
	if err := h.preApplySteps(app, spec); err != nil {
		return err
	}

	// func to apply
	return resultFunc(resultAbsFilePath, app.GetNamespace())
}

func (h *Helm) prepareHelmTemplate(overlay string, app templator.HelmApplication, spec *v1beta1.ToolsetSpec, valuesAbsFilePath string) error {

	logFields := map[string]interface{}{
		"application": app.GetName().String(),
		"overlay":     overlay,
		"action":      "preparetemplating",
	}
	monitor := h.monitor.WithFields(logFields)

	monitor.Debug("Generate values with toolsetSpec")
	values := app.SpecToHelmValues(monitor, spec)

	if helper.FileExists(valuesAbsFilePath) {
		if err := os.Remove(valuesAbsFilePath); err != nil {
			monitor.Error(err)
			return err
		}
	}

	if err := helper.AddStructToYaml(valuesAbsFilePath, values); err != nil {
		monitor.Error(err)
		return err
	}
	return nil
}

func (h *Helm) runHelmTemplate(overlay string, app templator.HelmApplication, valuesAbsFilePath, resultAbsFilePath string) error {
	logFields := map[string]interface{}{
		"application": app.GetName().String(),
		"overlay":     overlay,
		"action":      "templating",
	}
	monitor := h.monitor.WithFields(logFields)

	chartInfo := app.GetChartInfo()

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
