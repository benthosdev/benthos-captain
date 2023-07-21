package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PipelineFinalizer is the finalizer used by the pipeline controller to cleanup resources when the pipeline is being deleted.
	PipelineFinalizer = "pipeline.streaming.benthos.dev"
)

// BenthosPipelineSpec defines the desired state of BenthosPipeline
type BenthosPipelineSpec struct {
	// Currently this is set just to take a string. Ideally, we should be able to fetch the struct
	// as a package but currently `config` is internal. If we decide to fetch from a package, we will need to consider
	// Kubernetes API Versioning when it changes.

	// Config defines the Benthos configuration as a string.
	Config string `json:"config,omitempty"`

	// Replicas defines the amount of replicas to create for the Benthos deployment.
	Replicas int32 `json:"replicas,omitempty"`
}

// BenthosPipelineStatus defines the observed state of BenthosPipeline
type BenthosPipelineStatus struct {
	Ready bool   `json:"ready,omitempty"`
	Phase string `json:"phase,omitempty"`
	// AvailableReplicas is the amount of pods available from the deployment.
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=benthospipelines,scope=Namespaced
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="The current state the Benthos Pipeline."
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The current phase of the Benthos Pipeline."
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The desired amount of running Benthos replicas."
// +kubebuilder:printcolumn:name="Available",type="integer",JSONPath=".status.availableReplicas",description="The amount of available Benthos replicas."
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this resource"

// BenthosPipeline is the Schema for the benthospipelines API
type BenthosPipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BenthosPipelineSpec   `json:"spec,omitempty"`
	Status BenthosPipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BenthosPipelineList contains a list of BenthosPipeline
type BenthosPipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BenthosPipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BenthosPipeline{}, &BenthosPipelineList{})
}

// Currently these have been pulled directly from the Benthos repo. Ideally, we should be able to fetch these
// as a package but currently `config` is internal. If we decide to fetch from a package, we will need to consider
// Kubernetes API Versioning when it changes.
