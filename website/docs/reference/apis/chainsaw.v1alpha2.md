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
    

<p>ActionCheck contains check for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `check` | `policy/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## ActionCheckRef     {#chainsaw-kyverno-io-v1alpha2-ActionCheckRef}

**Appears in:**
    

<p>ActionCheckRef contains check reference options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha2-FileRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `resource` | `policy/v1alpha1.Any` |  |  | <p>Check provides a check used in assertions.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## ActionDryRun     {#chainsaw-kyverno-io-v1alpha2-ActionDryRun}

**Appears in:**
    

<p>ActionDryRun contains dry run options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |

## ActionEnv     {#chainsaw-kyverno-io-v1alpha2-ActionEnv}

**Appears in:**
    

<p>ActionEnv contains environment options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `env` | [`[]Binding`](#chainsaw-kyverno-io-v1alpha1-Binding) |  |  | <p>Env defines additional environment variables.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |

## ActionExpectations     {#chainsaw-kyverno-io-v1alpha2-ActionExpectations}

**Appears in:**
    

<p>ActionExpectations contains expectations for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## ActionFormat     {#chainsaw-kyverno-io-v1alpha2-ActionFormat}

**Appears in:**
    

<p>ActionFormat contains format for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `format` | [`Format`](#chainsaw-kyverno-io-v1alpha1-Format) |  |  | <p>Format determines the output format (json or yaml).</p> |

## ActionObject     {#chainsaw-kyverno-io-v1alpha2-ActionObject}

**Appears in:**
    

<p>ActionObject contains object selector options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectType` | [`ObjectType`](#chainsaw-kyverno-io-v1alpha1-ObjectType) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `ActionObjectSelector` | [`ActionObjectSelector`](#chainsaw-kyverno-io-v1alpha2-ActionObjectSelector) | :white_check_mark: | :white_check_mark: | *No description provided.* |

## ActionObjectSelector     {#chainsaw-kyverno-io-v1alpha2-ActionObjectSelector}

**Appears in:**
    
- [ActionObject](#chainsaw-kyverno-io-v1alpha2-ActionObject)

<p>ActionObjectSelector contains object selector options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectName` | [`ObjectName`](#chainsaw-kyverno-io-v1alpha1-ObjectName) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |

## ActionResourceRef     {#chainsaw-kyverno-io-v1alpha2-ActionResourceRef}

**Appears in:**
    

<p>ActionResourceRef contains resource reference options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha2-FileRef) | :white_check_mark: | :white_check_mark: | *No description provided.* |
| `resource` | [`meta/v1/unstructured.Unstructured`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#unstructured-unstructured-v1) |  |  | <p>Resource provides a resource to be applied.</p> |
| `template` | `bool` |  |  | <p>Template determines whether resources should be considered for templating.</p> |

## ActionTimeout     {#chainsaw-kyverno-io-v1alpha2-ActionTimeout}

**Appears in:**
    

<p>ActionTimeout contains timeout options for an action.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |

## CleanupOptions     {#chainsaw-kyverno-io-v1alpha2-CleanupOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>CleanupOptions contains the configuration used for cleaning up resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running a test.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |

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
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |

## DeletionOptions     {#chainsaw-kyverno-io-v1alpha2-DeletionOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>DeletionOptions contains the configuration used for deleting resources.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `propagation` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>Propagation decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation.</p> |

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

## ErrorOptions     {#chainsaw-kyverno-io-v1alpha2-ErrorOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>ErrorOptions contains the global error configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `catch` | [`[]CatchFinally`](#chainsaw-kyverno-io-v1alpha1-CatchFinally) |  |  | <p>Catch defines what the tests steps will execute when an error happens. This will be combined with catch handlers defined at the test and step levels.</p> |

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

## NamespaceOptions     {#chainsaw-kyverno-io-v1alpha2-NamespaceOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)
- [TestSpec](#chainsaw-kyverno-io-v1alpha2-TestSpec)

<p>NamespaceOptions contains the configuration used to allocate a namespace for each test.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `template` | `policy/v1alpha1.Any` |  |  | <p>Template defines a template to create the test namespace.</p> |

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
| `format` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha2-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML).</p> |
| `path` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `name` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |

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
| `TestStepSpec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: | :white_check_mark: | <p>TestStepSpec of the step.</p> |

## WaitFor     {#chainsaw-kyverno-io-v1alpha2-WaitFor}

**Appears in:**
    

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

  