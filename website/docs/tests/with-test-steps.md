# `TestStep`s based syntax

The `TestStep`s based syntax is more verbose than the manifests based one but offers more flexibility to influence how a test step runs.

It allows providing additional configuration per specific step, and makes it easier to reuse files accross multiple tests.

This document focuses on understanding how `TestStep`s work.

Keep in mind that a `TestStep`, like with the manifest based syntax, relies on [file naming conventions](with-manifests.md#file-naming-convention).
On the other hand it doesn't suffer the unsupported deletion limitation and can be combined with manifests based syntax when defining a step.

## The `TestStep` resource

A `TestStep` resource, like any Kubernetes resource has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the [step operations](what-is-a-test.md#operations) and other configuration elements belonging to the step being defined:

- **Timeout**: Dictates how long the test step should run before being marked as failed due to a timeout.
- **Delete**: Points out resources that need to be removed before this step gets executed. It ensures the desired state of the environment before the step runs.
- **Apply**: Denotes the Kubernetes resources or configurations that should be applied at this stage.
- **Assert**: Specifies the conditions that must be true for the step to pass. Essentially, it's where you set your expectations.
- **Error**: Lists the expected errors for this step. This is vital for cases where certain errors are anticipated and should be treated as part of the expected behavior.

The full structure of the `TestStep` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-TestStep)

## `TestStep` loading process

In addition to the [manifests loading process](with-manifests.md#manifests-loading-process), if the manifest being loading is a `TestStep` it is directly aggregated in the steps that compose the test (it would not make sense to consider the `TestStep` a resource to be applied in the target cluster).

Manifests that are not a `TestStep` are processed as usual and participate to the test step being loaded in the same way they do with [manifests based syntax](with-manifests.md#manifests-loading-process).

Note that it's only allowed to have a **single `TestStep` resource for a given test step**.

## Example

### 01-test-step.yaml

The manifest below contains a `TestStep` in a file called `01-test-step.yaml`. Chainsaw will load the `TestStep` in step `01`.
The `TestStep` defines a custom `timeout` for step `01` and references a config map manifest from a relative path.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: test-step-name
spec:
  # this timeout applies only to the step
  # it would not be possible to override the timeout
  # with a manifests based approach
  timeout: 10s
  # apply a configmap to the cluster
  # the path to the configmap is relative to the folder
  # containing the test, hence allow reusing manifests
  # across multiple tests
  apply:
  - file: ../resources/configmap.yaml
```

### 02-assert.yaml

The manifest below contains an assertion statement in a file called `02-assert.yaml`. Chainsaw will associate this manifest with an assert operation in step `02`.
Note that this file doesn't contain a `TestStep` resource.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
```

### 02-error.yaml

The manifest below contains a `TestStep` in a file called `02-error.yaml`. Chainsaw will load the `TestStep` and aggregate it in step `02`.
The `TestStep` defines a custom `timeout` for step `02` and references an error statement manifest from a relative path.

This illustrates how a `TestStep` resource can be combined with raw manifests to compose a test step in Chainsaw.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: test-step-name
spec:
  # this timeout applies only to the step
  # it would not be possible to override the timeout
  # with a manifests based approach
  timeout: 20s
  # evaluate an error statement against resources
  # present in the cluster
  error:
  - file: ../resources/configmap-error.yaml
```

## Conclusion

This test will create a config map in the first step.
The second step will both assert that the content of the config map contains the `foo: bar` data, and verify that the configmap does not contain the `lorem: ipsum` data.

The first step is made of a single `TestStep` resource, while the second step is the combination of a `TestStep` resource and a raw manifest.