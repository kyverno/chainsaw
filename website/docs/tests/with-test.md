# `Test`s based syntax

The `Test` based approach is the more verbose and explicit syntax.

It does not rely on file naming conventions, it makes it easy to reuse files accross multiple tests, it offers the flexibility to provide additional configuration at both the test and test step level and supports all operations.

## The `Test` resource

A `Test` resource, like any Kubernetes resource has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the steps of a test and other configuration elements belonging to the test being defined:

- **Timeout**: Determines how long the test should run before being marked as failed due to a timeout.
- **Skip**: A simple flag to decide if a particular test should be ignored during the test run.
- **Concurrent**: Determines whether the test should run concurrently with other tests.
- **SkipDelete**: Determines whether the resources created by the test should be deleted after the test is executed.
- **Namespace**: Specifies whether the test should run in a random ephemeral namespace or use the given namespace.
- **Steps**: An ordered collection of test steps to be executed during the test.

The full structure of the `Test` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test).

## `Test` loading process

A `Test` resource is self contained and fully represents a test.
Therefore the loading process is straightforward, Chainsaw loads the `Test` and adds it to the collection of tests to be processed.

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
  # these timeouts applies only to the test (and test steps if not overridden)
  timeout:
    apply: 10s
    assert: 10s
    error: 10s
  skip: false
  concurrent: false
  skipDelete: false
  steps:
  # first step
  # apply a configmap to the cluster
  # the path to the configmap is relative to the folder
  # containing the test, hence allow reusing manifests
  # across multiple tests
  - spec:
      apply:
      - file: ../resources/configmap.yaml
  # second step
  # execute assert statements against existing resources
  # in the cluster
  - spec:
      assert:
      - file: ../resources/configmap-assert.yaml
  # third step
  # execute error statements against existing resources
  # in the cluster
  - spec:
      error:
      - file: ../resources/configmap-error.yaml
```
