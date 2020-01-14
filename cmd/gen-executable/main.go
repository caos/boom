package main

import (
	"flag"
	"os"

	"github.com/caos/boom/internal/templator/helm/charts"

	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/kubebuilder"
	"github.com/caos/orbiter/logging/stdlib"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("gen-executables")
)

func main() {
	var toolsDirectoryPath string

	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "../../tools", "The local path where the tools folder should be")
	flag.Parse()

	logger := logcontext.Add(stdlib.New(os.Stdout))
	if *verbose {
		logger = logger.Verbose()
	}

	ctrl.SetLogger(kubebuilder.New(logger))

	if err := charts.FetchAll(logger, toolsDirectoryPath); err != nil {
		setupLog.Error(err, "unable to fetch charts")
		os.Exit(1)
	}
}
