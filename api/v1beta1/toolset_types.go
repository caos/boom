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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ToolsetSpec defines the desired state of Toolset
type ToolsetSpec struct {
	ForceApply                bool                       `json:"forceApply,omitempty" yaml:"forceApply,omitempty"`
	CurrentStateFolder        string                     `json:"currentStatePath,omitempty" yaml:"currentStatePath,omitempty"`
	PreApply                  *PreApply                  `json:"preApply,omitempty" yaml:"preApply,omitempty"`
	PostApply                 *PostApply                 `json:"postApply,omitempty" yaml:"postApply,omitempty"`
	PrometheusOperator        *PrometheusOperator        `json:"prometheus-operator,omitempty" yaml:"prometheus-operator"`
	LoggingOperator           *LoggingOperator           `json:"logging-operator,omitempty" yaml:"logging-operator"`
	PrometheusNodeExporter    *PrometheusNodeExporter    `json:"prometheus-node-exporter,omitempty" yaml:"prometheus-node-exporter"`
	PrometheusSystemdExporter *PrometheusSystemdExporter `json:"prometheus-systemd-exporter,omitempty" yaml:"prometheus-systemd-exporter"`
	Grafana                   *Grafana                   `json:"grafana,omitempty" yaml:"grafana"`
	Ambassador                *Ambassador                `json:"ambassador,omitempty" yaml:"ambassador"`
	KubeStateMetrics          *KubeStateMetrics          `json:"kube-state-metrics,omitempty" yaml:"kube-state-metrics"`
	Argocd                    *Argocd                    `json:"argocd,omitempty" yaml:"argocd"`
	Prometheus                *Prometheus                `json:"prometheus,omitempty" yaml:"prometheus"`
	Loki                      *Loki                      `json:"loki,omitempty" yaml:"loki"`
}

// ToolsetStatus defines the observed state of Toolset
type ToolsetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Toolset is the Schema for the toolsets API
type Toolset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec   *ToolsetSpec   `json:"spec,omitempty"`
	Status *ToolsetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ToolsetList contains a list of Toolset
type ToolsetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*Toolset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Toolset{}, &ToolsetList{})
}
