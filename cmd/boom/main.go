package main

import (
	"flag"
	"io/ioutil"
	"os"
	"time"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"

	"github.com/caos/boom/internal/app"
	"github.com/caos/boom/internal/clientgo"
	gitcrdconfig "github.com/caos/boom/internal/gitcrd/config"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/orb"

	gconfig "github.com/caos/boom/internal/bundle/application/applications/grafana/config"
)

func main() {

	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}

	var metricsAddr string
	var toolsDirectoryPath, dashboardsDirectoryPath string
	var gitOrbConfig, gitCrdPath, gitCrdURL, gitCrdPrivateKey, gitCrdDirectoryPath string
	var enableLeaderElection, localMode bool
	var intervalSeconds int
	var gitCrdEmail, gitCrdUser string
	var limitResources int64

	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.BoolVar(&localMode, "local-mode", false, "Disable the controller manager and only use the operator to handle gitcrds")

	flag.StringVar(&gitOrbConfig, "git-orbconfig", "", "The orbconfig path. If not provided, --git-crd-url and --git-crd-secret are used")

	flag.StringVar(&gitCrdURL, "git-crd-url", "https://github.com/stebenz/boom-crd.git", "The url for the git-repo to clone for the CRD")
	flag.StringVar(&gitCrdPrivateKey, "git-crd-private-key", "", "Path to private key required to clone the git-repo for the CRD")
	flag.StringVar(&gitCrdDirectoryPath, "git-crd-directory-path", "/tmp/crd", "Local path where the CRD git-repo will be cloned into")
	flag.StringVar(&gitCrdPath, "git-crd-path", "crd-test.yaml", "The path to the CRD in the cloned git-repo ")
	flag.StringVar(&gitCrdUser, "git-crd-user", "boom", "The name of the user used for pushing the current state in git")
	flag.StringVar(&gitCrdEmail, "git-crd-email", "boom@caos.ch", "The email of the user used for pushing the current state in git")

	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "/tmp/tools", "The local path where the tools folder should be")
	flag.StringVar(&dashboardsDirectoryPath, "dashboards-directory-path", "/dashboards", "The local path where the dashboards folder should be")

	flag.IntVar(&intervalSeconds, "intervalSeconds", 60, "defines the interval in which the reconiliation of the gitCrds runs")
	flag.Int64Var(&limitResources, "limit", 0, "Defines the limit which is used by the request for current state")
	flag.Parse()

	if *verbose {
		monitor = monitor.Verbose()
	}
	gconfig.DashboardsDirectoryPath = dashboardsDirectoryPath

	var gitCrdPrivateKeyBytes []byte

	if localMode {
		clientgo.InConfig = false
	}
	if limitResources != 0 {
		clientgo.Limit = limitResources
	}

	if gitOrbConfig != "" {
		orbconfig, err := orb.ParseOrbConfig(gitOrbConfig)
		if err != nil {
			monitor.Error(err)
			os.Exit(1)
		}

		gitCrdPrivateKeyBytes = []byte(orbconfig.Repokey)
		gitCrdURL = orbconfig.URL
	}

	if gitCrdPrivateKeyBytes == nil && gitCrdPrivateKey != "" {
		var err error
		gitCrdPrivateKeyBytes, err = ioutil.ReadFile(gitCrdPrivateKey)
		if err != nil {
			monitor.Error(errors.Wrap(err, "unable to read git crd private key"))
			os.Exit(1)
		}
	}

	// ctrl.SetLogger(monitor)

	appStruct := app.New(monitor, toolsDirectoryPath)

	var gitCrdError chan error
	if gitCrdPath != "" {
		gitcrdMonitor := monitor.WithFields(map[string]interface{}{"type": "gitcrd"})

		gitcrdConf := &gitcrdconfig.Config{
			Monitor:          gitcrdMonitor,
			CrdDirectoryPath: gitCrdDirectoryPath,
			CrdUrl:           gitCrdURL,
			PrivateKey:       gitCrdPrivateKeyBytes,
			CrdPath:          gitCrdPath,
			User:             gitCrdUser,
			Email:            gitCrdEmail,
		}

		if err := appStruct.AddGitCrd(gitcrdConf); err != nil {
			monitor.Error(errors.Wrap(err, "unable to start supervised crd"))
			os.Exit(1)
		}

		go func() {
			// TODO: use a function scoped error variable
			for {
				started := time.Now()
				goErr := appStruct.ReconcileGitCrds()
				recMonitor := monitor.WithFields(map[string]interface{}{
					"took": time.Since(started),
				})
				if goErr != nil {
					recMonitor.Error(goErr)
				}
				recMonitor.Info("Reconciling iteration done")
				time.Sleep(time.Duration(intervalSeconds) * time.Second)
			}
		}()

		go func() {
			for {
				started := time.Now()
				goErr := appStruct.WriteBackCurrentState()
				recMonitor := monitor.WithFields(map[string]interface{}{
					"took": time.Since(started),
				})
				if goErr != nil {
					recMonitor.Error(goErr)
				}
				recMonitor.Info("Current state iteration done")
				time.Sleep(time.Duration(intervalSeconds) * time.Second)
			}
		}()
	}

	cmd, err := kustomize.New("../../config/crd", true, false)
	if err != nil {
		monitor.Error(errors.Wrap(err, "unable to locate crd"))
		os.Exit(1)
	}

	err = errors.Wrapf(helper.Run(monitor, cmd.Build()), "Failed to apply crd")
	if err != nil {
		monitor.Error(errors.Wrap(err, "unable to apply crd"))
		os.Exit(1)
	}
	<-gitCrdError
}
