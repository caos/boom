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
	"context"
	"flag"
	"os"
	"time"

	"github.com/pkg/errors"

	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/kubebuilder"
	"github.com/caos/orbiter/logging/stdlib"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
	"github.com/caos/toolsop/controllers"
	"github.com/caos/toolsop/internal/app"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
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

func main() {
	var metricsAddr string
	var toolsDirectoryPath, toolsetsPath string
	var gitCrdPath, gitCrdUrl, gitCrdSecret, gitCrdDirectoryPath string
	var enableLeaderElection, localMode bool
	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&localMode, "local-mode", false,
		"Disable the controller manager and only use the operator to handle gitcrds")

	flag.StringVar(&gitCrdUrl, "git-crd-url", "git@github.com:caos/tools.git", "The url for the git-repo to clone for the CRD ")
	flag.StringVar(&gitCrdSecret, "git-crd-secret", "config/manager/secret/id_rsa-toolsop-tools-read", "Path to Secret to clone the git-repo for the CRD")
	flag.StringVar(&gitCrdDirectoryPath, "git-crd-directory-path", "/tmp/crd", "Local path where the CRD git-repo will be cloned into")
	flag.StringVar(&gitCrdPath, "git-crd-path", "crd/example/crd.yaml", "The path to the CRD in the cloned git-repo ")

	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "tools", "The local path where the tools folder should be")
	flag.StringVar(&toolsetsPath, "tools-toolset-path", "toolsets", "The path to the fold structue which defines the toolsets and their versions")
	flag.Parse()

	ctx := context.Background()

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

	if !localMode {
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
	}

	if gitCrdPath == "" {
		return
	}

	for {
		if err := app.AddGitCrd(gitCrdUrl, gitCrdSecret, gitCrdPath); err != nil {
			setupLog.Error(err, "unable to start supervised crd")
			os.Exit(1)
		}

		if err := app.ReconcileGitCrds(); err != nil {
			logger.Error(errors.Wrap(err, "unable to maintaining supervised crd"))
		}
		if err := app.CleanUp(); err != nil {
			setupLog.Error(err, "cleaning up failed")
			os.Exit(1)
		}
		time.Sleep(10 * time.Second)
	}
}
