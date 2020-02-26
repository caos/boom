package main

import (
	"flag"
	"os"

	"github.com/caos/boom/internal/templator/helm/chart/fetch"
	"github.com/caos/orbiter/mntr"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("gen-executables")
)

func main() {
	var toolsDirectoryPath string

	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "/tmp/tools", "The local path where the tools folder should be")
	flag.Parse()

	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}
	if *verbose {
		monitor = monitor.Verbose()
	}

	// ctrl.SetLogger(monitor)

	if err := fetch.All(monitor, toolsDirectoryPath); err != nil {
		setupLog.Error(err, "unable to fetch charts")
		os.Exit(1)
	}
}
