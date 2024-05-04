package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PipelineFinalizer is the finalizer used by the pipeline controller to cleanup resources when the pipeline is being deleted.
	PipelineFinalizer = "pipeline.captain.benthos.dev"
)

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// Currently this is set just to take a string. Ideally, we should be able to fetch the struct
	// as a package but currently `config` is internal. If we decide to fetch from a package, we will need to consider
	// Kubernetes API Versioning when it changes.

	// Config defines the Benthos configuration as a string.
	Config *apiextensionsv1.JSON `json:"config,omitempty"`

	// Replicas defines the amount of replicas to create for the Benthos deployment.
	Replicas int32 `json:"replicas,omitempty"`

	// Image defines the image and tag to use for the Benthos deployment.
	// +optional
	Image string `json:"image,omitempty"`

	// ConfigFiles Additional configuration, as Key/Value pairs, that will be mounted as files with the /config
	// directory on the pod. The key should be the file name and the value should be its content.
	ConfigFiles map[string]string `json:"configFiles,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	Ready bool   `json:"ready,omitempty"`
	Phase string `json:"phase,omitempty"`
	// AvailableReplicas is the amount of pods available from the deployment.
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pipelines,scope=Namespaced
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="The current state the Benthos Pipeline."
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The current phase of the Benthos Pipeline."
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The desired amount of running Benthos replicas."
// +kubebuilder:printcolumn:name="Available",type="integer",JSONPath=".status.availableReplicas",description="The amount of available Benthos replicas."
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this resource"

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

// Currently these have been pulled directly from the Benthos repo. Ideally, we should be able to fetch these
// as a package but currently `config` is internal. If we decide to fetch from a package, we will need to consider
// Kubernetes API Versioning when it changes.
