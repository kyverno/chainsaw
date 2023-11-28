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

<p>Test is the resource that contains aa test definition.</p>


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
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be applied.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## `Assert`     {#chainsaw-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Assert represents a test condition that is expected to hold true
during the testing process.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | <p>FileRef provides a reference to the file containing the assertion.</p> |

## `Catch`     {#chainsaw-kyverno-io-v1alpha1-Catch}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Catch defines actions to be executed on failure.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `podLogs` | [`PodLogs`](#chainsaw-kyverno-io-v1alpha1-PodLogs) |  |  | <p>PodLogs determines the pod logs collector to execute.</p> |
| `events` | [`Events`](#chainsaw-kyverno-io-v1alpha1-Events) |  |  | <p>Events determines the events collector to execute.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |

## `Command`     {#chainsaw-kyverno-io-v1alpha1-Command}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Command describes a command to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `entrypoint` | `string` | :white_check_mark: |  | <p>Entrypoint is the command entry point to run.</p> |
| `args` | `[]string` |  |  | <p>Args is the command arguments.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `check` | `github.com/kyverno/kyverno-json/pkg/apis/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## `ConfigurationSpec`     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Global timeouts configuration. Applies to all tests/test steps if not overridden.</p> |
| `testDirs` | `[]string` |  |  | <p>Directories containing test cases to run.</p> |
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running the tests (implies SkipClusterDelete).</p> |
| `failFast` | `bool` |  |  | <p>FailFast determines whether the test should stop upon encountering the first failure.</p> |
| `parallel` | `int` |  |  | <p>The maximum number of tests to run at once.</p> |
| `reportFormat` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha1-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML|nil) nil == no report. maps to report.Type, however we don't want generated.deepcopy to have reference to it.</p> |
| `reportName` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |
| `namespace` | `string` |  |  | <p>Namespace defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `fullName` | `bool` |  |  | <p>FullName makes use of the full test case folder path instead of the folder name.</p> |
| `excludeTestRegex` | `string` |  |  | <p>ExcludeTestRegex is used to exclude tests based on a regular expression.</p> |
| `includeTestRegex` | `string` |  |  | <p>IncludeTestRegex is used to include tests based on a regular expression.</p> |
| `repeatCount` | `int` |  |  | <p>RepeatCount indicates how many times the tests should be executed.</p> |
| `testFile` | `string` |  |  | <p>TestFile is the name of the file containing the test to run.</p> |
| `forceTerminationGracePeriod` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.</p> |

## `Create`     {#chainsaw-kyverno-io-v1alpha1-Create}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Create represents a set of resources that should be created.
If a resource already exists in the cluster it will fail.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRefOrResource` | [`FileRefOrResource`](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource) | :white_check_mark: | :white_check_mark: | <p>FileRefOrResource provides a reference to the file containing the resources to be created.</p> |
| `dryRun` | `bool` |  |  | <p>DryRun determines whether the file should be applied in dry run mode.</p> |
| `expect` | [`[]Expectation`](#chainsaw-kyverno-io-v1alpha1-Expectation) |  |  | <p>Expect defines a list of matched checks to validate the operation outcome.</p> |

## `Delete`     {#chainsaw-kyverno-io-v1alpha1-Delete}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Delete is a reference to an object that should be deleted</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ref` | [`ObjectReference`](#chainsaw-kyverno-io-v1alpha1-ObjectReference) | :white_check_mark: |  | <p>ObjectReference determines objects to be deleted.</p> |
| `check` | `github.com/kyverno/kyverno-json/pkg/apis/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## `Error`     {#chainsaw-kyverno-io-v1alpha1-Error}

**Appears in:**
    
- [Operation](#chainsaw-kyverno-io-v1alpha1-Operation)

<p>Error represents an anticipated error condition that may arise during testing.
Instead of treating such an error as a test failure, it acknowledges it as expected.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | <p>FileRef provides a reference to the file containing the expected error.</p> |

## `Events`     {#chainsaw-kyverno-io-v1alpha1-Events}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>Events defines how to collects events.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `namespace` | `string` |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | `string` |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |

## `Expectation`     {#chainsaw-kyverno-io-v1alpha1-Expectation}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)

<p>Expectation represents a check to be applied on the result of an operation
with a match filter to determine if the verification should be considered.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `match` | `github.com/kyverno/kyverno-json/pkg/apis/v1alpha1.Any` |  |  | <p>Match defines the matching statement.</p> |
| `check` | `github.com/kyverno/kyverno-json/pkg/apis/v1alpha1.Any` | :white_check_mark: |  | <p>Match defines the matching statement.</p> |

## `FileRef`     {#chainsaw-kyverno-io-v1alpha1-FileRef}

**Appears in:**
    
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)
- [FileRefOrResource](#chainsaw-kyverno-io-v1alpha1-FileRefOrResource)

<p>FileRef represents a file reference.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `file` | `string` | :white_check_mark: |  | <p>File is the path to the referenced file.</p> |

## `FileRefOrResource`     {#chainsaw-kyverno-io-v1alpha1-FileRefOrResource}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Create](#chainsaw-kyverno-io-v1alpha1-Create)

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
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `podLogs` | [`PodLogs`](#chainsaw-kyverno-io-v1alpha1-PodLogs) |  |  | <p>PodLogs determines the pod logs collector to execute.</p> |
| `events` | [`Events`](#chainsaw-kyverno-io-v1alpha1-Events) |  |  | <p>Events determines the events collector to execute.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |

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
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the operation. Overrides the global timeout set in the Configuration.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |
| `apply` | [`Apply`](#chainsaw-kyverno-io-v1alpha1-Apply) |  |  | <p>Apply represents resources that should be applied for this test step. This can include things like configuration settings or any other resources that need to be available during the test.</p> |
| `assert` | [`Assert`](#chainsaw-kyverno-io-v1alpha1-Assert) |  |  | <p>Assert represents an assertion to be made. It checks whether the conditions specified in the assertion hold true.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `create` | [`Create`](#chainsaw-kyverno-io-v1alpha1-Create) |  |  | <p>Create represents a creation operation.</p> |
| `delete` | [`Delete`](#chainsaw-kyverno-io-v1alpha1-Delete) |  |  | <p>Delete represents a creation operation.</p> |
| `error` | [`Error`](#chainsaw-kyverno-io-v1alpha1-Error) |  |  | <p>Error represents the expected errors for this test step. If any of these errors occur, the test will consider them as expected; otherwise, they will be treated as test failures.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |

## `PodLogs`     {#chainsaw-kyverno-io-v1alpha1-PodLogs}

**Appears in:**
    
- [Catch](#chainsaw-kyverno-io-v1alpha1-Catch)
- [Finally](#chainsaw-kyverno-io-v1alpha1-Finally)

<p>PodLogs defines how to collects pod logs.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
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
| `content` | `string` |  |  | <p>Content defines a shell script (run with "sh -c ...").</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `check` | `github.com/kyverno/kyverno-json/pkg/apis/v1alpha1.Any` |  |  | <p>Check is an assertion tree to validate the operation outcome.</p> |

## `TestSpec`     {#chainsaw-kyverno-io-v1alpha1-TestSpec}

**Appears in:**
    
- [Test](#chainsaw-kyverno-io-v1alpha1-Test)

<p>TestSpec contains the test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.</p> |
| `skip` | `bool` |  |  | <p>Skip determines whether the test should skipped.</p> |
| `concurrent` | `bool` |  |  | <p>Concurrent determines whether the test should run concurrently with other tests.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the test should be deleted after the test is executed.</p> |
| `namespace` | `string` |  |  | <p>Namespace determines whether the test should run in a random ephemeral namespace or not.</p> |
| `steps` | [`[]TestSpecStep`](#chainsaw-kyverno-io-v1alpha1-TestSpecStep) | :white_check_mark: |  | <p>Steps defining the test.</p> |

## `TestSpecStep`     {#chainsaw-kyverno-io-v1alpha1-TestSpecStep}

**Appears in:**
    
- [TestSpec](#chainsaw-kyverno-io-v1alpha1-TestSpec)

<p>TestSpecStep contains the test step definition used in a test spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name of the step.</p> |
| `Spec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: | :white_check_mark: | <p>Spec of the step.</p> |

## `TestStepSpec`     {#chainsaw-kyverno-io-v1alpha1-TestStepSpec}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)
- [TestSpecStep](#chainsaw-kyverno-io-v1alpha1-TestSpecStep)

<p>TestStepSpec defines the desired state and behavior for each test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.</p> |
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
| `error` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Error defines the timeout for the error operation</p> |
| `delete` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Delete defines the timeout for the delete operation</p> |
| `cleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Cleanup defines the timeout for the cleanup operation</p> |
| `exec` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) | :white_check_mark: |  | <p>Exec defines the timeout for exec operations</p> |

  