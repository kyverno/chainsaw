---
title: chainsaw (v1alpha2)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha2
auto_generated: true
---
<p>Package v1alpha2 contains API Schema definitions for the v1alpha2 API group.</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha2-Configuration)
  
## Configuration     {#chainsaw-kyverno-io-v1alpha2-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha2` |
| `kind` | `string` | :white_check_mark: | | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec) | :white_check_mark: |  | <p>Configuration spec.</p> |

## CleanupOptions     {#chainsaw-kyverno-io-v1alpha2-CleanupOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

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
| `timeouts` | [`DefaultTimeouts`](#chainsaw-kyverno-io-v1alpha1-DefaultTimeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |

## DeletionOptions     {#chainsaw-kyverno-io-v1alpha2-DeletionOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

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

## NamespaceOptions     {#chainsaw-kyverno-io-v1alpha2-NamespaceOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>NamespaceOptions contains the configuration used to allocate a namespace for each test.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `compiler` | `policy/v1alpha1.Compiler` |  |  | <p>Compiler defines the default compiler to use when evaluating expressions.</p> |
| `template` | [`Projection`](#chainsaw-kyverno-io-v1alpha1-Projection) |  |  | <p>Template defines a template to create the test namespace.</p> |

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

## TemplatingOptions     {#chainsaw-kyverno-io-v1alpha2-TemplatingOptions}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>TemplatingOptions contains the templating configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled determines whether resources should be considered for templating.</p> |

  