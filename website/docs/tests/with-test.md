# `Test`s based syntax

The `Test` based approach is the more verbose and explicit syntax.

It does not rely on file naming conventions, it makes it easy to reuse files accross multiple tests, it offers the flexibility to provide additional configuration at both the test and test step level and supports all operations.

## The `Test` resource

A `Test` resource, like any Kubernetes resource has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the steps of a test and other configuration elements belonging to the test being defined.

The full structure of the `Test` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test).

## `Test` loading process

A `Test` resource is self contained and fully represents a test.
Therefore the loading process is straightforward, Chainsaw loads the `Test` and adds it to the collection of tests to be processed.

!!! note

    For the time being, the `Test` based approach requires the file name to match `chainsaw-test.yaml`. This will be configurable by a command flag in the future.

## Example

### chainsaw-test.yaml

The manifest below contains a `Test` in a file called `chainsaw-test.yaml`.
Chainsaw will load the `Test` and steps defined in its `spec` section.
The `Test` defines a custom `timeout` for the whole test.
Note that the `timeout` could have been overridden in specific steps if needed.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test-name
spec:
  skip: false
  concurrent: false
  skipDelete: false
  # these timeouts are applied per operation
  timeouts:
    apply: 10s
    assert: 10s
    error: 10s
  steps:
  # first step
  # apply a configmap to the cluster
  # the path to the configmap is relative to the folder
  # containing the test, hence allow reusing manifests
  # across multiple tests
  - try:
    - apply:
        file: ../resources/configmap.yaml
  # second step
  # execute assert statements against existing resources
  # in the cluster
  - try:
    - assert:
        file: ../resources/configmap-assert.yaml
  # third step
  # execute error statements against existing resources
  # in the cluster
  - try:
    - error:
        file: ../resources/configmap-error.yaml
```
