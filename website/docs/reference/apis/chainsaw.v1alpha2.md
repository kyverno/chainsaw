---
title: chainsaw (v1alpha2)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha2
auto_generated: true
---
<p>Package v1alpha2 contains API Schema definitions for the v1alpha2 API group.</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha2-Configuration)
- [Test](#chainsaw-kyverno-io-v1alpha2-Test)
  
## Configuration     {#chainsaw-kyverno-io-v1alpha2-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha2` |
| `kind` | `string` | :white_check_mark: | | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec) | :white_check_mark: |  | <p>Configuration spec.</p> |

## Test     {#chainsaw-kyverno-io-v1alpha2-Test}

<p>Test is the resource that contains a test definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha2` |
| `kind` | `string` | :white_check_mark: | | `Test` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestSpec`](#chainsaw-kyverno-io-v1alpha2-TestSpec) | :white_check_mark: |  | <p>Test spec.</p> |

## ActionCheck     {#chainsaw-kyverno-io-v1alpha2-ActionCheck}

**Appears in:**
    
- [Command](#chainsaw-kyverno-io-v1alpha2-Command)
- [Script](#chainsaw-kyverno-io-v1alpha2-Script)

<p>ActionCheck contains check for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `check` | `policy/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## ActionCheckRef     {#chainsaw-kyverno-io-v1alpha2-ActionCheckRef}

**Appears in:**
    
- [Assert](#chainsaw-kyverno-io-v1alpha2-Assert)
- [Error](#chainsaw-kyverno-io-v1alpha2-Error)

<p>ActionCheckRef contains check reference options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha2-FileRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `resource` | `policy/v1alpha1.Any` |  |  | <p>Check provides a check used in assertions.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## ActionDryRun     {#chainsaw-kyverno-io-v1alpha2-ActionDryRun}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha2-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha2-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha2-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha2-Update)

<p>ActionDryRun contains dry run options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |

## ActionEnv     {#chainsaw-kyverno-io-v1alpha2-ActionEnv}

**Appears in:**
    
- [Command](#chainsaw-kyverno-io-v1alpha2-Command)
- [Script](#chainsaw-kyverno-io-v1alpha2-Script)

<p>ActionEnv contains environment options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |

## ActionExpectations     {#chainsaw-kyverno-io-v1alpha2-ActionExpectations}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha2-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha2-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha2-Delete)
- [Patch](#chainsaw-kyverno-io-v1alpha2-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha2-Update)

<p>ActionExpectations contains expectations for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## ActionFormat     {#chainsaw-kyverno-io-v1alpha2-ActionFormat}

**Appears in:**
    
- [Events](#chainsaw-kyverno-io-v1alpha2-Events)
- [Get](#chainsaw-kyverno-io-v1alpha2-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha2-Wait)

<p>ActionFormat contains format for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## ActionObject     {#chainsaw-kyverno-io-v1alpha2-ActionObject}

**Appears in:**
    
- [Describe](#chainsaw-kyverno-io-v1alpha2-Describe)
- [Get](#chainsaw-kyverno-io-v1alpha2-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha2-Wait)

<p>ActionObject contains object selector options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha2-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## ActionObjectSelector     {#chainsaw-kyverno-io-v1alpha2-ActionObjectSelector}

**Appears in:**
    
- [ActionObject](#chainsaw-kyverno-io-v1alpha2-ActionObject)
- [Events](#chainsaw-kyverno-io-v1alpha2-Events)
- [PodLogs](#chainsaw-kyverno-io-v1alpha2-PodLogs)

<p>ActionObjectSelector contains object selector options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectName` | [`ObjectName`](#chainsaw-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |

## ActionResourceRef     {#chainsaw-kyverno-io-v1alpha2-ActionResourceRef}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha2-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha2-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha2-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha2-Update)

<p>ActionResourceRef contains resource reference options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha2-FileRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `resource` | [`meta/v1/unstructured.Unstructured`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#unstructured-unstructured-v1) |  |  | <p>Resource provides a resource to be applied.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## ActionTimeout     {#chainsaw-kyverno-io-v1alpha2-ActionTimeout}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha2-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha2-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha2-Command)
- [Create](#chainsaw-kyverno-io-v1alpha2-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha2-Delete)
- [Describe](#chainsaw-kyverno-io-v1alpha2-Describe)
- [Error](#chainsaw-kyverno-io-v1alpha2-Error)
- [Events](#chainsaw-kyverno-io-v1alpha2-Events)
- [Get](#chainsaw-kyverno-io-v1alpha2-Get)
- [Patch](#chainsaw-kyverno-io-v1alpha2-Patch)
- [PodLogs](#chainsaw-kyverno-io-v1alpha2-PodLogs)
- [Script](#chainsaw-kyverno-io-v1alpha2-Script)
- [Update](#chainsaw-kyverno-io-v1alpha2-Update)
- [Wait](#chainsaw-kyverno-io-v1alpha2-Wait)

<p>ActionTimeout contains timeout options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |

## Apply     {#chainsaw-kyverno-io-v1alpha2-Apply}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Apply represents a set of configurations or resources that
should be applied during testing.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha2-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha2-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha2-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Assert     {#chainsaw-kyverno-io-v1alpha2-Assert}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Assert represents a test condition that is expected to hold true
during the testing process.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionCheckRef` | [`ActionCheckRef`](#chainsaw-kyverno-io-v1alpha2-ActionCheckRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## CleanupOptions     {#chainsaw-kyverno-io-v1alpha2-CleanupOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>CleanupOptions contains the configuration used for cleaning up resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running a test.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |

## Command     {#chainsaw-kyverno-io-v1alpha2-Command}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Command describes a command to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionCheck` | [`ActionCheck`](#chainsaw-kyverno-io-v1alpha2-ActionCheck) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionEnv` | [`ActionEnv`](#chainsaw-kyverno-io-v1alpha2-ActionEnv) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `entrypoint` | `string` | :white_check_mark: |  | <p>Entrypoint is the command entry point to run.</p> |
| `args` | `[]string` |  |  | <p>Args is the command arguments.</p> |

## ConfigurationSpec     {#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha2-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `cleanup` | [`CleanupOptions`](#chainsaw-kyverno-io-v1alpha2-CleanupOptions) |  |  | <p>Cleanup contains cleanup configuration.</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `deletion` | [`DeletionOptions`](#chainsaw-kyverno-io-v1alpha2-DeletionOptions) |  |  | <p>Deletion contains the global deletion configuration.</p> |
| `discovery` | [`DiscoveryOptions`](#chainsaw-kyverno-io-v1alpha2-DiscoveryOptions) |  |  | <p>Discovery contains tests discovery configuration.</p> |
| `error` | [`ErrorOptions`](#chainsaw-kyverno-io-v1alpha2-ErrorOptions) |  |  | <p>Error contains the global error configuration.</p> |
| `execution` | [`ExecutionOptions`](#chainsaw-kyverno-io-v1alpha2-ExecutionOptions) |  |  | <p>Execution contains tests execution configuration.</p> |
| `namespace` | [`NamespaceOptions`](#chainsaw-kyverno-io-v1alpha2-NamespaceOptions) |  |  | <p>Namespace contains properties for the namespace to use for tests.</p> |
| `report` | [`ReportOptions`](#chainsaw-kyverno-io-v1alpha2-ReportOptions) |  |  | <p>Report contains properties for the report.</p> |
| `templating` | [`TemplatingOptions`](#chainsaw-kyverno-io-v1alpha2-TemplatingOptions) |  |  | <p>Templating contains the templating config.</p> |
| `timeouts` | [`DefaultTimeouts`](#chainsaw-kyverno-io-v1alpha1-DefaultTimeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |

## Create     {#chainsaw-kyverno-io-v1alpha2-Create}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Create represents a set of resources that should be created.
If a resource already exists in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha2-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha2-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha2-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Delete     {#chainsaw-kyverno-io-v1alpha2-Delete}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Delete is a reference to an object that should be deleted</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha2-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `file` | `string` |  |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |
| `ref` | [`ObjectReference`](#chainsaw-kyverno-io-v1alpha2-ObjectReference) |  |  | <p>Ref determines objects to be deleted.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in the Configuration, the Test and the TestStep.</p> |

## DeletionOptions     {#chainsaw-kyverno-io-v1alpha2-DeletionOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>DeletionOptions contains the configuration used for deleting resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `propagation` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>Propagation decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation.</p> |

## Describe     {#chainsaw-kyverno-io-v1alpha2-Describe}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Describe defines how to describe resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionObject` | [`ActionObject`](#chainsaw-kyverno-io-v1alpha2-ActionObject) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `showEvents` | `bool` |  |  | <p>Show Events indicates whether to include related events.</p> |

## DiscoveryOptions     {#chainsaw-kyverno-io-v1alpha2-DiscoveryOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>DiscoveryOptions contains the discovery configuration used when discovering tests in folders.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `excludeTestRegex` | `string` |  |  | <p>ExcludeTestRegex is used to exclude tests based on a regular expression.</p> |
| `includeTestRegex` | `string` |  |  | <p>IncludeTestRegex is used to include tests based on a regular expression.</p> |
| `testFile` | `string` |  |  | <p>TestFile is the name of the file containing the test to run. If no extension is provided, chainsaw will try with .yaml first and .yml if needed.</p> |
| `fullName` | `bool` |  |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |

## Error     {#chainsaw-kyverno-io-v1alpha2-Error}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Error represents an anticipated error condition that may arise during testing.
Instead of treating such an error as a test failure, it acknowledges it as expected.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionCheckRef` | [`ActionCheckRef`](#chainsaw-kyverno-io-v1alpha2-ActionCheckRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## ErrorOptions     {#chainsaw-kyverno-io-v1alpha2-ErrorOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>ErrorOptions contains the global error configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `catch` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Catch defines what the tests steps will execute when an error happens. This will be combined with catch handlers defined at the test and step levels.</p> |

## Events     {#chainsaw-kyverno-io-v1alpha2-Events}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Events defines how to collect events.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionFormat` | [`ActionFormat`](#chainsaw-kyverno-io-v1alpha2-ActionFormat) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha2-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## ExecutionOptions     {#chainsaw-kyverno-io-v1alpha2-ExecutionOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>ExecutionOptions determines how tests are run.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `failFast` | `bool` |  |  | <p>FailFast determines whether the test should stop upon encountering the first failure.</p> |
| `parallel` | `int` |  |  | <p>The maximum number of tests to run at once.</p> |
| `repeatCount` | `int` |  |  | <p>RepeatCount indicates how many times the tests should be executed.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |

## FileRef     {#chainsaw-kyverno-io-v1alpha2-FileRef}

**Appears in:**
    
- [ActionCheckRef](#chainsaw-kyverno-io-v1alpha2-ActionCheckRef)
- [ActionResourceRef](#chainsaw-kyverno-io-v1alpha2-ActionResourceRef)

<p>FileRef represents a file reference.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `file` | `string` | :white_check_mark: |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |

## Get     {#chainsaw-kyverno-io-v1alpha2-Get}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Get defines how to get resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionFormat` | [`ActionFormat`](#chainsaw-kyverno-io-v1alpha2-ActionFormat) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObject` | [`ActionObject`](#chainsaw-kyverno-io-v1alpha2-ActionObject) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## NamespaceOptions     {#chainsaw-kyverno-io-v1alpha2-NamespaceOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>NamespaceOptions contains the configuration used to allocate a namespace for each test.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `template` | `policy/v1alpha1.Any` |  |  | <p>Template defines a template to create the test namespace.</p> |

## ObjectReference     {#chainsaw-kyverno-io-v1alpha2-ObjectReference}

**Appears in:**
    
- [Delete](#chainsaw-kyverno-io-v1alpha2-Delete)

<p>ObjectReference represents one or more objects with a specific apiVersion and kind.
For a single object name and namespace are used to identify the object.
For multiple objects use labels.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ObjectName` | [`ObjectName`](#chainsaw-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `labelSelector` | [`meta/v1.LabelSelector`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#labelselector-v1-meta) |  |  | <p>Label selector to match objects to delete</p> |

## Operation     {#chainsaw-kyverno-io-v1alpha2-Operation}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha2-TestStepSpec)
- [TryOperation](#chainsaw-kyverno-io-v1alpha2-TryOperation)

<p>Operation defines operation elements.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `OperationAction` | [`OperationAction`](#chainsaw-kyverno-io-v1alpha2-OperationAction) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `OperationBindings` | [`OperationBindings`](#chainsaw-kyverno-io-v1alpha2-OperationBindings) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `OperationClusters` | [`OperationClusters`](#chainsaw-kyverno-io-v1alpha2-OperationClusters) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `OperationOutputs` | [`OperationOutputs`](#chainsaw-kyverno-io-v1alpha2-OperationOutputs) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `description` | `string` |  |  | <p>Description contains a description of the operation.</p> |

## OperationAction     {#chainsaw-kyverno-io-v1alpha2-OperationAction}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha2-Operation)

<p>OperationAction defines an operation action, only one action should be specified per operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apply` | [`Apply`](#chainsaw-kyverno-io-v1alpha2-Apply) |  |  | <p>Apply represents resources that should be applied for this test step. This can include things like configuration settings or any other resources that need to be available during the test.</p> |
| `assert` | [`Assert`](#chainsaw-kyverno-io-v1alpha2-Assert) |  |  | <p>Assert represents an assertion to be made. It checks whether the conditions specified in the assertion hold true.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha2-Command) |  |  | <p>Command defines a command to run.</p> |
| `create` | [`Create`](#chainsaw-kyverno-io-v1alpha2-Create) |  |  | <p>Create represents a creation operation.</p> |
| `delete` | [`Delete`](#chainsaw-kyverno-io-v1alpha2-Delete) |  |  | <p>Delete represents a deletion operation.</p> |
| `describe` | [`Describe`](#chainsaw-kyverno-io-v1alpha2-Describe) |  |  | <p>Describe determines the resource describe collector to execute.</p> |
| `error` | [`Error`](#chainsaw-kyverno-io-v1alpha2-Error) |  |  | <p>Error represents the expected errors for this test step. If any of these errors occur, the test will consider them as expected; otherwise, they will be treated as test failures.</p> |
| `events` | [`Events`](#chainsaw-kyverno-io-v1alpha2-Events) |  |  | <p>Events determines the events collector to execute.</p> |
| `get` | [`Get`](#chainsaw-kyverno-io-v1alpha2-Get) |  |  | <p>Get determines the resource get collector to execute.</p> |
| `patch` | [`Patch`](#chainsaw-kyverno-io-v1alpha2-Patch) |  |  | <p>Patch represents a patch operation.</p> |
| `podLogs` | [`PodLogs`](#chainsaw-kyverno-io-v1alpha2-PodLogs) |  |  | <p>PodLogs determines the pod logs collector to execute.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha2-Script) |  |  | <p>Script defines a script to run.</p> |
| `sleep` | [`Sleep`](#chainsaw-kyverno-io-v1alpha2-Sleep) |  |  | <p>Sleep defines zzzz.</p> |
| `update` | [`Update`](#chainsaw-kyverno-io-v1alpha2-Update) |  |  | <p>Update represents an update operation.</p> |
| `wait` | [`Wait`](#chainsaw-kyverno-io-v1alpha2-Wait) |  |  | <p>Wait determines the resource wait collector to execute.</p> |

## OperationBindings     {#chainsaw-kyverno-io-v1alpha2-OperationBindings}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha2-Operation)

<p>OperationBindings contains bindings options for an operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |

## OperationClusters     {#chainsaw-kyverno-io-v1alpha2-OperationClusters}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha2-Operation)

<p>OperationClusters contains clusters options for an operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |

## OperationOutputs     {#chainsaw-kyverno-io-v1alpha2-OperationOutputs}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha2-Operation)

<p>OperationOutputs contains outputs options for an operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |

## Patch     {#chainsaw-kyverno-io-v1alpha2-Patch}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Patch represents a set of resources that should be patched.
If a resource doesn't exist yet in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha2-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha2-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha2-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## PodLogs     {#chainsaw-kyverno-io-v1alpha2-PodLogs}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>PodLogs defines how to collect pod logs.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha2-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `container` | `string` |  |  | <p>Container in pod to get logs from else --all-containers is used.</p> |
| `tail` | `int` |  |  | <p>Tail is the number of last lines to collect from pods. If omitted or zero, then the default is 10 if you use a selector, or -1 (all) if you use a pod name. This matches default behavior of `kubectl logs`.</p> |

## ReportFormatType     {#chainsaw-kyverno-io-v1alpha2-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [ReportOptions](#chainsaw-kyverno-io-v1alpha2-ReportOptions)

## ReportOptions     {#chainsaw-kyverno-io-v1alpha2-ReportOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>ReportOptions contains the configuration used for reporting.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `format` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha2-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML|JUNIT-TEST|JUNIT-STEP|JUNIT-OPERATION).</p> |
| `path` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `name` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |

## Script     {#chainsaw-kyverno-io-v1alpha2-Script}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Script describes a script to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionCheck` | [`ActionCheck`](#chainsaw-kyverno-io-v1alpha2-ActionCheck) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionEnv` | [`ActionEnv`](#chainsaw-kyverno-io-v1alpha2-ActionEnv) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `content` | `string` |  |  | <p>Content defines a shell script (run with "sh -c ...").</p> |

## Sleep     {#chainsaw-kyverno-io-v1alpha2-Sleep}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Sleep represents a duration while nothing happens.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `duration` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Duration is the delay used for sleeping.</p> |

## TemplatingOptions     {#chainsaw-kyverno-io-v1alpha2-TemplatingOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>TemplatingOptions contains the templating configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled determines whether resources should be considered for templating.</p> |

## TestExecutionOptions     {#chainsaw-kyverno-io-v1alpha2-TestExecutionOptions}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>TestExecutionOptions determines how tests are run.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `concurrent` | `bool` |  |  | <p>Concurrent determines whether the test should run concurrently with other tests.</p> |
| `skip` | `bool` |  |  | <p>Skip determines whether the test should skipped.</p> |
| `terminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>TerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |

## TestSpec     {#chainsaw-kyverno-io-v1alpha2-TestSpec}

**Appears in:**
    
- [Test](#chainsaw-kyverno-io-v1alpha2-Test)

<p>TestSpec contains the test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `cleanup` | [`CleanupOptions`](#chainsaw-kyverno-io-v1alpha2-CleanupOptions) |  |  | <p>Cleanup contains cleanup configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `execution` | [`TestExecutionOptions`](#chainsaw-kyverno-io-v1alpha2-TestExecutionOptions) |  |  | <p>Execution contains tests execution configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `deletion` | [`DeletionOptions`](#chainsaw-kyverno-io-v1alpha2-DeletionOptions) |  |  | <p>Deletion contains the global deletion configuration.</p> |
| `description` | `string` |  |  | <p>Description contains a description of the test.</p> |
| `error` | [`ErrorOptions`](#chainsaw-kyverno-io-v1alpha2-ErrorOptions) |  |  | <p>Error contains the global error configuration.</p> |
| `namespace` | [`NamespaceOptions`](#chainsaw-kyverno-io-v1alpha2-NamespaceOptions) |  |  | <p>Namespace contains properties for the namespace to use for tests.</p> |
| `steps` | [`[]TestStep`](#chainsaw-kyverno-io-v1alpha2-TestStep) | :white_check_mark: |  | <p>Steps defining the test.</p> |
| `templating` | [`TemplatingOptions`](#chainsaw-kyverno-io-v1alpha2-TemplatingOptions) |  |  | <p>Templating contains the templating config.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.</p> |

## TestStep     {#chainsaw-kyverno-io-v1alpha2-TestStep}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>TestStep contains the test step definition used in a test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name of the step.</p> |
| `TestStepSpec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha2-TestStepSpec) | :white_check_mark: | :white_check_mark: | <p>TestStepSpec of the step.</p> |

## TestStepSpec     {#chainsaw-kyverno-io-v1alpha2-TestStepSpec}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha2-TestStep)

<p>TestStepSpec defines the desired state and behavior for each test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the test step.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in both the Configuration and the Test.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `try` | [`[]TryOperation`](#chainsaw-kyverno-io-v1alpha2-TryOperation) | :white_check_mark: |  | <p>Try defines what the step will try to execute.</p> |
| `catch` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha2-Operation) |  |  | <p>Catch defines what the step will execute when an error happens.</p> |
| `finally` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha2-Operation) |  |  | <p>Finally defines what the step will execute after the step is terminated.</p> |
| `cleanup` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha2-Operation) |  |  | <p>Cleanup defines what will be executed after the test is terminated.</p> |

## TryOperation     {#chainsaw-kyverno-io-v1alpha2-TryOperation}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha2-TestStepSpec)

<p>TryOperation defines operation elements.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `Operation` | [`Operation`](#chainsaw-kyverno-io-v1alpha2-Operation) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## Update     {#chainsaw-kyverno-io-v1alpha2-Update}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Update represents a set of resources that should be updated.
If a resource does not exist in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionDryRun` | [`ActionDryRun`](#chainsaw-kyverno-io-v1alpha2-ActionDryRun) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionExpectations` | [`ActionExpectations`](#chainsaw-kyverno-io-v1alpha2-ActionExpectations) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionResourceRef` | [`ActionResourceRef`](#chainsaw-kyverno-io-v1alpha2-ActionResourceRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## Wait     {#chainsaw-kyverno-io-v1alpha2-Wait}

**Appears in:**
    
- [OperationAction](#chainsaw-kyverno-io-v1alpha2-OperationAction)

<p>Wait specifies how to perform wait operations on resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ActionTimeout` | [`ActionTimeout`](#chainsaw-kyverno-io-v1alpha2-ActionTimeout) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionFormat` | [`ActionFormat`](#chainsaw-kyverno-io-v1alpha2-ActionFormat) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObject` | [`ActionObject`](#chainsaw-kyverno-io-v1alpha2-ActionObject) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `for` | [`WaitFor`](#chainsaw-kyverno-io-v1alpha2-WaitFor) | :white_check_mark: |  | <p>WaitFor specifies the condition to wait for.</p> |

## WaitFor     {#chainsaw-kyverno-io-v1alpha2-WaitFor}

**Appears in:**
    
- [Wait](#chainsaw-kyverno-io-v1alpha2-Wait)

<p>WaitFor specifies the condition to wait for.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `deletion` | [`WaitForDeletion`](#chainsaw-kyverno-io-v1alpha2-WaitForDeletion) |  |  | <p>Deletion specifies parameters for waiting on a resource's deletion.</p> |
| `condition` | [`WaitForCondition`](#chainsaw-kyverno-io-v1alpha2-WaitForCondition) |  |  | <p>Condition specifies the condition to wait for.</p> |
| `jsonPath` | [`WaitForJsonPath`](#chainsaw-kyverno-io-v1alpha2-WaitForJsonPath) |  |  | <p>JsonPath specifies the json path condition to wait for.</p> |

## WaitForCondition     {#chainsaw-kyverno-io-v1alpha2-WaitForCondition}

**Appears in:**
    
- [WaitFor](#chainsaw-kyverno-io-v1alpha2-WaitFor)

<p>WaitForCondition represents parameters for waiting on a specific condition of a resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name defines the specific condition to wait for, e.g., "Available", "Ready".</p> |
| `value` | `string` |  |  | <p>Value defines the specific condition status to wait for, e.g., "True", "False".</p> |

## WaitForDeletion     {#chainsaw-kyverno-io-v1alpha2-WaitForDeletion}

**Appears in:**
    
- [WaitFor](#chainsaw-kyverno-io-v1alpha2-WaitFor)

<p>WaitForDeletion represents parameters for waiting on a resource's deletion.</p>


## WaitForJsonPath     {#chainsaw-kyverno-io-v1alpha2-WaitForJsonPath}

**Appears in:**
    
- [WaitFor](#chainsaw-kyverno-io-v1alpha2-WaitFor)

<p>WaitForJsonPath represents parameters for waiting on a json path of a resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `path` | `string` | :white_check_mark: |  | <p>Path defines the json path to wait for, e.g. '{.status.phase}'.</p> |
| `value` | `string` | :white_check_mark: |  | <p>Value defines the expected value to wait for, e.g., "Running".</p> |

  