/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/caos/orbiter/mntr"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/controllers"
	"github.com/caos/boom/internal/app"
	"github.com/caos/boom/internal/clientgo"
	gitcrdconfig "github.com/caos/boom/internal/gitcrd/config"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kustomize"
	"github.com/caos/boom/internal/orb"

	// +kubebuilder:scaffold:imports

	gconfig "github.com/caos/boom/internal/bundle/application/applications/grafana/config"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = toolsetsv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var toolsDirectoryPath, dashboardsDirectoryPath string
	var gitOrbConfig, gitCrdPath, gitCrdURL, gitCrdPrivateKey, gitCrdDirectoryPath string
	var enableLeaderElection, localMode bool
	var intervalSeconds int
	var gitCrdEmail, gitCrdUser string
	var metrics bool
	var metricsport string

	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.BoolVar(&localMode, "local-mode", false, "Disable the controller manager and only use the operator to handle gitcrds")

	flag.StringVar(&gitOrbConfig, "git-orbconfig", "", "The orbconfig path. If not provided, --git-crd-url and --git-crd-secret are used")

	flag.StringVar(&gitCrdURL, "git-crd-url", "https://github.com/stebenz/boom-crd.git", "The url for the git-repo to clone for the CRD")
	flag.StringVar(&gitCrdPrivateKey, "git-crd-private-key", "", "Path to private key required to clone the git-repo for the CRD")
	flag.StringVar(&gitCrdDirectoryPath, "git-crd-directory-path", "/tmp/crd", "Local path where the CRD git-repo will be cloned into")
	flag.StringVar(&gitCrdPath, "git-crd-path", "crd.yaml", "The path to the CRD in the cloned git-repo ")
	flag.StringVar(&gitCrdUser, "git-crd-user", "boom", "The name of the user used for pushing the current state in git")
	flag.StringVar(&gitCrdEmail, "git-crd-email", "boom@caos.ch", "The email of the user used for pushing the current state in git")

	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "/tmp/tools", "The local path where the tools folder should be")
	flag.StringVar(&dashboardsDirectoryPath, "dashboards-directory-path", "/dashboards", "The local path where the dashboards folder should be")

	flag.IntVar(&intervalSeconds, "intervalSeconds", 60, "defines the interval in which the reconiliation of the gitCrds runs")

	flag.BoolVar(&metrics, "metrics", false, "Defines if a metrics endpoint should be exposed")
	flag.StringVar(&metricsport, "metricsport", "2112", "Port with which the metrics endpoint will get exposed")
	flag.Parse()

	gconfig.DashboardsDirectoryPath = dashboardsDirectoryPath

	var gitCrdPrivateKeyBytes []byte

	if localMode {
		clientgo.InConfig = false
	}

	if gitOrbConfig != "" {
		orb, err := orb.ParseOrbConfig(gitOrbConfig)
		if err != nil {
			setupLog.Error(err, "unable to parse orbconfig")
			os.Exit(1)
		}

		gitCrdPrivateKeyBytes = []byte(orb.Repokey)
		gitCrdURL = orb.URL
	}

	if gitCrdPrivateKeyBytes == nil && gitCrdPrivateKey != "" {
		var err error
		gitCrdPrivateKeyBytes, err = ioutil.ReadFile(gitCrdPrivateKey)
		if err != nil {
			setupLog.Error(err, "unable to read git crd private key")
			os.Exit(1)
		}
	}

	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}
	if *verbose {
		monitor = monitor.Verbose()
	}

	// ctrl.SetLogger(monitor)

	app, err := app.New(monitor, toolsDirectoryPath, dashboardsDirectoryPath)
	if err != nil {
		setupLog.Error(err, "unable to start app")
		os.Exit(1)
	}

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

		if err := app.AddGitCrd(gitcrdConf); err != nil {
			setupLog.Error(err, "unable to start supervised crd")
			os.Exit(1)
		}

		go func() {
			// TODO: use a function scoped error variable
			for {
				started := time.Now()
				goErr := app.ReconcileGitCrds()
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
				goErr := app.WriteBackCurrentState()
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

	if metrics {
		http.Handle("/metrics", promhttp.Handler())
		address := strings.Join([]string{":", metricsport}, "")
		go func() {
			if err := http.ListenAndServe(address, nil); err != nil {
				setupLog.Error(err, "error while serving metrics endpoint")
				os.Exit(1)
			}

			monitor.WithFields(map[string]interface{}{
				"port":     metricsport,
				"endpoint": "/metrics",
			}).Info("Started metrics")
		}()
	}

	if !localMode {
		cmd, err := kustomize.New("/crd", true, false)
		if err != nil {
			setupLog.Error(err, "unable to locate crd")
			os.Exit(1)
		}

		err = errors.Wrapf(helper.Run(monitor, cmd.Build()), "Failed to apply crd")
		if err != nil {
			setupLog.Error(err, "unable to apply crd")
			os.Exit(1)
		}

		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme:             scheme,
			MetricsBindAddress: metricsAddr,
			LeaderElection:     enableLeaderElection,
			Port:               9443,
		})
		if err != nil {
			setupLog.Error(err, "unable to start manager")
			os.Exit(1)
		}

		if err = (&controllers.ToolsetReconciler{
			App:    app,
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Toolset"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Toolset")
			os.Exit(1)
		}
		// +kubebuilder:scaffold:builder

		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}

		setupLog.Info("starting manager")
	} else {
		cmd, err := kustomize.New("../../config/crd", true, false)
		if err != nil {
			setupLog.Error(err, "unable to locate crd")
			os.Exit(1)
		}

		err = errors.Wrapf(helper.Run(monitor, cmd.Build()), "Failed to apply crd")
		if err != nil {
			setupLog.Error(err, "unable to apply crd")
			os.Exit(1)
		}
	}

	<-gitCrdError
}
