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
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	// PrometheusOperator defines the desired state for the prometheus operator
	PrometheusOperator     *PrometheusOperator     `json:"prometheus-operator,omitempty" yaml:"prometheus-operator"`
	LoggingOperator        *LoggingOperator        `json:"logging-operator,omitempty" yaml:"logging-operator"`
	PrometheusNodeExporter *PrometheusNodeExporter `json:"prometheus-node-exporter,omitempty" yaml:"prometheus-node-exporter"`
}

type PrometheusOperator struct {
	Prefix         string          `json:"prefix,omitempty"`
	Namespace      string          `json:"namespace,omitempty"`
	Prometheus     *Prometheus     `json:"prometheus,omitempty"`
	ServiceMonitor *ServiceMonitor `json:"servicemonitor,omitempty"`
}

type ServiceMonitor struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Interval    string            `json:"interval,omitempty"`
	Relabelings []*Relabeling     `json:"relabelings,omitempty"`
}

type Relabeling struct {
	Action       string   `json:"action,omitempty"`
	Regex        string   `json:"regex,omitempty"`
	Replacement  string   `json:"replacement,omitempty"`
	SourceLabels []string `json:"sourcelabels,omitempty"`
	TargetLabel  string   `json:"targetlabel,omitempty"`
}

type Prometheus struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	RemoteWrite []*RemoteWrite    `json:"remotewrite,omitempty"`
}

type RemoteWrite struct {
	URL string `json:"url,omitempty"`
}

type LoggingOperator struct {
	Prefix    string `json:"prefix,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type PrometheusNodeExporter struct {
	Prefix    string `json:"prefix,omitempty"`
	Namespace string `json:"namespace,omitempty"`
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
	metav1.ObjectMeta `json:"metadata,omitempty"`

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
