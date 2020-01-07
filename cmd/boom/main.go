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
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"

	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/kubebuilder"
	"github.com/caos/orbiter/logging/stdlib"
	"github.com/pkg/errors"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/controllers"
	"github.com/caos/boom/internal/app"
	"github.com/caos/boom/internal/helper"
	"github.com/caos/boom/internal/kustomize"
	// +kubebuilder:scaffold:imports
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

type Orb struct {
	URL     string
	Repokey string
}

func main() {
	var metricsAddr string
	var toolsDirectoryPath, toolsetsPath string
	var gitOrbConfig, gitCrdPath, gitCrdURL, gitCrdPrivateKey, gitCrdDirectoryPath string
	var enableLeaderElection, localMode bool
	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&localMode, "local-mode", false,
		"Disable the controller manager and only use the operator to handle gitcrds")

	flag.StringVar(&gitOrbConfig, "git-orbconfig", "", "The orbconfig path. If not provided, --git-crd-url and --git-crd-secret are used")
	flag.StringVar(&gitCrdURL, "git-crd-url", "https://github.com/stebenz/boom-crd.git", "The url for the git-repo to clone for the CRD")
	flag.StringVar(&gitCrdPrivateKey, "git-crd-private-key", "", "Path to private key required to clone the git-repo for the CRD")
	flag.StringVar(&gitCrdDirectoryPath, "git-crd-directory-path", "/tmp/crd", "Local path where the CRD git-repo will be cloned into")
	flag.StringVar(&gitCrdPath, "git-crd-path", "crd.yaml", "The path to the CRD in the cloned git-repo ")

	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "../../tools", "The local path where the tools folder should be")
	flag.StringVar(&toolsetsPath, "toolsq-toolset-path", "toolsets", "The path to the fold structue which defines the toolsets and their versions")
	flag.Parse()

	var gitCrdPrivateKeyBytes []byte

	if gitOrbConfig != "" {
		gitOrbConfig, err := ioutil.ReadFile(gitOrbConfig)
		if err != nil {
			setupLog.Error(err, "unable to read orbconfig")
			os.Exit(1)
		}

		orb := Orb{}
		if err := yaml.Unmarshal(gitOrbConfig, &orb); err != nil {
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

	logger := logcontext.Add(stdlib.New(os.Stdout))
	if *verbose {
		logger = logger.Verbose()
	}

	ctrl.SetLogger(kubebuilder.New(logger))

	app, err := app.New(logger, toolsDirectoryPath, gitCrdDirectoryPath, toolsetsPath)
	if err != nil {
		setupLog.Error(err, "unable to start app")
		os.Exit(1)
	}

	var gitCrdError chan error
	if gitCrdPath != "" {
		if err := app.AddGitCrd(gitCrdURL, gitCrdPrivateKeyBytes, gitCrdPath); err != nil {
			setupLog.Error(err, "unable to start supervised crd")
			os.Exit(1)
		}

		go func() {
			// TODO: use a function scoped error variable
			for {
				started := time.Now()
				goErr := app.ReconcileGitCrds()
				recLogger := logger.WithFields(map[string]interface{}{
					"took": time.Since(started),
				})
				if goErr != nil {
					recLogger.Error(goErr)
					gitCrdError <- goErr
				}
				recLogger.Info("Iteration done")
				time.Sleep(10 * time.Second)
			}
		}()
	}

	if !localMode {
		cmd, err := kustomize.New("/crd")
		if err != nil {
			setupLog.Error(err, "unable to locate crd")
			os.Exit(1)
		}

		err = errors.Wrapf(helper.Run(logger, cmd.Build()), "Failed to apply crd")
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
		cmd, err := kustomize.New("../../config/crd")
		if err != nil {
			setupLog.Error(err, "unable to locate crd")
			os.Exit(1)
		}

		err = errors.Wrapf(helper.Run(logger, cmd.Build()), "Failed to apply crd")
		if err != nil {
			setupLog.Error(err, "unable to apply crd")
			os.Exit(1)
		}

		<-gitCrdError
	}

}
