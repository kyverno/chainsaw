# Explicit approach

The explicit is a bit more verbose than the conventional one but offers far more flexibility and features:

- It does not rely on file naming conventions for operations ordering
- It encourages file reuse across tests, reducing duplication and maintenance
- It offers the flexibility to provide additional configurations like timeouts, complex logic, etc...
- It supports all [operations](../operations/index.md) without restrictions

## The Test resource

A `Test` resource, like any other Kubernetes resource, has an `apiVersion`, `kind` and `metadata` section.

It also comes with a `spec` section used to declaratively represent the test logic, steps and operations, as well as other configuration elements belonging to the test being defined.

!!! tip "Reference documentation"
    The full structure of `Test` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test).

## Example

### chainsaw-test.yaml

The `Test` below illustrates a simple test. Chainsaw will load the `Test` and steps defined in its `spec` section.

It's worth noting that:

- The test defines its own `timeouts`
- It also states that this test should not be executed in parallel with other tests
- It has multiple steps, most of them reference files that can be used in other tests if needed
- It uses an arbitrary shell script

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  # state that this test should not be executed in parallel with other tests
  concurrent: false
  # timeouts for this specific test
  timeouts:
    apply: 10s
    assert: 10s
    error: 10s
  steps:
  # step 1
  # apply a configmap to the cluster
  # the path to the configmap is relative to the folder
  # containing the test, hence allow reusing manifests
  # across multiple tests
  - try:
    - apply:
        file: ../resources/configmap.yaml
  # step 2
  # execute assert statements against existing resources
  # in the cluster
  - try:
    - assert:
        file: ../resources/configmap-assert.yaml
  # step 3
  # execute error statements against existing resources
  # in the cluster
  - try:
    - error:
        file: ../resources/configmap-error.yaml
  # step 4
  # execute an arbitrary shell script
  - try:
    - script:
        content: echo "goodbye"
```

## Conclusion

While this test is simple, it illustrates the differences with the [conventional approach](./conventional.md).

The purpose here is only to present the explicit approach and there are a lot more features to discuss, we will cover them in the next sections.
