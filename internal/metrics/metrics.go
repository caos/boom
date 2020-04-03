package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

type metricsStruct struct {
	gitClone               *prometheus.CounterVec
	crdFormat              *prometheus.CounterVec
	currentStateWrite      *prometheus.CounterVec
	currentStateRead       *prometheus.CounterVec
	reconcilingBundle      *prometheus.CounterVec
	reconcilingApplication *prometheus.CounterVec
}

var (
	metrics = &metricsStruct{
		gitClone:               newCounterVec("caos_boom_git_clone", "Counter how many times git repositories were cloned", "result", "url"),
		crdFormat:              newCounterVec("caos_boom_crd_format", "Counter how many failures there were with the crd unmarshalling", "result", "url", "path", "reason"),
		currentStateWrite:      newCounterVec("caos_boom_current_state_write", "Counter how many times the current state was written", "result", "action", "url", "path"),
		currentStateRead:       newCounterVec("caos_boom_current_state_read", "Counter how many times the current state was read", "result", "action"),
		reconcilingBundle:      newCounterVec("caos_boom_reconciling_bundle", "Counter how many times the bundle was reconciled", "result", "bundle"),
		reconcilingApplication: newCounterVec("caos_boom_reconciling_application", "Counter how many times a application was reconciled", "result", "application", "templator", "deploy"),
	}
)

func newCounterVec(name string, help string, labels ...string) *prometheus.CounterVec {
	counterVec := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	},
		labels,
	)

	err := prometheus.Register(counterVec)
	_, ok := err.(prometheus.AlreadyRegisteredError)
	if err != nil && !ok {
		return counterVec
	}

	return counterVec
}

func SuccessfulGitClone(url string) {
	metrics.gitClone.With(prometheus.Labels{
		"result": "success",
		"url":    url,
	}).Inc()
}

func FailedGitClone(url string) {
	metrics.gitClone.With(prometheus.Labels{
		"result": "failure",
		"url":    url,
	}).Inc()
}

func WrongCRDFormat(url string, path string) {
	metrics.crdFormat.With(prometheus.Labels{
		"result": "failure",
		"url":    url,
		"path":   path,
		"reason": "structure",
	}).Inc()
}

func UnsupportedAPIGroup(url string, path string) {
	metrics.crdFormat.With(prometheus.Labels{
		"result": "failure",
		"url":    url,
		"path":   path,
		"reason": "apiGroup",
	}).Inc()

}

func UnsupportedVersion(url string, path string) {
	metrics.crdFormat.With(prometheus.Labels{
		"result": "failure",
		"url":    url,
		"path":   path,
		"reason": "version",
	}).Inc()
}

func WroteCurrentState(url string, path string) {
	metrics.currentStateWrite.With(prometheus.Labels{
		"result": "success",
		"url":    url,
		"path":   path,
		"action": "write",
	}).Inc()
}

func FailedWritingCurrentState(url string, path string) {
	metrics.currentStateWrite.With(prometheus.Labels{
		"result": "failure",
		"url":    url,
		"path":   path,
		"action": "write",
	}).Inc()
}

func SuccessfulReadingCurrentState() {
	metrics.currentStateRead.With(prometheus.Labels{
		"result": "success",
		"action": "read",
	}).Inc()
}

func FailedReadingCurrentState() {
	metrics.currentStateRead.With(prometheus.Labels{
		"result": "failure",
		"action": "read",
	}).Inc()
}

func SuccessfulReconcilingBundle(bundle string) {
	metrics.reconcilingBundle.With(prometheus.Labels{
		"result": "success",
		"bundle": bundle,
	}).Inc()
}

func FailureReconcilingBundle(bundle string) {
	metrics.reconcilingBundle.With(prometheus.Labels{
		"result": "failure",
		"bundle": bundle,
	}).Inc()
}

func SuccessfulReconcilingApplication(app string, templator string, deploy bool) {
	metrics.reconcilingApplication.With(prometheus.Labels{
		"result":      "success",
		"application": app,
		"deploy":      strconv.FormatBool(deploy),
		"templator":   templator,
	}).Inc()
}

func FailureReconcilingApplication(app string, templator string, deploy bool) {
	metrics.reconcilingApplication.With(prometheus.Labels{
		"result":      "failure",
		"application": app,
		"deploy":      strconv.FormatBool(deploy),
		"templator":   templator,
	}).Inc()
}
