# `Test` based syntax

The `Test` based syntax is more verbose than the [manifests based syntax](./manifests-based.md) but offers more flexibility and features:

- Does not rely on file naming conventions for operations ordering
- Allows to easily reuse files accross multiple tests
- Offers the flexibility to provide additional configuration at the test, step and operation level
- Supports all [operations](../operations/index.md) and [collectors](../collectors/index.md)

## The `Test` resource

A `Test` resource, like any Kubernetes resource has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the steps of a test and other configuration elements belonging to the test being defined.

!!! tip "Reference documentation"
    The full structure of the `Test` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test).

### Test steps

A `Test` is mostly made of test steps. Test steps are detailed in a [dedicated documentation](../steps/index.md).

## Example

### chainsaw-test

The manifest below contains a `Test` in a file called `chainsaw-test.yaml` (or `chainsaw-test.yml`).

Chainsaw will load the `Test` and steps defined in its `spec` section.

The `Test` uses a custom `timeouts` for the whole test. Note that `timeouts` could have been overridden in specific steps if needed.

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

### Specifying a Custom Test File

If you have your test defined in a different file, you can specify it when running Chainsaw:

```bash
chainsaw test . --test-file <custom-test-file-name>.yaml
```

If you don't provide a file extension, chainsaw will search for a file with the `.yaml` extension first and the `.yml` extension if not found.
If you provide a file extension, chainsaw will only search for a file with the extension you provided.

## Raw Resource Support

Chainsaw now includes the raw resource feature, allowing direct specification of Kubernetes resources within the test definitions.

This feature offers a more streamlined approach for defining resources, especially useful for simpler test scenarios or for cases where resource definitions need to be reused or slightly modified across different tests.

### Example Raw Resource

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test-name
spec:
  skip: false
  concurrent: false
  skipDelete: false
  timeouts:
    apply: 10s
    assert: 10s
    error: 10s
  steps:
  # first step applies a configmap directly to the cluster
  - try:
    - apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: chainsaw-quick-start
          data:
            foo: bar
  # second step executes assert statements against existing resources
  - try:
    - assert:
        file: ../resources/configmap-assert.yaml
  # third step executes error statements against existing resources
  - try:
    - error:
        file: ../resources/configmap-error.yaml
```
## URL Support for File References

Chainsaw has support for URLs in file references for assert, apply, error and similar operations.

This feature enhances the reach of Chainsaw by allowing users to reference files hosted on remote locations, such as GitHub raw URLs or other web URLs, directly within their test definitions.

### Example URL File Reference

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test-name
spec:
  skip: false
  concurrent: false
  skipDelete: false
  timeouts:
    apply: 10s
    assert: 10s
    error: 10s
  steps:
  # first step
  # apply a Kubernetes manifest hosted at a GitHub raw URL
  - try:
    - apply:
        file: https://raw.githubusercontent.com/username/repo/branch/path/to/configmap.yaml
  # second step
  # execute assert statements against existing resources
  # using a file hosted on another web URL
  - try:
    - assert:
        file: https://example.com/path/to/configmap-assert.yaml
  # third step
  # execute error statements against existing resources
  - try:
    - error:
        file: https://mywebsite.com/path/to/configmap-error.yaml
```
