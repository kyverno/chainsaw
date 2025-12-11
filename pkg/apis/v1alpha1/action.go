package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ActionBindings contains bindings options for an action.
type ActionBindings struct {
	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`
}

// ActionCheck contains check for an action.
type ActionCheck struct {
	// Check is an assertion tree to validate the operation outcome.
	// +optional
	Check *Check `json:"check,omitempty"`
}

// ActionCheckRef contains check reference options for an action.
// +kubebuilder:not:={required:{file,resource}}
type ActionCheckRef struct {
	FileRef `json:",inline"`

	// Check provides a check used in assertions.
	// +optional
	Check *Projection `json:"resource,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`
}

// ActionClusters contains clusters options for an action.
type ActionClusters struct {
	// Cluster defines the target cluster (will be inherited if not specified).
	// +optional
	Cluster *string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`
}

// ActionDryRun contains dry run options for an action.
type ActionDryRun struct {
	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`
}

// ActionEnv contains environment options for an action.
type ActionEnv struct {
	// Env defines additional environment variables.
	// +optional
	Env []Binding `json:"env,omitempty"`

	// SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`

	// SkipCommandOutput removes the command from the output logs.
	// +optional
	SkipCommandOutput bool `json:"skipCommandOutput,omitempty"`
}

// ActionExpectations contains expectations for an action.
type ActionExpectations struct {
	// Expect defines a list of matched checks to validate the operation outcome.
	// +optional
	Expect []Expectation `json:"expect,omitempty"`
}

// ActionFormat contains format for an action.
type ActionFormat struct {
	// Format determines the output format (json or yaml).
	// +optional
	Format Format `json:"format,omitempty"`
}

type ActionInlineResource struct {
	// Resource provides a resource to be applied.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Resource *unstructured.Unstructured `json:"resource,omitempty"`
}

// ActionObject contains object selector options for an action.
type ActionObject struct {
	ObjectType           `json:",inline"`
	ActionObjectSelector `json:",inline"`
}

// ActionObjectSelector contains object selector options for an action.
// +kubebuilder:not:={required:{name,selector}}
type ActionObjectSelector struct {
	ObjectName `json:",inline"`

	// Selector defines labels selector.
	// +optional
	Selector Expression `json:"selector,omitempty"`
}

// ActionOutputs contains outputs options for an action.
type ActionOutputs struct {
	// Outputs defines output bindings.
	// +optional
	Outputs []Output `json:"outputs,omitempty"`
}

// FileRef represents a file reference.
type FileRef struct {
	// File is the path to the referenced file. This can be a direct path to a file
	// or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML
	// files within the "manifest" directory.
	File Expression `json:"file,omitempty"`
}

// ActionResourceRef contains resource reference options for an action.
// +kubebuilder:not:={required:{file,resource}}
type ActionResourceRef struct {
	FileRef `json:",inline"`
	// Resource provides a resource to be applied.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Resource *unstructured.Unstructured `json:"resource,omitempty"`
	// ActionInlineResource `json:",inline"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`
}

// ActionTimeout contains timeout options for an action.
type ActionTimeout struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	ActionBindings     `json:",inline"`
	ActionClusters     `json:",inline"`
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionOutputs      `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// Assert represents a test condition that is expected to hold true
// during the testing process.
type Assert struct {
	ActionBindings `json:",inline"`
	ActionCheckRef `json:",inline"`
	ActionClusters `json:",inline"`
	ActionTimeout  `json:",inline"`
}

// Command describes a command to run as a part of a test step.
type Command struct {
	ActionBindings `json:",inline"`
	ActionCheck    `json:",inline"`
	ActionClusters `json:",inline"`
	ActionEnv      `json:",inline"`
	ActionOutputs  `json:",inline"`
	ActionTimeout  `json:",inline"`

	// Entrypoint is the command entry point to run.
	Entrypoint string `json:"entrypoint"`

	// Args is the command arguments.
	// +optional
	Args []string `json:"args,omitempty"`

	// WorkDir is the working directory for command.
	// +optional
	WorkDir *string `json:"workDir,omitempty"`
}

// Create represents a set of resources that should be created.
// If a resource already exists in the cluster it will fail.
type Create struct {
	ActionBindings     `json:",inline"`
	ActionClusters     `json:",inline"`
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionOutputs      `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// Delete is a reference to an object that should be deleted
// +kubebuilder:not:={required:{file,ref}}
type Delete struct {
	ActionBindings     `json:",inline"`
	ActionClusters     `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionTimeout      `json:",inline"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// File is the path to the referenced file. This can be a direct path to a file
	// or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML
	// files within the "manifest" directory.
	// +optional
	File Expression `json:"file,omitempty"`

	// Ref determines objects to be deleted.
	// +optional
	Ref *ObjectReference `json:"ref,omitempty"`

	// DeletionPropagationPolicy decides if a deletion will propagate to the dependents of
	// the object, and how the garbage collector will handle the propagation.
	// Overrides the deletion propagation policy set in the Configuration, the Test and the TestStep.
	// +optional
	// +kubebuilder:validation:Enum:=Orphan;Background;Foreground
	DeletionPropagationPolicy *metav1.DeletionPropagation `json:"deletionPropagationPolicy,omitempty"`
}

