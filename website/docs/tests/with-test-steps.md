# `TestStep`s based syntax

The `TestStep`s based syntax is more verbose than the manifests based one but offers more flexibility to influence how a test step runs.

It allows providing additional configuration per specific step, and makes it easier to reuse files accross multiple tests.

This document focuses on understanding how `TestStep`s work.

Keep in mind that a `TestStep`, like with the manifest based syntax, relies on [file naming conventions](with-manifests.md#file-naming-convention).
On the other hand it doesn't suffer the unsupported deletion limitation and can be combined with manifests based syntax when defining a step.

## The `TestStep` resource

A `TestStep` resource, like any Kubernetes resource has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the [step operations](what-is-a-test.md#operations) and other configuration elements belonging to the step being defined.

The full structure of the `TestStep` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-TestStep).

## `TestStep` loading process

In addition to the [manifests loading process](with-manifests.md#manifests-loading-process), if the manifest being loading is a `TestStep` it is directly aggregated in the steps that compose the test (it would not make sense to consider the `TestStep` a resource to be applied in the target cluster).

Manifests that are not a `TestStep` are processed as usual and participate to the test step being loaded in the same way they do with [manifests based syntax](with-manifests.md#manifests-loading-process).

Note that it's only allowed to have a **single `TestStep` resource for a given test step**.

## Raw Resource Support

Chainsaw now allows the specification of Kubernetes resources directly within the TestStep definition. This raw resource feature enhances flexibility by allowing inline resource definitions, particularly useful for concise or reusable configurations.

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
  skipDelete: true
  # these timeouts are applied per operation
  # it would not be possible to override the timeouts
  # with a manifests based approach
  timeouts:
    apply: 45s
  try:
  # apply a configmap to the cluster
  # the path to the configmap is relative to the folder
  # containing the test, hence allow reusing manifests
  # across multiple tests
  - apply:
      file: ../resources/configmap.yaml
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
  # these timeouts are applied per operation
  # it would not be possible to override the timeouts
  # with a manifests based approach
  timeouts:
    error: 20s
  try:
  # evaluate an error statement against resources
  # present in the cluster
  - error:
      file: ../resources/configmap-error.yaml
```

### 02-assert-exec.yaml

This manifest contains both an assertion and an exec operation, showing how you can mix operations in a single step.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: test-step-name-02
spec:
  try:
  - assert:
      file: ../resources/configmap-assert.yaml
  - command:
      entrypoint: "echo"
      args:
      - "Hello Chainsaw"

```

## Example Raw Resource

### 01-test-step.yaml

This `TestStep` defines a custom `timeout` and applies a `ConfigMap` directly within the step using the `raw resource` feature.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: test-step-name
spec:
  skipDelete: true
  timeouts:
    apply: 45s
  try:
  - apply:
      resource:
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: chainsaw-quick-start
        data:
          foo: bar
```

### 02-assert.yaml

This manifest contains an assertion statement for a `ConfigMap` and does not include a `TestStep` resource.

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

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: test-step-name
spec:
  timeouts:
    error: 20s
  try:
  - error:
      file: ../resources/configmap-error.yaml
```
## URL File Reference Support

Chainsaw supports URL file references in TestStep operations like apply, assert, and others. This feature allows users to specify files hosted on external sources, such as GitHub raw URLs or other accessible web URLs, directly within their test steps. 

## Example URL File Reference

### 03-url-apply-test-step.yaml

This TestStep demonstrates the use of a URL file reference in an apply operation. Chainsaw will fetch the file from the provided URL and execute the apply operation using the fetched content.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: url-apply-test-step
spec:
  try:
  - apply:
      file: https://raw.githubusercontent.com/user/repo/branch/path/to/external-configmap.yaml
```

### 04-url-assert-test-step.yaml

The manifest below contains a TestStep where the assert operation references a file hosted on an external URL. Chainsaw will load the content from the URL for the assertion.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: url-assert-test-step
spec:
  try:
  - assert:
      file: https://example.com/path/to/assert-configmap.yaml
```
## Conclusion

This test will create a config map in the first step.
The second step will both assert that the content of the config map contains the `foo: bar` data, and verify that the configmap does not contain the `lorem: ipsum` data.

The first step is made of a single `TestStep` resource, while the second step is the combination of a `TestStep` resource and a raw manifest.
