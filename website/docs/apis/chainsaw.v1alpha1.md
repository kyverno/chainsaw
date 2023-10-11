---
title: chainsaw (v1alpha1)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the v1alpha1 API group</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)
  
## `Configuration`     {#chainsaw-kyverno-io-v1alpha1-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec) | :white_check_mark: | <p>Configuration spec.</p> |

## `ConfigurationSpec`     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  | <p>Timeout per test step.</p> |
| `testDirs` | `[]string` |  | <p>Directories containing test cases to run.</p> |
| `skipDelete` | `bool` |  | <p>If set, do not delete the resources after running the tests (implies SkipClusterDelete).</p> |
| `stopOnFirstFailure` | `bool` |  | <p>StopOnFirstFailure determines whether the test should stop upon encountering the first failure.</p> |
| `parallel` | `int` | :white_check_mark: | <p>The maximum number of tests to run at once.</p> |
| `reportFormat` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha1-ReportFormatType) |  | <p>ReportFormat determines test report format (JSON|XML|nil) nil == no report. maps to report.Type, however we don't want generated.deepcopy to have reference to it.</p> |
| `reportName` | `string` |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |
| `namespace` | `string` |  | <p>Namespace defines the namespace to use for tests.</p> |
| `suppress` | `[]string` |  | <p>Suppress is used to suppress logs.</p> |
| `fullName` | `bool` |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |
| `skipTestRegex` | `string` |  | <p>SkipTestRegex is used to skip tests based on a regular expression.</p> |

## `ReportFormatType`     {#chainsaw-kyverno-io-v1alpha1-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

  