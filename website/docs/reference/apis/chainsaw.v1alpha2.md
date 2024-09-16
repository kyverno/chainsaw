---
title: chainsaw (v1alpha2)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the v1alpha1 API group.</p>
---
title: chainsaw (v1alpha2)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha2
auto_generated: true
---
<p>Package v1alpha2 contains API Schema definitions for the v1alpha2 API group.</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)
- [Configuration](#chainsaw-kyverno-io-v1alpha2-Configuration)
- [Test](#chainsaw-kyverno-io-v1alpha2-Test)
  
## Configuration     {#chainsaw-kyverno-io-v1alpha1-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec) | :white_check_mark: |  | <p>Configuration spec.</p> |

## Test     {#chainsaw-kyverno-io-v1alpha1-Test}

<p>Test is the resource that contains a test definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `Test` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestSpec`](#chainsaw-kyverno-io-v1alpha1-TestSpec) | :white_check_mark: |  | <p>Test spec.</p> |

## Apply     {#chainsaw-kyverno-io-v1alpha1-Apply}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Apply represents a set of configurations or resources that
should be applied during testing.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the resources to be applied.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## Assert     {#chainsaw-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Assert represents a test condition that is expected to hold true
during the testing process.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `FileRefOrCheck` | [`FileRefOrCheck`](#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck) | :white_check_mark: | :white_check_mark: | <p>FileRefOrAssert provides a reference to the assertion.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## Binding     {#chainsaw-kyverno-io-v1alpha1-Binding}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [Output](#chainsaw-kyverno-io-v1alpha1-Output)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>Binding represents a key/value set as a binding in an executing test.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name the name of the binding.</p> |
| `value` | `policy/v1alpha1.Any` | :white_check_mark: |  | <p>Value value of the binding.</p> |

## Catch     {#chainsaw-kyverno-io-v1alpha1-Catch}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)
- [ErrorOptions](#chainsaw-kyverno-io-v1alpha2-ErrorOptions)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Catch defines actions to be executed on failure.</p>


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
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [PodLogs](#chainsaw-kyverno-io-v1alpha1-PodLogs)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>Clusters defines a cluster map.</p>


## Command     {#chainsaw-kyverno-io-v1alpha1-Command}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Command describes a command to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `entrypoint` | `string` | :white_check_mark: |  | <p>Entrypoint is the command entry point to run.</p> |
| `args` | `[]string` |  |  | <p>Args is the command arguments.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `check` | `policy/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## Condition     {#chainsaw-kyverno-io-v1alpha1-Condition}

**Appears in:**
    
- [For](#chainsaw-kyverno-io-v1alpha1-For)

<p>Condition represents parameters for waiting on a specific condition of a resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name defines the specific condition to wait for, e.g., "Available", "Ready".</p> |
| `value` | `string` |  |  | <p>Value defines the specific condition status to wait for, e.g., "True", "False".</p> |

## ConfigurationSpec     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running the tests (implies SkipClusterDelete).</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `failFast` | `bool` |  |  | <p>FailFast determines whether the test should stop upon encountering the first failure.</p> |
| `parallel` | `int` |  |  | <p>The maximum number of tests to run at once.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation.</p> |
| `reportFormat` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha1-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML|nil) nil == no report. maps to report.Type, however we don't want generated.deepcopy to have reference to it.</p> |
| `reportPath` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `reportName` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |
| `namespace` | `string` |  |  | <p>Namespace defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `namespaceTemplate` | `policy/v1alpha1.Any` |  |  | <p>NamespaceTemplate defines a template to create the test namespace.</p> |
| `fullName` | `bool` |  |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |
| `excludeTestRegex` | `string` |  |  | <p>ExcludeTestRegex is used to exclude tests based on a regular expression.</p> |
| `includeTestRegex` | `string` |  |  | <p>IncludeTestRegex is used to include tests based on a regular expression.</p> |
| `repeatCount` | `int` |  |  | <p>RepeatCount indicates how many times the tests should be executed.</p> |
| `testFile` | `string` |  |  | <p>TestFile is the name of the file containing the test to run. If no extension is provided, chainsaw will try with .yaml first and .yml if needed.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `catch` | [`[]Catch`](#chainsaw-kyverno-io-v1alpha1-Catch) |  |  | <p>Catch defines what the tests steps will execute when an error happens. This will be combined with catch handlers defined at the test and step levels.</p> |

## Create     {#chainsaw-kyverno-io-v1alpha1-Create}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Create represents a set of resources that should be created.
If a resource already exists in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be created.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## Delete     {#chainsaw-kyverno-io-v1alpha1-Delete}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Delete is a reference to an object that should be deleted</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `file` | `string` |  |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |
| `ref` | [`ObjectReference`](#chainsaw-kyverno-io-v1alpha1-ObjectReference) |  |  | <p>Ref determines objects to be deleted.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in the Configuration, the Test and the TestStep.</p> |

## Deletion     {#chainsaw-kyverno-io-v1alpha1-Deletion}

**Appears in:**
    
- [For](#chainsaw-kyverno-io-v1alpha1-For)

<p>Deletion represents parameters for waiting on a resource's deletion.</p>


## Describe     {#chainsaw-kyverno-io-v1alpha1-Describe}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Describe defines how to describe resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `ResourceReference` | [`ResourceReference`](#chainsaw-kyverno-io-v1alpha1-ResourceReference) | :white_check_mark: | :white_check_mark: | <p>ResourceReference referenced resource type.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) |  | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `showEvents` | `bool` |  |  | <p>Show Events indicates whether to include related events.</p> |

## Error     {#chainsaw-kyverno-io-v1alpha1-Error}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Error represents an anticipated error condition that may arise during testing.
Instead of treating such an error as a test failure, it acknowledges it as expected.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `FileRefOrCheck` | [`FileRefOrCheck`](#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck) | :white_check_mark: | :white_check_mark: | <p>FileRefOrAssert provides a reference to the expected error.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## Events     {#chainsaw-kyverno-io-v1alpha1-Events}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Events defines how to collect events.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) |  | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## Expectation     {#chainsaw-kyverno-io-v1alpha1-Expectation}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>Expectation represents a check to be applied on the result of an operation
with a match filter to determine if the verification should be considered.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `match` | `policy/v1alpha1.Any` |  |  | <p>Match defines the matching statement.</p> |
| `check` | `policy/v1alpha1.Any` | :white_check_mark: |  | <p>Check defines the verification statement.</p> |

## FileRef     {#chainsaw-kyverno-io-v1alpha1-FileRef}

**Appears in:**
    
- [FileRefOrCheck](#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck)
- [FileRefOrResource](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource)

<p>FileRef represents a file reference.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `file` | `string` | :white_check_mark: |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |

## FileRefOrCheck     {#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck}

**Appears in:**
    
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)

<p>FileRefOrCheck represents a file reference or resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) |  | :white_check_mark: | <p>FileRef provides a reference to the file containing the resources to be applied.</p> |
| `resource` | `policy/v1alpha1.Any` |  |  | <p>Check provides a check used in assertions.</p> |

## FileRefOrResource     {#chainsaw-kyverno-io-v1alpha1-FileRefOrResource}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>FileRefOrResource represents a file reference or resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) |  | :white_check_mark: | <p>FileRef provides a reference to the file containing the resources to be applied.</p> |
| `resource` | [`meta/v1/unstructured.Unstructured`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#unstructured-unstructured-v1) |  |  | <p>Resource provides a resource to be applied.</p> |

## Finally     {#chainsaw-kyverno-io-v1alpha1-Finally}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Finally defines actions to be executed at the end of a test.</p>


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

## For     {#chainsaw-kyverno-io-v1alpha1-For}

**Appears in:**
    
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>For specifies the condition to wait for.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `deletion` | [`Deletion`](#chainsaw-kyverno-io-v1alpha1-Deletion) |  |  | <p>Deletion specifies parameters for waiting on a resource's deletion.</p> |
| `condition` | [`Condition`](#chainsaw-kyverno-io-v1alpha1-Condition) |  |  | <p>Condition specifies the condition to wait for.</p> |
| `jsonPath` | [`JsonPath`](#chainsaw-kyverno-io-v1alpha1-JsonPath) |  |  | <p>JsonPath specifies the json path condition to wait for.</p> |

## Format     {#chainsaw-kyverno-io-v1alpha1-Format}

(Alias of `string`)

**Appears in:**
    
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>Format determines the output format (json or yaml).</p>


## Get     {#chainsaw-kyverno-io-v1alpha1-Get}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Get defines how to get resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `ResourceReference` | [`ResourceReference`](#chainsaw-kyverno-io-v1alpha1-ResourceReference) | :white_check_mark: | :white_check_mark: | <p>ResourceReference referenced resource type.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) |  | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## JsonPath     {#chainsaw-kyverno-io-v1alpha1-JsonPath}

**Appears in:**
    
- [For](#chainsaw-kyverno-io-v1alpha1-For)

<p>JsonPath represents parameters for waiting on a json path of a resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `path` | `string` | :white_check_mark: |  | <p>Path defines the json path to wait for, e.g. '{.status.phase}'.</p> |
| `value` | `string` | :white_check_mark: |  | <p>Value defines the expected value to wait for, e.g., "Running".</p> |

## ObjectLabelsSelector     {#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector}

**Appears in:**
    
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [PodLogs](#chainsaw-kyverno-io-v1alpha1-PodLogs)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ObjectLabelsSelector represents a strategy to select objects.
For a single object name and namespace are used to identify the object.
For multiple objects use selector.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `namespace` | `string` |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | `string` |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |

## ObjectReference     {#chainsaw-kyverno-io-v1alpha1-ObjectReference}

**Appears in:**
    
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)

<p>ObjectReference represents one or more objects with a specific apiVersion and kind.
For a single object name and namespace are used to identify the object.
For multiple objects use labels.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | <p>ObjectType determines the type of referenced objects.</p> |
| `ObjectSelector` | [`ObjectSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectSelector) | :white_check_mark: | :white_check_mark: | <p>ObjectSelector determines the selection process of referenced objects.</p> |

## ObjectSelector     {#chainsaw-kyverno-io-v1alpha1-ObjectSelector}

**Appears in:**
    
- [ObjectReference](#chainsaw-kyverno-io-v1alpha1-ObjectReference)

<p>ObjectSelector represents a strategy to select objects.
For a single object name and namespace are used to identify the object.
For multiple objects use labels.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `namespace` | `string` |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | `string` |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |
| `labels` | `map[string]string` |  |  | <p>Label selector to match objects to delete</p> |

## ObjectType     {#chainsaw-kyverno-io-v1alpha1-ObjectType}

**Appears in:**
    
- [ObjectReference](#chainsaw-kyverno-io-v1alpha1-ObjectReference)

<p>ObjectType represents a specific apiVersion and kind.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: |  | <p>API version of the referent.</p> |
| `kind` | `string` | :white_check_mark: |  | <p>Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</p> |

## Operation     {#chainsaw-kyverno-io-v1alpha1-Operation}

**Appears in:**
    
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
| `error` | [`Error`](#chainsaw-kyverno-io-v1alpha1-Error) |  |  | <p>Error represents the expected errors for this test step. If any of these errors occur, the test will consider them as expected; otherwise, they will be treated as test failures.</p> |
| `patch` | [`Patch`](#chainsaw-kyverno-io-v1alpha1-Patch) |  |  | <p>Patch represents a patch operation.</p> |
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
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [Update](#chainsaw-kyverno-io-v1alpha1-Update)

<p>Output represents an output binding with a match to determine if the binding must be considered or not.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `Binding` | [`Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) | :white_check_mark: | :white_check_mark: | <p>Binding determines the binding to create when the match succeeds.</p> |
| `match` | `policy/v1alpha1.Any` |  |  | <p>Match defines the matching statement.</p> |

## Patch     {#chainsaw-kyverno-io-v1alpha1-Patch}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Patch represents a set of resources that should be patched.
If a resource doesn't exist yet in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be patched.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## PodLogs     {#chainsaw-kyverno-io-v1alpha1-PodLogs}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>PodLogs defines how to collect pod logs.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) |  | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `container` | `string` |  |  | <p>Container in pod to get logs from else --all-containers is used.</p> |
| `tail` | `int` |  |  | <p>Tail is the number of last lines to collect from pods. If omitted or zero, then the default is 10 if you use a selector, or -1 (all) if you use a pod name. This matches default behavior of `kubectl logs`.</p> |

## ReportFormatType     {#chainsaw-kyverno-io-v1alpha1-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

## ResourceReference     {#chainsaw-kyverno-io-v1alpha1-ResourceReference}

**Appears in:**
    
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ResourceReference represents a resource (API).</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: |  | <p>API version of the referent.</p> |
| `kind` | `string` | :white_check_mark: |  | <p>Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</p> |

## Script     {#chainsaw-kyverno-io-v1alpha1-Script}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Script describes a script to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `content` | `string` |  |  | <p>Content defines a shell script (run with "sh -c ...").</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `check` | `policy/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## Sleep     {#chainsaw-kyverno-io-v1alpha1-Sleep}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Sleep represents a duration while nothing happens.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `duration` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Duration is the delay used for sleeping.</p> |

## TestSpec     {#chainsaw-kyverno-io-v1alpha1-TestSpec}

**Appears in:**
    
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)

<p>TestSpec contains the test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the test.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `skip` | `bool` |  |  | <p>Skip determines whether the test should skipped.</p> |
| `concurrent` | `bool` |  |  | <p>Concurrent determines whether the test should run concurrently with other tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the test should be deleted after the test is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `namespace` | `string` |  |  | <p>Namespace determines whether the test should run in a random ephemeral namespace or not.</p> |
| `namespaceTemplate` | `policy/v1alpha1.Any` |  |  | <p>NamespaceTemplate defines a template to create the test namespace.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `steps` | [`[]TestStep`](#chainsaw-kyverno-io-v1alpha1-TestStep) | :white_check_mark: |  | <p>Steps defining the test.</p> |
| `catch` | [`[]Catch`](#chainsaw-kyverno-io-v1alpha1-Catch) |  |  | <p>Catch defines what the steps will execute when an error happens. This will be combined with catch handlers defined at the step level.</p> |
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
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `try` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha1-Operation) | :white_check_mark: |  | <p>Try defines what the step will try to execute.</p> |
| `catch` | [`[]Catch`](#chainsaw-kyverno-io-v1alpha1-Catch) |  |  | <p>Catch defines what the step will execute when an error happens.</p> |
| `finally` | [`[]Finally`](#chainsaw-kyverno-io-v1alpha1-Finally) |  |  | <p>Finally defines what the step will execute after the step is terminated.</p> |
| `cleanup` | [`[]Finally`](#chainsaw-kyverno-io-v1alpha1-Finally) |  |  | <p>Cleanup defines what will be executed after the test is terminated.</p> |

## Timeouts     {#chainsaw-kyverno-io-v1alpha1-Timeouts}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)
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
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `outputs` | [`[]Output`](#chainsaw-kyverno-io-v1alpha1-Output) |  |  | <p>Outputs defines output bindings.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be created.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## Wait     {#chainsaw-kyverno-io-v1alpha1-Wait}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Wait specifies how to perform wait operations on resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Specifies how long to wait for the condition to be met before timing out.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster where the wait operation will be performed (default cluster will be used if not specified).</p> |
| `clusters` | [`Clusters`](#chainsaw-kyverno-io-v1alpha1-Clusters) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `ResourceReference` | [`ResourceReference`](#chainsaw-kyverno-io-v1alpha1-ResourceReference) | :white_check_mark: | :white_check_mark: | <p>ResourceReference referenced resource type.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) |  | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `for` | [`For`](#chainsaw-kyverno-io-v1alpha1-For) | :white_check_mark: |  | <p>For specifies the condition to wait for.</p> |
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

  
  
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

  