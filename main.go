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

	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	toolsetsv1beta1 "github.com/caos/toolsop/api/v1beta1"
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
	var toolsDirectoryPath, toolsSecret, toolsUrl, toolsetsPath string
	var gitCrdPath, gitCrdUrl, gitCrdSecret, gitCrdDirectoryPath string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.StringVar(&gitCrdUrl, "git-crd-url", "git@github.com:caos/tools.git", "The url for the git-repo to clone for the CRD ")
	flag.StringVar(&gitCrdSecret, "git-crd-secret", "./secretdata/ssh-keys/id_rsa-toolsop-tools-read", "Path to Secret to clone the git-repo for the CRD")
	flag.StringVar(&gitCrdDirectoryPath, "git-crd-directory-path", "/tmp/crd", "Local path where the CRD git-repo will be cloned into")
	flag.StringVar(&gitCrdPath, "git-crd-path", "crd/example/crd.yaml", "The path to the CRD in the cloned git-repo ")

	flag.StringVar(&toolsUrl, "tools-url", "git@github.com:caos/tools.git", "The URL from where the tools-repo should be cloned from")
	flag.StringVar(&toolsDirectoryPath, "tools-directory-path", "/tmp/tools", "The local path where the tools-repo should be cloned to")
	flag.StringVar(&toolsSecret, "tools-secret", "./secretdata/ssh-keys/id_rsa-toolsop-tools-read", "The secret which get used to clone the tools-repo")
	flag.StringVar(&toolsetsPath, "tools-toolset-path", "./toolsets/toolsets.yaml", "The path to the yaml which defined the toolsets and their versions")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	// mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
	// 	Scheme:             scheme,
	// 	MetricsBindAddress: metricsAddr,
	// 	LeaderElection:     enableLeaderElection,
	// 	Port:               9443,
	// })
	// if err != nil {
	// 	setupLog.Error(err, "unable to start manager")
	// 	os.Exit(1)
	// }

	ctx := context.Background()

	app, err := app.New(toolsDirectoryPath, gitCrdDirectoryPath, toolsetsPath, toolsUrl, toolsSecret)
	if err != nil {
		setupLog.Error(err, "unable to start app")
		os.Exit(1)
	}
	// if err = (&controllers.ToolsetReconciler{
	// 	App:    app,
	// 	Client: mgr.GetClient(),
	// 	Log:    ctrl.Log.WithName("controllers").WithName("Toolset"),
	// 	Scheme: mgr.GetScheme(),
	// }).SetupWithManager(mgr); err != nil {
	// 	setupLog.Error(err, "unable to create controller", "controller", "Toolset")
	// 	os.Exit(1)
	// }
	// +kubebuilder:scaffold:builder

	errChan := make(chan error)
	if gitCrdPath != "" {
		if err := app.AddGitCrd(gitCrdUrl, gitCrdSecret, gitCrdPath); err != nil {
			setupLog.Error(err, "unable to start supervised crd")
			os.Exit(1)
		}

		ctxChild, cancel := context.WithCancel(ctx)

		go func() {
			<-ctxChild.Done()
			cancel()
		}()

		go func() {
			for err == nil {
				err = app.ReconcileGitCrds()
				time.Sleep(10 * time.Second)
			}

			setupLog.Error(err, "unable to maintaining supervised crd")
			errChan <- err
		}()
	}

	setupLog.Info("starting manager")
	// if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
	// 	setupLog.Error(err, "problem running manager")
	// 	os.Exit(1)
	// }

	<-errChan
	app.CleanUp()
	setupLog.Info("stopped")
}
