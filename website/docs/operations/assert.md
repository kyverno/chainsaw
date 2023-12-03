# Assert

The `assert` operation allows you to specify conditions that should hold true for a successful test. For example, after applying certain resources, you might want to ensure that a particular pod is running or a service is accessible.

!!! info

    Assertions in Chainsaw are based on **assertion trees**.

    Assertion trees is a solution to declaratively represent complex conditions like partial array comparisons or complex operations against an incoming data structure.

    Assertion trees are compatible with standard assertions that exist in tools like KUTTL but can do a lot more.
    Please see the [assertion trees documentation](https://kyverno.github.io/kyverno-json/policies/asserts/) in kyverno-json for details.

!!! tip "Reference documentation"
    The full structure of the `Assert` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Assert).

## Usage in `Test`

Below is an example of using `assert` in a `Test` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - assert:
            file: ../resources/configmap-assert.yaml
        # ...
    ```

## Usage in `TestStep`

Below is an example of using `assert` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - assert:
          file: ../resources/configmap-assert.yaml
      # ...
    ```
