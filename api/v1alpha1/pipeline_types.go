/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type HTTPConfig struct {
}

type InputConfig struct {
	Name string `json:"name,omitempty"`
}

type BufferConfig struct {
}

type PipelineConfig struct {
}

type OutputConfig struct {
}

type ManagerConfig struct {
}

type LoggerConfig struct {
}

type MetricsConfig struct {
}

type TracerConfig struct {
}

type SystemCloseTimeoutConfig struct {
}

type TestsConfig struct {
}

type Config struct {
	HTTP               HTTPConfig               `json:"http" yaml:"http"`
	Input              InputConfig              `json:"input" yaml:"input"`
	Buffer             BufferConfig             `json:"buffer" yaml:"buffer"`
	Pipeline           PipelineConfig           `json:"pipeline" yaml:"pipeline"`
	Output             OutputConfig             `json:"output" yaml:"output"`
	Manager            ManagerConfig            `json:"resources" yaml:"resources"`
	Logger             LoggerConfig             `json:"logger" yaml:"logger"`
	Metrics            MetricsConfig            `json:"metrics" yaml:"metrics"`
	Tracer             TracerConfig             `json:"tracer" yaml:"tracer"`
	SystemCloseTimeout SystemCloseTimeoutConfig `json:"shutdown_timeout" yaml:"shutdown_timeout"`
	Tests              TestsConfig              `json:"tests,omitempty" yaml:"tests,omitempty"`
}

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Workers defines the number of workers
	Workers int `json:"workers,omitempty"`

	// Config defines the pipeline config
	Config Config `json:"config,omitempty"`
}

type State string

const (
	Degraded State = "Degraded"
	Running        = "Running"
	Paused         = "Paused"
	Failed         = "Failed"
)

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// State
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pipeline is the Schema for the pipelines API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
