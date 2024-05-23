# Assert

The `assert` operation allows you to specify conditions that should hold true for a successful test.

For example, after applying resources, you might want to ensure that a particular pod is running or a service is accessible.

!!! info "Assertion trees"

    Assertions in Chainsaw are based on **assertion trees**.

    Assertion trees are a solution to declaratively represent complex conditions like partial array comparisons or complex operations against an incoming data structure.

    Assertion trees are compatible with standard assertions that exist in tools like KUTTL but can do a lot more.
    Please see the [assertion trees documentation](https://kyverno.github.io/kyverno-json/latest/policies/asserts/) in kyverno-json for details.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Assert` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Assert).
    - This operation supports [bindings](../tests/general/bindings.md).

## Usage examples

Below is an example of using `assert` in a `Test` resource.

!!! example "Using a specific file for assertions"

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
            file: ../resources/deployment-assert.yaml
        # ...
    ```

!!! example "Using file path expressions for assertions"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example-multi
    spec:
      steps:
      - try:
        # ...
        - assert:
            file: "../assertions/*.yaml"
        # ...
    ```

!!! example "Using an URL"

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
            file: https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/valid.yaml
        # ...
    ```

!!! example "Using an inline assertion tree"

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
            resource:
              apiVersion: v1
              kind: Deployment
              metadata:
                name: foo
              spec:
                (replicas > 3): true
        # ...
    ```