// Describe defines how to describe resources.
type Describe struct {
	ActionClusters `json:",inline"`
	ActionObject   `json:",inline"`
	ActionTimeout  `json:",inline"`

	// Show Events indicates whether to include related events.
	// +optional
	ShowEvents *bool `json:"showEvents,omitempty"`
}

// Error represents an anticipated error condition that may arise during testing.
// Instead of treating such an error as a test failure, it acknowledges it as expected.
type Error struct {
	ActionBindings `json:",inline"`
	ActionCheckRef `json:",inline"`
	ActionClusters `json:",inline"`
	ActionTimeout  `json:",inline"`
}

// Events defines how to collect events.
type Events struct {
	ActionClusters       `json:",inline"`
	ActionFormat         `json:",inline"`
	ActionObjectSelector `json:",inline"`
	ActionTimeout        `json:",inline"`
}

// Get defines how to get resources.
type Get struct {
	ActionClusters `json:",inline"`
	ActionFormat   `json:",inline"`
	ActionObject   `json:",inline"`
	ActionTimeout  `json:",inline"`
}

// Patch represents a set of resources that should be patched.
// If a resource doesn't exist yet in the cluster it will fail.
type Patch struct {
	ActionBindings     `json:",inline"`
	ActionClusters     `json:",inline"`
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionOutputs      `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// PodLogs defines how to collect pod logs.
type PodLogs struct {
	ActionClusters       `json:",inline"`
	ActionObjectSelector `json:",inline"`
	ActionTimeout        `json:",inline"`

	// Container in pod to get logs from else --all-containers is used.
	// +optional
	Container Expression `json:"container,omitempty"`

	// Tail is the number of last lines to collect from pods. If omitted or zero,
	// then the default is 10 if you use a selector, or -1 (all) if you use a pod name.
	// This matches default behavior of `kubectl logs`.
	// +optional
	Tail *int `json:"tail,omitempty"`
}

// Proxy defines how to get resources.
type Proxy struct {
	ActionClusters `json:",inline"`
	ActionOutputs  `json:",inline"`
	ActionTimeout  `json:",inline"`
	ObjectName     `json:",inline"`
	ObjectType     `json:",inline"`

	// TargetPort defines the target port to proxy the request.
	// +optional
	TargetPort Expression `json:"port,omitempty"`

	// TargetPath defines the target path to proxy the request.
	// +optional
	TargetPath Expression `json:"path,omitempty"`
}

// Script describes a script to run as a part of a test step.
type Script struct {
	ActionBindings `json:",inline"`
	ActionCheck    `json:",inline"`
	ActionClusters `json:",inline"`
	ActionEnv      `json:",inline"`
	ActionOutputs  `json:",inline"`
	ActionTimeout  `json:",inline"`

	// Content defines a shell script (run with "$shell $shellArgs ...").
	// +optional
	Content string `json:"content,omitempty"`

	// Shell defines the host shell (run with "... $shellArgs $content").
	Shell *string `json:"shell,omitempty"`

	// ShellArgs defines arguments for the host shell (run with "$shell ... $content").
	ShellArgs []string `json:"shellArgs,omitempty"`

	// WorkDir is the working directory for script.
	// +optional
	WorkDir *string `json:"workDir,omitempty"`
}

// Sleep represents a duration while nothing happens.
type Sleep struct {
	// Duration is the delay used for sleeping.
	Duration metav1.Duration `json:"duration"`
}

// Update represents a set of resources that should be updated.
// If a resource does not exist in the cluster it will fail.
type Update struct {
	ActionBindings     `json:",inline"`
	ActionClusters     `json:",inline"`
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionOutputs      `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// Wait specifies how to perform wait operations on resources.
type Wait struct {
	ActionTimeout  `json:",inline"`
	ActionFormat   `json:",inline"`
	ActionClusters `json:",inline"`
	ActionObject   `json:",inline"`

	// WaitFor specifies the condition to wait for.
	WaitFor `json:"for"`
}

// WaitFor specifies the condition to wait for.
type WaitFor struct {
	// Deletion specifies parameters for waiting on a resource's deletion.
	// +optional
	Deletion *WaitForDeletion `json:"deletion,omitempty"`

	// Condition specifies the condition to wait for.
	// +optional
	Condition *WaitForCondition `json:"condition,omitempty"`

	// JsonPath specifies the json path condition to wait for.
	// +optional
	JsonPath *WaitForJsonPath `json:"jsonPath,omitempty"`
}

// WaitForCondition represents parameters for waiting on a specific condition of a resource.
type WaitForCondition struct {
	// Name defines the specific condition to wait for, e.g., "Available", "Ready".
	Name Expression `json:"name"`

	// Value defines the specific condition status to wait for, e.g., "True", "False".
	// +optional
	Value *Expression `json:"value,omitempty"`
}

// WaitForDeletion represents parameters for waiting on a resource's deletion.
type WaitForDeletion struct{}

// WaitForJsonPath represents parameters for waiting on a json path of a resource.
type WaitForJsonPath struct {
	// Path defines the json path to wait for, e.g. '{.status.phase}'.
	Path Expression `json:"path"`

	// Value defines the expected value to wait for, e.g., "Running".
	// +optional
	Value *Expression `json:"value,omitempty"`
}
