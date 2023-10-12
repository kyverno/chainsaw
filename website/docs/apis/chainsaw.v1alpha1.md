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


| Field | Type | Required | Description |
|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec) |  | <p>Configuration spec.</p> |

## `Test`     {#chainsaw-kyverno-io-v1alpha1-Test}

<p>Test is the resource that contains aa test definition.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | `Test` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestSpec`](#chainsaw-kyverno-io-v1alpha1-TestSpec) | :white_check_mark: | <p>Test spec.</p> |

## `TestStep`     {#chainsaw-kyverno-io-v1alpha1-TestStep}

<p>TestStep is the resource that contains the testStep used to run tests.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | `TestStep` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  | <p>Standard object's metadata.</p> |
| `spec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: | <p>TestStep spec.</p> |

## `Apply`     {#chainsaw-kyverno-io-v1alpha1-Apply}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

| Field | Type | Required | Description |
|---|---|---|---|
| `file` | `string` | :white_check_mark: | <p>File containing the manifest to be applied.</p> |

## `Assert`     {#chainsaw-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

| Field | Type | Required | Description |
|---|---|---|---|
| `file` | `string` | :white_check_mark: | <p>File containing the assertion manifest.</p> |

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

## `Error`     {#chainsaw-kyverno-io-v1alpha1-Error}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

| Field | Type | Required | Description |
|---|---|---|---|
| `file` | `string` | :white_check_mark: | <p>File containing the assertion manifest. It is expected that the assertion fails.</p> |

## `ReportFormatType`     {#chainsaw-kyverno-io-v1alpha1-ReportFormatType}

(Alias of `string`)

**Appears in:**
    
- [ConfigurationSpec](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec)

## `TestSpec`     {#chainsaw-kyverno-io-v1alpha1-TestSpec}

**Appears in:**
    
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)

<p>TestSpec contains the test spec.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `steps` | [`[]TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: | <p>Steps defining the test.</p> |

## `TestStepSpec`     {#chainsaw-kyverno-io-v1alpha1-TestStepSpec}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)

| Field | Type | Required | Description |
|---|---|---|---|
| `assert` | [`[]Assert`](#chainsaw-kyverno-io-v1alpha1-Assert) |  | *No description provided.* |
| `apply` | [`[]Apply`](#chainsaw-kyverno-io-v1alpha1-Apply) |  | *No description provided.* |
| `error` | [`[]Error`](#chainsaw-kyverno-io-v1alpha1-Error) |  | *No description provided.* |

  