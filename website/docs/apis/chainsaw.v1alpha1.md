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
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Apply represents a set of configurations or resources that
should be applied during testing.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | <p>FileRef provides a reference to the file containing the</p> |
| `shouldFail` | `bool` |  |  | <p>ShouldFail determines whether applying the file is expected to fail.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## `Assert`     {#chainsaw-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Assert represents a test condition that is expected to hold true
during the testing process.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | <p>FileRef provides a reference to the file containing the assertion.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## `Collect`     {#chainsaw-kyverno-io-v1alpha1-Collect}

**Appears in:**
    
- [OnFailure](#chainsaw-kyverno-io-v1alpha1-OnFailure)

<p>Collect defines a set of collectors.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `podLogs` | [`PodLogs`](#chainsaw-kyverno-io-v1alpha1-PodLogs) |  |  | <p>PodLogs determines the pod logs collector to execute.</p> |
| `events` | [`Events`](#chainsaw-kyverno-io-v1alpha1-Events) |  |  | <p>Events determines the events collector to execute.</p> |

## `Command`     {#chainsaw-kyverno-io-v1alpha1-Command}

**Appears in:**
    
- [Exec](#chainsaw-kyverno-io-v1alpha1-Exec)

<p>Command describes a command to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `entrypoint` | `string` | :white_check_mark: |  | <p>Entrypoint is the command entry point to run.</p> |
| `args` | `[]string` |  |  | <p>Args is the command arguments.</p> |

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

## `Delete`     {#chainsaw-kyverno-io-v1alpha1-Delete}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Delete is a reference to an object that should be deleted</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `ObjectReference` | [`ObjectReference`](#chainsaw-kyverno-io-v1alpha1-ObjectReference) | :white_check_mark: | :white_check_mark: | <p>ObjectReference determines objects to be deleted.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## `Error`     {#chainsaw-kyverno-io-v1alpha1-Error}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Error represents an anticipated error condition that may arise during testing.
Instead of treating such an error as a test failure, it acknowledges it as expected.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `FileRef` | [`FileRef`](#chainsaw-kyverno-io-v1alpha1-FileRef) | :white_check_mark: | :white_check_mark: | <p>FileRef provides a reference to the file containing the expected error.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## `Events`     {#chainsaw-kyverno-io-v1alpha1-Events}

**Appears in:**
    
- [Collect](#chainsaw-kyverno-io-v1alpha1-Collect)

<p>Events defines how to collects events.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `namespace` | `string` |  |  | <p>Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/</p> |
| `name` | `string` |  |  | <p>Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</p> |
| `selector` | `string` |  |  | <p>Selector defines labels selector.</p> |

## `Exec`     {#chainsaw-kyverno-io-v1alpha1-Exec}

**Appears in:**
    
- [ExecOperation](#chainsaw-kyverno-io-v1alpha1-ExecOperation)
- [OnFailure](#chainsaw-kyverno-io-v1alpha1-OnFailure)

<p>Exec describes a command or script.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeout` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>Timeout for the command. Overrides the global timeout set in the Configuration.</p> |
| `command` | [`Command`](#chainsaw-kyverno-io-v1alpha1-Command) |  |  | <p>Command defines a command to run.</p> |
| `script` | [`Script`](#chainsaw-kyverno-io-v1alpha1-Script) |  |  | <p>Script defines a script to run.</p> |

## `ExecOperation`     {#chainsaw-kyverno-io-v1alpha1-ExecOperation}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>Exec describes a command and/or script operation.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `Exec` | [`Exec`](#chainsaw-kyverno-io-v1alpha1-Exec) |  | :white_check_mark: | <p>Exec defines the command and/or script.</p> |
| `skipLogOutput` | `bool` |  |  | <p>SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.</p> |
| `continueOnError` | `bool` |  |  | <p>ContinueOnError determines whether a test should continue or not in case the operation was not successful. Even if the test continues executing, it will still be reported as failed.</p> |

## `FileRef`     {#chainsaw-kyverno-io-v1alpha1-FileRef}

**Appears in:**
    
- [Apply](#chainsaw-kyverno-io-v1alpha1-Apply)
- [Assert](#chainsaw-kyverno-io-v1alpha1-Assert)
- [Error](#chainsaw-kyverno-io-v1alpha1-Error)

<p>FileRef represents a file reference.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `file` | `string` | :white_check_mark: |  | <p>File is the path to the referenced file.</p> |

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

## `OnFailure`     {#chainsaw-kyverno-io-v1alpha1-OnFailure}

**Appears in:**
    
- [TestStepSpec](#chainsaw-kyverno-io-v1alpha1-TestStepSpec)

<p>OnFailure defines actions to be executed on failure.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `collect` | [`Collect`](#chainsaw-kyverno-io-v1alpha1-Collect) |  |  | <p>Collect define the collectors to run.</p> |
| `exec` | [`Exec`](#chainsaw-kyverno-io-v1alpha1-Exec) |  |  | <p>Exec define the commands and/or scripts to run.</p> |

## `PodLogs`     {#chainsaw-kyverno-io-v1alpha1-PodLogs}

**Appears in:**
    
- [Collect](#chainsaw-kyverno-io-v1alpha1-Collect)

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
    
- [Exec](#chainsaw-kyverno-io-v1alpha1-Exec)

<p>Script describes a script to run as a part of a test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `content` | `string` |  |  | <p>Content defines a shell script (run with "sh -c ...").</p> |

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
| `spec` | [`TestStepSpec`](#chainsaw-kyverno-io-v1alpha1-TestStepSpec) | :white_check_mark: |  | <p>Spec of the step.</p> |

## `TestStepSpec`     {#chainsaw-kyverno-io-v1alpha1-TestStepSpec}

**Appears in:**
    
- [TestStep](#chainsaw-kyverno-io-v1alpha1-TestStep)
- [TestSpecStep](#chainsaw-kyverno-io-v1alpha1-TestSpecStep)

<p>TestStepSpec defines the desired state and behavior for each test step.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `timeouts` | [`Timeouts`](#chainsaw-kyverno-io-v1alpha1-Timeouts) |  |  | <p>Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.</p> |
| `skipDelete` | `bool` |  |  | <p>SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.</p> |
| `assert` | [`[]Assert`](#chainsaw-kyverno-io-v1alpha1-Assert) |  |  | <p>Assert represents the assertions to be made for this test step. It checks whether the conditions specified in each assertion hold true.</p> |
| `apply` | [`[]Apply`](#chainsaw-kyverno-io-v1alpha1-Apply) |  |  | <p>Apply lists the resources that should be applied for this test step. This can include things like configuration settings or any other resources that need to be available during the test.</p> |
| `error` | [`[]Error`](#chainsaw-kyverno-io-v1alpha1-Error) |  |  | <p>Error lists the expected errors for this test step. If any of these errors occur, the test will consider them as expected; otherwise, they will be treated as test failures.</p> |
| `delete` | [`[]Delete`](#chainsaw-kyverno-io-v1alpha1-Delete) |  |  | <p>Delete provides a list of objects that should be deleted before this test step is executed. This helps in ensuring that the environment is set up correctly before the test step runs.</p> |
| `exec` | [`[]ExecOperation`](#chainsaw-kyverno-io-v1alpha1-ExecOperation) |  |  | <p>Exec provides a list of commands and/or scripts that should be executed as a part of this test step.</p> |
| `onFailure` | [`[]OnFailure`](#chainsaw-kyverno-io-v1alpha1-OnFailure) |  |  | <p>OnFailure defines actions to be executed in case of step failure.</p> |

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

  