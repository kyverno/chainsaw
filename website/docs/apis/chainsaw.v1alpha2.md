---
title: chainsaw (v1alpha2)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha2
auto_generated: true
---
<p>Package v1alpha2 contains API Schema definitions for the v1alpha2 API group.</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha2-Configuration)
  
## `Configuration`     {#chainsaw-kyverno-io-v1alpha2-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `chainsaw.kyverno.io/v1alpha2` |
| `kind` | `string` | :white_check_mark: | | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec) | :white_check_mark: |  | <p>Configuration spec.</p> |

## `Cleanup`     {#chainsaw-kyverno-io-v1alpha2-Cleanup}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>Cleanup contains the cleanup configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running a test.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |

## `ConfigurationSpec`     {#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha2-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `catch` | [`[]Catch`](#chainsaw-kyverno-io-v1alpha1-Catch) |  |  | <p>Catch defines what the tests steps will execute when an error happens. This will be combined with catch handlers defined at the test and step levels.</p> |
| `cleanup` | [`Cleanup`](#chainsaw-kyverno-io-v1alpha2-Cleanup) |  |  | <p>Cleanup contains cleanup configuration.</p> |
| `clusters` | [`map[string]Cluster`](#chainsaw-kyverno-io-v1alpha1-Cluster) |  |  | <p>Clusters holds a registry to clusters to support multi-cluster tests.</p> |
| `discovery` | [`Discovery`](#chainsaw-kyverno-io-v1alpha2-Discovery) |  |  | <p>Discovery contains tests discovery configuration.</p> |
| `execution` | [`Execution`](#chainsaw-kyverno-io-v1alpha2-Execution) |  |  | <p>Execution contains tests execution configuration.</p> |
| `namespace` | [`Namespace`](#chainsaw-kyverno-io-v1alpha2-Namespace) |  |  | <p>Namespace contains properties for the namespace to use for tests.</p> |
| `deletionPropagationPolicy` | [`meta/v1.DeletionPropagation`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#deletionpropagation-v1-meta) |  |  | <p>DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation.</p> |
| `report` | [`Report`](#chainsaw-kyverno-io-v1alpha2-Report) |  |  | <p>Report contains properties for the report.</p> |
| `templating` | [`Templating`](#chainsaw-kyverno-io-v1alpha2-Templating) |  |  | <p>Templating contains the templating config.</p> |
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |

## `Discovery`     {#chainsaw-kyverno-io-v1alpha2-Discovery}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>Discovery options contain the discovery configuration used when discovering tests in folders.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `excludeTestRegex` | `string` |  |  | <p>ExcludeTestRegex is used to exclude tests based on a regular expression.</p> |
| `includeTestRegex` | `string` |  |  | <p>IncludeTestRegex is used to include tests based on a regular expression.</p> |
| `testFile` | `string` |  |  | <p>TestFile is the name of the file containing the test to run. If no extension is provided, chainsaw will try with .yaml first and .yml if needed.</p> |
| `fullName` | `bool` |  |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |

## `Execution`     {#chainsaw-kyverno-io-v1alpha2-Execution}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>Execution options determine how tests are run.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `failFast` | `bool` |  |  | <p>FailFast determines whether the test should stop upon encountering the first failure.</p> |
| `parallel` | `int` |  |  | <p>The maximum number of tests to run at once.</p> |
| `repeatCount` | `int` |  |  | <p>RepeatCount indicates how many times the tests should be executed.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |

## `Namespace`     {#chainsaw-kyverno-io-v1alpha2-Namespace}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>Namespace contains info about the namespace used for testing.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `template` | `policy/v1alpha1.Any` |  |  | <p>Template defines a template to create the test namespace.</p> |

## `Report`     {#chainsaw-kyverno-io-v1alpha2-Report}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>Report contains info about the report.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `format` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha2-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML).</p> |
| `path` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `name` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |

## `ReportFormatType`     {#chainsaw-kyverno-io-v1alpha2-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [Report](#chainsaw-kyverno-io-v1alpha2-Report)

## `Templating`     {#chainsaw-kyverno-io-v1alpha2-Templating}

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha2-ConfigurationSpec)

<p>Templating contains the templating configuration.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled determines whether resources should be considered for templating.</p> |

  