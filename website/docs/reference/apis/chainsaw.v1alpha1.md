---
title: chainsaw (v1alpha1)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the v1alpha1 API group.</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)
- [StepTemplate](#chainsaw-kyverno-io-v1alpha1-StepTemplate)
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)
  
## Configuration     {#chainsaw-kyverno-io-v1alpha1-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec) | :white_check_mark: |  | <p>Configuration spec.</p> |

## StepTemplate     {#chainsaw-kyverno-io-v1alpha1-StepTemplate}

<p>StepTemplate is the resource that contains a step definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `StepTemplate` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`StepTemplateSpec`](#chainsaw-kyverno-io-v1alpha1-StepTemplateSpec) | :white_check_mark: |  | <p>Test step spec.</p> |

## Test     {#chainsaw-kyverno-io-v1alpha1-Test}

<p>Test is the resource that contains a test definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `Test` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestSpec`](#chainsaw-kyverno-io-v1alpha1-TestSpec) | :white_check_mark: |  | <p>Test spec.</p> |

## ActionBindings     {#chainsaw-kyverno-io-v1alpha1-ActionBindings}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>ActionBindings contains bindings options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |

## ActionCheck     {#chainsaw-kyverno-io-v1alpha1-ActionCheck}

**Appears in:**
    
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)

<p>ActionCheck contains check for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `check` | `policy/v1alpha1.AssertionTree` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## ActionCheckRef     {#chainsaw-kyverno-io-v1alpha1-ActionCheckRef}

**Appears in:**
    
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)

<p>ActionCheckRef contains check reference options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `resource` | [`Projection`](#chainsaw-kyverno-io-v1alpha1-Projection) |  |  | <p>Check provides a check used in assertions.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## ActionClusters     {#chainsaw-kyverno-io-v1alpha1-ActionClusters}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [PodLogs](#chainsaw-kyverno-io-v1alpha1-PodLogs)
- [Proxy](#chainsaw-kyverno-io-v1alpha1-Proxy)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ActionClusters contains clusters options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (will be inherited if not specified).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |

## ActionDryRun     {#chainsaw-kyverno-io-v1alpha1-ActionDryRun}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>ActionDryRun contains dry run options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |

## ActionEnv     {#chainsaw-kyverno-io-v1alpha1-ActionEnv}

**Appears in:**
    
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)

<p>ActionEnv contains environment options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |

## ActionExpectations     {#chainsaw-kyverno-io-v1alpha1-ActionExpectations}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>ActionExpectations contains expectations for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## ActionFormat     {#chainsaw-kyverno-io-v1alpha1-ActionFormat}

**Appears in:**
    
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ActionFormat contains format for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## ActionObject     {#chainsaw-kyverno-io-v1alpha1-ActionObject}

**Appears in:**
    
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ActionObject contains object selector options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha1-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## ActionObjectSelector     {#chainsaw-kyverno-io-v1alpha1-ActionObjectSelector}

**Appears in:**
    
- [ActionObject](#chainsaw-kyverno-io-v1alpha1-ActionObject)
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [PodLogs](#chainsaw-kyverno-io-v1alpha1-PodLogs)

<p>ActionObjectSelector contains object selector options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectName` | [`ObjectName`](#chainsaw-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `selector` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>Selector defines labels selector.</p> |

## ActionOutputs     {#chainsaw-kyverno-io-v1alpha1-ActionOutputs}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Proxy](#chainsaw-kyverno-io-v1alpha1-Proxy)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>ActionOutputs contains outputs options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |

## ActionResourceRef     {#chainsaw-kyverno-io-v1alpha1-ActionResourceRef}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>ActionResourceRef contains resource reference options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `resource` | [`meta/v1/unstructured.Unstructured`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#unstructured-unstructured-v1) |  |  | <p>Resource provides a resource to be applied.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## ActionTimeout     {#chainsaw-kyverno-io-v1alpha1-ActionTimeout}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [PodLogs](#chainsaw-kyverno-io-v1alpha1-PodLogs)
- [Proxy](#chainsaw-kyverno-io-v1alpha1-Proxy)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ActionTimeout contains timeout options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |

## Apply     {#chainsaw-kyverno-io-v1alpha1-Apply}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Apply represents a set of configurations or resources that
should be applied during testing.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha1-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha1-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha1-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Assert     {#chainsaw-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Assert represents a test condition that is expected to hold true
during the testing process.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionCheckRef` | [`ActionCheckRef`](#chainsaw-kyverno-io-v1alpha1-ActionCheckRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Binding     {#chainsaw-kyverno-io-v1alpha1-Binding}

**Appears in:**
    
- [ActionBindings](#chainsaw-kyverno-io-v1alpha1-ActionBindings)
- [ActionEnv](#chainsaw-kyverno-io-v1alpha1-ActionEnv)
- [Output](#chainsaw-kyverno-io-v1alpha1-Output)
- [Scenario](#chainsaw-kyverno-io-v1alpha1-Scenario)
- [StepTemplateSpec](#chainsaw-kyverno-io-v1alpha1-StepTemplateSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)
- [With](#chainsaw-kyverno-io-v1alpha1-With)

<p>Binding represents a key/value set as a binding in an executing test.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) | :white_check_mark: |  | <p>Name the name of the binding.</p> |
| `compiler` | `policy/v1alpha1.Compiler` |  |  | <p>Compiler defines the default compiler to use when evaluating expressions.</p> |
| `value` | [`Projection`](#chainsaw-kyverno-io-v1alpha1-Projection) | :white_check_mark: |  | <p>Value value of the binding.</p> |

## CatchFinally     {#chainsaw-kyverno-io-v1alpha1-CatchFinally}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)
- [StepTemplateSpec](#chainsaw-kyverno-io-v1alpha1-StepTemplateSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>CatchFinally defines actions to be executed in catch, finally and cleanup blocks.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the operation.</p> |
| `podLogs` | [`PodLogs`](#chainsaw-kyverno-io-v1alpha1-PodLogs) |  |  | <p>PodLogs determines the pod logs collector to execute.</p> |
| `events` | [`Events`](#chainsaw-kyverno-io-v1alpha1-Events) |  |  | <p>Events determines the events collector to execute.</p> |
| `describe` | [`Describe`](#chainsaw-kyverno-io-v1alpha1-Describe) |  |  | <p>Describe determines the resource describe collector to execute.</p> |
| `wait` | [`Wait`](#chainsaw-kyverno-io-v1alpha1-Wait) |  |  | <p>Wait determines the resource wait collector to execute.</p> |
| `get` | [`Get`](#chainsaw-kyverno-io-v1alpha1-Get) |  |  | <p>Get determines the resource get collector to execute.</p> |
| `delete` | [`Delete`](#chainsaw-kyverno-io-v1alpha1-Delete) |  |  | <p>Delete represents a deletion operation.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |
| `sleep` | [`Sleep`](#chainsaw-kyverno-io-v1alpha1-Sleep) |  |  | <p>Sleep defines zzzz.</p> |

## Clusters     {#chainsaw-kyverno-io-v1alpha1-Clusters}

(Alias of `map[string]github.com/kyverno/chainsaw/pkg/apis/v1alpha1.Cluster`)

**Appears in:**
    
- [ActionClusters](#chainsaw-kyverno-io-v1alpha1-ActionClusters)
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Clusters defines a cluster map.</p>


## Command     {#chainsaw-kyverno-io-v1alpha1-Command}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Command describes a command to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionCheck` | [`ActionCheck`](#chainsaw-kyverno-io-v1alpha1-ActionCheck) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionEnv` | [`ActionEnv`](#chainsaw-kyverno-io-v1alpha1-ActionEnv) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `entrypoint` | `string` | :white_check_mark: |  | <p>Entrypoint is the command entry point to run.</p> |
| `args` | `[]string` |  |  | <p>Args is the command arguments.</p> |
| `workDir` | `string` |  |  | <p>WorkDir is the working directory for command.</p> |

## ConfigurationSpec     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeouts` | [`DefaultTimeouts`](#chainsaw-kyverno-io-v1alpha1-DefaultTimeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running the tests (implies SkipClusterDelete).</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `compiler` | `policy/v1alpha1.Compiler` |  |  | <p>Compiler defines the default compiler to use when evaluating expressions.</p> |
| `failFast` | `bool` |  |  | <p>FailFast determines whether the test should stop upon encountering the first failure.</p> |
| `parallel` | `int` |  |  | <p>The maximum number of tests to run at once.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation.</p> |
| `reportFormat` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha1-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML|JUNIT-TEST|JUNIT-STEP|JUNIT-OPERATION|nil) nil == no report. maps to report.Type, however we don't want generated.deepcopy to have reference to it.</p> |
| `reportPath` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `reportName` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |
| `namespace` | `string` |  |  | <p>Namespace defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `namespaceTemplateCompiler` | `policy/v1alpha1.Compiler` |  |  | <p>NamespaceTemplateCompiler defines the default compiler to use when evaluating expressions.</p> |
| `namespaceTemplate` | [`Projection`](#chainsaw-kyverno-io-v1alpha1-Projection) |  |  | <p>NamespaceTemplate defines a template to create the test namespace.</p> |
| `fullName` | `bool` |  |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |
| `excludeTestRegex` | `string` |  |  | <p>ExcludeTestRegex is used to exclude tests based on a regular expression.</p> |
| `includeTestRegex` | `string` |  |  | <p>IncludeTestRegex is used to include tests based on a regular expression.</p> |
| `repeatCount` | `int` |  |  | <p>RepeatCount indicates how many times the tests should be executed.</p> |
| `testFile` | `string` |  |  | <p>TestFile is the name of the file containing the test to run. If no extension is provided, chainsaw will try with .yaml first and .yml if needed.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `catch` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Catch defines what the tests steps will execute when an error happens. This will be combined with catch handlers defined at the test and step levels.</p> |

## Create     {#chainsaw-kyverno-io-v1alpha1-Create}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Create represents a set of resources that should be created.
If a resource already exists in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha1-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha1-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha1-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## DefaultTimeouts     {#chainsaw-kyverno-io-v1alpha1-DefaultTimeouts}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

<p>DefaultTimeouts contains defautl timeouts per operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apply` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Apply defines the timeout for the apply operation</p> |
| `assert` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Assert defines the timeout for the assert operation</p> |
| `cleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Cleanup defines the timeout for the cleanup operation</p> |
| `delete` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Delete defines the timeout for the delete operation</p> |
| `error` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Error defines the timeout for the error operation</p> |
| `exec` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Exec defines the timeout for exec operations</p> |

## Delete     {#chainsaw-kyverno-io-v1alpha1-Delete}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Delete is a reference to an object that should be deleted</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha1-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `file` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |
| `ref` | [`ObjectReference`](#chainsaw-kyverno-io-v1alpha1-ObjectReference) |  |  | <p>Ref determines objects to be deleted.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in the Configuration, the Test and the TestStep.</p> |

## Describe     {#chainsaw-kyverno-io-v1alpha1-Describe}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Describe defines how to describe resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObject` | [`ActionObject`](#chainsaw-kyverno-io-v1alpha1-ActionObject) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `showEvents` | `bool` |  |  | <p>Show Events indicates whether to include related events.</p> |

## Error     {#chainsaw-kyverno-io-v1alpha1-Error}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Error represents an anticipated error condition that may arise during testing.
Instead of treating such an error as a test failure, it acknowledges it as expected.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionCheckRef` | [`ActionCheckRef`](#chainsaw-kyverno-io-v1alpha1-ActionCheckRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Events     {#chainsaw-kyverno-io-v1alpha1-Events}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Events defines how to collect events.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionFormat` | [`ActionFormat`](#chainsaw-kyverno-io-v1alpha1-ActionFormat) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha1-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Expectation     {#chainsaw-kyverno-io-v1alpha1-Expectation}

**Appears in:**
    
- [ActionExpectations](#chainsaw-kyverno-io-v1alpha1-ActionExpectations)

<p>Expectation represents a check to be applied on the result of an operation
with a match filter to determine if the verification should be considered.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `match` | `policy/v1alpha1.AssertionTree` |  |  | <p>Match defines the matching statement.</p> |
| `check` | `policy/v1alpha1.AssertionTree` | :white_check_mark: |  | <p>Check defines the verification statement.</p> |

## Expression     {#chainsaw-kyverno-io-v1alpha1-Expression}

(Alias of `string`)

**Appears in:**
    
- [ActionObjectSelector](#chainsaw-kyverno-io-v1alpha1-ActionObjectSelector)
- [Binding](#chainsaw-kyverno-io-v1alpha1-Binding)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [FileRef](#chainsaw-kyverno-io-v1alpha1-FileRef)
- [ObjectName](#chainsaw-kyverno-io-v1alpha1-ObjectName)
- [ObjectType](#chainsaw-kyverno-io-v1alpha1-ObjectType)
- [PodLogs](#chainsaw-kyverno-io-v1alpha1-PodLogs)
- [Proxy](#chainsaw-kyverno-io-v1alpha1-Proxy)
- [WaitForCondition](#chainsaw-kyverno-io-v1alpha1-WaitForCondition)
- [WaitForJsonPath](#chainsaw-kyverno-io-v1alpha1-WaitForJsonPath)

<p>Expression defines an expression to be used in string fields.</p>


## FileRef     {#chainsaw-kyverno-io-v1alpha1-FileRef}

**Appears in:**
    
- [ActionCheckRef](#chainsaw-kyverno-io-v1alpha1-ActionCheckRef)
- [ActionResourceRef](#chainsaw-kyverno-io-v1alpha1-ActionResourceRef)

<p>FileRef represents a file reference.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `file` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) | :white_check_mark: |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |

## Format     {#chainsaw-kyverno-io-v1alpha1-Format}

(Alias of `string`)

**Appears in:**
    
- [ActionFormat](#chainsaw-kyverno-io-v1alpha1-ActionFormat)

<p>Format determines the output format (json or yaml).</p>


## Get     {#chainsaw-kyverno-io-v1alpha1-Get}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Get defines how to get resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionFormat` | [`ActionFormat`](#chainsaw-kyverno-io-v1alpha1-ActionFormat) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObject` | [`ActionObject`](#chainsaw-kyverno-io-v1alpha1-ActionObject) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## ObjectName     {#chainsaw-kyverno-io-v1alpha1-ObjectName}

**Appears in:**
    
- [ActionObjectSelector](#chainsaw-kyverno-io-v1alpha1-ActionObjectSelector)
- [ObjectReference](#chainsaw-kyverno-io-v1alpha1-ObjectReference)
- [Proxy](#chainsaw-kyverno-io-v1alpha1-Proxy)

<p>ObjectName represents an object namespace and name.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `namespace` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |

## ObjectReference     {#chainsaw-kyverno-io-v1alpha1-ObjectReference}

**Appears in:**
    
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)

<p>ObjectReference represents one or more objects with a specific apiVersion and kind.
For a single object name and namespace are used to identify the object.
For multiple objects use labels.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ObjectName` | [`ObjectName`](#chainsaw-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `labels` | `map[string]string` |  |  | <p>Label selector to match objects to delete</p> |

## ObjectType     {#chainsaw-kyverno-io-v1alpha1-ObjectType}

**Appears in:**
    
- [ActionObject](#chainsaw-kyverno-io-v1alpha1-ActionObject)
- [ObjectReference](#chainsaw-kyverno-io-v1alpha1-ObjectReference)
- [Proxy](#chainsaw-kyverno-io-v1alpha1-Proxy)

<p>ObjectType represents a specific apiVersion and kind.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) | :white_check_mark: |  | <p>API version of the referent.</p> |
| `kind` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) | :white_check_mark: |  | <p>Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</p> |

## Operation     {#chainsaw-kyverno-io-v1alpha1-Operation}

**Appears in:**
    
- [StepTemplateSpec](#chainsaw-kyverno-io-v1alpha1-StepTemplateSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Operation defines a single operation, only one action is permitted for a given operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `OperationBase` | [`OperationBase`](#chainsaw-kyverno-io-v1alpha1-OperationBase) |  | :white_check_mark: | <p>OperationBase defines common elements to all operations.</p> |
| `apply` | [`Apply`](#chainsaw-kyverno-io-v1alpha1-Apply) |  |  | <p>Apply represents resources that should be applied for this test step. This can include things like configuration settings or any other resources that need to be available during the test.</p> |
| `assert` | [`Assert`](#chainsaw-kyverno-io-v1alpha1-Assert) |  |  | <p>Assert represents an assertion to be made. It checks whether the conditions specified in the assertion hold true.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `create` | [`Create`](#chainsaw-kyverno-io-v1alpha1-Create) |  |  | <p>Create represents a creation operation.</p> |
| `delete` | [`Delete`](#chainsaw-kyverno-io-v1alpha1-Delete) |  |  | <p>Delete represents a deletion operation.</p> |
| `describe` | [`Describe`](#chainsaw-kyverno-io-v1alpha1-Describe) |  |  | <p>Describe determines the resource describe collector to execute.</p> |
| `error` | [`Error`](#chainsaw-kyverno-io-v1alpha1-Error) |  |  | <p>Error represents the expected errors for this test step. If any of these errors occur, the test will consider them as expected; otherwise, they will be treated as test failures.</p> |
| `events` | [`Events`](#chainsaw-kyverno-io-v1alpha1-Events) |  |  | <p>Events determines the events collector to execute.</p> |
| `get` | [`Get`](#chainsaw-kyverno-io-v1alpha1-Get) |  |  | <p>Get determines the resource get collector to execute.</p> |
| `patch` | [`Patch`](#chainsaw-kyverno-io-v1alpha1-Patch) |  |  | <p>Patch represents a patch operation.</p> |
| `podLogs` | [`PodLogs`](#chainsaw-kyverno-io-v1alpha1-PodLogs) |  |  | <p>PodLogs determines the pod logs collector to execute.</p> |
| `proxy` | [`Proxy`](#chainsaw-kyverno-io-v1alpha1-Proxy) |  |  | <p>Proxy runs a proxy request.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |
| `sleep` | [`Sleep`](#chainsaw-kyverno-io-v1alpha1-Sleep) |  |  | <p>Sleep defines zzzz.</p> |
| `update` | [`Update`](#chainsaw-kyverno-io-v1alpha1-Update) |  |  | <p>Update represents an update operation.</p> |
| `wait` | [`Wait`](#chainsaw-kyverno-io-v1alpha1-Wait) |  |  | <p>Wait determines the resource wait collector to execute.</p> |

## OperationBase     {#chainsaw-kyverno-io-v1alpha1-OperationBase}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>OperationBase defines common elements to all operations.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the operation.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## Output     {#chainsaw-kyverno-io-v1alpha1-Output}

**Appears in:**
    
- [ActionOutputs](#chainsaw-kyverno-io-v1alpha1-ActionOutputs)

<p>Output represents an output binding with a match to determine if the binding must be considered or not.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `Binding` | [`Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) | :white_check_mark: | :white_check_mark: | <p>Binding determines the binding to create when the match succeeds.</p> |
| `match` | `policy/v1alpha1.AssertionTree` |  |  | <p>Match defines the matching statement.</p> |

## Patch     {#chainsaw-kyverno-io-v1alpha1-Patch}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Patch represents a set of resources that should be patched.
If a resource doesn't exist yet in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha1-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha1-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha1-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## PodLogs     {#chainsaw-kyverno-io-v1alpha1-PodLogs}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>PodLogs defines how to collect pod logs.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha1-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `container` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>Container in pod to get logs from else --all-containers is used.</p> |
| `tail` | `int` |  |  | <p>Tail is the number of last lines to collect from pods. If omitted or zero, then the default is 10 if you use a selector, or -1 (all) if you use a pod name. This matches default behavior of `kubectl logs`.</p> |

## Projection     {#chainsaw-kyverno-io-v1alpha1-Projection}

**Appears in:**
    
- [ActionCheckRef](#chainsaw-kyverno-io-v1alpha1-ActionCheckRef)
- [Binding](#chainsaw-kyverno-io-v1alpha1-Binding)
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)

<p>Projection can be any type.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|

## Proxy     {#chainsaw-kyverno-io-v1alpha1-Proxy}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Proxy defines how to get resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ObjectName` | [`ObjectName`](#chainsaw-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `port` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>TargetPort defines the target port to proxy the request.</p> |
| `path` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>TargetPath defines the target path to proxy the request.</p> |

## ReportFormatType     {#chainsaw-kyverno-io-v1alpha1-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

## Scenario     {#chainsaw-kyverno-io-v1alpha1-Scenario}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)

<p>Scenario defines per scenario bindings.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines binding key/values.</p> |

## Script     {#chainsaw-kyverno-io-v1alpha1-Script}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Script describes a script to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionCheck` | [`ActionCheck`](#chainsaw-kyverno-io-v1alpha1-ActionCheck) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionEnv` | [`ActionEnv`](#chainsaw-kyverno-io-v1alpha1-ActionEnv) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `content` | `string` |  |  | <p>Content defines a shell script (run with "sh -c ...").</p> |
| `workDir` | `string` |  |  | <p>WorkDir is the working directory for script.</p> |

## Sleep     {#chainsaw-kyverno-io-v1alpha1-Sleep}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Sleep represents a duration while nothing happens.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `duration` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Duration is the delay used for sleeping.</p> |

## StepTemplateSpec     {#chainsaw-kyverno-io-v1alpha1-StepTemplateSpec}

**Appears in:**
    
- [StepTemplate](#chainsaw-kyverno-io-v1alpha1-StepTemplate)

<p>StepTemplateSpec defines the spec of a step template.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `try` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha1-Operation) | :white_check_mark: |  | <p>Try defines what the step will try to execute.</p> |
| `catch` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Catch defines what the step will execute when an error happens.</p> |
| `finally` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Finally defines what the step will execute after the step is terminated.</p> |
| `cleanup` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Cleanup defines what will be executed after the test is terminated.</p> |

## TestSpec     {#chainsaw-kyverno-io-v1alpha1-TestSpec}

**Appears in:**
    
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)

<p>TestSpec contains the test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the test.</p> |
| `failFast` | `bool` |  |  | <p>FailFast determines whether the test should stop upon encountering the first failure.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (will be inherited if not specified).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `skip` | `bool` |  |  | <p>Skip determines whether the test should skipped.</p> |
| `concurrent` | `bool` |  |  | <p>Concurrent determines whether the test should run concurrently with other tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the test should be deleted after the test is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `compiler` | `policy/v1alpha1.Compiler` |  |  | <p>Compiler defines the default compiler to use when evaluating expressions.</p> |
| `namespace` | `string` |  |  | <p>Namespace determines whether the test should run in a random ephemeral namespace or not.</p> |
| `namespaceTemplate` | [`Projection`](#chainsaw-kyverno-io-v1alpha1-Projection) |  |  | <p>NamespaceTemplate defines a template to create the test namespace.</p> |
| `namespaceTemplateCompiler` | `policy/v1alpha1.Compiler` |  |  | <p>NamespaceTemplateCompiler defines the default compiler to use when evaluating expressions.</p> |
| `scenarios` | [`[]Scenario`](#chainsaw-kyverno-io-v1alpha1-Scenario) |  |  | <p>Scenarios defines test scenarios.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `steps` | [`[]TestStep`](#chainsaw-kyverno-io-v1alpha1-TestStep) | :white_check_mark: |  | <p>Steps defining the test.</p> |
| `catch` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Catch defines what the steps will execute when an error happens. This will be combined with catch handlers defined at the step level.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in the Configuration.</p> |

## TestStep     {#chainsaw-kyverno-io-v1alpha1-TestStep}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)

<p>TestStep contains the test step definition used in a test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name of the step.</p> |
| `use` | [`Use`](#chainsaw-kyverno-io-v1alpha1-Use) | :white_check_mark: |  | <p>Use defines a reference to a step template.</p> |
| `TestStepSpec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: | :white_check_mark: | <p>TestStepSpec of the step.</p> |

## TestStepSpec     {#chainsaw-kyverno-io-v1alpha1-TestStepSpec}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)

<p>TestStepSpec defines the desired state and behavior for each test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the test step.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in both the Configuration and the Test.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (will be inherited if not specified).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `try` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha1-Operation) |  |  | <p>Try defines what the step will try to execute.</p> |
| `catch` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Catch defines what the step will execute when an error happens.</p> |
| `finally` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Finally defines what the step will execute after the step is terminated.</p> |
| `cleanup` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Cleanup defines what will be executed after the test is terminated.</p> |

## Timeouts     {#chainsaw-kyverno-io-v1alpha1-Timeouts}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Timeouts contains timeouts per operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apply` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Apply defines the timeout for the apply operation</p> |
| `assert` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Assert defines the timeout for the assert operation</p> |
| `cleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Cleanup defines the timeout for the cleanup operation</p> |
| `delete` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Delete defines the timeout for the delete operation</p> |
| `error` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Error defines the timeout for the error operation</p> |
| `exec` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Exec defines the timeout for exec operations</p> |

## Update     {#chainsaw-kyverno-io-v1alpha1-Update}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Update represents a set of resources that should be updated.
If a resource does not exist in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionBindings` | [`ActionBindings`](#chainsaw-kyverno-io-v1alpha1-ActionBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha1-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha1-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionOutputs` | [`ActionOutputs`](#chainsaw-kyverno-io-v1alpha1-ActionOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha1-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Use     {#chainsaw-kyverno-io-v1alpha1-Use}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)

<p>Use defines a reference to a step template.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `template` | `string` | :white_check_mark: |  | <p>Template references a step template.</p> |
| `with` | [`With`](#chainsaw-kyverno-io-v1alpha1-With) |  |  | <p>With defines arguments passed to the step template.</p> |

## Wait     {#chainsaw-kyverno-io-v1alpha1-Wait}

**Appears in:**
    
- [CatchFinally](#chainsaw-kyverno-io-v1alpha1-CatchFinally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Wait specifies how to perform wait operations on resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha1-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionFormat` | [`ActionFormat`](#chainsaw-kyverno-io-v1alpha1-ActionFormat) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionClusters` | [`ActionClusters`](#chainsaw-kyverno-io-v1alpha1-ActionClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObject` | [`ActionObject`](#chainsaw-kyverno-io-v1alpha1-ActionObject) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `for` | [`WaitFor`](#chainsaw-kyverno-io-v1alpha1-WaitFor) | :white_check_mark: |  | <p>WaitFor specifies the condition to wait for.</p> |

## WaitFor     {#chainsaw-kyverno-io-v1alpha1-WaitFor}

**Appears in:**
    
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>WaitFor specifies the condition to wait for.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `deletion` | [`WaitForDeletion`](#chainsaw-kyverno-io-v1alpha1-WaitForDeletion) |  |  | <p>Deletion specifies parameters for waiting on a resource's deletion.</p> |
| `condition` | [`WaitForCondition`](#chainsaw-kyverno-io-v1alpha1-WaitForCondition) |  |  | <p>Condition specifies the condition to wait for.</p> |
| `jsonPath` | [`WaitForJsonPath`](#chainsaw-kyverno-io-v1alpha1-WaitForJsonPath) |  |  | <p>JsonPath specifies the json path condition to wait for.</p> |

## WaitForCondition     {#chainsaw-kyverno-io-v1alpha1-WaitForCondition}

**Appears in:**
    
- [WaitFor](#chainsaw-kyverno-io-v1alpha1-WaitFor)

<p>WaitForCondition represents parameters for waiting on a specific condition of a resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) | :white_check_mark: |  | <p>Name defines the specific condition to wait for, e.g., "Available", "Ready".</p> |
| `value` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>Value defines the specific condition status to wait for, e.g., "True", "False".</p> |

## WaitForDeletion     {#chainsaw-kyverno-io-v1alpha1-WaitForDeletion}

**Appears in:**
    
- [WaitFor](#chainsaw-kyverno-io-v1alpha1-WaitFor)

<p>WaitForDeletion represents parameters for waiting on a resource's deletion.</p>


## WaitForJsonPath     {#chainsaw-kyverno-io-v1alpha1-WaitForJsonPath}

**Appears in:**
    
- [WaitFor](#chainsaw-kyverno-io-v1alpha1-WaitFor)

<p>WaitForJsonPath represents parameters for waiting on a json path of a resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `path` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) | :white_check_mark: |  | <p>Path defines the json path to wait for, e.g. '{.status.phase}'.</p> |
| `value` | [`Expression`](#chainsaw-kyverno-io-v1alpha1-Expression) |  |  | <p>Value defines the expected value to wait for, e.g., "Running".</p> |

## With     {#chainsaw-kyverno-io-v1alpha1-With}

**Appears in:**
    
- [Use](#chainsaw-kyverno-io-v1alpha1-Use)

<p>With defines arguments passed to step templates.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |

  