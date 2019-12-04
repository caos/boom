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

package controllers

import (
	"context"

	"github.com/caos/boom/internal/app"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
)

// ToolsetReconciler reconciles a Toolset object
type ToolsetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	App    *app.App
}

// +kubebuilder:rbac:groups=toolsets.boom.caos.ch,resources=toolsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=toolsets.boom.caos.ch,resources=toolsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=*,verbs=*
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=*,verbs=*
// +kubebuilder:rbac:groups=apps,resources=*,verbs=*
// +kubebuilder:rbac:groups=batch,resources=*,verbs=*
// +kubebuilder:rbac:groups=extensions,resources=*,verbs=*
// +kubebuilder:rbac:groups=logging.banzaicloud.io,resources=*,verbs=*
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=*,verbs=*
// +kubebuilder:rbac:groups=policy,resources=*,verbs=*
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=*,verbs=*
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=*,verbs=*
// +kubebuilder:rbac:groups=apiregistration.k8s.io,resources=*,verbs=*

func (r *ToolsetReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("toolset", req.NamespacedName)

	getToolset := func(instance runtime.Object) error {
		if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
			log.Error(err, "unable to fetch Toolset")
			return err
		}
		r.Log.Info("crd successfully loaded")
		return nil
	}

	if err := r.App.ReconcileCrd("v1beta1", req.NamespacedName.String(), getToolset); err != nil {
		log.Error(err, "unable to reconcile Toolset")
	}
	r.Log.Info("Toolset sucessfully reconciled")

	return ctrl.Result{}, nil
}

func (r *ToolsetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsetsv1beta1.Toolset{}).
		Complete(r)
}
