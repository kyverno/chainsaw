# `Test`s based syntax

The `Test` based approach is the more verbose and explicit syntax.

It does not rely on file naming conventions, it makes it easy to reuse files accross multiple tests, it offers the flexibility to provide additional configuration at both the test and test step level and supports all operations.

## The `Test` resource

A `Test` resource, like any Kubernetes resource has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the steps of a test and other configuration elements belonging to the test being defined.

The full structure of the `Test` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test).

## `Test` loading process

A `Test` resource is self-contained and fully represents a test. Chainsaw loads the Test and adds it to the collection of tests to be processed. By default, Chainsaw expects the file name to match `chainsaw-test.yaml`. However, this can be customized using the `--test-file` command flag.

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

### Specifying a Custom Test File

If you have your test defined in a different file, you can specify it when running Chainsaw:

```bash
chainsaw test --test-dir . --test-file=<custom-test-file-name>.yaml
```

## Raw Resource Support

Chainsaw now includes the raw resource feature, allowing direct specification of Kubernetes resources within the test definitions. This feature offers a more streamlined approach for defining resources, especially useful for simpler test scenarios or for cases where resource definitions need to be reused or slightly modified across different tests.

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

Chainsaw support for URLs in file references for assert, apply, error and similar operations. This feature enhances the reach of Chainsaw by allowing users to reference files hosted on remote locations, such as GitHub raw URLs or other web URLs, directly within their test definitions.

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
