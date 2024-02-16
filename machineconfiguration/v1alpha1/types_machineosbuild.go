package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineOSBuild describes a build process managed by the MCO
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type MachineOSBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the configuration of the machine os build
	// +kubebuilder:validation:Required
	Spec MachineOSBuildSpec `json:"spec"`

	// status describes the lst observed state of this machine os build
	// +optional
	Status MachineOSBuildStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineOSBuildList describes all of the Builds on the system
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type MachineOSBuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MachineOSBuild `json:"items"`
}

// MachineOSBuildSpec describes user-configurable options as well as information about a build process.
type MachineOSBuildSpec struct {
	// machineConfigPool is the pool which the build is for
	// +kubebuilder:validation:Required
	MachineConfigPool MachineConfigPoolReference `json:"machineConfigPool"`
	// buildInputs is where user options for the build live
	// +kubebuilder:validation:Required
	BuildInputs BuildInputs `json:"buildInputs"`
	// currentConfig is the currently running config on the MCP
	// +optional
	CurrentConfig string `json:"currentConfig"`
	// desiredConfig is the desired config we want to build an image for. If currentConfig and desiredConfig are not the same, we need to build an image.
	// +kubebuilder:validation:Required
	DesiredConfig string `json:"desiredConfig"`
}

// MachineOSBuildStatus describes the state of a build and other helpful information.
type MachineOSBuildStatus struct {
	// conditions are state related conditions for the build. Valid types are:
	// BuildPrepared, Building, BuildFailed, BuildInterrupted, BuildRestarted, and Ready
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// renderedMachineOSImage describes the machineOsImage object created to track image specific information
	// machineConfig is a reference to the MC that was used to build this image.
	// +optional
	RenderedMachineOSImage MachineOSImageReference `json:"renderedMachineOSImage"`
	// startTime describes when this build began
	// +kubebuilder:validation:Required
	StartTime *metav1.Time `json:"startTime"`
	// endTime describes when the build ended
	// +optional
	EndTime *metav1.Time `json:"endTime"`
	// buildHistory contains previous iterations of failed or interrupted builds related to this one.
	// This will keep track of failed or interrupted build names
	// I am thinking we keep track of their name and why they failed
	// +listType=set
	// +optional
	BuildHistory []PriorMachineOSBuilds `json:"buildHistory,omitempty"`
}

// BuildProgess highlights some of the key phases of a build to be tracked in Conditions.
type BuildProgress string

const (
	// prepared indicates that the build has finished preparing. A build is prepared
	// by gathering the build inputs, validating them, and making sure we can do an update as specified.
	MachineOSBuildPrepared BuildProgress = "Prepared"
	// building indicates that the build has been kicked off with the specified image builder
	MachineOSBuiling BuildProgress = "Building"
	// failed indicates that during the build or preparation process, the build failed.
	MachineOSBuildFailed BuildProgress = "Failed"
	// interrupted indicates that the user stopped the build process by modifying part of the build config
	MachineOSBuildInterrupted BuildProgress = "Interrupted"
	// restarted indicates that this build has been either interrupted or failed, and now a new
	// iteration of the process has begun.
	MachineOSBuildRestarted BuildProgress = "Restarted"
	// ready indicates that the build has completed and the image is ready to roll out.
	MachineOSReady BuildProgress = "Ready"
)

// PriorMachineOSBuilds contains information about related builds
type PriorMachineOSBuilds struct {
	// name is the name of the build
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// buildFailure contains an optional message of why this build ended prematurely.
	// +optional
	BuildFailure string `json:"buildFailure"`
}

// Refers to the name of a MachineConfigPool (e.g., "worker", "infra", etc.):
type MachineConfigPoolReference struct {
	// name is the name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// Refers to the name of a rendered MachineConfig (e.g., "rendered-worker-ec40d2965ff81bce7cd7a7e82a680739", etc.):
type RenderedMachineConfigReference struct {
	// name is the name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// Refers to the name of a (future) MachineOSImage (e.g., "worker-os-image-167651b10ec98af17971d6a47df9e22f", etc.):
type MachineOSImageReference struct {
	// name is the name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// Refers to the name of an image registry push/pull secret needed in the build process.
type ImageSecretObjectReference struct {
	// name is the name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// BuildInputs holds all of the information needed to trigger a build
type BuildInputs struct {
	// imageBuilderType specifies the backend to be used to build the image.
	// Valid options are: OpenShiftImageBuilder, PodImageBuilder, and Default (OpenShiftImageBuilder)
	ImageBuilderType MachineOSImageBuilderType `json:"imageBuilderType"`
	// baseOSImageURL is the base OSImage we use to build our custom image.
	// +kubebuilder:validation:Required
	BaseOSImageURL string `json:"baseOSImageURL"`
	// baseImagePullSecret is the secret used to pull the base image.
	// +kubebuilder :validation:Required
	BaseImagePullSecret string `json:"baseImagePullSecret"`
	// finalImagePushSecret is the secret used to connect to a user registry.
	// +kubebuilder:validation:Required
	FinalImagePushSecret ImageSecretObjectReference `json:"finalImagePushSecret"`
	// finalImagePullSecret is the secret used to pull the final produced image.
	// +kubebuilder:validation:Required
	FinalImagePullSecret ImageSecretObjectReference `json:"finalImagePullSecret"`
	// containerFile describes the custom data the user has specified to build into the image.
	// +optional
	Containerfile []byte `json:"containerFile"`
	// finalImagePullspec describes the location of the final image.
	// +kubebuilder:validation:Pattern=`^https://`
	// +kubebuilder:validation:Required
	FinalImagePullspec string `json:"finalImagePullspec"`
}

type MachineOSImageBuilderType string

const (
	OCPBuilder     MachineOSImageBuilderType = "OpenShiftImageBuilder"
	PodBuilder     MachineOSImageBuilderType = "PodImageBuilder"
	DefaultBuilder MachineOSImageBuilderType = "Default"
)
