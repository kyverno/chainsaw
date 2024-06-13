package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ActionCheck contains check for an action.
type ActionCheck struct {
	// Check is an assertion tree to validate the operation outcome.
	// +optional
	Check *Check `json:"check,omitempty"`
}

// ActionCheckRef contains check reference options for an action.
type ActionCheckRef struct {
	FileRef `json:",inline"`

	// Check provides a check used in assertions.
	// +optional
	Check *Check `json:"resource,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`
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
	// +kubebuilder:validation:XEmbeddedResource
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
type ActionObjectSelector struct {
	ObjectName `json:",inline"`

	// Selector defines labels selector.
	// +optional
	Selector string `json:"selector,omitempty"`
}

// FileRef represents a file reference.
type FileRef struct {
	// File is the path to the referenced file. This can be a direct path to a file
	// or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML
	// files within the "manifest" directory.
	File string `json:"file,omitempty"`
}

// ActionResourceRef contains resource reference options for an action.
type ActionResourceRef struct {
	FileRef `json:",inline"`
	// Resource provides a resource to be applied.
	// +kubebuilder:validation:XEmbeddedResource
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
// +k8s:conversion-gen=false
type Apply struct {
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// Assert represents a test condition that is expected to hold true
// during the testing process.
// +k8s:conversion-gen=false
type Assert struct {
	ActionCheckRef `json:",inline"`
	ActionTimeout  `json:",inline"`
}

// Command describes a command to run as a part of a test step.
// +k8s:conversion-gen=false
type Command struct {
	ActionCheck   `json:",inline"`
	ActionEnv     `json:",inline"`
	ActionTimeout `json:",inline"`

	// Entrypoint is the command entry point to run.
	Entrypoint string `json:"entrypoint"`

	// Args is the command arguments.
	// +optional
	Args []string `json:"args,omitempty"`
}

// Create represents a set of resources that should be created.
// If a resource already exists in the cluster it will fail.
// +k8s:conversion-gen=false
type Create struct {
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// Delete is a reference to an object that should be deleted
// +k8s:conversion-gen=false
type Delete struct {
	ActionExpectations `json:",inline"`
	ActionTimeout      `json:",inline"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// File is the path to the referenced file. This can be a direct path to a file
	// or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML
	// files within the "manifest" directory.
	// +optional
	File string `json:"file,omitempty"`

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
// +k8s:conversion-gen=false
type Describe struct {
	ActionObject  `json:",inline"`
	ActionTimeout `json:",inline"`

	// Show Events indicates whether to include related events.
	// +optional
	ShowEvents *bool `json:"showEvents,omitempty"`
}

// Error represents an anticipated error condition that may arise during testing.
// Instead of treating such an error as a test failure, it acknowledges it as expected.
// +k8s:conversion-gen=false
type Error struct {
	ActionCheckRef `json:",inline"`
	ActionTimeout  `json:",inline"`
}

// Events defines how to collect events.
// +k8s:conversion-gen=false
type Events struct {
	ActionFormat         `json:",inline"`
	ActionObjectSelector `json:",inline"`
	ActionTimeout        `json:",inline"`
}

// Get defines how to get resources.
// +k8s:conversion-gen=false
type Get struct {
	ActionFormat  `json:",inline"`
	ActionObject  `json:",inline"`
	ActionTimeout `json:",inline"`
}

// Patch represents a set of resources that should be patched.
// If a resource doesn't exist yet in the cluster it will fail.
// +k8s:conversion-gen=false
type Patch struct {
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// PodLogs defines how to collect pod logs.
// +k8s:conversion-gen=false
type PodLogs struct {
	ActionObjectSelector `json:",inline"`
	ActionTimeout        `json:",inline"`

	// Container in pod to get logs from else --all-containers is used.
	// +optional
	Container string `json:"container,omitempty"`

	// Tail is the number of last lines to collect from pods. If omitted or zero,
	// then the default is 10 if you use a selector, or -1 (all) if you use a pod name.
	// This matches default behavior of `kubectl logs`.
	// +optional
	Tail *int `json:"tail,omitempty"`
}

// Script describes a script to run as a part of a test step.
// +k8s:conversion-gen=false
type Script struct {
	ActionCheck   `json:",inline"`
	ActionEnv     `json:",inline"`
	ActionTimeout `json:",inline"`

	// Content defines a shell script (run with "sh -c ...").
	// +optional
	Content string `json:"content,omitempty"`
}

// Sleep represents a duration while nothing happens.
// +k8s:conversion-gen=false
type Sleep struct {
	// Duration is the delay used for sleeping.
	Duration metav1.Duration `json:"duration"`
}

// Update represents a set of resources that should be updated.
// If a resource does not exist in the cluster it will fail.
// +k8s:conversion-gen=false
type Update struct {
	ActionDryRun       `json:",inline"`
	ActionExpectations `json:",inline"`
	ActionResourceRef  `json:",inline"`
	ActionTimeout      `json:",inline"`
}

// Wait specifies how to perform wait operations on resources.
// +k8s:conversion-gen=false
type Wait struct {
	ActionTimeout `json:",inline"`
	ActionFormat  `json:",inline"`
	ActionObject  `json:",inline"`

	// WaitFor specifies the condition to wait for.
	WaitFor `json:"for"`
}

// WaitFor specifies the condition to wait for.
// +k8s:conversion-gen=false
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
// +k8s:conversion-gen=false
type WaitForCondition struct {
	// Name defines the specific condition to wait for, e.g., "Available", "Ready".
	Name string `json:"name"`

	// Value defines the specific condition status to wait for, e.g., "True", "False".
	// +optional
	Value *string `json:"value,omitempty"`
}

// WaitForDeletion represents parameters for waiting on a resource's deletion.
// +k8s:conversion-gen=false
type WaitForDeletion struct{}

// WaitForJsonPath represents parameters for waiting on a json path of a resource.
// +k8s:conversion-gen=false
type WaitForJsonPath struct {
	// Path defines the json path to wait for, e.g. '{.status.phase}'.
	Path string `json:"path"`

	// Value defines the expected value to wait for, e.g., "Running".
	Value string `json:"value"`
}
