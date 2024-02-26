---
title: chainsaw (v1alpha1)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the v1alpha1 API group.</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)
  
## `Configuration`     {#chainsaw-kyverno-io-v1alpha1-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec) | :white_check_mark: |  | <p>Configuration spec.</p> |

## `Test`     {#chainsaw-kyverno-io-v1alpha1-Test}

<p>Test is the resource that contains a test definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `Test` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestSpec`](#chainsaw-kyverno-io-v1alpha1-TestSpec) | :white_check_mark: |  | <p>Test spec.</p> |

## `TestStep`     {#chainsaw-kyverno-io-v1alpha1-TestStep}

<p>TestStep is the resource that contains the testStep used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `TestStep` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: |  | <p>TestStep spec.</p> |

## `Apply`     {#chainsaw-kyverno-io-v1alpha1-Apply}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Apply represents a set of configurations or resources that
should be applied during testing.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the resources to be applied.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## `Assert`     {#chainsaw-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Assert represents a test condition that is expected to hold true
during the testing process.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `FileRefOrCheck` | [`FileRefOrCheck`](#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck) | :white_check_mark: | :white_check_mark: | <p>FileRefOrAssert provides a reference to the assertion.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## `Binding`     {#chainsaw-kyverno-io-v1alpha1-Binding}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Command](#chainsaw-kyverno-io-v1alpha1-Command)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)
- [Script](#chainsaw-kyverno-io-v1alpha1-Script)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Binding represents a key/value set as a binding in an executing test.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name the name of the binding.</p> |
| `value` | `policy/v1alpha1.Any` | :white_check_mark: |  | <p>Value value of the binding.</p> |

## `Catch`     {#chainsaw-kyverno-io-v1alpha1-Catch}

**Appears in:**
    
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

## `Cluster`     {#chainsaw-kyverno-io-v1alpha1-Cluster}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `kubeconfig` | `string` | :white_check_mark: |  | <p>Kubeconfig is the path to the referenced file.</p> |
| `context` | `string` |  |  | <p>Context is the name of the context to use.</p> |

## `Command`     {#chainsaw-kyverno-io-v1alpha1-Command}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Command describes a command to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `entrypoint` | `string` | :white_check_mark: |  | <p>Entrypoint is the command entry point to run.</p> |
| `args` | `[]string` |  |  | <p>Args is the command arguments.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `check` | `policy/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## `ConfigurationSpec`     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}

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
| `reportFormat` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha1-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML|nil) nil == no report. maps to report.Type, however we don't want generated.deepcopy to have reference to it.</p> |
| `reportPath` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `reportName` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |
| `namespace` | `string` |  |  | <p>Namespace defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `namespaceTemplate` | `policy/v1alpha1.Any` |  |  | <p>NamespaceTemplate defines a template to create the test namespace.</p> |
| `fullName` | `bool` |  |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |
| `excludeTestRegex` | `string` |  |  | <p>ExcludeTestRegex is used to exclude tests based on a regular expression.</p> |
| `includeTestRegex` | `string` |  |  | <p>IncludeTestRegex is used to include tests based on a regular expression.</p> |
| `repeatCount` | `int` |  |  | <p>RepeatCount indicates how many times the tests should be executed.</p> |
| `testFile` | `string` |  |  | <p>TestFile is the name of the file containing the test to run.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |
| `clusters` | [`map[string]Cluster`](#chainsaw-kyverno-io-v1alpha1-Cluster) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |

## `Create`     {#chainsaw-kyverno-io-v1alpha1-Create}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Create represents a set of resources that should be created.
If a resource already exists in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be created.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## `Delete`     {#chainsaw-kyverno-io-v1alpha1-Delete}

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
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `ref` | [`ObjectReference`](#chainsaw-kyverno-io-v1alpha1-ObjectReference) | :white_check_mark: |  | <p>ObjectReference determines objects to be deleted.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## `Describe`     {#chainsaw-kyverno-io-v1alpha1-Describe}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Describe defines how to describe resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `resource` | `string` | :white_check_mark: |  | <p>Resource type.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) | :white_check_mark: | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `showEvents` | `bool` |  |  | <p>Show Events indicates whether to include related events.</p> |

## `Error`     {#chainsaw-kyverno-io-v1alpha1-Error}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Error represents an anticipated error condition that may arise during testing.
Instead of treating such an error as a test failure, it acknowledges it as expected.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `FileRefOrCheck` | [`FileRefOrCheck`](#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck) | :white_check_mark: | :white_check_mark: | <p>FileRefOrAssert provides a reference to the expected error.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## `Events`     {#chainsaw-kyverno-io-v1alpha1-Events}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Events defines how to collect events.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) | :white_check_mark: | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## `Expectation`     {#chainsaw-kyverno-io-v1alpha1-Expectation}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)

<p>Expectation represents a check to be applied on the result of an operation
with a match filter to determine if the verification should be considered.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `match` | `policy/v1alpha1.Any` |  |  | <p>Match defines the matching statement.</p> |
| `check` | `policy/v1alpha1.Any` | :white_check_mark: |  | <p>Check defines the verification statement.</p> |

## `FileRef`     {#chainsaw-kyverno-io-v1alpha1-FileRef}

**Appears in:**
    
- [FileRefOrCheck](#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck)
- [FileRefOrResource](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource)

<p>FileRef represents a file reference.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `file` | `string` | :white_check_mark: |  | <p>File is the path to the referenced file. This can be a direct path to a file or an expression that matches multiple files, such as "manifest/*.yaml" for all YAML files within the "manifest" directory.</p> |

## `FileRefOrCheck`     {#chainsaw-kyverno-io-v1alpha1-FileRefOrCheck}

**Appears in:**
    
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)

<p>FileRefOrCheck represents a file reference or resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) |  | :white_check_mark: | <p>FileRef provides a reference to the file containing the resources to be applied.</p> |
| `resource` | `policy/v1alpha1.Any` |  |  | <p>Check provides a check used in assertions.</p> |

## `FileRefOrResource`     {#chainsaw-kyverno-io-v1alpha1-FileRefOrResource}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)
- [Patch](#chainsaw-kyverno-io-v1alpha1-Patch)

<p>FileRefOrResource represents a file reference or resource.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) |  | :white_check_mark: | <p>FileRef provides a reference to the file containing the resources to be applied.</p> |
| `resource` | [`meta/v1/unstructured.Unstructured`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#unstructured-unstructured-v1) |  |  | <p>Resource provides a resource to be applied.</p> |

## `Finally`     {#chainsaw-kyverno-io-v1alpha1-Finally}

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

## `Format`     {#chainsaw-kyverno-io-v1alpha1-Format}

(Alias of `string`)

**Appears in:**
    
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)

<p>Format determines the output format (json or yaml).</p>


## `Get`     {#chainsaw-kyverno-io-v1alpha1-Get}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Get defines how to get resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `resource` | `string` | :white_check_mark: |  | <p>Resource type.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) | :white_check_mark: | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## `ObjectLabelsSelector`     {#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector}

**Appears in:**
    
- [Describe](#chainsaw-kyverno-io-v1alpha1-Describe)
- [Events](#chainsaw-kyverno-io-v1alpha1-Events)
- [Get](#chainsaw-kyverno-io-v1alpha1-Get)
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>ObjectLabelsSelector represents a strategy to select objects.
For a single object name and namespace are used to identify the object.
For multiple objects use selector.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `namespace` | `string` |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | `string` |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |

## `ObjectReference`     {#chainsaw-kyverno-io-v1alpha1-ObjectReference}

**Appears in:**
    
- [Delete](#chainsaw-kyverno-io-v1alpha1-Delete)

<p>ObjectReference represents one or more objects with a specific apiVersion and kind.
For a single object name and namespace are used to identify the object.
For multiple objects use labels.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectSelector` | [`ObjectSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectSelector) | :white_check_mark: | :white_check_mark: | <p>ObjectSelector determines the selection process of referenced objects.</p> |
| `apiVersion` | `string` | :white_check_mark: |  | <p>API version of the referent.</p> |
| `kind` | `string` | :white_check_mark: |  | <p>Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds</p> |

## `ObjectSelector`     {#chainsaw-kyverno-io-v1alpha1-ObjectSelector}

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

## `Operation`     {#chainsaw-kyverno-io-v1alpha1-Operation}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Operation defines a single operation, only one action is permitted for a given operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the operation.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |
| `apply` | [`Apply`](#chainsaw-kyverno-io-v1alpha1-Apply) |  |  | <p>Apply represents resources that should be applied for this test step. This can include things like configuration settings or any other resources that need to be available during the test.</p> |
| `assert` | [`Assert`](#chainsaw-kyverno-io-v1alpha1-Assert) |  |  | <p>Assert represents an assertion to be made. It checks whether the conditions specified in the assertion hold true.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `create` | [`Create`](#chainsaw-kyverno-io-v1alpha1-Create) |  |  | <p>Create represents a creation operation.</p> |
| `delete` | [`Delete`](#chainsaw-kyverno-io-v1alpha1-Delete) |  |  | <p>Delete represents a deletion operation.</p> |
| `error` | [`Error`](#chainsaw-kyverno-io-v1alpha1-Error) |  |  | <p>Error represents the expected errors for this test step. If any of these errors occur, the test will consider them as expected; otherwise, they will be treated as test failures.</p> |
| `patch` | [`Patch`](#chainsaw-kyverno-io-v1alpha1-Patch) |  |  | <p>Patch represents a patch operation.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |
| `sleep` | [`Sleep`](#chainsaw-kyverno-io-v1alpha1-Sleep) |  |  | <p>Sleep defines zzzz.</p> |

## `Patch`     {#chainsaw-kyverno-io-v1alpha1-Patch}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Patch represents a set of resources that should be patched.
If a resource doesn't exist yet in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be patched.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## `PodLogs`     {#chainsaw-kyverno-io-v1alpha1-PodLogs}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>PodLogs defines how to collect pod logs.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `namespace` | `string` |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | `string` |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |
| `container` | `string` |  |  | <p>Container in pod to get logs from else --all-containers is used.</p> |
| `tail` | `int` |  |  | <p>Tail is the number of last lines to collect from pods. If omitted or zero, then the default is 10 if you use a selector, or -1 (all) if you use a pod name. This matches default behavior of `kubectl logs`.</p> |

## `ReportFormatType`     {#chainsaw-kyverno-io-v1alpha1-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

## `Script`     {#chainsaw-kyverno-io-v1alpha1-Script}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Script describes a script to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `content` | `string` |  |  | <p>Content defines a shell script (run with "sh -c ...").</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `check` | `policy/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## `Sleep`     {#chainsaw-kyverno-io-v1alpha1-Sleep}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Sleep represents a duration while nothing happens.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `duration` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Duration is the delay used for sleeping.</p> |

## `TestSpec`     {#chainsaw-kyverno-io-v1alpha1-TestSpec}

**Appears in:**
    
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)

<p>TestSpec contains the test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the test.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `skip` | `bool` |  |  | <p>Skip determines whether the test should skipped.</p> |
| `concurrent` | `bool` |  |  | <p>Concurrent determines whether the test should run concurrently with other tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the test should be deleted after the test is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `namespace` | `string` |  |  | <p>Namespace determines whether the test should run in a random ephemeral namespace or not.</p> |
| `namespaceTemplate` | `policy/v1alpha1.Any` |  |  | <p>NamespaceTemplate defines a template to create the test namespace.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `steps` | [`[]TestSpecStep`](#chainsaw-kyverno-io-v1alpha1-TestSpecStep) | :white_check_mark: |  | <p>Steps defining the test.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |

## `TestSpecStep`     {#chainsaw-kyverno-io-v1alpha1-TestSpecStep}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)

<p>TestSpecStep contains the test step definition used in a test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name of the step.</p> |
| `TestStepSpec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: | :white_check_mark: | <p>TestStepSpec of the step.</p> |

## `TestStepSpec`     {#chainsaw-kyverno-io-v1alpha1-TestStepSpec}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)
- [TestSpecStep](#chainsaw-kyverno-io-v1alpha1-TestSpecStep)

<p>TestStepSpec defines the desired state and behavior for each test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `description` | `string` |  |  | <p>Description contains a description of the test step.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |
| `bindings` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Bindings defines additional binding key/values.</p> |
| `try` | [`[]Operation`](#chainsaw-kyverno-io-v1alpha1-Operation) | :white_check_mark: |  | <p>Try defines what the step will try to execute.</p> |
| `catch` | [`[]Catch`](#chainsaw-kyverno-io-v1alpha1-Catch) |  |  | <p>Catch defines what the step will execute when an error happens.</p> |
| `finally` | [`[]Finally`](#chainsaw-kyverno-io-v1alpha1-Finally) |  |  | <p>Finally defines what the step will execute after the step is terminated.</p> |

## `Timeouts`     {#chainsaw-kyverno-io-v1alpha1-Timeouts}

**Appears in:**
    
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

## `Wait`     {#chainsaw-kyverno-io-v1alpha1-Wait}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Wait specifies how to perform wait operations on resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Specifies how long to wait for the condition to be met before timing out.</p> |
| `cluster` | `string` |  |  | <p>Cluster defines the target cluster where the wait operation will be performed (default cluster will be used if not specified).</p> |
| `resource` | `string` | :white_check_mark: |  | <p>Resource type on which the wait operation will be applied.</p> |
| `ObjectLabelsSelector` | [`ObjectLabelsSelector`](#chainsaw-kyverno-io-v1alpha1-ObjectLabelsSelector) | :white_check_mark: | :white_check_mark: | <p>ObjectLabelsSelector determines the selection process of referenced objects.</p> |
| `condition` | `string` |  |  | <p>Condition represents the specific condition to wait for. Example: "Available", "Ready", etc.</p> |
| `waitType` | [`WaitType`](#chainsaw-kyverno-io-v1alpha1-WaitType) | :white_check_mark: |  | <p>WaitType specifies the type of wait operation: "condition" or "delete".</p> |
| `allNamespaces` | `bool` |  |  | <p>AllNamespaces indicates whether to wait for resources in all namespaces.</p> |

## `WaitType`     {#chainsaw-kyverno-io-v1alpha1-WaitType}

(Alias of `string`)

**Appears in:**
    
- [Wait](#chainsaw-kyverno-io-v1alpha1-Wait)

<p>WaitType represents the type of wait operation.</p>


  